// SPDX-License-Identifier: Apache-2.0

// Package config provides support for configuring the system.
//
// BaseConfig provides the data structures that represent the
// configuration and can be used to load configuration from a
// TOML file.
//
// Config consumes the BaseConfig and returns implementations
// of the various interfaces used by the application (in a
// dependency-injection style).
package config
