// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
)

// encodePasswordCmd represents the encodePassword command
var encodePasswordCmd = &cobra.Command{
	Use:   "encode-password",
	Short: "Encode a password for storage in the database",
	Run: func(cmd *cobra.Command, args []string) {
		for _, pwd := range args {
			hash := sha256.Sum256([]byte(pwd))
			b64 := base64.StdEncoding.EncodeToString(hash[:])
			fmt.Printf("%s: %s\n", pwd, b64)
		}
	},
}

func init() {
	authCmd.AddCommand(encodePasswordCmd)
}
