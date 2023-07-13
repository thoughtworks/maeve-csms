package cmd

import (
	"github.com/spf13/cobra"
)

// contractCmd represents the contract command
var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Support for contract certificates",
}

func init() {
	rootCmd.AddCommand(contractCmd)
}
