// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

// GetCertificateId returns the SHA-256 thumbprint of the certificate
func GetCertificateId(pemData string) (string, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return "", fmt.Errorf("failed to decode certificate chain")
	}
	if block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("expected certificate, got %s", block.Type)
	}
	hash := sha256.Sum256(block.Bytes)
	certificateId := hex.EncodeToString(hash[:])

	return certificateId, nil
}
