// SPDX-License-Identifier: Apache-2.0
//go:build e2e

package test_driver

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"sync"
	"testing"
	"time"
)

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func shutdownBrokerConnection(t *testing.T, client mqtt.Client, wg *sync.WaitGroup) {
	defer client.Disconnect(0)

	t.Log("unplugging...")
	wg.Add(1)
	tok := client.Publish("everest_external/nodered/1/carsim/cmd/modify_charging_session", 0, false, "unplug")
	if tok.WaitTimeout(1 * time.Second); tok.Error() != nil {
		t.Errorf("failed to publish unplug to topic: %v", tok.Error())
	}
	timedOut := waitTimeout(wg, 5*time.Second)
	if timedOut {
		t.Errorf("timed out waiting for unplug to complete")
	}
}

func setupBrokerConnection(t *testing.T, wg *sync.WaitGroup) (mqtt.Client, func()) {
	brokerAddr := os.Getenv("MQTT_BROKER_ADDR")
	if brokerAddr == "" {
		brokerAddr = "127.0.0.1:1884"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddr)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.WaitTimeout(2*time.Second) && token.Error() != nil {
		t.Fatalf("failed to connect to MQTT broker: %v", token.Error())
	}

	return client, func() {
		shutdownBrokerConnection(t, client, wg)
	}
}

func TestRFIDCharge(t *testing.T) {
	wg := sync.WaitGroup{}

	client, shutdown := setupBrokerConnection(t, &wg)
	defer shutdown()

	state := "Unset"

	tok := client.Subscribe("everest_external/nodered/1/state/state_string", 0, func(client mqtt.Client, msg mqtt.Message) {
		t.Logf("charger status changed to %s\n", msg.Payload())

		switch string(msg.Payload()) {
		case "Wait for Auth":
			state = string(msg.Payload())
			wg.Done()
		case "Charging":
			state = string(msg.Payload())
			wg.Done()
		case "Idle":
			state = string(msg.Payload())
			wg.Done()
		}
	})
	if tok.WaitTimeout(1 * time.Second); tok.Error() != nil {
		t.Fatalf("failed to subscribe to topic: %v", tok.Error())
	}

	wg.Add(1)
	t.Logf("plug in...")
	tok = client.Publish("everest_external/nodered/1/carsim/cmd/execute_charging_session", 0, false,
		"sleep 1;iec_wait_pwr_ready;sleep 1;draw_power_regulated 16,3;sleep 5")

	if tok.WaitTimeout(1 * time.Second); tok.Error() != nil {
		t.Fatalf("failed to publish start charge to topic: %v", tok.Error())
	}

	timedOut := waitTimeout(&wg, 5*time.Second)
	if timedOut {
		t.Errorf("timed out waiting for ready to authorise")
	}
	if state != "Wait for Auth" {
		t.Fatalf("expected state to be Wait for Auth, got %s", state)
	}

	wg.Add(1)
	t.Logf("authorise charge...")
	tok = client.Publish("everest_api/dummy_token_provider/cmd/provide", 0, false,
		fmt.Sprintln("{\"id_token\":\"DEADBEEF\",\"authorization_type\":\"RFID\",\"prevalidated\":false,\"connectors\":[1]}"))
	if tok.WaitTimeout(1 * time.Second); tok.Error() != nil {
		t.Fatalf("failed to publish authorise to topic: %v", tok.Error())
	}

	timedOut = waitTimeout(&wg, 5*time.Second)
	if timedOut {
		t.Errorf("timed out waiting for authorise to complete")
	}
	if state != "Charging" {
		t.Fatalf("expected state to be Charging, got %s", state)
	}
}

func TestISO15118Charge(t *testing.T) {
	t.Skip("skipping test in short mode.")
	wg := sync.WaitGroup{}

	client, shutdown := setupBrokerConnection(t, &wg)
	defer shutdown()

	state := "Unset"

	tok := client.Subscribe("everest_external/nodered/1/state/state_string", 0, func(client mqtt.Client, msg mqtt.Message) {
		t.Logf("charger status changed to %s\n", msg.Payload())

		switch string(msg.Payload()) {
		case "PrepareCharging":
			state = string(msg.Payload())
			wg.Done()
		case "Idle":
			state = string(msg.Payload())
			wg.Done()
		}
	})
	if tok.WaitTimeout(1 * time.Second); tok.Error() != nil {
		t.Fatalf("failed to subscribe to topic: %v", tok.Error())
	}

	wg.Add(1)
	t.Logf("plug in...")
	tok = client.Publish("everest_external/nodered/1/carsim/cmd/execute_charging_session", 0, false,
		"sleep 1;iso_wait_slac_matched;iso_start_v2g_session Contract,AC_three_phase_core;iso_wait_pwr_ready;iso_draw_power_regulated 16,3;sleep 36000")
	if tok.WaitTimeout(1 * time.Second); tok.Error() != nil {
		t.Fatalf("failed to publish start charge to topic: %v", tok.Error())
	}

	timedOut := waitTimeout(&wg, 15*time.Second)
	if timedOut {
		t.Errorf("timed out waiting for charge to start")
	}
	if state != "PrepareCharging" {
		t.Fatalf("expected state to be PrepareCharging, got %s", state)
	}
}
