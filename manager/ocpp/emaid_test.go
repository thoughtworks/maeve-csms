package ocpp_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"math/rand"
	"strings"
	"testing"
)

func TestNormalizeEmaidWithEmaidWithHyphens(t *testing.T) {
	norm, err := ocpp.NormalizeEmaid("GB-TWK-012345678-V")
	require.NoError(t, err)
	assert.Equal(t, "GBTWK012345678V", norm)
}

func TestNormalizeEmaidWithLowerCaseEmaid(t *testing.T) {
	norm, err := ocpp.NormalizeEmaid("gbtwk123456789")
	require.NoError(t, err)
	assert.Equal(t, "GBTWK123456789B", norm)
}

func TestNormalizeEmaidWithoutCheckDigit(t *testing.T) {
	emaids := []string{
		"GBTWK012345678V",
		"CSKTH5U8TC90A1S",
		"IHRFRNPCZVPPVEW",
	}

	for _, emaid := range emaids {
		t.Run(emaid, func(t *testing.T) {
			norm, err := ocpp.NormalizeEmaid(emaid[0:14])
			require.NoError(t, err)
			assert.Equal(t, emaid[14], norm[14])
		})
	}
}

func TestNormalizeInvalidEmaid(t *testing.T) {
	_, err := ocpp.NormalizeEmaid("GB-TWK-01234567")
	assert.ErrorContains(t, err, "emaid GB-TWK-01234567 is invalid")
}

var alpha = "ABCDEFGHIJKLMNOPQRSTUVWXZY"
var alphaNumeric = "0123456789" + alpha

func generateEmaid() string {
	countryCode := fmt.Sprintf("%c%c", alpha[rand.Int()%len(alpha)], alpha[rand.Int()%len(alpha)])
	partyId := fmt.Sprintf("%c%c%c", alpha[rand.Int()%len(alpha)], alpha[rand.Int()%len(alpha)], alpha[rand.Int()%len(alpha)])

	var id string
	for i := 0; i < 9; i++ {
		id += fmt.Sprintf("%c", alphaNumeric[rand.Int()%len(alphaNumeric)])
	}

	return fmt.Sprintf("%s%s%s", countryCode, partyId, id)
}

func TestNormalizeRandomEmaid(t *testing.T) {
	for i := 0; i < 15; i++ {
		emaid := generateEmaid()
		t.Run(emaid, func(t *testing.T) {
			norm, err := ocpp.NormalizeEmaid(emaid)
			// validate no error
			require.NoError(t, err)
			// validate check digit is from allowed alphabet
			assert.True(t, strings.Contains(alphaNumeric, string(norm[14])))
		})
	}
}
