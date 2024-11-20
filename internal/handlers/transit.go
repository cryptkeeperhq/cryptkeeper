package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	enginetransit "github.com/cryptkeeperhq/cryptkeeper/internal/engine/transit"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
)

// Request and Response Types
type CreateKeyRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateKeyResponse struct {
	KeyID        string    `json:"key_id"`
	CreationDate time.Time `json:"creation_date"`
}

type EncryptRequest struct {
	KeyID     string `json:"key_id"`
	Plaintext string `json:"plaintext"`
}

type EncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
}

type DecryptRequest struct {
	KeyID      string `json:"key_id"`
	Ciphertext string `json:"ciphertext"`
}

type DecryptResponse struct {
	Plaintext string `json:"plaintext"`
}

type SignRequest struct {
	KeyID   string `json:"key_id"`
	Message string `json:"message"`
}

type SignResponse struct {
	Signature string `json:"signature"`
}

type VerifyRequest struct {
	KeyID     string `json:"key_id"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type VerifyResponse struct {
	Verified bool `json:"verified"`
}

type HmacRequest struct {
	KeyID   string `json:"key_id"`
	Message string `json:"message"`
}
type HmacResponse struct {
	HMAC string `json:"hmac"`
}

type HmacVerifyRequest struct {
	KeyID   string `json:"key_id"`
	Message string `json:"message"`
	HMAC    string `json:"hmac"`
}

func (h *Handler) ListTransitKeys(w http.ResponseWriter, r *http.Request) {
	var secrets []models.Secret

	err := h.DB.Model(&secrets).
		ColumnExpr("secret.*, p.path, secret.metadata->>'key_type' AS key_type").
		Join("JOIN paths p ON p.id = secret.path_id").
		Where("p.engine_type = ?", "transit").
		Where("secret.version = (SELECT MAX(version) FROM secrets AS s WHERE s.id = secret.id)").
		Select()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, secrets)
}

func handleSecretAndPath(h *Handler, keyID string, version int) (*enginetransit.Handler, error) {
	secret, err := db.GetTransitEncryptionKey(keyID, version)
	if err != nil {
		return nil, err
	}

	// Decrypt the template using Path Key
	decryptedPathKeyHandle, err := h.getPathHandle(secret.PathID)
	if err != nil {
		return nil, err
	}

	// plaintextValue, err := h.CryptoOps.DecryptSecretValue(secret.EncryptedDEK, secret.EncryptedValue, decryptedPathKeyHandle)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println("plaintextValue", plaintextValue)

	// buf := bytes.NewReader([]byte(plaintextValue))
	// reader := keyset.NewBinaryReader(buf)

	// // Read and return the keyset handle
	// secretHandle, err := keyset.ReadWithNoSecrets(reader)
	// if err != nil {
	// 	return nil, err
	// }

	// var template tinkpb.KeyTemplate
	// if err := proto.Unmarshal([]byte(plaintextValue), &template); err != nil {
	// 	fmt.Println("ERROR!!!", "failed to unmarshal KeyTemplate")
	// 	return nil, fmt.Errorf("failed to unmarshal KeyTemplate: %v", err)
	// }

	// fmt.Println("template", template)
	// secretHandle, err := keyset.NewHandle(&template)

	pathAead, err := aead.New(decryptedPathKeyHandle)
	if err != nil {
		return nil, err
	}

	decryptedDek, err := pathAead.Decrypt(secret.EncryptedDEK, nil)
	if err != nil {
		return nil, err
	}

	secretHandle, _ := keyset.Read(keyset.NewBinaryReader(bytes.NewReader(decryptedDek)), pathAead)

	return enginetransit.NewHandler(secretHandle, secret.KeyType, version)

}

// Handlers
func (h *Handler) Encrypt(w http.ResponseWriter, r *http.Request) {
	var req EncryptRequest
	if decodeJSONBody(w, r, &req) != nil {
		return
	}

	transitHandler, err := handleSecretAndPath(h, req.KeyID, 0)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	plainText, err := base64.URLEncoding.DecodeString(req.Plaintext)
	if err != nil {
		http.Error(w, "invalid base64 message", http.StatusNotFound)
		return
	}

	ciphertext, err := transitHandler.Encrypt(plainText, []byte{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, EncryptResponse{Ciphertext: fmt.Sprintf("%d:%s", transitHandler.Version, base64.StdEncoding.EncodeToString(ciphertext))})
}

func (h *Handler) Decrypt(w http.ResponseWriter, r *http.Request) {
	var req DecryptRequest
	if decodeJSONBody(w, r, &req) != nil {
		return
	}

	val := strings.Split(req.Ciphertext, ":")
	if len(val) < 2 {
		http.Error(w, "invalid encrypted value", http.StatusBadRequest)
		return
	}
	version, cipherText := val[0], val[1]

	transitHandler, err := handleSecretAndPath(h, req.KeyID, utils.ToInt(version))
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	cipherTextBytes, _ := base64.StdEncoding.DecodeString(cipherText)
	plaintext, err := transitHandler.Decrypt(cipherTextBytes, []byte{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, DecryptResponse{Plaintext: base64.StdEncoding.EncodeToString(plaintext)})
}

func (h *Handler) Sign(w http.ResponseWriter, r *http.Request) {
	var req SignRequest
	if decodeJSONBody(w, r, &req) != nil {
		return
	}

	transitHandler, err := handleSecretAndPath(h, req.KeyID, 0)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	plainText, err := base64.URLEncoding.DecodeString(req.Message)
	if err != nil {
		http.Error(w, "invalid base64 message", http.StatusNotFound)
		return
	}

	signature, err := transitHandler.Sign(plainText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, SignResponse{Signature: fmt.Sprintf("%d:%s", transitHandler.Version, hex.EncodeToString(signature))})
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if decodeJSONBody(w, r, &req) != nil {
		return
	}

	val := strings.Split(req.Signature, ":")
	if len(val) < 2 {
		http.Error(w, "invalid encrypted value", http.StatusBadRequest)
		return
	}
	version, cipherText := val[0], val[1]

	transitHandler, err := handleSecretAndPath(h, req.KeyID, utils.ToInt(version))
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	plainText, err := base64.URLEncoding.DecodeString(req.Message)
	if err != nil {
		http.Error(w, "invalid base64 message", http.StatusNotFound)
		return
	}

	signatureBytes, _ := hex.DecodeString(cipherText)
	err = transitHandler.Verify(plainText, signatureBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, VerifyResponse{Verified: true})
}

func (h *Handler) Hmac(w http.ResponseWriter, r *http.Request) {
	var req HmacRequest
	if decodeJSONBody(w, r, &req) != nil {
		return
	}

	transitHandler, err := handleSecretAndPath(h, req.KeyID, 0)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	plainText, err := base64.URLEncoding.DecodeString(req.Message)
	if err != nil {
		http.Error(w, "invalid base64 message", http.StatusNotFound)
		return
	}

	ciphertext, err := transitHandler.ComputeHMAC(plainText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, HmacResponse{
		HMAC: fmt.Sprintf("%d:%s", transitHandler.Version, base64.StdEncoding.EncodeToString(ciphertext)),
	})
}

func (h *Handler) HmacVerify(w http.ResponseWriter, r *http.Request) {
	var req HmacVerifyRequest
	if decodeJSONBody(w, r, &req) != nil {
		return
	}

	val := strings.Split(req.HMAC, ":")
	if len(val) < 2 {
		http.Error(w, "invalid encrypted value", http.StatusBadRequest)
		return
	}
	version, cipherText := val[0], val[1]

	transitHandler, err := handleSecretAndPath(h, req.KeyID, utils.ToInt(version))
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	plainText, err := base64.URLEncoding.DecodeString(req.Message)
	if err != nil {
		http.Error(w, "invalid base64 message", http.StatusNotFound)
		return
	}

	hmacBytes, _ := base64.StdEncoding.DecodeString(cipherText)
	err = transitHandler.VerifyHMAC(plainText, hmacBytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, VerifyResponse{
		Verified: true,
	})
}

// Utility to get decrypted path handle
func (h *Handler) getPathHandle(pathId string) (*keyset.Handle, error) {
	var path models.Path
	if err := h.DB.Model(&path).Where("id = ?", pathId).Select(); err != nil {
		return nil, err
	}

	decryptedPathKeyHandle, err := h.CryptoOps.DecryptPathKey(path.KeyData)
	return decryptedPathKeyHandle, err
}
