// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// hashCertificateCmd represents the hashCertificate command
var hashCertificateCmd = &cobra.Command{
	Use:   "hash-certificate",
	Short: "Hash a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		for _, pemFile := range args {
			//#nosec G304 - only files specified by the person running the application will be loaded
			b, err := os.ReadFile(pemFile)
			if err != nil {
				fmt.Printf("%s: %v\n", pemFile, err)
				continue
			}
			block, _ := pem.Decode(b)
			if block == nil {
				fmt.Printf("%s: no PEM data found\n", pemFile)
				continue
			}
			if block.Type != "CERTIFICATE" {
				fmt.Printf("%s: expected CERTIFICATE, got %s\n", pemFile, block.Type)
				continue
			}
			hash := sha256.Sum256(block.Bytes)
			b64Hash := base64.URLEncoding.EncodeToString(hash[:])

			fmt.Printf("%s: %s\n", pemFile, b64Hash)
		}
	},
}

func init() {
	authCmd.AddCommand(hashCertificateCmd)
}
