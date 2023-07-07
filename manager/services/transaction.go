// SPDX-License-Identifier: Apache-2.0

package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

type Transaction struct {
	ChargeStationId   string       `json:"chargeStationId"`
	TransactionId     string       `json:"transactionId"`
	IdToken           string       `json:"idToken"`
	TokenType         string       `json:"tokenType"`
	MeterValues       []MeterValue `json:"meterValues"`
	StartSeqNo        int          `json:"startSeqNo"`
	EndedSeqNo        int          `json:"endedSeqNo"`
	UpdatedSeqNoCount int          `json:"updatedSeqNoCount"`
	Offline           bool         `json:"offline"`
}

type MeterValue struct {
	SampledValues []SampledValue `json:"sampledValue"`
	Timestamp     string         `json:"timestamp"`
}

type SampledValue struct {
	Context       *string        `json:"context"`
	Location      *string        `json:"location"`
	Measurand     *string        `json:"measurand"`
	Phase         *string        `json:"phase"`
	UnitOfMeasure *UnitOfMeasure `json:"unitOfMeasure"`
	Value         float64        `json:"value"`
}

type UnitOfMeasure struct {
	Unit      string `json:"unit"`
	Multipler int    `json:"multipler"`
}

func (t Transaction) String() string {
	return fmt.Sprintf("[%s] token=%s(%s) start=%d update=%d end=%d meterValues=%v",
		key(t.ChargeStationId, t.TransactionId),
		t.IdToken,
		t.TokenType,
		t.StartSeqNo,
		t.UpdatedSeqNoCount,
		t.EndedSeqNo,
		t.MeterValues)
}

type TransactionStore interface {
	Transactions() ([]*Transaction, error)
	FindTransaction(chargeStationId, transactionId string) (*Transaction, error)
	CreateTransaction(chargeStationId, transactionId, idToken, tokenType string, meterValue []MeterValue, seqNo int, offline bool) error
	UpdateTransaction(chargeStationId, transactionId string, meterValue []MeterValue) error
	EndTransaction(chargeStationId, transactionId, idToken, tokenType string, meterValue []MeterValue, seqNo int) error
}

type InMemoryTransactionStore struct {
	sync.Mutex

	transactions map[string]*Transaction
}

func key(chargeStationId, transactionId string) string {
	return fmt.Sprintf("%s:%s", chargeStationId, transactionId)
}

func (i *InMemoryTransactionStore) getTransaction(chargeStationId, transactionId string) *Transaction {
	transaction := i.transactions[key(chargeStationId, transactionId)]
	return transaction
}

func (i *InMemoryTransactionStore) updateTransaction(transaction *Transaction) {
	key := key(transaction.ChargeStationId, transaction.TransactionId)
	i.transactions[key] = transaction
}

func (i *InMemoryTransactionStore) Transactions() ([]*Transaction, error) {
	i.Lock()
	defer i.Unlock()

	transactions := make([]*Transaction, 0, len(i.transactions))

	for _, transaction := range i.transactions {
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (i *InMemoryTransactionStore) FindTransaction(chargeStationId, transactionId string) (*Transaction, error) {
	i.Lock()
	defer i.Unlock()
	return i.getTransaction(chargeStationId, transactionId), nil
}

func (i *InMemoryTransactionStore) CreateTransaction(chargeStationId, transactionId, idToken, tokenType string, meterValues []MeterValue, seqNo int, offline bool) error {
	i.Lock()
	defer i.Unlock()
	transaction := i.getTransaction(chargeStationId, transactionId)
	if transaction != nil {
		transaction.IdToken = idToken
		transaction.TokenType = tokenType
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.StartSeqNo = seqNo
		transaction.Offline = offline
	} else {
		transaction = &Transaction{
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
		i.updateTransaction(transaction)
	}
	return nil
}

func (i *InMemoryTransactionStore) UpdateTransaction(chargeStationId, transactionId string, meterValues []MeterValue) error {
	i.Lock()
	defer i.Unlock()
	transaction := i.getTransaction(chargeStationId, transactionId)
	if transaction == nil {
		transaction = &Transaction{
			ChargeStationId:   chargeStationId,
			TransactionId:     transactionId,
			MeterValues:       meterValues,
			UpdatedSeqNoCount: 1,
		}
		i.updateTransaction(transaction)
	} else {
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.UpdatedSeqNoCount++
	}
	return nil
}

func (i *InMemoryTransactionStore) EndTransaction(chargeStationId, transactionId, idToken, tokenType string, meterValues []MeterValue, seqNo int) error {
	i.Lock()
	defer i.Unlock()
	transaction := i.getTransaction(chargeStationId, transactionId)

	if transaction == nil {
		transaction = &Transaction{
			ChargeStationId: chargeStationId,
			TransactionId:   transactionId,
			IdToken:         idToken,
			TokenType:       tokenType,
			MeterValues:     meterValues,
			EndedSeqNo:      seqNo,
		}
		i.updateTransaction(transaction)
	} else {
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.EndedSeqNo = seqNo
	}
	return nil
}

func NewInMemoryTransactionStore() *InMemoryTransactionStore {
	store := new(InMemoryTransactionStore)
	store.transactions = make(map[string]*Transaction)
	return store
}

type RedisTransactionStore struct {
	sync.Mutex
	client *redis.Client
}

func NewRedisTransactionStore(address string) *RedisTransactionStore {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Network:  "tcp",
		Password: "",
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil
	}

	return &RedisTransactionStore{
		Mutex:  sync.Mutex{},
		client: client,
	}
}

func (r *RedisTransactionStore) getTransaction(key string) (*Transaction, error) {
	var transaction Transaction

	value, err := r.client.Get(key).Result()
	if err != redis.Nil {
		err = json.Unmarshal([]byte(value), &transaction)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling: %v", err)
		}

		return &transaction, nil
	}

	return nil, nil
}

func (r *RedisTransactionStore) updateStore(key string, transaction *Transaction) error {
	newValue, err := json.Marshal(&transaction)
	if err != nil {
		return fmt.Errorf("marshalling transaction: %w", err)
	}

	err = r.client.Set(key, newValue, 0).Err()
	if err != nil {
		return fmt.Errorf("storing value: %w", err)
	}

	return nil
}

func (r *RedisTransactionStore) Transactions() ([]*Transaction, error) {
	r.Lock()
	defer r.Unlock()

	keys, err := r.client.Keys("*").Result()
	if err != nil {
		return nil, fmt.Errorf("getting transaction: %w", err)
	}

	transactions := make([]*Transaction, 0, len(keys))
	for _, key := range keys {
		transaction, err := r.getTransaction(key)
		if err != nil {
			fmt.Printf("Not a transaction: [%s] %v\n", key, transaction)
			continue
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *RedisTransactionStore) FindTransaction(chargeStationId, transactionId string) (*Transaction, error) {
	r.Lock()
	defer r.Unlock()

	key := key(chargeStationId, transactionId)
	return r.getTransaction(key)
}

func (r *RedisTransactionStore) CreateTransaction(chargeStationId, transactionId, idToken, tokenType string, meterValues []MeterValue, seqNo int, offline bool) error {
	r.Lock()
	defer r.Unlock()

	key := key(chargeStationId, transactionId)
	transaction, err := r.getTransaction(key)

	if err != nil {
		return fmt.Errorf("getting transaction: %w", err)
	}

	if transaction != nil {
		transaction.IdToken = idToken
		transaction.TokenType = tokenType
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.StartSeqNo = seqNo
		transaction.Offline = offline
	} else {
		transaction = &Transaction{
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
	}

	return r.updateStore(key, transaction)
}

func (r *RedisTransactionStore) UpdateTransaction(chargeStationId, transactionId string, meterValues []MeterValue) error {
	r.Lock()
	defer r.Unlock()

	key := key(chargeStationId, transactionId)
	transaction, err := r.getTransaction(key)

	if err != nil {
		return fmt.Errorf("getting transaction: %w", err)
	}

	if transaction == nil {
		transaction = &Transaction{
			ChargeStationId:   chargeStationId,
			TransactionId:     transactionId,
			MeterValues:       meterValues,
			UpdatedSeqNoCount: 1,
		}
	} else {
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.UpdatedSeqNoCount++
	}

	return r.updateStore(key, transaction)
}

func (r *RedisTransactionStore) EndTransaction(chargeStationId, transactionId, idToken, tokenType string, meterValues []MeterValue, seqNo int) error {
	r.Lock()
	defer r.Unlock()

	key := key(chargeStationId, transactionId)
	transaction, err := r.getTransaction(key)

	if err != nil {
		return fmt.Errorf("getting transaction: %w", err)
	}

	if transaction == nil {
		transaction = &Transaction{
			ChargeStationId: chargeStationId,
			TransactionId:   transactionId,
			IdToken:         idToken,
			TokenType:       tokenType,
			MeterValues:     meterValues,
			EndedSeqNo:      seqNo,
		}
	} else {
		transaction.MeterValues = append(transaction.MeterValues, meterValues...)
		transaction.EndedSeqNo = seqNo
	}

	return r.updateStore(key, transaction)
}
