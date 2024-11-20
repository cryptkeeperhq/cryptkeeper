package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	"github.com/cryptkeeperhq/cryptkeeper/internal/crypt"
	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/messaging"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "github.com/cryptkeeperhq/cryptkeeper/internal/kafka"
	"github.com/cryptkeeperhq/cryptkeeper/internal/keycloak"
	"github.com/cryptkeeperhq/cryptkeeper/internal/messaging/channel"
	"github.com/cryptkeeperhq/cryptkeeper/internal/messaging/kafka"
	"github.com/cryptkeeperhq/cryptkeeper/internal/messaging/mhandlers"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	// "github.com/cryptkeeperhq/cryptkeeper/internal/scheduler"
)

const (
	AccessLevelList   = "list"
	AccessLevelRead   = "read"
	AccessLevelCreate = "create"
	AccessLevelUpdate = "update"
	AccessLevelDelete = "delete"
	AccessLevelOwner  = "owner"
)

var (
	policyTopic = "policy_updates"
	secretTopic = "secret_updates"
)

type Handler struct {
	DB        *pg.DB
	Config    config.Config
	CryptoOps crypt.CryptographicOperations
	// Z         zanzibar.Zanzibar
	Producer messaging.Producer
	Consumer messaging.Consumer
	KAuth    *keycloak.KAuth
	JWTKey   []byte
	// Scheduler *scheduler.Scheduler
}

func setupMessaging(config *config.Config, db *pg.DB) (messaging.Producer, messaging.Consumer) {
	useKafka := config.Kafka.Enabled

	var producer messaging.Producer
	var consumer messaging.Consumer

	handlers := map[string]messaging.MessageHandler{
		"policy_updates": &mhandlers.PolicyUpdateHandler{DB: db},
	}

	if useKafka {
		// Initialize Kafka producer and consumer
		kafkaBrokers := []string{config.Kafka.Broker}
		var err error
		producer, err = kafka.NewKafkaProducer(kafkaBrokers)
		if err != nil {
			log.Fatalf("Failed to initialize Kafka producer: %v", err)
		}

		consumer, err = kafka.NewKafkaConsumer(kafkaBrokers, handlers)
		if err != nil {
			log.Fatalf("Failed to initialize Kafka consumer: %v", err)
		}

	} else {
		// Initialize channel-based producer and consumer
		prod := channel.NewChannelProducer()
		producer = prod
		consumer = channel.NewChannelConsumer(prod, handlers)
	}

	consumer.Consume([]string{policyTopic, secretTopic})
	return producer, consumer
}

func setupZanzibar() {
	// // Initialize SpiceDB Zanzibar
	// zanzibarClient, err := zanzibar.NewZanzibar(config)
	// if err != nil {
	// 	log.Fatalf("Failed to create Zanzibar client: %v", err)
	// }
	// err = zanzibarClient.WriteSchema()
	// if err != nil {
	// 	log.Fatalf("Failed to create schema: %v", err)
	// }
	// config.Logger.Debug("Connected to Zanibar")

	// // zanzibarClient.AddUserToGroup("dev_user", "devs")
	// // err = zanzibarClient.AddPathGroupPermissions("/kv", "devs", []string{"list", "read", "create", "update", "delete", "rotate"})
	// // if err != nil {
	// // 	log.Println("Failed to create schema: %v", err)
	// // }
	// fmt.Println("+++++++++++++")
	// zanzibarClient.VerifyGroupPermissions("/kv", "devs", "list")
	// fmt.Println("+++++++++++++")
	// zanzibarClient.VerifyGroupMembership("dev_user", "devs")
	// fmt.Println("+++++++++++++")
	// zanzibarClient.VerifyPathPermissions("/kv", "dev_user", "read")

	// return
}

func Init(config *config.Config) *Handler {

	// Initialize Crypto Operations Software or HSM
	cryptOps, err := crypt.New(config, "software")
	if err != nil {
		log.Fatalf("Failed to initialize cryptographic operations: %v", err)
	}
	config.Logger.Debug("Crypto Initialized...")

	// Initialize Postgres DB
	db := db.Init(config)
	// defer db.Close()
	config.Logger.Debug(fmt.Sprintf("Connected to database [%s] on [%s:%d]", config.Database.Name, config.Database.Host, config.Database.Port))

	// Initialize Scheduler
	// scheduler := scheduler.Init(db, cryptOps, config)
	// scheduler.Start()
	// config.Logger.Debug("Scheduler Running...")

	var kAuth *keycloak.KAuth
	if config.Auth.SSOEnabled {
		kAuth, _ = keycloak.Init()
	}

	producer, consumer := setupMessaging(config, db)
	config.Logger.Debug("Messaging Initialized...")

	// Initialize API Handlers
	// h := &Handler{DB: db, Config: *config, CryptoOps: cryptOps, Z: *zanzibarClient, Producer: *producer, JWTKey: []byte(config.Server.JWTKey), KAuth: kAuth}
	h := &Handler{DB: db, Config: *config, CryptoOps: cryptOps, Producer: producer, Consumer: consumer, JWTKey: []byte(config.Server.JWTKey), KAuth: kAuth}
	return h
}

// TLSAuthMiddleware is a middleware to handle TLS authentication
func TLSAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil {
			clientCerts := r.TLS.PeerCertificates
			if len(clientCerts) > 0 {
				clientCert := clientCerts[0] // Get the first client certificate
				// Log the subject and serial number of the client certificate
				log.Printf("Client Certificate Subject: %s, Serial Number: %s", clientCert.Subject.String(), clientCert.SerialNumber.String())
			} else {
				log.Println("No client certificate provided")
				http.Error(w, "Client certificate required", http.StatusUnauthorized)
				return
			}
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func respondWithOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg": "OK"}`))
}

func respondDenied(w http.ResponseWriter) {
	http.Error(w, "Access Denied", http.StatusForbidden)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

func (h *Handler) NewHandler() *mux.Router {

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	router := r.PathPrefix("/api").Subrouter()

	// Unauthenticated Routers
	router.HandleFunc("/register", h.RegisterUser).Methods("POST")
	router.HandleFunc("/auth/keycloak", h.AuthKeycloak).Methods("POST")
	router.HandleFunc("/auth/user", h.AuthenticateUser).Methods("POST")
	router.HandleFunc("/auth/role", h.AuthenticateAppRole).Methods("POST")
	router.Handle("/access-shared-link/{linkID}", http.HandlerFunc(h.AccessSharedLink)).Methods("GET")

	// Transit Engine Router
	transitRouter := router.PathPrefix("/transit").Subrouter()
	transitRouter.Use(TLSAuthMiddleware)
	transitRouter.Handle("/keys", http.HandlerFunc(h.ListTransitKeys)).Methods("GET")
	// transitRouter.Handle("/fpe-encrypt", http.HandlerFunc(h.FPEEncrypt)).Methods("POST")
	transitRouter.Handle("/encrypt", http.HandlerFunc(h.Encrypt)).Methods("POST")
	transitRouter.Handle("/decrypt", http.HandlerFunc(h.Decrypt)).Methods("POST")
	// transitRouter.Handle("/fpe-decrypt", http.HandlerFunc(h.FPEDecrypt)).Methods("POST")
	transitRouter.Handle("/sign", http.HandlerFunc(h.Sign)).Methods("POST")
	transitRouter.Handle("/verify", http.HandlerFunc(h.Verify)).Methods("POST")
	transitRouter.Handle("/hmac", http.HandlerFunc(h.Hmac)).Methods("POST")
	transitRouter.Handle("/hmac/verify", http.HandlerFunc(h.HmacVerify)).Methods("POST")

	// Authenticated Routers
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(h.Authenticate)

	workflowRouter := authRouter.PathPrefix("/workflows").Subrouter()
	workflowRouter.HandleFunc("", http.HandlerFunc(h.GetWorkflows)).Methods("GET")
	workflowRouter.HandleFunc("", http.HandlerFunc(h.SaveOrCreateWorkflow)).Methods("POST")
	workflowRouter.Handle("/workflow/{id}", http.HandlerFunc(h.GetWorkflow)).Methods("GET")
	workflowRouter.HandleFunc("/events", h.GetEvents).Methods("GET")
	workflowRouter.HandleFunc("/execute", h.ExecuteWorkflow).Methods("POST")

	sealUnsealRouter := authRouter.PathPrefix("/admin").Subrouter()
	sealUnsealRouter.Handle("/seal/status", http.HandlerFunc(h.SealStatusHandler)).Methods("GET")
	sealUnsealRouter.Handle("/seal", http.HandlerFunc(h.SealHandler)).Methods("POST")
	sealUnsealRouter.Handle("/unseal", http.HandlerFunc(h.UnsealHandler)).Methods("POST")

	dashboardRouter := authRouter.PathPrefix("/dashboard").Subrouter()
	dashboardRouter.Handle("/summary", http.HandlerFunc(h.GetDashboardSummary)).Methods("GET")
	dashboardRouter.Handle("/recent-activity", http.HandlerFunc(h.GetRecentActivity)).Methods("GET")

	// Generic endpoints
	authRouter.Handle("/templates", http.HandlerFunc(h.GetTemplates)).Methods("GET")
	authRouter.Handle("/secrets/scan", http.HandlerFunc(h.ScanForSecrets)).Methods("POST")
	authRouter.Handle("/notifications", http.HandlerFunc(h.GetNotifications)).Methods("GET")
	authRouter.Handle("/audit-logs", http.HandlerFunc(h.GetAuditLogs)).Methods("GET")
	authRouter.Handle("/search-secrets", http.HandlerFunc(h.SearchSecrets)).Methods("GET")

	// User and Group management router
	usersRouter := authRouter.PathPrefix("/users").Subrouter()
	usersRouter.Handle("", http.HandlerFunc(h.ListUsers)).Methods("GET")
	usersRouter.Handle("", http.HandlerFunc(h.CreateUser)).Methods("POST")
	usersRouter.Handle("/{userID}/groups", http.HandlerFunc(h.ListUserGroups)).Methods("GET")

	groupsRouter := authRouter.PathPrefix("/groups").Subrouter()
	groupsRouter.Handle("", http.HandlerFunc(h.ListGroups)).Methods("GET")
	groupsRouter.Handle("", http.HandlerFunc(h.CreateGroup)).Methods("POST")
	groupsRouter.Handle("/{groupID}/users", http.HandlerFunc(h.ListGroupUsers)).Methods("GET")
	groupsRouter.Handle("/add_user", http.HandlerFunc(h.AddUserToGroup)).Methods("POST")
	groupsRouter.Handle("/remove_user", http.HandlerFunc(h.RemoveUserFromGroup)).Methods("POST")

	// Other APIs
	userRouter := authRouter.PathPrefix("/user").Subrouter()
	userRouter.Handle("/paths", http.HandlerFunc(h.ListUserPaths)).Methods("GET")

	// Path Router
	pathsRouter := authRouter.PathPrefix("/paths").Subrouter()
	pathsRouter.Handle("", http.HandlerFunc(h.ListAllPaths)).Methods("GET")
	pathsRouter.Handle("", http.HandlerFunc(h.CreatePath)).Methods("POST")
	pathsRouter.Handle("", http.HandlerFunc(h.UpdatePath)).Methods("PUT")
	pathsRouter.Handle("/{pathID}", http.HandlerFunc(h.GetPath)).Methods("GET")
	pathsRouter.Handle("/{pathID}/permissions", http.HandlerFunc(h.GetPathPermissions)).Methods("GET")
	pathsRouter.Handle("/{pathID}/policy", http.HandlerFunc(h.GetPathPolicy)).Methods("GET")
	pathsRouter.Handle("/{pathID}/deleted", http.HandlerFunc(h.GetDeletedSecrets)).Methods("GET")
	// secretsRouter.Handle("/restore", h.CheckPermission(AccessLevelDelete, http.HandlerFunc(h.RestoreDeletedSecret))).Methods("POST")
	pathsRouter.Handle("/{pathID}/deleted/{secretID}/restore", http.HandlerFunc(h.RestoreDeletedSecret)).Methods("POST")

	// Secrets Router
	secretsRouter := authRouter.PathPrefix("/secrets").Subrouter()
	secretsRouter.Handle("", h.CheckPermission(AccessLevelCreate, http.HandlerFunc(h.CreateSecret))).Methods("POST")
	secretsRouter.Handle("", h.CheckPermission(AccessLevelList, http.HandlerFunc(h.GetSecrets))).Methods("GET")
	secretsRouter.Handle("/secret/{id}", http.HandlerFunc(h.GetSecret)).Methods("GET")
	// secretsRouter.Handle("/accesses", http.HandlerFunc(h.GetSecretAccesses)).Methods("GET")
	secretsRouter.Handle("/version", h.CheckPermission(AccessLevelRead, http.HandlerFunc(h.GetSecretVersion))).Methods("GET")
	secretsRouter.Handle("/history", h.CheckPermission(AccessLevelList, http.HandlerFunc(h.GetSecretHistory))).Methods("GET")
	secretsRouter.Handle("/lineage", h.CheckPermission(AccessLevelList, http.HandlerFunc(h.GetSecretLineage))).Methods("GET")

	secretsRouter.Handle("/delete", h.CheckPermission(AccessLevelDelete, http.HandlerFunc(h.DeleteSecret))).Methods("DELETE")
	secretsRouter.Handle("/rotate", h.CheckPermission(AccessLevelCreate, http.HandlerFunc(h.RotateSecret))).Methods("POST")
	secretsRouter.Handle("/metadata", h.CheckPermission(AccessLevelUpdate, http.HandlerFunc(h.UpdateSecretMetadata))).Methods("PUT")
	//TODO: This fails as we don't have few details in deletion table
	// secretsRouter.Handle("/restore", h.CheckPermission(AccessLevelDelete, http.HandlerFunc(h.RestoreDeletedSecret))).Methods("POST")
	// secretsRouter.Handle("/assign_access", h.CheckPermission(AccessLevelUpdate, http.HandlerFunc(h.AassignSecretAccess))).Methods("POST")
	secretsRouter.Handle("/share", h.CheckPermission(AccessLevelRead, http.HandlerFunc(h.CreateSharedLink))).Methods("POST")

	// Secrets Approval Router
	approvalRequestRouter := authRouter.PathPrefix("/approval-requests").Subrouter()
	approvalRequestRouter.Handle("", http.HandlerFunc(h.CreateApprovalRequest)).Methods("POST")
	approvalRequestRouter.Handle("", http.HandlerFunc(h.ListApprovalRequests)).Methods("GET")
	approvalRequestRouter.Handle("/approve", h.CheckPermission(AccessLevelCreate, http.HandlerFunc(h.ApproveRequest))).Methods("POST")
	approvalRequestRouter.Handle("/reject", h.CheckPermission(AccessLevelCreate, http.HandlerFunc(h.RejectRequest))).Methods("POST")

	// Path Policies Router
	policiesRouter := authRouter.PathPrefix("/policies").Subrouter()
	policiesRouter.Handle("", http.HandlerFunc(h.GetPolicies)).Methods("GET")
	policiesRouter.Handle("", http.HandlerFunc(h.SavePolicy)).Methods("POST")
	policiesRouter.HandleFunc("/{id}", h.DeletePolicy).Methods("DELETE")
	policiesRouter.HandleFunc("/audit-logs", h.GetPolicyAuditLogs).Methods("GET")

	// PKI Router
	pkiRouter := authRouter.PathPrefix("/pki").Subrouter()
	// Get called from Secret Details page
	pkiRouter.Handle("/download-certificate", http.HandlerFunc(h.DownloadCertificate)).Methods("GET")
	pkiRouter.Handle("/download-ca", http.HandlerFunc(h.DownloadCA)).Methods("GET")

	// PKI Admin page
	pkiRouter.Handle("/ca", http.HandlerFunc(h.getCAs)).Methods("GET")
	pkiRouter.Handle("/ca", http.HandlerFunc(h.addCA)).Methods("POST")
	pkiRouter.Handle("/template", http.HandlerFunc(h.getTemplates)).Methods("GET")
	pkiRouter.Handle("/template", http.HandlerFunc(h.addTemplate)).Methods("POST")
	// pkiRouter.Handle("/request_certificate", http.HandlerFunc(h.issueCertificate)).Methods("POST")

	// Get called from Admin PKI page
	// pkiRouter.Handle("/download_certificate", http.HandlerFunc(h.downloadCertificate)).Methods("GET")
	// pkiRouter.Handle("/download_private_key", http.HandlerFunc(h.downloadPrivateKey)).Methods("GET")

	// App Roles Router
	appRolesRouter := authRouter.PathPrefix("/approles").Subrouter()
	appRolesRouter.Handle("", http.HandlerFunc(h.GetAppRoles)).Methods("GET")
	appRolesRouter.Handle("", http.HandlerFunc(h.CreateAppRole)).Methods("POST")

	// router.Handle("/assign-policy", utils.Authenticate(http.HandlerFunc(h.AssignPolicyToPath))).Methods("POST")

	return r
}

// TODO: change this to userName and make things easier
func (h *Handler) logAction(username string, action string, secretID *string, details map[string]interface{}) error {
	logEntry := models.AuditLog{
		Username:  username,
		Action:    action,
		SecretID:  secretID,
		Details:   details,
		Timestamp: time.Now(),
	}
	_, err := h.DB.Model(&logEntry).Insert()
	return err
}
