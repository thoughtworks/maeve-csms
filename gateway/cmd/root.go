// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csms",
	Short: "Charge Station Management System",
	Long: `Provides a Charge Station Management System that is horizontally
scalable. There are two core components, the gateway accepts
connections from charge stations and forwards messages to/from
an MQTT broker. The manager reads messages from the MQTT broker,
implements any logic and determines the appropriate response.
The manager can also initiate a request to the charge station
and receive the response.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
