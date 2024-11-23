package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var DB *pg.DB
var logger *slog.Logger

func Get() *pg.DB {
	return DB
}

func Init(config *config.Config) *pg.DB {
	if config == nil {
		log.Fatalf("Config cannot be nil")
	}
	validateConfig(config)
	logger = config.Logger

	DB = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Database.Host, config.Database.Port),
		User:     config.Database.User,
		Password: config.Database.Password,
		Database: config.Database.Name,
		PoolSize: 10,
	})

	DB.AddQueryHook(dbLogger{})
	createSchema(DB)

	logger.Debug("Database connection established and schema created.")
	return DB
}

func Close() error {
	if DB != nil {
		logger.Debug("Closing database connection.")
		return DB.Close()
	}
	return nil
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	logger.Debug(string(query))
	return nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.Path)(nil),
		(*models.Secret)(nil),
		(*models.User)(nil),
		(*models.Group)(nil),
		(*models.Certificate)(nil),
		(*models.UserGroup)(nil),
		(*models.SecretDeletion)(nil),
		(*models.Notification)(nil),
		(*models.SharedLink)(nil),
		(*models.AuditLog)(nil),
		(*models.ApprovalRequest)(nil),
		(*models.AccessLog)(nil),
		(*models.RootCA)(nil),
		(*models.SubCA)(nil),
		(*models.CertificateAuthority)(nil),
		(*models.CertificateTemplate)(nil),
		(*models.AppRole)(nil),
		(*models.Policy)(nil),
		(*models.UserPolicy)(nil),
		(*models.GroupPolicy)(nil),
		(*models.AppPolicy)(nil),
		(*models.PolicyAuditLog)(nil),
		(*models.SchedulerMetadata)(nil),
		(*models.Workflow)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			logger.Error("Failed to create table for model %T: %v", model, err)
			os.Exit(1)
		}
	}
	logger.Debug("Schema created successfully.")
	return nil
}

func validateConfig(config *config.Config) {
	if config.Database.Host == "" || config.Database.Port == 0 || config.Database.User == "" || config.Database.Password == "" || config.Database.Name == "" {
		log.Fatalf("Database configuration is incomplete.")
	}
}
