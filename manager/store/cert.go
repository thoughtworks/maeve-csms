// SPDX-License-Identifier: Apache-2.0

package store

import "context"

type CertificateStore interface {
	SetCertificate(ctx context.Context, pemCertificate string) error
	LookupCertificate(ctx context.Context, certificateHash string) (string, error)
	DeleteCertificate(ctx context.Context, certificateHash string) error
}
