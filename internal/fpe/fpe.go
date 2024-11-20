package fpe

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/tink-crypto/tink-go/v2/core/registry"
	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
	ubiqfpe "gitlab.com/ubiqsecurity/ubiq-fpe-go"
	"google.golang.org/protobuf/proto"

	fpepb "github.com/cryptkeeperhq/cryptkeeper/internal/fpe/proto"
)

// Define constants for FPE
const (
	inputCharacterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tweakSize         = 7
	fpeTypeURL        = "type.googleapis.com/fpe.FF31FPE"
)

// FPEPrimitive interface representing FPE operations
type FPEPrimitive interface {
	EncryptFPE(plaintext string) (string, error)
	DecryptFPE(ciphertext string) (string, error)
}

// FF31FPE struct representing an FPE primitive using FF3-1
type FF31FPE struct {
	fpeClient *ubiqfpe.FF3_1
}

func New(handle *keyset.Handle) (FPEPrimitive, error) {
	// Write keyset to an in-memory buffer
	var buf bytes.Buffer
	memKeyset := keyset.NewBinaryWriter(&buf)
	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
	if err := insecurecleartextkeyset.Write(handle, memKeyset); err != nil {
		return nil, fmt.Errorf("failed to write keyset: %v", err)
	}

	// Read the keyset back from the buffer
	reader := keyset.NewBinaryReader(&buf)
	keysetData, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read keyset: %v", err)
	}

	// Retrieve primary key data
	var keyData []byte
	for _, key := range keysetData.Key {
		if key.KeyId == keysetData.PrimaryKeyId {
			keyData = key.KeyData.Value
			break
		}
	}
	if keyData == nil {
		return nil, fmt.Errorf("no primary key found in keyset")
	}

	// Initialize FPEKeyManager and create primitive
	km := &FPEKeyManager{}
	primitive, err := km.Primitive(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to create FPE primitive: %v", err)
	}
	return primitive.(FPEPrimitive), nil
}

// NewFF31FPE creates a new FF3-1 FPE primitive
func NewFF31FPE(masterKey, tweak []byte) (*FF31FPE, error) {
	inputAlphabet, err := ubiqfpe.NewAlphabet(inputCharacterSet)
	if err != nil {
		return nil, fmt.Errorf("failed to create alphabet: %v", err)
	}

	fpeClient, err := ubiqfpe.NewFF3_1(masterKey, tweak, inputAlphabet.Len(), inputCharacterSet)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize FF3-1 client: %v", err)
	}

	return &FF31FPE{fpeClient: fpeClient}, nil
}

// EncryptFPE encrypts the plaintext using FPE
func (f *FF31FPE) EncryptFPE(plaintext string) (string, error) {
	return f.fpeClient.Encrypt(plaintext, nil)
}

// DecryptFPE decrypts the ciphertext using FPE
func (f *FF31FPE) DecryptFPE(ciphertext string) (string, error) {
	return f.fpeClient.Decrypt(ciphertext, nil)
}

// FPEKeyManager manages FPE keys and generates FPE primitives
type FPEKeyManager struct{}

// Primitive constructs an FPE primitive from the key data
func (km *FPEKeyManager) Primitive(keyData []byte) (interface{}, error) {
	if len(keyData) < 39 {
		return nil, fmt.Errorf("keyData is too short to contain both master key and tweak")
	}

	var fpeKey fpepb.FPEKey
	if err := proto.Unmarshal(keyData, &fpeKey); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key format: %v", err)
	}

	fmt.Println(fpeKey.Version)
	fmt.Println(fpeKey.MasterKey, len(fpeKey.MasterKey))
	fmt.Println(fpeKey.Tweak, len(fpeKey.Tweak))

	// Validate the tweak length to prevent runtime errors
	if len(fpeKey.Tweak) != 7 {
		return nil, fmt.Errorf("invalid tweak length: expected 7 bytes, got %d", len(fpeKey.Tweak))
	}

	return NewFF31FPE(fpeKey.MasterKey, fpeKey.Tweak)
}

// DoesSupport checks if this KeyManager supports the given type URL
func (km *FPEKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == fpeTypeURL
}

// NewKey generates new key data for FPE
func (km *FPEKeyManager) NewKey(keyFormat []byte) (proto.Message, error) {
	var fpeKeyFormat fpepb.FPEKeyFormat
	if err := proto.Unmarshal(keyFormat, &fpeKeyFormat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key format: %v", err)
	}

	fmt.Println("KEY SIZE", fpeKeyFormat.KeySize)

	keySize := int(fpeKeyFormat.GetKeySize())
	if keySize == 0 {
		keySize = 32
	}

	masterKey, err := generateRandomKeyData(keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	// tweak := generateTweak(masterKey)
	tweak := generateRandomTweak(7)

	fmt.Printf("Master key length: %d, Tweak length: %d\n", len(masterKey), len(tweak))

	fpeKey := &fpepb.FPEKey{
		Version:   0,
		MasterKey: masterKey,
		Tweak:     tweak,
	}

	fmt.Printf("FPEKey MasterKey length: %d, Tweak length: %d\n", len(fpeKey.MasterKey), len(fpeKey.Tweak))

	return fpeKey, nil
}

// NewKeyData generates new key data for FPE
func (km *FPEKeyManager) NewKeyData(keyFormat []byte) (*tinkpb.KeyData, error) {
	key, err := km.NewKey(keyFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new key: %v", err)
	}

	serializedKey, err := proto.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize key: %v", err)
	}

	return &tinkpb.KeyData{
		TypeUrl:         fpeTypeURL,
		Value:           serializedKey,
		KeyMaterialType: tinkpb.KeyData_SYMMETRIC,
	}, nil
}

// KeyMaterialType specifies the type of key material used
func (km *FPEKeyManager) KeyMaterialType() tinkpb.KeyData_KeyMaterialType {
	return tinkpb.KeyData_SYMMETRIC
}

// TypeURL returns the type URL of the FPE key
func (km *FPEKeyManager) TypeURL() string {
	return fpeTypeURL
}

// generateRandomKeyData generates random key data of the specified length
func generateRandomKeyData(length int) ([]byte, error) {
	keyData := make([]byte, length)
	if _, err := rand.Read(keyData); err != nil {
		return nil, fmt.Errorf("failed to generate random key data: %v", err)
	}
	return keyData, nil
}

// generateTweak generates a random tweak
func generateTweak(keyData []byte) []byte {
	h := hmac.New(sha256.New, keyData)
	h.Write(keyData)
	fullTweak := h.Sum(nil)
	return fullTweak[:tweakSize]
}

// generateRandomTweak generates a tweak of the specified length (7 bytes for FF3)
func generateRandomTweak(length int) []byte {
	tweak := make([]byte, length)
	if _, err := rand.Read(tweak); err != nil {
		// Handle error appropriately if tweak generation fails
		return nil
	}
	return tweak
}

// FPEKeyTemplate returns a KeyTemplate for creating FPE keys in Tink
func FPEKeyTemplate() *tinkpb.KeyTemplate {
	return createFPEKeyTemplate(32, tinkpb.OutputPrefixType_RAW)

}

func createFPEKeyTemplate(keySize uint32, outputPrefixType tinkpb.OutputPrefixType) *tinkpb.KeyTemplate {
	format := &fpepb.FPEKeyFormat{
		KeySize: keySize,
	}

	serializedFormat, err := proto.Marshal(format)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal key format: %s", err))
	}

	return &tinkpb.KeyTemplate{
		TypeUrl:          fpeTypeURL,
		Value:            serializedFormat,
		OutputPrefixType: outputPrefixType,
	}
}

// RegisterFPEKeyManager registers the FPEKeyManager with Tink's registry
func RegisterFPEKeyManager() error {
	return registry.RegisterKeyManager(new(FPEKeyManager))
}
