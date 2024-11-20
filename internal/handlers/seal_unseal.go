package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/hashicorp/vault/shamir"
)

var (
	mu sync.Mutex
)

type SealRequest struct {
	Shares []string `json:"shares"`
}

// CombineShares combines the shares to reconstruct the master key
func CombineShares(shares [][]byte) ([]byte, error) {
	masterKey, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}
	return masterKey, nil
}

func (h *Handler) SealStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": false,
	})

}

func (h *Handler) SealHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// isSealed = true

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Vault sealed"))
}

func (h *Handler) UnsealHandler(w http.ResponseWriter, r *http.Request) {
	var req SealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Combine the shares to reconstruct the master key
	shares := [][]byte{}
	for _, v := range req.Shares {
		shareString, err := hex.DecodeString(v)
		if err != nil {
			http.Error(w, "Failed to unseal: "+err.Error(), http.StatusInternalServerError)
			return
		}
		shares = append(shares, shareString)
	}

	fmt.Println(shares)

	// masterKey, err := CombineShares(shares)
	// if err != nil {
	// 	http.Error(w, "Failed to unseal: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// masterKeyString := base64.StdEncoding.EncodeToString(masterKey)
	// fmt.Printf("Reconstructed Key: %x\n", masterKey)
	// err = crypto.InitKeyManagement(masterKeyString)
	// if err != nil {
	// 	http.Error(w, "Failed to unseal: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Vault unsealed"))
}
