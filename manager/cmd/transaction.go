// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

// transactionCmd represents the transaction command
var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Interact with the transaction store",
	Long:  `Interact with the transaction store.`,
}

func init() {
	rootCmd.AddCommand(transactionCmd)
}
