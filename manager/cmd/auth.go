package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Support configuring charge station authentication",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
