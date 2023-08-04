// SPDX-License-Identifier: Apache-2.0

package config

type InMemoryStorageConfig struct{}

type FirestoreStorageConfig struct {
	ProjectId string `mapstructure:"project_id" toml:"project_id"`
}

type StorageConfig struct {
	Type             string                  `mapstructure:"type" toml:"type"`
	FirestoreStorage *FirestoreStorageConfig `mapstructure:"firestore,omitempty" toml:"firestore,omitempty"`
	InMemoryStorage  *InMemoryStorageConfig  `mapstructure:"in_memory,omitempty" toml:"in_memory,omitempty"`
}
