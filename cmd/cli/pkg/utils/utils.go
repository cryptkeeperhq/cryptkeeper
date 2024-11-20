package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetPath(name string, token string) int64 {
	req, err := http.NewRequest("GET", "http://localhost:8000/paths?path="+name, nil)
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

	var paths []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&paths); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		os.Exit(1)
	}

	if len(paths) == 0 {
		fmt.Println("no path by that name")
	}

	if len(paths) > 1 {
		fmt.Println("multiple paths by the name")
	}

	path := paths[0]

	// data, err := json.MarshalIndent(paths, "", "  ")
	// if err != nil {
	// 	fmt.Printf("Error marshaling secret: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(string(data))

	// pathID, err := strconv.ParseInt(path["id"].(int), 10, 64)
	// if err != nil {
	// 	fmt.Printf("Error marshaling secret: %v\n", err)
	// 	os.Exit(1)
	// }
	return int64(path["id"].(float64))

}

func ToInt64(i string) int64 {
	v, _ := strconv.ParseInt(i, 10, 64)
	return v
}

func ToInt(i string) int {
	v, _ := strconv.Atoi(i)
	return v
}
