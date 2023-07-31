// SPDX-License-Identifier: Apache-2.0

package firestore

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type certificate struct {
	PemCertificate string `firestore:"pem"`
}

func (s *Store) SetCertificate(ctx context.Context, pemCertificate string) error {
	certificateHash, err := getPEMCertificateHash(pemCertificate)
	if err != nil {
		return err
	}
	csRef := s.client.Doc(fmt.Sprintf("Certificate/%s", certificateHash))
	_, err = csRef.Set(ctx, &certificate{
		PemCertificate: pemCertificate,
	})
	if err != nil {
		return err
	}
	return nil
}

func getPEMCertificateHash(pemCertificate string) (string, error) {
	var cert *x509.Certificate
	block, _ := pem.Decode([]byte(pemCertificate))
	if block != nil {
		if block.Type == "CERTIFICATE" {
			var err error
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("pem block does not contain certificate, but %s", block.Type)
		}
	} else {
		return "", fmt.Errorf("pem block not found")
	}

	hash := sha256.Sum256(cert.Raw)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash[:])
	return b64Hash, nil
}

func (s *Store) LookupCertificate(ctx context.Context, certificateHash string) (string, error) {
	csRef := s.client.Doc(fmt.Sprintf("Certificate/%s", certificateHash))
	snap, err := csRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", nil
		}
		return "", fmt.Errorf("lookup certificate %s: %w", certificateHash, err)
	}
	var cert certificate
	if err = snap.DataTo(&cert); err != nil {
		return "", fmt.Errorf("lookup certificate %s: %w", certificateHash, err)
	}
	return cert.PemCertificate, nil
}

func (s *Store) DeleteCertificate(ctx context.Context, certificateHash string) error {
	csRef := s.client.Doc(fmt.Sprintf("Certificate/%s", certificateHash))
	_, err := csRef.Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return fmt.Errorf("delete certificate %s: %w", certificateHash, err)
	}

	return nil
}
