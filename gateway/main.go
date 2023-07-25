// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/thoughtworks/maeve-csms/gateway/cmd"
	"golang.org/x/exp/slog"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	cmd.Execute()
}
