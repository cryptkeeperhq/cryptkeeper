package transit

import (
	"errors"
	"fmt"

	"github.com/tink-crypto/tink-go/v2/signature"
)

func (h *Handler) Sign(message []byte) ([]byte, error) {
	if h.KeyType != "ed25519" && h.KeyType != "ecdsa-p256" && h.KeyType != "ecdsa-p384" && h.KeyType != "ecdsa-p521" && h.KeyType != "rsa-2048" && h.KeyType != "rsa-3072" && h.KeyType != "rsa-4096" {
		return nil, errors.New("unsupported key type for signing")
	}

	signer, err := signature.NewSigner(h.Handle)
	if err != nil {
		return nil, fmt.Errorf("error getting signer: %v", err)
	}

	signature, err := signer.Sign(message)
	if err != nil {
		return nil, fmt.Errorf("error signing message: %v", err)
	}

	return signature, nil
}

func (h *Handler) Verify(message []byte, signatureValue []byte) error {
	if h.KeyType != "ed25519" && h.KeyType != "ecdsa-p256" && h.KeyType != "ecdsa-p384" && h.KeyType != "ecdsa-p521" && h.KeyType != "rsa-2048" && h.KeyType != "rsa-3072" && h.KeyType != "rsa-4096" {
		return errors.New("unsupported key type for verification")
	}

	// Extract public keyset handle for verification
	pubKeysetHandle, err := h.Handle.Public()
	if err != nil {
		return fmt.Errorf("error getting public keyset handle: %v", err)
	}

	verifier, err := signature.NewVerifier(pubKeysetHandle)
	if err != nil {
		return fmt.Errorf("error getting verifier: %v", err)
	}

	err = verifier.Verify(signatureValue, message)
	if err != nil {
		return fmt.Errorf("error verifying signature: %v", err)
	}

	return nil
}
