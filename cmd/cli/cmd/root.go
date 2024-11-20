package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger *slog.Logger

var cfgFile string

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "cryptkeeper-cli",
	Short: "A CLI for CryptKeeper",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration for CryptKeeper CLI",
	Run: func(cmd *cobra.Command, args []string) {
		var url, token string

		fmt.Print("Enter CryptKeeper URL: ")
		fmt.Scan(&url)
		// fmt.Print("Enter Token: ")
		// fmt.Scan(&token)

		viper.Set("url", url)
		// viper.Set("token", token)

		if err := viper.WriteConfigAs("cryptkeeper.yaml"); err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}

		fmt.Println("Configuration saved successfully.")
		fmt.Printf("To use the token, run: export CRYPTKEEPER_TOKEN=%s\n", token)
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show-config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		token := viper.GetString("token")

		fmt.Printf("CryptKeeper URL: %s\n", url)
		fmt.Printf("Token: %s\n", token)
	},
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("cryptkeeper")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cryptkeeper.yaml)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(showConfigCmd)
}

// func main() {
// 	if err := rootCmd.Execute(); err != nil {
// 		log.Fatalf("Error: %v", err)
// 	}
// }

// var rootCmd = &cobra.Command{
// 	Use:   "cryptkeeper-cli",
// 	Short: "A CLI for managing secrets with CryptKeeper",
// }

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func init() {
// 	opts := &slog.HandlerOptions{
// 		Level:     slog.LevelInfo,
// 		AddSource: true,
// 	}

// 	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

// }
