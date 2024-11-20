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

var getSecretVersion string

var getSecretCmd = &cobra.Command{
	Use:   "get [path] [key]",
	Short: "Retrieve a secret",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		key := args[1]

		token := os.Getenv("CRYPTKEEPER_TOKEN")
		if token == "" {
			fmt.Println("Did you login? CRYPTKEEPER_TOKEN environment variable is missing")
			os.Exit(1)
		}

		// pathID := utils.GetPath(path, token)

		if getSecretVersion != "" {
			logger.Debug("msg", "Getting Version ", getSecretVersion)
		}

		url := fmt.Sprintf("http://localhost:8000/api/secrets/version?path=%s&key=%s&version=%s", path, key, getSecretVersion)
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

		var secret map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			os.Exit(1)
		}

		// data, err := json.MarshalIndent(secret, "", "  ")
		// if err != nil {
		// 	fmt.Printf("Error marshaling secret: %v\n", err)
		// 	os.Exit(1)
		// }

		// fmt.Println(string(data))

		printData := [][]string{
			// []string{"Key", secret["key"].(string)},
			// []string{"Path", secret["path"].(string)},
		}

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
				printData = append(printData, []string{k, value})
			case float64:
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

		table.Render()

	},
}

func init() {
	getSecretCmd.Flags().StringVar(&getSecretVersion, "version", "", "Get a specific version of the secret")
	rootCmd.AddCommand(getSecretCmd)
}
