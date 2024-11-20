package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var roleID string
var secretID string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate using AppRole",
	Long:  `Authenticate using AppRole and get a token for subsequent API calls.`,
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&roleID, "role_id", "", "Role ID for AppRole authentication")
	loginCmd.Flags().StringVar(&secretID, "secret_id", "", "Secret ID for AppRole authentication")

	loginCmd.MarkFlagRequired("role_id")
	loginCmd.MarkFlagRequired("secret_id")
}

func login() {
	url := "http://localhost:8000" + "/api/auth/role"
	payload := map[string]string{
		"role_id":   roleID,
		"secret_id": secretID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		os.Exit(1)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Error decoding response body:", err)
		os.Exit(1)
	}

	token, ok := result["token"].(string)
	if !ok {
		fmt.Println("Error: no token found in response")
		os.Exit(1)
	}

	// viper.Set("token", token)
	// if err := viper.WriteConfig(); err != nil {
	// 	fmt.Println("Error saving token to config:", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("Login successful, token generated.")
	// fmt.Println(token)

	data := [][]string{
		[]string{"Token", token},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Token", "Value"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetCaption(true, "Login successful, token generated.")
	table.Render() // Send output

}
