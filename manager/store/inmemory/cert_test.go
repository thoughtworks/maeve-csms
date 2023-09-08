// SPDX-License-Identifier: Apache-2.0

package inmemory_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"math/big"
	"testing"
	"time"
)

func TestSetAndLookupAndDeleteCertificate(t *testing.T) {
	cert := generateCertificate(t)

	pemCertificate := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))

	store := inmemory.NewStore(clock.RealClock{})

	err := store.SetCertificate(context.Background(), pemCertificate)
	require.NoError(t, err)

	hash := sha256.Sum256(cert.Raw)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash[:])

	got, err := store.LookupCertificate(context.Background(), b64Hash)
	require.NoError(t, err)

	assert.Equal(t, pemCertificate, got)

	err = store.DeleteCertificate(context.Background(), b64Hash)
	require.NoError(t, err)

	got, err = store.LookupCertificate(context.Background(), b64Hash)
	require.NoError(t, err)

	assert.Equal(t, "", got)
}

func generateCertificate(t *testing.T) *x509.Certificate {
	keyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Thoughtworks"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &keyPair.PublicKey, keyPair)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(derBytes)
	require.NoError(t, err)

	return cert
}
