// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	handlers16 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"strconv"
)

// convertIdCmd represents the convertId command
var convertIdCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert an OCPP 1.6 transaction id into the stored transaction id",
	Long: `OCPP 1.6 transaction ids are integer values: but the transaction
store uses UUIDs to store transactions. This command converts the provided
OCPP 1.6 transaction id into a transaction store UUID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tbl := table.New("Transaction Id", "Transaction UUID")

		for _, arg := range args {
			transactionId, err := strconv.ParseInt(arg, 10, 32)
			if err != nil {
				tbl.Print()
				return fmt.Errorf("converting %s to an integer: %w", arg, err)
			}
			transactionUuid := handlers16.ConvertToUUID(int(transactionId))

			tbl.AddRow(transactionId, transactionUuid)
		}

		tbl.Print()

		return nil
	},
}

func init() {
	transactionCmd.AddCommand(convertIdCmd)
}
