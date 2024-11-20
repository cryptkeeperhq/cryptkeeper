package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host   string
		Port   int
		JWTKey string `mapstructure:"jwt_key"`
	}

	TLS struct {
		Enabled  bool
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
		CaFile   string `mapstructure:"ca"`
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Security struct {
		MasterKey string `mapstructure:"master_key"`
	}

	Zanzibar struct {
		ApiKey   string `mapstructure:"api_key"`
		Endpoint string
	}

	Auth struct {
		SSOEnabled bool `mapstructure:"sso_enabled"`
	}

	Kafka struct {
		Broker  string
		Enabled bool `mapstructure:"enabled"`
	}

	Sync struct {
		Platform string
		Vault    struct {
			Address string
			Token   string
		}
		AWS struct {
			Region          string
			AccessKeyID     string
			SecretAccessKey string
		}
		Azure struct {
			TenantID     string
			ClientID     string
			ClientSecret string
			KeyVaultName string
		}
	}

	Logger *slog.Logger
}

func LoadConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("/etc/cryptkeeper/")
	viper.AddConfigPath("$HOME/.cryptkeeper")

	// Load environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// If the config file is not found, continue without an error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, err
		}
		// Otherwise, proceed with environment variables
	}

	// Set default values from config file
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Overwrite with environment variables if they are set
	if host := viper.GetString("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := viper.GetInt("DB_PORT"); port != 0 {
		config.Database.Port = port
	}
	if user := viper.GetString("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := viper.GetString("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if name := viper.GetString("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if sslmode := viper.GetString("DB_SSLMODE"); sslmode != "" {
		config.Database.SSLMode = sslmode
	}

	if masterKey := viper.GetString("MASTER_KEY"); masterKey != "" {
		config.Security.MasterKey = masterKey
	}

	if value := viper.GetString("SPICEDB_ENDPOINT"); value != "" {
		config.Zanzibar.Endpoint = value
	}

	if value := viper.GetString("SPICEDB_API_KEY"); value != "" {
		config.Zanzibar.ApiKey = value
	}

	if value := viper.GetString("KAFKA_BROKER"); value != "" {
		config.Kafka.Broker = value
	}

	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	config.Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

	return &config, nil
}
