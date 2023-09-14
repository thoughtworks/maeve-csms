// SPDX-License-Identifier: Apache-2.0

package inmemory

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"k8s.io/utils/clock"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/thoughtworks/maeve-csms/manager/store"
)

// Store is an in-memory implementation of the store.Engine interface. As everything
// is stored in memory it is not stateless and cannot be used if running >1 manager
// instances. It is primarily provided to support unit testing.
type Store struct {
	sync.Mutex
	clock                            clock.PassiveClock
	chargeStationAuth                map[string]*store.ChargeStationAuth
	chargeStationSettings            map[string]*store.ChargeStationSettings
	chargeStationInstallCertificates map[string]*store.ChargeStationInstallCertificates
	chargeStationRuntimeDetails      map[string]*store.ChargeStationRuntimeDetails
	tokens                           map[string]*store.Token
	transactions                     map[string]*store.Transaction
	certificates                     map[string]string
	registrations                    map[string]*store.OcpiRegistration
	partyDetails                     map[string]*store.OcpiParty
	locations                        map[string]*store.Location
}

func NewStore(clock clock.PassiveClock) *Store {
	return &Store{
		clock:                            clock,
		chargeStationAuth:                make(map[string]*store.ChargeStationAuth),
		chargeStationSettings:            make(map[string]*store.ChargeStationSettings),
		chargeStationInstallCertificates: make(map[string]*store.ChargeStationInstallCertificates),
		chargeStationRuntimeDetails:      make(map[string]*store.ChargeStationRuntimeDetails),
		tokens:                           make(map[string]*store.Token),
		transactions:                     make(map[string]*store.Transaction),
		certificates:                     make(map[string]string),
		registrations:                    make(map[string]*store.OcpiRegistration),
		partyDetails:                     make(map[string]*store.OcpiParty),
		locations:                        make(map[string]*store.Location),
	}
}

func (s *Store) SetChargeStationAuth(_ context.Context, chargeStationId string, auth *store.ChargeStationAuth) error {
	s.Lock()
	defer s.Unlock()
	s.chargeStationAuth[chargeStationId] = auth
	return nil
}

func (s *Store) LookupChargeStationAuth(_ context.Context, chargeStationId string) (*store.ChargeStationAuth, error) {
	s.Lock()
	defer s.Unlock()
	return s.chargeStationAuth[chargeStationId], nil
}

func (s *Store) UpdateChargeStationSettings(_ context.Context, chargeStationId string, settings *store.ChargeStationSettings) error {
	s.Lock()
	defer s.Unlock()
	set := s.chargeStationSettings[chargeStationId]
	if set == nil {
		set = &store.ChargeStationSettings{
			ChargeStationId: chargeStationId,
			Settings:        make(map[string]*store.ChargeStationSetting, len(settings.Settings)),
		}
		maps.Copy(set.Settings, settings.Settings)
	} else {
		for k, v := range settings.Settings {
			set.Settings[k] = v
		}
	}
	s.chargeStationSettings[chargeStationId] = set
	return nil
}

func (s *Store) LookupChargeStationSettings(_ context.Context, chargeStationId string) (*store.ChargeStationSettings, error) {
	s.Lock()
	defer s.Unlock()
	return s.chargeStationSettings[chargeStationId], nil
}

func (s *Store) ListChargeStationSettings(_ context.Context, pageSize int, previousChargeStationId string) ([]*store.ChargeStationSettings, error) {
	s.Lock()
	defer s.Unlock()

	keys := maps.Keys(s.chargeStationSettings)
	sort.Strings(keys)

	i, found := slices.BinarySearch(keys, previousChargeStationId)
	if !found {
		i = 0
	} else {
		i++
	}

	var settings []*store.ChargeStationSettings
	max := int(math.Min(float64(i+pageSize), float64(len(keys))))
	for _, k := range keys[i:max] {
		settings = append(settings, s.chargeStationSettings[k])
	}
	return settings, nil
}

func (s *Store) UpdateChargeStationInstallCertificates(_ context.Context, chargeStationId string, certificates *store.ChargeStationInstallCertificates) error {
	s.Lock()
	defer s.Unlock()
	certs := s.chargeStationInstallCertificates[chargeStationId]
	if certs == nil {
		certs = &store.ChargeStationInstallCertificates{
			ChargeStationId: chargeStationId,
			Certificates:    slices.Clone(certificates.Certificates),
		}
	} else {
		var newCerts []*store.ChargeStationInstallCertificate
		for _, v := range certificates.Certificates {
			matched := false
			for _, c := range certs.Certificates {
				if v.CertificateId == c.CertificateId {
					c.CertificateData = v.CertificateData
					c.CertificateInstallationStatus = v.CertificateInstallationStatus
					c.CertificateType = v.CertificateType
					matched = true
					break
				}
			}
			if !matched {
				newCerts = append(newCerts, v)
			}
		}
		certs.Certificates = append(certs.Certificates, newCerts...)
	}
	s.chargeStationInstallCertificates[chargeStationId] = certs
	return nil
}

func (s *Store) LookupChargeStationInstallCertificates(_ context.Context, chargeStationId string) (*store.ChargeStationInstallCertificates, error) {
	s.Lock()
	defer s.Unlock()
	return s.chargeStationInstallCertificates[chargeStationId], nil
}

func (s *Store) ListChargeStationInstallCertificates(_ context.Context, pageSize int, previousChargeStationId string) ([]*store.ChargeStationInstallCertificates, error) {
	s.Lock()
	defer s.Unlock()

	keys := maps.Keys(s.chargeStationInstallCertificates)
	sort.Strings(keys)

	i, found := slices.BinarySearch(keys, previousChargeStationId)
	if !found {
		i = 0
	} else {
		i++
	}

	var installCertificates []*store.ChargeStationInstallCertificates
	max := int(math.Min(float64(i+pageSize), float64(len(keys))))
	for _, k := range keys[i:max] {
		installCertificates = append(installCertificates, s.chargeStationInstallCertificates[k])
	}
	return installCertificates, nil
}

func (s *Store) SetChargeStationRuntimeDetails(_ context.Context, chargeStationId string, details *store.ChargeStationRuntimeDetails) error {
	s.Lock()
	defer s.Unlock()
	s.chargeStationRuntimeDetails[chargeStationId] = details
	return nil
}

func (s *Store) LookupChargeStationRuntimeDetails(_ context.Context, chargeStationId string) (*store.ChargeStationRuntimeDetails, error) {
	s.Lock()
	defer s.Unlock()
	return s.chargeStationRuntimeDetails[chargeStationId], nil
}

func (s *Store) SetToken(_ context.Context, token *store.Token) error {
	s.Lock()
	defer s.Unlock()
	token.LastUpdated = time.Now().UTC().Format(time.RFC3339)
	s.tokens[token.Uid] = token
	return nil
}

func (s *Store) LookupToken(_ context.Context, tokenUid string) (*store.Token, error) {
	s.Lock()
	defer s.Unlock()
	return s.tokens[tokenUid], nil
}

func (s *Store) ListTokens(_ context.Context, offset int, limit int) ([]*store.Token, error) {
	s.Lock()
	defer s.Unlock()
	var tokens []*store.Token
	count := 0
	for _, token := range s.tokens {
		if count >= offset && count < offset+limit {
			tokens = append(tokens, token)
		}
		count++
	}
	if tokens == nil {
		tokens = make([]*store.Token, 0)
	}
	return tokens, nil
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
	b64Hash := base64.RawURLEncoding.EncodeToString(hash[:])
	return b64Hash, nil
}

func (s *Store) LookupCertificate(_ context.Context, certificateHash string) (string, error) {
	s.Lock()
	defer s.Unlock()

	return s.certificates[certificateHash], nil
}

func (s *Store) DeleteCertificate(_ context.Context, certificateHash string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.certificates, certificateHash)

	return nil
}

func (s *Store) SetRegistrationDetails(_ context.Context, token string, registration *store.OcpiRegistration) error {
	s.Lock()
	defer s.Unlock()

	s.registrations[token] = registration

	return nil
}

func (s *Store) GetRegistrationDetails(_ context.Context, token string) (*store.OcpiRegistration, error) {
	s.Lock()
	defer s.Unlock()
	return s.registrations[token], nil
}

func (s *Store) DeleteRegistrationDetails(_ context.Context, token string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.registrations, token)

	return nil
}

func (s *Store) SetPartyDetails(_ context.Context, partyDetails *store.OcpiParty) error {
	s.Lock()
	defer s.Unlock()

	recordId := fmt.Sprintf("%s:%s:%s", partyDetails.Role, partyDetails.CountryCode, partyDetails.PartyId)

	s.partyDetails[recordId] = partyDetails

	return nil
}

func (s *Store) GetPartyDetails(_ context.Context, role, countryCode, partyId string) (*store.OcpiParty, error) {
	s.Lock()
	defer s.Unlock()

	recordId := fmt.Sprintf("%s:%s:%s", role, countryCode, partyId)

	return s.partyDetails[recordId], nil
}

func (s *Store) ListPartyDetailsForRole(_ context.Context, role string) ([]*store.OcpiParty, error) {
	s.Lock()
	defer s.Unlock()
	var parties []*store.OcpiParty
	for _, party := range s.partyDetails {
		if party.Role == role {
			parties = append(parties, party)
		}
	}
	if parties == nil {
		parties = make([]*store.OcpiParty, 0)
	}
	return parties, nil
}

func (s *Store) SetLocation(_ context.Context, location *store.Location) error {
	s.Lock()
	defer s.Unlock()

	s.locations[location.Id] = location

	return nil
}

func (s *Store) LookupLocation(_ context.Context, locationId string) (*store.Location, error) {
	s.Lock()
	defer s.Unlock()

	return s.locations[locationId], nil
}

func (s *Store) ListLocations(_ context.Context, offset int, limit int) ([]*store.Location, error) {
	s.Lock()
	defer s.Unlock()
	var locations []*store.Location
	count := 0
	for _, location := range s.locations {
		if count >= offset && count < offset+limit {
			locations = append(locations, location)
		}
		count++
	}
	if locations == nil {
		locations = make([]*store.Location, 0)
	}
	return locations, nil
}
