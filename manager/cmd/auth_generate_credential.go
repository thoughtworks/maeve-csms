package cmd

import (
	"crypto/rand"
	"fmt"
	"github.com/spf13/cobra"
	"math/big"
)

// generateCredentialCmd represents the generateCredentialCmd command
var generateCredentialCmd = &cobra.Command{
	Use:   "generate-credential",
	Short: "Generates a credential for use as initial token in OCPI registration process",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := generateRandomString()
		if err != nil {
			return err
		}
		fmt.Println(result)

		return nil
	},
}

func generateRandomString() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, 64)
	for i := 0; i < 64; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func init() {
	authCmd.AddCommand(generateCredentialCmd)
}
