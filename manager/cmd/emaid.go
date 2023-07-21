// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

// emaidCmd represents the emaid command
var emaidCmd = &cobra.Command{
	Use:   "emaid",
	Short: "Support for working with eMAIDs",
}

func init() {
	rootCmd.AddCommand(emaidCmd)
}
