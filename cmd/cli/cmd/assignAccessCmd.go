package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var assignAccessCmd = &cobra.Command{
	Use:   "assign-access [secret_id] [access_level]",
	Short: "Assign access to a secret",
	Long:  `Assign access to a secret for a user or group`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		secretID := args[0]
		accessLevel := args[1]

		userID, _ := cmd.Flags().GetInt64("user-id")
		groupID, _ := cmd.Flags().GetInt64("group-id")

		if userID == 0 && groupID == 0 {
			fmt.Println("Either --user-id or --group-id must be specified")
			os.Exit(1)
		}

		assignAccess(secretID, userID, groupID, accessLevel)
	},
}

func init() {
	assignAccessCmd.Flags().Int64("user-id", 0, "User ID to assign access")
	assignAccessCmd.Flags().Int64("group-id", 0, "Group ID to assign access")
	rootCmd.AddCommand(assignAccessCmd)
}

func assignAccess(secretID string, userID, groupID int64, accessLevel string) {
	access := map[string]interface{}{
		"secret_id":    secretID,
		"access_level": accessLevel,
	}

	if userID != 0 {
		access["user_id"] = userID
	}
	if groupID != 0 {
		access["group_id"] = groupID
	}

	jsonData, err := json.Marshal(access)
	if err != nil {
		fmt.Println("Error marshalling access:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/assign_access", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Error assigning access: %s\n", resp.Status)
		os.Exit(1)
	}

	fmt.Println("Access assigned successfully.")
}
