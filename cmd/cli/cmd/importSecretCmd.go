package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Secret struct {
	Path         string                 `json:"path"`
	Key          string                 `json:"key" yaml:"key"`
	Value        string                 `json:"value" yaml:"value"`
	Metadata     map[string]interface{} `json:"metadata" yaml:"metadata"`
	IsMultiValue bool                   `json:"is_multi_value"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	IsOneTime    bool                   `json:"is_one_time"`
	MultiValue   map[string]interface{} `json:"multi_value"`
}

// func flattenMap(data map[string]interface{}, prefix string) map[string]interface{} {
// 	flatMap := make(map[string]interface{})

// 	for k, v := range data {
// 		key := k
// 		if prefix != "" {
// 			key = prefix + "." + k
// 		}

// 		logger.Debug("Map type", "Type", reflect.TypeOf(v))

// 		switch value := v.(type) {
// 		case map[interface{}]interface{}:
// 			for k1, v1 := range value {
// 				nestedMap := flattenMap(v1.(map[string]interface{}), k1.(string))
// 				for nk, nv := range nestedMap {
// 					flatMap[nk] = nv
// 				}
// 			}
// 		case map[string]interface{}:
// 			nestedMap := flattenMap(value, key)
// 			for nk, nv := range nestedMap {
// 				flatMap[nk] = nv
// 			}
// 		default:
// 			flatMap[key] = value
// 		}
// 	}

// 	return flatMap
// }

func Flatten2(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}

	for k, v := range src {
		// log.Println("Map type", "Type", reflect.TypeOf(v))

		switch child := v.(type) {
		case map[string]interface{}:
			Flatten2(prefix+k, child, dest)
		case map[interface{}]interface{}:
			// fmt.Println(len(child), child)
			// for i := 0; i < len(child); i++ {
			for kk, vv := range child {
				jsonValue, _ := json.Marshal(vv)
				dest[prefix+k+"."+kk.(string)] = string(jsonValue)

				// dest[prefix+k+"."+kk.(string)] = vv
				// Flatten2(prefix+k+prefix+kk.(string), vv.(map[string]interface{}), dest)
			}
		case []interface{}:
			for i := 0; i < len(child); i++ {
				dest[prefix+k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			dest[prefix+k] = v
		}
	}
}

func createSecret(secret Secret) error {

	token := os.Getenv("TOKEN")
	if token == "" {
		fmt.Println("TOKEN environment variable is required")
		os.Exit(1)
	}

	data, err := json.Marshal(secret)
	if err != nil {
		fmt.Printf("Error marshaling secret: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000/api/secrets?path=%s&key=%s", secret.Path, secret.Key), bytes.NewBuffer(data))
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

	fmt.Println("Secret created successfully")

	return nil
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func importSecrets(path string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	// var data map[string]interface{}
	var data interface{}

	if err := yaml.Unmarshal(byteValue, &data); err != nil {
		if err := json.Unmarshal(byteValue, &data); err != nil {
			return fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}

	data = convert(data)

	// if b, err := json.Marshal(data); err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("Output: %s\n", b)
	// }

	// secrets := make(map[string]interface{})

	// secrets := flattenMap(data, "")
	// Flatten2("", data, secrets)

	for key, value := range data.(map[string]interface{}) {
		// jsonValue, err := json.Marshal(value)
		// if err != nil {
		// 	continue
		// }
		// fmt.Println(string(jsonValue))

		// fmt.Printf("%s = %s\n", secret.Key, string(jsonValue))

		secret := Secret{
			Path:         path,
			Key:          key,
			MultiValue:   value.(map[string]interface{}),
			IsMultiValue: true,
		}

		if err := createSecret(secret); err != nil {
			return fmt.Errorf("failed to create secret for key %s: %w", secret.Key, err)
		}
	}

	return nil
}

var importSecretsCmd = &cobra.Command{
	Use:   "import [path] [file]",
	Short: "Import secrets from config files",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		file := args[1]

		importSecrets(path, file)

		// token := os.Getenv("TOKEN")
		// if token == "" {
		// 	fmt.Println("TOKEN environment variable is required")
		// 	os.Exit(1)
		// }

		// var expiresAt *time.Time
		// if createSecretExpiresAt != "" {
		// 	t, err := time.Parse(time.RFC3339, createSecretExpiresAt)
		// 	if err != nil {
		// 		fmt.Printf("Error parsing expires_at: %v\n", err)
		// 		os.Exit(1)
		// 	}
		// 	expiresAt = &t
		// }

		// m := make(map[string]interface{})
		// err := json.Unmarshal([]byte(metadata), &m)
		// if err != nil {
		// 	fmt.Printf("Error marshaling Metadata: %v\n", err)
		// 	os.Exit(1)
		// }

		// // pathID := utils.GetPath(path, token)

		// secret := map[string]interface{}{
		// 	// "path_id":     pathID,
		// 	"path":        path,
		// 	"key":         key,
		// 	"value":       value,
		// 	"expires_at":  expiresAt,
		// 	"is_one_time": isOneTime,
		// 	"metadata":    m,
		// }

		// log.Println(secret)

		// data, err := json.Marshal(secret)
		// if err != nil {
		// 	fmt.Printf("Error marshaling secret: %v\n", err)
		// 	os.Exit(1)
		// }

		// req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000/api/secrets?path=%s&key=%s", path, key), bytes.NewBuffer(data))
		// if err != nil {
		// 	fmt.Printf("Error creating request: %v\n", err)
		// 	os.Exit(1)
		// }

		// req.Header.Set("Authorization", "Bearer "+token)
		// req.Header.Set("Content-Type", "application/json")

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Printf("Error making request: %v\n", err)
		// 	os.Exit(1)
		// }
		// defer resp.Body.Close()

		// if resp.StatusCode != http.StatusOK {
		// 	body, _ := io.ReadAll(resp.Body)
		// 	fmt.Printf("Error creating secret: %s\n", body)
		// 	os.Exit(1)
		// }

		// fmt.Println("Secret created successfully")

		// if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		// 	fmt.Printf("Error decoding response: %v\n", err)
		// 	os.Exit(1)
		// }

		// formattedData, err := json.MarshalIndent(secret, "", "  ")
		// if err != nil {
		// 	fmt.Printf("Error marshaling secret: %v\n", err)
		// 	os.Exit(1)
		// }

		// fmt.Println(string(formattedData))
	},
}

func init() {

	// importSecretsCmd.Flags().StringVar(&createSecretExpiresAt, "expires-at", "", "Expiration time in RFC3339 format")
	// importSecretsCmd.Flags().BoolVar(&isOneTime, "is-one-time", false, "One time use secret only")
	// createSecretCmd.Flags().StringVar(&metadata, "metadata", "{}", "Optional metadata for secret")
	rootCmd.AddCommand(importSecretsCmd)
}
