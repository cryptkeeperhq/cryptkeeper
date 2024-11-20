package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	"github.com/cryptkeeperhq/cryptkeeper/internal/crypt"
	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/go-pg/pg/v10"
)

const (
	notificationThreshold = 100 // Maximum notifications per run
)

type Scheduler struct {
	DB            *pg.DB
	Config        *config.Config
	CryptoOps     crypt.CryptographicOperations
	ticker        *time.Ticker
	stopChan      chan struct{}
	wg            sync.WaitGroup
	checkInterval time.Duration
	lastCheckedAt time.Time
}

func main() {
	// Load Config
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	// Initialize Postgres DB
	db := db.Init(config)
	defer db.Close()
	config.Logger.Debug(fmt.Sprintf("Connected to database [%s] on [%s:%d]", config.Database.Name, config.Database.Host, config.Database.Port))

	// Initialize Crypto Operations Software or HSM
	cryptOps, err := crypt.New(config, "software")
	if err != nil {
		log.Fatalf("Failed to initialize cryptographic operations: %v", err)
	}
	config.Logger.Debug("Crypto Initialized...")

	// Initialize Scheduler
	scheduler := Init(db, cryptOps, config)
	scheduler.Start()
	config.Logger.Debug("Scheduler Running...")

	// Wait for exit signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	// Graceful shutdown
	config.Logger.Debug("Shutting down server...")

	// Stop Scheduler
	scheduler.Stop()
}

func Init(db *pg.DB, c crypt.CryptographicOperations, config *config.Config) *Scheduler {
	s := &Scheduler{
		DB:            db,
		CryptoOps:     c,
		Config:        config,
		stopChan:      make(chan struct{}),
		checkInterval: 5 * time.Second, // TODO: externalize to config
	}

	// Initialize lastCheckedAt
	var metadata models.SchedulerMetadata
	err := db.Model(&metadata).First()
	if err != nil {
		metadata.LastCheckedAt = time.Now().Add(-24 * time.Hour) // Default to 24 hours ago if not set
		db.Model(&metadata).Insert()
	}

	s.lastCheckedAt = metadata.LastCheckedAt

	return s
}

func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(s.checkInterval)
	s.wg.Add(1)
	go s.run()
}

func (s *Scheduler) run() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ticker.C:
			s.Config.Logger.Debug("SCHEDULER RUNNING")
			s.checkExpiringSecrets()
			s.checkExpiredSecrets()
			s.deleteExpiredSecrets()
			s.rotateSecrets()
			s.updateLastCheckedAt()
		case <-s.stopChan:
			s.ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	s.Config.Logger.Debug("Scheduler stopped.")
}

func (s *Scheduler) updateLastCheckedAt() {
	_, err := s.DB.Model(&models.SchedulerMetadata{LastCheckedAt: time.Now()}).Column("last_checked_at").Where("id = ?", 1).Update()
	if err != nil {
		s.Config.Logger.Error("Error updating last checked time", "error", err)
	}
}

func (s *Scheduler) rotateSecrets() {
	var secrets []models.Secret

	// Subquery to get the latest version for each secret
	subquery := s.DB.Model((*models.Secret)(nil)).
		Column("key").
		ColumnExpr("MAX(version) AS max_version").
		Group("key")

	err := s.DB.Model(&secrets).Distinct().
		Column("secret.*").
		Join("JOIN (?) AS sq ON sq.key = secret.key AND sq.max_version = secret.version", subquery).
		Where("rotation_interval IS NOT NULL").
		Where("secret.updated_at > ?", s.lastCheckedAt). // Filter by last checked time
		Select()

	if err != nil {
		s.Config.Logger.Error("Error fetching secrets for rotation", "error", err)
		return
	}

	now := time.Now()
	for _, secret := range secrets {
		parsedDuration, err := time.ParseDuration(secret.RotationInterval)
		if err != nil {
			s.Config.Logger.Error("Error parsing rotation interval", "error", err, "secretID", secret.ID)
			continue
		}

		if secret.LastRotatedAt == nil || now.Sub(*secret.LastRotatedAt) >= parsedDuration {
			// Fetch the path to determine the engine type
			var path models.Path
			err := s.DB.Model(&path).Where("id = ?", secret.PathID).Select()
			if err != nil {
				s.Config.Logger.Error("Error fetching path", "error", err, "secretID", secret.ID)
				continue
			}

			// Decrypt the path key
			decryptedPathKeyHandle, err := s.CryptoOps.DecryptPathKey(path.KeyData)
			if err != nil {
				s.Config.Logger.Error("Failed to decrypt path key", "error", err, "pathID", path.ID)
				continue
			}

			// Decrypt the secret value
			decryptedValue, err := s.CryptoOps.DecryptSecretValue(secret.EncryptedDEK, secret.EncryptedValue, decryptedPathKeyHandle)
			if err != nil {
				s.Config.Logger.Error("Error decrypting secret value", "error", err, "secretID", secret.ID)
				continue
			}

			secret.Value = string(decryptedValue)

			// Generate a new value to be stored
			newSecret, err := db.WriteSecret(secret.CreatedBy, secret, s.CryptoOps)
			if err != nil {
				s.Config.Logger.Error("Error rotating secret", "error", err, "secretID", secret.ID)
				continue
			}

			message := fmt.Sprintf("Secret Rotated: %d / %s", newSecret.PathID, newSecret.Key)
			s.sendNotification(newSecret.ID, message)
		}
	}
}

func (s *Scheduler) checkExpiringSecrets() {
	var secrets []models.Secret
	err := s.DB.Model(&secrets).
		Where("expires_at IS NOT NULL").
		Where("expires_at > ?", time.Now()).
		Where("expires_at < ?", time.Now().Add(24*time.Hour)).
		Where("secret.updated_at > ?", s.lastCheckedAt). // Filter by last checked time
		Limit(notificationThreshold).                    // Limit the number of notifications
		Select()
	if err != nil {
		s.Config.Logger.Error("Error querying expiring secrets", "error", err)
		return
	}

	for _, secret := range secrets {
		message := fmt.Sprintf("Secret is expiring soon: %d %s (%d)", secret.PathID, secret.Key, secret.Version)
		s.sendNotification(secret.ID, message)
	}
}

func (s *Scheduler) checkExpiredSecrets() {
	var secrets []models.Secret
	err := s.DB.Model(&secrets).
		Where("expires_at IS NOT NULL").
		Where("expires_at <= ?", time.Now()).
		Where("secret.updated_at > ?", s.lastCheckedAt). // Filter by last checked time
		Limit(notificationThreshold).                    // Limit the number of notifications
		Select()
	if err != nil {
		s.Config.Logger.Error("Error querying expired secrets", "error", err)
		return
	}

	for _, secret := range secrets {
		message := fmt.Sprintf("Secret has expired: %d %s (%d)", secret.PathID, secret.Key, secret.Version)
		s.sendNotification(secret.ID, message)
	}
}

func (s *Scheduler) deleteExpiredSecrets() {
	now := time.Now()
	var expiredSecrets []models.Secret
	err := s.DB.Model(&expiredSecrets).
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Where("secret.updated_at > ?", s.lastCheckedAt). // Filter by last checked time
		Limit(notificationThreshold).                    // Limit the number of deletions
		Select()
	if err != nil && err != pg.ErrNoRows {
		s.Config.Logger.Error("Error querying expired secrets", "error", err)
		return
	}

	for _, secret := range expiredSecrets {
		_, err = s.DB.Model(&secret).Where("id = ?", secret.ID).Delete()
		if err != nil {
			s.Config.Logger.Error("Error deleting expired secret", "error", err, "secretID", secret.ID)
			message := fmt.Sprintf("Error deleting expired secret: %d %s (%d)", secret.PathID, secret.Key, secret.Version)
			s.sendNotification(secret.ID, message)
		} else {
			s.Config.Logger.Info("Deleted expired secret", "secretID", secret.ID, "pathID", secret.PathID)
			message := fmt.Sprintf("Deleted expired secret: %d %s (%d)", secret.PathID, secret.Key, secret.Version)
			s.sendNotification(secret.ID, message)
		}
	}
}

func (s *Scheduler) sendNotification(secretID string, message string) {
	s.Config.Logger.Info("Notification sent", "secretID", secretID, "message", message)

	notification := models.Notification{
		SecretID:  secretID,
		Message:   message,
		CreatedAt: time.Now(),
	}

	_, err := s.DB.Model(&notification).Insert()
	if err != nil {
		s.Config.Logger.Error("Error inserting notification", "error", err, "secretID", secretID)
	}
}
