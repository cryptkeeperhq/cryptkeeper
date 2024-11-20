package transit

import (
	"errors"
	"fmt"

	"github.com/tink-crypto/tink-go/v2/tink"
)

func (h *Handler) ComputeHMAC(message []byte) ([]byte, error) {
	if h.KeyType != "hmac" {
		return nil, errors.New("unsupported key type for HMAC")
	}

	m, err := h.getPrimitive()
	if err != nil {
		return nil, fmt.Errorf("error getting MAC primitive: %v", err)
	}

	hmacValue, err := m.(tink.MAC).ComputeMAC(message)
	if err != nil {
		return nil, fmt.Errorf("error computing HMAC: %v", err)
	}

	return hmacValue, nil
}

func (h *Handler) VerifyHMAC(message []byte, hmacValue []byte) error {
	if h.KeyType != "hmac" {
		return errors.New("unsupported key type for HMAC")
	}

	m, err := h.getPrimitive()
	if err != nil {
		return fmt.Errorf("error getting MAC primitive: %v", err)
	}

	err = m.(tink.MAC).VerifyMAC(hmacValue, message)
	if err != nil {
		return fmt.Errorf("error verifying HMAC: %v", err)
	}

	return nil
}
