// SPDX-License-Identifier: Apache-2.0

package services

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"io"
	"os"
)

type LocalSource interface {
	GetData(ctx context.Context) (string, error)
}

type StringSource struct {
	Data string
}

func (s StringSource) GetData(_ context.Context) (string, error) {
	return s.Data, nil
}

type FileSource struct {
	FileName string
}

func (f FileSource) GetData(_ context.Context) (string, error) {
	file, err := os.Open(f.FileName)
	if err != nil {
		return "", fmt.Errorf("opening file %s: %v", f.FileName, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %v", f.FileName, err)
	}
	return string(data), nil
}

type GoogleSecretSource struct {
	SecretName string
}

func (g GoogleSecretSource) GetData(ctx context.Context) (string, error) {
	var client *secretmanager.Client
	client, err := secretmanager.NewClient(ctx)
	defer func() {
		_ = client.Close()
	}()
	if err != nil {
		return "", fmt.Errorf("creating secretmanager client: %v", err)
	}
	var resp *secretmanagerpb.AccessSecretVersionResponse
	resp, err = client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: g.SecretName,
	})
	if err != nil {
		return "", fmt.Errorf("reading secret %s: %v", g.SecretName, err)
	}

	return string(resp.GetPayload().GetData()), nil
}
