// SPDX-License-Identifier: Apache-2.0

package config

type InMemoryStorageConfig struct{}

type FirestoreStorageConfig struct {
	ProjectId string `mapstructure:"project_id" toml:"project_id" validate:"required"`
}

type StorageConfig struct {
	Type             string                  `mapstructure:"type" toml:"type" validate:"required,oneof=firestore in_memory"`
	FirestoreStorage *FirestoreStorageConfig `mapstructure:"firestore,omitempty" toml:"firestore,omitempty" validate:"required_if=Type firestore"`
	InMemoryStorage  *InMemoryStorageConfig  `mapstructure:"in_memory,omitempty" toml:"in_memory,omitempty"`
}
