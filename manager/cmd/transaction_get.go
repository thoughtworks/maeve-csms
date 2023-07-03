package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
)

var long = `Get transactions from the transaction store.
Requires an even number of arguments in the format "cs-id transaction-id".`

// serveCmd represents the serve command
var getTransactionsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get transaction details from the store",
	Long:  long,
	RunE: func(cmd *cobra.Command, args []string) error {
		transactionStore := services.NewRedisTransactionStore(redisAddr)
		if transactionStore == nil {
			return errors.New("unable to connect to transaction store at address " + redisAddr)
		}

		if len(args) == 0 {
			transactions, err := transactionStore.Transactions()
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
			transaction, err := transactionStore.FindTransaction(args[i], args[i+1])
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

func formatTransaction(transaction *services.Transaction) (string, error) {
	b, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func init() {
	transactionCmd.AddCommand(getTransactionsCmd)

	getTransactionsCmd.Flags().StringVarP(&redisAddr, "redis-addr", "r", "127.0.0.1:6379",
		"The address of the Redis store, e.g. 127.0.0.1:6379")
}
