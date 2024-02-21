// SPDX-License-Identifier: Apache-2.0

package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestStringSource(t *testing.T) {
	source := StringSource{
		Data: "hello world",
	}

	data, err := source.GetData(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, "hello world", data)
}

func TestFileSource(t *testing.T) {
	source := FileSource{
		FileName: "testdata/root_ca.pem",
	}

	data, err := source.GetData(context.TODO())
	require.NoError(t, err)
	assert.Contains(t, data, "-----BEGIN CERTIFICATE-----")
}

func TestGoogleCloudSecretSource(t *testing.T) {
	secretName := os.Getenv("TEST_GOOGLE_CLOUD_SECRET_NAME")
	if secretName == "" {
		t.Skip("no test google cloud secrets configured")
	}
	t.Logf("Using secret %s", secretName)

	source := GoogleSecretSource{
		SecretName: secretName,
	}

	data, err := source.GetData(context.TODO())
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}
