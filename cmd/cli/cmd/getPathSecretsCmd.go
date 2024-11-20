package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// var getSecretVersion string

var getPathSecretsCmd = &cobra.Command{
	Use:   "secrets [path]",
	Short: "Retrieve Secrets for a given path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		token := os.Getenv("CRYPTKEEPER_TOKEN")
		if token == "" {
			fmt.Println("TOKEN environment variable is required")
			os.Exit(1)
		}

		// pathID := utils.GetPath(path, token)

		if getSecretVersion != "" {
			logger.Debug("msg", "Getting Version ", getSecretVersion)
		}

		url := fmt.Sprintf("http://localhost:8000/api/secrets?path=%s", path)
		logger.Debug("msg", "url", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error getting secret: %s\n", body)
			os.Exit(1)
		}

		var paths []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&paths); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Path", "Key", "Version", "One Time?", "Multi Value", "Created At"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetBorder(true)
		table.SetCenterSeparator("+")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, v := range paths {
			table.Append([]string{v["path"].(string), v["key"].(string), fmt.Sprintf("%v", v["version"]), fmt.Sprintf("%v", v["is_one_time"]), fmt.Sprintf("%v", v["is_multi_value"]), v["created_at"].(string)})
		}

		table.Render()

	},
}

func init() {
	// getSecretCmd.Flags().StringVar(&getSecretVersion, "version", "", "Get a specific version of the secret")
	rootCmd.AddCommand(getPathSecretsCmd)
}
