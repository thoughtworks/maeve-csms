// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"fmt"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
	"net"
	"net/url"
	"testing"
)

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "127.0.0.1:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

// NewBroker creates a local MQTT broker that can be used for testing.
func NewBroker(t *testing.T) (*mqtt.Server, *url.URL) {
	server := mqtt.New(nil)

	err := server.AddHook(new(auth.AllowHook), nil)
	if err != nil {
		t.Fatalf("adding auth hook: %v", err)
	}

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("getting free port: %v", err)
	}

	ws := listeners.NewWebsocket("broker1", fmt.Sprintf("127.0.0.1:%d", port), nil)
	err = server.AddListener(ws)
	if err != nil {
		t.Fatalf("adding ws listener: %v", err)
	}

	addr, err := url.Parse(fmt.Sprintf("ws://%s", ws.Address()))
	if err != nil {
		t.Fatalf("parsing broker url: %v", err)
	}

	return server, addr
}
