// SPDX-License-Identifier: Apache-2.0

package store

import "context"

type Transaction struct {
	ChargeStationId   string       `firestore:"chargeStationId"`
	TransactionId     string       `firestore:"transactionId"`
	IdToken           string       `firestore:"idToken"`
	TokenType         string       `firestore:"tokenType"`
	MeterValues       []MeterValue `firestore:"meterValues"`
	StartSeqNo        int          `firestore:"startSeqNo"`
	EndedSeqNo        int          `firestore:"endedSeqNo"`
	UpdatedSeqNoCount int          `firestore:"updatedSeqNoCount"`
	Offline           bool         `firestore:"offline"`
}

type MeterValue struct {
	SampledValues []SampledValue `firestore:"sampledValue"`
	Timestamp     string         `firestore:"timestamp"`
}

type SampledValue struct {
	Context       *string        `firestore:"context"`
	Location      *string        `firestore:"location"`
	Measurand     *string        `firestore:"measurand"`
	Phase         *string        `firestore:"phase"`
	UnitOfMeasure *UnitOfMeasure `firestore:"unitOfMeasure"`
	Value         float64        `firestore:"value"`
}

type UnitOfMeasure struct {
	Unit      string `firestore:"unit"`
	Multipler int    `firestore:"multipler"`
}

type TransactionStore interface {
	Transactions(ctx context.Context) ([]*Transaction, error)
	FindTransaction(ctx context.Context, chargeStationId, transactionId string) (*Transaction, error)
	CreateTransaction(ctx context.Context, chargeStationId, transactionId, idToken, tokenType string, meterValue []MeterValue, seqNo int, offline bool) error
	UpdateTransaction(ctx context.Context, chargeStationId, transactionId string, meterValue []MeterValue) error
	EndTransaction(ctx context.Context, chargeStationId, transactionId, idToken, tokenType string, meterValue []MeterValue, seqNo int) error
}
