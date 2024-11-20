package transit

import (
	"errors"

	commonpb "github.com/tink-crypto/tink-go/v2/proto/common_go_proto"
	rsapb "github.com/tink-crypto/tink-go/v2/proto/rsa_ssa_pkcs1_go_proto"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"

	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/daead"
	"github.com/tink-crypto/tink-go/v2/mac"
	"github.com/tink-crypto/tink-go/v2/signature"

	"google.golang.org/protobuf/proto"

	"github.com/cryptkeeperhq/cryptkeeper/internal/fpe"
)

type Handler struct {
	KeyType string
	Handle  *keyset.Handle
	Version int
}

type Transit interface {
	GenerateKeyTemplate() (*tinkpb.KeyTemplate, error)

	Encrypt(plainText []byte, associatedData []byte) ([]byte, error)
	Decrypt(cipherText []byte, associatedData []byte) ([]byte, error)

	ComputeHMAC(message []byte) ([]byte, error)
	VerifyHMAC(message []byte, hmacValue []byte) error

	Sign(message []byte) ([]byte, error)
	Verify(message []byte, signatureValue []byte) error
}

func NewHandler(handle *keyset.Handle, keyType string, version int) (*Handler, error) {
	return &Handler{
		Version: version,
		Handle:  handle,
		KeyType: keyType,
	}, nil
}

func (h *Handler) GenerateKeyTemplate() (*tinkpb.KeyTemplate, error) {
	var keyTemplate *tinkpb.KeyTemplate
	switch h.KeyType {
	case "fpe":
		err := fpe.RegisterFPEKeyManager()
		if err != nil {
			panic(err)
		}
		keyTemplate = fpe.FPEKeyTemplate()
	case "aes128-gcm96":
		keyTemplate = aead.AES128GCMKeyTemplate()
	case "aes256-gcm96":
		keyTemplate = aead.AES256GCMKeyTemplate()
	case "chacha20-poly1305":
		keyTemplate = aead.ChaCha20Poly1305KeyTemplate()
	case "ed25519":
		keyTemplate = signature.ED25519KeyTemplate()
	case "ecdsa-p256":
		keyTemplate = signature.ECDSAP256KeyTemplate()
	case "ecdsa-p384":
		keyTemplate = signature.ECDSAP384SHA512KeyTemplate()
	case "ecdsa-p521":
		keyTemplate = signature.ECDSAP521KeyTemplate()
	case "rsa-2048":
		keyTemplate = createRSASignatureKeyTemplate(2048)
	case "rsa-3072":
		keyTemplate = createRSASignatureKeyTemplate(3072)
	case "rsa-4096":
		keyTemplate = createRSASignatureKeyTemplate(4096)
	case "hmac":
		keyTemplate = mac.HMACSHA256Tag128KeyTemplate()
	case "aes256S-iv":
		keyTemplate = daead.AESSIVKeyTemplate()
	default:
		return nil, errors.New("unsupported key type")
	}

	return keyTemplate, nil
}

func createRSASignatureKeyTemplate(modulusSize int) *tinkpb.KeyTemplate {
	format := &rsapb.RsaSsaPkcs1KeyFormat{
		Params: &rsapb.RsaSsaPkcs1Params{
			HashType: commonpb.HashType_SHA256,
		},
		ModulusSizeInBits: uint32(modulusSize),
		PublicExponent:    []byte{0x01, 0x00, 0x01}, // 65537
	}
	serializedFormat, _ := proto.Marshal(format)
	return &tinkpb.KeyTemplate{
		TypeUrl:          "type.googleapis.com/google.crypto.tink.RsaSsaPkcs1PrivateKey",
		OutputPrefixType: tinkpb.OutputPrefixType_TINK,
		Value:            serializedFormat,
	}
}
