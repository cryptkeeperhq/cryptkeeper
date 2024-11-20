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

var doTransitCmd = &cobra.Command{
	Use:   "transit [operation] [key] [data]",
	Short: "Perform a transit operation (encrypt, decrypt, hmac, sign, verify)",
	// Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		operation := args[0]
		key := args[1]
		inputData := args[2]

		token := os.Getenv("CRYPTKEEPER_TOKEN")
		if token == "" {
			fmt.Println("Did you login? CRYPTKEEPER_TOKEN environment variable is missing")
			os.Exit(1)
		}

		// Set up request URL and payload based on operation
		var url string
		request := map[string]interface{}{
			"key_id": key,
		}

		switch operation {
		case "encrypt":
			url = "http://localhost:8000/api/transit/encrypt"
			request["plaintext"] = inputData
		case "decrypt":
			url = "http://localhost:8000/api/transit/decrypt"
			request["ciphertext"] = inputData
		case "hmac":
			url = "http://localhost:8000/api/transit/hmac"
			request["message"] = inputData
		case "hmac-verify":
			url = "http://localhost:8000/api/transit/hmac/verify"
			request["message"] = inputData
			// Assume the signature is provided in an additional argument
			if len(args) < 4 {
				fmt.Println("Signature is required for verify operation")
				os.Exit(1)
			}
			request["hmac"] = args[3]
		case "sign":
			url = "http://localhost:8000/api/transit/sign"
			request["message"] = inputData
		case "verify":
			url = "http://localhost:8000/api/transit/verify"
			request["message"] = inputData
			// Assume the signature is provided in an additional argument
			if len(args) < 4 {
				fmt.Println("Signature is required for verify operation")
				os.Exit(1)
			}
			request["signature"] = args[3]
		default:
			fmt.Printf("Unsupported operation: %s\n", operation)
			os.Exit(1)
		}

		// Marshal the request payload
		data, err := json.Marshal(request)
		if err != nil {
			fmt.Printf("Error marshaling request: %v\n", err)
			os.Exit(1)
		}

		// Make the HTTP request
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

		// Handle non-OK response
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error response from server: %s\n", body)
			os.Exit(1)
		}

		// Decode response
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			os.Exit(1)
		}

		// Print results based on operation
		var printData [][]string
		switch operation {
		case "encrypt":
			printData = [][]string{{"CipherText", response["ciphertext"].(string)}}
		case "decrypt":
			printData = [][]string{{"PlainText", response["plaintext"].(string)}}
		case "hmac":
			printData = [][]string{{"HMAC", response["hmac"].(string)}}
		case "hmac-verify":
			verified := "false"
			if response["verified"].(bool) {
				verified = "true"
			}
			printData = [][]string{{"Verified", verified}}
		case "sign":
			printData = [][]string{{"Signature", response["signature"].(string)}}
		case "verify":
			verified := "false"
			if response["verified"].(bool) {
				verified = "true"
			}
			printData = [][]string{{"Verified", verified}}
		}

		// Render the output table
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
	rootCmd.AddCommand(doTransitCmd)
}
