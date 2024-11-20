package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rotateSecretExpiresAt string

var rotateSecretCmd = &cobra.Command{
	Use:   "rotate [path] [key] [value]",
	Short: "Rotate a secret",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		key := args[1]
		value := args[2]

		token := os.Getenv("CRYPTKEEPER_TOKEN")
		if token == "" {
			fmt.Println("TOKEN environment variable is required")
			os.Exit(1)
		}

		var expiresAt *time.Time
		if rotateSecretExpiresAt != "" {
			t, err := time.Parse(time.RFC3339, rotateSecretExpiresAt)
			if err != nil {
				fmt.Printf("Error parsing expires_at: %v\n", err)
				os.Exit(1)
			}
			expiresAt = &t
		}

		secret := map[string]interface{}{
			"path":       path,
			"key":        key,
			"value":      value,
			"expires_at": expiresAt,
		}

		data, err := json.Marshal(secret)
		if err != nil {
			fmt.Printf("Error marshaling secret: %v\n", err)
			os.Exit(1)
		}

		url := fmt.Sprintf("http://localhost:8000/api/secrets/rotate?path=%s&key=%s", path, key)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error rotating secret: %s\n", body)
			os.Exit(1)
		}

		fmt.Println("Secret rotated successfully")
	},
}

func init() {
	rotateSecretCmd.Flags().StringVar(&rotateSecretExpiresAt, "expires-at", "", "Expiration time in RFC3339 format")
	rootCmd.AddCommand(rotateSecretCmd)
}
