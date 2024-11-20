package crypt

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/tink-crypto/tink-go/v2/aead"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	pkcs11 "github.com/miekg/pkcs11"
	"github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
	"google.golang.org/protobuf/proto"
)

var (
	libraryPath = "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	slotLabel   = "ForKMS"
	keyLabel    = "MyKEKLabel"
	hsmPin      = "1234"
)

// softhsm2-util --init-token --slot 0 --label "MyKEKLabel" --pin 1234 --so-pin 0000
// softhsm2-util --init-token --slot 5 --label "MyKEKLabel" --pin 1234 --so-pin 0000
// The token has been initialized and is reassigned to slot 868617254

type HSM struct {
	ctx       *pkcs11.Ctx
	session   pkcs11.SessionHandle
	slotLabel string
}

func NewHSM(config *config.Config) (*HSM, error) {
	h := &HSM{}
	h.ctx = pkcs11.New(libraryPath) // Path to your SoftHSM library
	err := h.ctx.Initialize()
	if err != nil {
		return nil, err
	}

	slot, err := h.getSlot()
	if err != nil {
		return nil, err
	}

	pin := os.Getenv("SOFTHSM_PIN")
	session, err := h.ctx.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, err
	}
	err = h.ctx.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		return nil, err
	}

	h.session = session
	h.slotLabel = slotLabel
	return h, nil
}

// Generate a new Path Key (DEK)
func (t *HSM) GeneratePathKey() (*keyset.Handle, error) {
	pathKeyHandle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	return pathKeyHandle, err
}

func (h *HSM) EncryptPathKey(pathKey *keyset.Handle) ([]byte, error) {
	// Prepare to write the keyset into a memory buffer
	memKeyset := &keyset.MemReaderWriter{}
	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
	if err := insecurecleartextkeyset.Write(pathKey, memKeyset); err != nil {
		fmt.Println("Failed to write keyset:", err)
		return nil, err
	}

	// Serialize the keyset (stored in memKeyset) to a byte slice
	dekBuf, err := proto.Marshal(memKeyset.Keyset)
	if err != nil {
		fmt.Println("Failed to marshal keyset:", err)
		return nil, err
	}

	handle, err := h.findKey(keyLabel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initialize failed: %s\n", err)
		return nil, err
	}

	encryptedDekBuf, err := h.encrypt(handle, dekBuf)

	return encryptedDekBuf, err
}

func (h *HSM) DecryptPathKey(encryptedPathKey []byte) (*keyset.Handle, error) {
	handle, err := h.findKey(keyLabel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initialize failed: %s\n", err)
		return nil, err
	}

	decryptedDek, err := h.decrypt(handle, encryptedPathKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[decryptedDek] Decrypt failed: %s\n", err)
		return nil, err
	}

	// Deserialize the keyset
	ks := tink_go_proto.Keyset{}
	if err := proto.Unmarshal(decryptedDek, &ks); err != nil {
		log.Fatalf("Failed to unmarshal keyset: %v", err)
		return nil, err
	}

	// Use keyset.NewHandle to create a new keyset handle from the keyset
	keysetReader := &keyset.MemReaderWriter{Keyset: &ks}
	kh, err := insecurecleartextkeyset.Read(keysetReader)
	if err != nil {
		log.Fatalf("Failed to create keyset handle: %v", err)
		return nil, err
	}

	return kh, nil
}

func (h *HSM) EncryptSecretValue(input string, pathKeyHandle *keyset.Handle) ([]byte, []byte, error) {
	pathAead, _ := aead.New(pathKeyHandle)

	// Generate Key Handle
	secretKeyHandle, _ := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	var buffer bytes.Buffer
	writer := keyset.NewBinaryWriter(&buffer)
	if err := secretKeyHandle.Write(writer, pathAead); err != nil {
		return nil, nil, err
	}

	// New AEAD primitive
	secretAead, err := aead.New(secretKeyHandle)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the value
	ciphertext, err := secretAead.Encrypt([]byte(input), nil)
	if err != nil {
		return nil, nil, err
	}
	// Lastly encrypt the DEK using Path Key
	serializedSecretKey := buffer.Bytes()
	encryptedSecretKey, err := pathAead.Encrypt(serializedSecretKey, nil)
	if err != nil {
		return nil, nil, err
	}

	return encryptedSecretKey, ciphertext, err

}

func (h *HSM) DecryptSecretValue(encryptedDEK, encryptedValue []byte, pathKeyHandle *keyset.Handle) (string, error) {
	pathAead, err := aead.New(pathKeyHandle)
	if err != nil {
		return "", err
	}

	// Firstly decrypt the DEK using Path Key
	decryptedDek, err := pathAead.Decrypt(encryptedDEK, nil)
	if err != nil {
		return "", err
	}

	// Read the decrypted DEK into Key Handle
	var buffer bytes.Buffer
	reader := keyset.NewBinaryReader(&buffer)
	buffer.Write(decryptedDek)
	secretKeyHandle, err := keyset.Read(reader, pathAead)
	if err != nil {
		return "", err
	}

	// New AEAD primitive
	secretAead, err := aead.New(secretKeyHandle)
	if err != nil {
		return "", err
	}

	decryptedValue, err := secretAead.Decrypt(encryptedValue, nil)

	return string(decryptedValue), err
}

// FindKey finds a key by label.
func (client *HSM) findKey(label string) (pkcs11.ObjectHandle, error) {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
	}

	if err := client.ctx.FindObjectsInit(client.session, template); err != nil {
		return 0, fmt.Errorf("find objects init failed: %w", err)
	}
	obj, _, err := client.ctx.FindObjects(client.session, 1)
	if err != nil {
		return 0, fmt.Errorf("find objects failed: %w", err)
	}
	if err := client.ctx.FindObjectsFinal(client.session); err != nil {
		return 0, fmt.Errorf("find objects final failed: %w", err)
	}
	if len(obj) == 0 {
		return 0, fmt.Errorf("key not found")
	}
	return obj[0], nil
}

// Encrypt encrypts data using the specified key.
func (client *HSM) encrypt(key pkcs11.ObjectHandle, data []byte) ([]byte, error) {
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, make([]byte, 16))} // AES CBC mode with PKCS7 padding
	err := client.ctx.EncryptInit(client.session, mechanism, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EncryptInit failed: %s\n", err)
		return nil, err
	}

	encryptedDek, err := client.ctx.Encrypt(client.session, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encrypt failed: %s\n", err)
		return nil, err
	}
	return encryptedDek, nil
}

// Decrypt decrypts data using the specified key.
func (client *HSM) decrypt(key pkcs11.ObjectHandle, data []byte) ([]byte, error) {
	// Assuming obj[0] is our KEK, use it to decrypt the encrypted DEK
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, make([]byte, 16))} // Ensure correct IV
	err := client.ctx.DecryptInit(client.session, mechanism, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DecryptInit failed: %s\n", err)
		return nil, err
	}

	decryptedDek, err := client.ctx.Decrypt(client.session, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Decrypt failed: %s\n", err)
		return nil, err
	}

	return decryptedDek, err
}

func (h *HSM) getSlot() (uint, error) {
	slots, err := h.ctx.GetSlotList(true)

	if err != nil {
		return 0, err
	}
	if len(slots) == 0 {
		return 0, errors.New("no slots available")
	}
	return slots[0], nil
}