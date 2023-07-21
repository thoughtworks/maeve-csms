// SPDX-License-Identifier: Apache-2.0

package inmemory

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"sync"

	"github.com/thoughtworks/maeve-csms/manager/store"
)

// Store is an in-memory implementation of the store.Engine interface. As everything
// is stored in memory it is not stateless and cannot be used if running >1 manager
// instances. It is primarily provided to support unit testing.
type Store struct {
	sync.Mutex
	chargeStationAuth map[string]*store.ChargeStationAuth
	tokens            map[string]*store.Token
	transactions      map[string]*store.Transaction
	certificates      map[string]string
}

func NewStore() *Store {
	return &Store{
		chargeStationAuth: make(map[string]*store.ChargeStationAuth),
		tokens:            make(map[string]*store.Token),
		transactions:      make(map[string]*store.Transaction),
		certificates:      make(map[string]string),
	}
}

func (s *Store) SetChargeStationAuth(ctx context.Context, chargeStationId string, auth *store.ChargeStationAuth) error {
	s.Lock()
	defer s.Unlock()
	s.chargeStationAuth[chargeStationId] = auth
	return nil
}

func (s *Store) LookupChargeStationAuth(ctx context.Context, chargeStationId string) (*store.ChargeStationAuth, error) {
	s.Lock()
	defer s.Unlock()
	return s.chargeStationAuth[chargeStationId], nil
}

func (s *Store) SetToken(ctx context.Context, token *store.Token) error {
	s.Lock()
	defer s.Unlock()
	s.tokens[token.Uid] = token
	return nil
}

func (s *Store) LookupToken(ctx context.Context, tokenUid string) (*store.Token, error) {
	s.Lock()
	defer s.Unlock()
	return s.tokens[tokenUid], nil
}

func transactionKey(chargeStationId, transactionId string) string {
	return fmt.Sprintf("%s:%s", chargeStationId, transactionId)
}

func (s *Store) getTransaction(chargeStationId, transactionId string) *store.Transaction {
	transaction := s.transactions[transactionKey(chargeStationId, transactionId)]
	return transaction
}

func (s *Store) updateTransaction(transaction *store.Transaction) {
	key := transactionKey(transaction.ChargeStationId, transaction.TransactionId)
	s.transactions[key] = transaction
}

func (s *Store) Transactions(_ context.Context) ([]*store.Transaction, error) {
	s.Lock()
	defer s.Unlock()

	transactions := make([]*store.Transaction, 0, len(s.transactions))

	for _, transaction := range s.transactions {
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Store) FindTransaction(_ context.Context, chargeStationId, transactionId string) (*store.Transaction, error) {
	s.Lock()
	defer s.Unlock()
	return s.getTransaction(chargeStationId, transactionId), nil
}

func (s *Store) CreateTransaction(_ context.Context, chargeStationId, transactionId, idToken, tokenType string, meterValues []store.MeterValue, seqNo int, offline bool) error {
	s.Lock()
	defer s.Unlock()
	transaction := s.getTransaction(chargeStationId, transactionId)
	if transaction != nil {
		transaction.IdToken = idToken
		transaction.TokenType = tokenType
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.StartSeqNo = seqNo
		transaction.Offline = offline
	} else {
		transaction = &store.Transaction{
			ChargeStationId:   chargeStationId,
			TransactionId:     transactionId,
			IdToken:           idToken,
			TokenType:         tokenType,
			MeterValues:       meterValues,
			StartSeqNo:        seqNo,
			EndedSeqNo:        0,
			UpdatedSeqNoCount: 0,
			Offline:           offline,
		}
		s.updateTransaction(transaction)
	}
	return nil
}

func (s *Store) UpdateTransaction(_ context.Context, chargeStationId, transactionId string, meterValues []store.MeterValue) error {
	s.Lock()
	defer s.Unlock()
	transaction := s.getTransaction(chargeStationId, transactionId)
	if transaction == nil {
		transaction = &store.Transaction{
			ChargeStationId:   chargeStationId,
			TransactionId:     transactionId,
			MeterValues:       meterValues,
			UpdatedSeqNoCount: 1,
		}
		s.updateTransaction(transaction)
	} else {
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.UpdatedSeqNoCount++
	}
	return nil
}

func (s *Store) EndTransaction(_ context.Context, chargeStationId, transactionId, idToken, tokenType string, meterValues []store.MeterValue, seqNo int) error {
	s.Lock()
	defer s.Unlock()
	transaction := s.getTransaction(chargeStationId, transactionId)

	if transaction == nil {
		transaction = &store.Transaction{
			ChargeStationId: chargeStationId,
			TransactionId:   transactionId,
			IdToken:         idToken,
			TokenType:       tokenType,
			MeterValues:     meterValues,
			EndedSeqNo:      seqNo,
		}
		s.updateTransaction(transaction)
	} else {
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.EndedSeqNo = seqNo
	}
	return nil
}

func (s *Store) SetCertificate(_ context.Context, pemCertificate string) error {
	s.Lock()
	defer s.Unlock()

	b64Hash, err := getPEMCertificateHash(pemCertificate)
	if err != nil {
		return err
	}

	s.certificates[b64Hash] = pemCertificate

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
	b64Hash := base64.URLEncoding.EncodeToString(hash[:])
	return b64Hash, nil
}

func (s *Store) LookupCertificate(_ context.Context, certificateHash string) (string, error) {
	s.Lock()
	defer s.Unlock()

	return s.certificates[certificateHash], nil
}

func (s *Store) DeleteCertificate(ctx context.Context, certificateHash string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.certificates, certificateHash)

	return nil
}
