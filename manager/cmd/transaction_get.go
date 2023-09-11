// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/utils/clock"

	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
)

var (
	storageEngine string
	gcloudProject string
)

var long = `Get transactions from the transaction store.
Requires an even number of arguments in the format "cs-id transaction-id".`

// serveCmd represents the serve command
var getTransactionsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get transaction details from the store",
	Long:  long,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		transactionStore, err := firestore.NewStore(ctx, gcloudProject, clock.RealClock{})
		if err != nil {
			return fmt.Errorf("creating transaction store: %w", err)
		}

		if len(args) == 0 {
			transactions, err := transactionStore.Transactions(ctx)
			if err != nil {
				return fmt.Errorf("getting transactions: %w", err)
			}

			for _, transaction := range transactions {
				if err != nil {
					fmt.Printf("formatting transaction: %v\n", err)
				}
				fmt.Printf("transaction: %s %s\n", transaction.ChargeStationId, transaction.TransactionId)
			}
			return nil
		}

		if len(args)%2 != 0 {
			return fmt.Errorf("incorrect number of arguments: %v", args)
		}

		for i := 0; i < len(args); i += 2 {
			transaction, err := transactionStore.FindTransaction(ctx, args[i], args[i+1])
			if err != nil {
				fmt.Println("transaction not found:", args)
			}
			formattedTransaction, err := formatTransaction(transaction)
			if err != nil {
				fmt.Printf("formatting transaction: %v\n", err)
			}
			fmt.Printf("transaction: %s\n", formattedTransaction)
		}

		return nil
	},
}

func formatTransaction(transaction *store.Transaction) (string, error) {
	b, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func init() {
	transactionCmd.AddCommand(getTransactionsCmd)

	transactionCmd.Flags().StringVarP(&storageEngine, "storage-engine", "s", "firestore",
		"The storage engine to use for persistence, one of [firestore, inmemory]")
	transactionCmd.Flags().StringVar(&gcloudProject, "gcloud-project", "*detect-project-id*",
		"The google cloud project that hosts the firestore instance (if chosen storage-engine)")
}
