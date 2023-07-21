// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
)

// normalizeCmd represents the checkdigit command
var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize an eMAID",
	Long:  "Processes each argument and returns its normalized value",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			norm, err := ocpp.NormalizeEmaid(arg)
			if err != nil {
				fmt.Printf("%s: %v\n", arg, err)
			} else {
				fmt.Printf("%s: %s\n", arg, norm)
			}
		}
	},
}

func init() {
	emaidCmd.AddCommand(normalizeCmd)
}
