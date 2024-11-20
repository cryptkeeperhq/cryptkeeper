package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/cmd/cli/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var shareSecretExpiresAt string

var shareSecretCmd = &cobra.Command{
	Use:   "share [path] [key] [version]",
	Short: "Share secret",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		key := args[1]
		version := args[2]

		token := os.Getenv("CRYPTKEEPER_TOKEN")
		if token == "" {
			fmt.Println("TOKEN environment variable is required")
			os.Exit(1)
		}

		var expiresAt *time.Time
		if shareSecretExpiresAt != "" {
			t, err := time.Parse(time.RFC3339, shareSecretExpiresAt)
			if err != nil {
				fmt.Printf("Error parsing expires_at: %v\n", err)
				os.Exit(1)
			}
			expiresAt = &t
		}

		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(metadata), &m)
		if err != nil {
			fmt.Printf("Error marshaling Metadata: %v\n", err)
			os.Exit(1)
		}

		// pathID := utils.GetPath(path, token)

		secret := map[string]interface{}{
			"path":       path,
			"key":        key,
			"version":    utils.ToInt(version),
			"expires_at": expiresAt,
		}

		data, err := json.Marshal(secret)
		if err != nil {
			fmt.Printf("Error marshaling secret: %v\n", err)
			os.Exit(1)
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000/api/secrets/share?path=%s&key=%s", path, key), bytes.NewBuffer(data))
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
			fmt.Printf("Error creating secret: %s\n", body)
			os.Exit(1)
		}

		// fmt.Println("Secret created successfully")

		if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			os.Exit(1)
		}

		// formattedData, err := json.MarshalIndent(secret, "", "  ")
		// if err != nil {
		// 	fmt.Printf("Error marshaling secret: %v\n", err)
		// 	os.Exit(1)
		// }

		// fmt.Println(string(formattedData))

		printData := [][]string{
			// []string{"Key", secret["key"].(string)},
			// []string{"Path", secret["path"].(string)},
		}
		// printData = append(printData, []string{"Version", secret["version"].(string)})

		for k, v := range secret {
			if k == "encrypted_dek" || k == "encrypted_value" {
				continue
			}
			switch value := v.(type) {
			case map[interface{}]interface{}:
				// for k1, v1 := range value {
				// 	nestedMap := flattenMap(v1.(map[string]interface{}), k1.(string))
				// 	for nk, nv := range nestedMap {
				// 		flatMap[nk] = nv
				// 	}
				// }
			case map[string]interface{}:
				// nestedMap := flattenMap(value, key)
				// for nk, nv := range nestedMap {
				// 	flatMap[nk] = nv
				// }
			case string:
				if v != "" {
					printData = append(printData, []string{k, value})
				}
			case float64:
				printData = append(printData, []string{k, fmt.Sprintf("%v", value)})
			case int:
				printData = append(printData, []string{k, fmt.Sprintf("%v", value)})
			case int64:
				printData = append(printData, []string{k, fmt.Sprintf("%v", value)})
			default:
				// flatMap[key] = value
			}

		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Key", "Value"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetBorder(true)
		table.SetCenterSeparator("+")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		table.AppendBulk(printData)
		table.SetCaption(true, "Secret created successfully.")

		table.Render()
	},
}

func init() {
	shareSecretCmd.Flags().StringVar(&shareSecretExpiresAt, "expires-at", "", "Expiration time in RFC3339 format")
	rootCmd.AddCommand(shareSecretCmd)
}
