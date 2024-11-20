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

var getSecretHistoryCmd = &cobra.Command{
	Use:   "versions [path] [key]",
	Short: "Get all versions for a secret",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		key := args[1]

		token := os.Getenv("CRYPTKEEPER_TOKEN")
		if token == "" {
			fmt.Println("Did you login? CRYPTKEEPER_TOKEN environment variable is missing")
			os.Exit(1)
		}

		url := fmt.Sprintf("http://localhost:8000/api/secrets/history?path=%s&key=%s", path, key)
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

		var secrets []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Path", "Key", "Version", "Created At"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetBorder(true)
		table.SetCenterSeparator("+")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, v := range secrets {
			table.Append([]string{v["path"].(string), v["key"].(string), fmt.Sprintf("%v", v["version"]), v["created_at"].(string)})
		}

		table.Render()

	},
}

func init() {

	rootCmd.AddCommand(getSecretHistoryCmd)
}
