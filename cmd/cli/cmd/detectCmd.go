package cmd

import (
	"fmt"
	"os"

	"github.com/cryptkeeperhq/cryptkeeper/internal/detection"
	"github.com/spf13/cobra"
)

// detectCmd represents the detect command
var detectCmd = &cobra.Command{
	Use:   "detect [directory]",
	Short: "Detect secrets in a specified directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := args[0]
		results, err := detection.ScanDirectory(dirPath)
		if err != nil {
			fmt.Printf("Error scanning directory: %v\n", err)
			os.Exit(1)
		}

		if len(results) > 0 {
			fmt.Println("Secrets detected:")
			for file, matches := range results {
				fmt.Printf("File: %s\n", file)
				for _, match := range matches {
					fmt.Printf("  %s\n", match)
				}
			}
			os.Exit(1) // Indicate that secrets were found
		} else {
			fmt.Println("No secrets detected.")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
