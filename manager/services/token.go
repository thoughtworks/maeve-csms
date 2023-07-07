// SPDX-License-Identifier: Apache-2.0

package services

import (
	"fmt"
)

type Token struct {
	Type string
	Uid  string
	// TODO: see OCPI for more details
}

type TokenStore interface {
	FindToken(typ, id string) (*Token, error)
}

type InMemoryTokenStore struct {
	Tokens map[string]*Token
}

func (s InMemoryTokenStore) FindToken(typ, id string) (*Token, error) {
	if typ == "" {
		for _, tok := range s.Tokens {
			if tok.Uid == id {
				return tok, nil
			}
		}
	}
	return s.Tokens[fmt.Sprintf("%s:%s", typ, id)], nil
}
