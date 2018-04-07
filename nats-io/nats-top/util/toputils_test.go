package toputils

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"nats-io/gnatsd/server"
	gnatsd "nats-io/gnatsd/test"
)

// Borrowed from gnatsd tests
const GNATSD_PORT = 11422

func runMonitorServer(monitorPort int) *server.Server {
	resetPreviousHTTPConnections()
	opts := gnatsd.DefaultTestOptions
	opts.Host = "127.0.0.1"
	opts.Port = GNATSD_PORT
	opts.HTTPPort = monitorPort

	return gnatsd.RunServer(&opts)
}

func resetPreviousHTTPConnections() {
	http.DefaultTransport = &http.Transport{}
}

func TestFetchingStatz(t *testing.T) {
	engine := &Engine{}
	engine.Uri = fmt.Sprintf("http://%s:%d", "127.0.0.1", server.DEFAULT_HTTP_PORT)
	engine.HttpClient = &http.Client{}

	s := runMonitorServer(server.DEFAULT_HTTP_PORT)
	defer s.Shutdown()

	var varz *server.Varz
	result, err := engine.Request("/varz")
	if err != nil {
		t.Fatalf("Failed getting /varz: %v", err)
	}

	if varzVal, ok := result.(*server.Varz); ok {
		varz = varzVal
	}

	// At the very least it is guaranteed that we have one core
	got := varz.Cores
	if got < 1 {
		t.Fatalf("Could not monitor number of cores. got: %v", got)
	}

	// Create simple subscription to gnatsd port to show subscriptions
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", GNATSD_PORT))
		if err != nil {
			t.Fatalf("could not create subcription to NATS: %s", err)
		}
		fmt.Fprintf(conn, "SUB hello.world  90\r\n")
		time.Sleep(5 * time.Second)
		conn.Close()
	}()
	time.Sleep(1 * time.Second)

	var connz *server.Connz
	result, err = engine.Request("/connz")
	if err != nil {
		t.Fatalf("Failed getting /connz: %v", err)
	}

	if connzVal, ok := result.(*server.Connz); ok {
		connz = connzVal
	}

	// Check that we got connections
	got = len(connz.Conns)
	if got <= 0 {
		t.Fatalf("Could not monitor with subscriptions option. expected non-nil conns, got: %v", got)
	}

	engine.DisplaySubs = true
	result, err = engine.Request("/connz")
	if err != nil {
		t.Fatalf("Failed getting /connz: %v", err)
	}

	if connzVal, ok := result.(*server.Connz); ok {
		connz = connzVal
	}

	// Check that we got subscriptions
	got = len(connz.Conns[0].Subs)
	if got <= 0 {
		t.Fatalf("Could not monitor with client subscriptions. expected client with subscriptions, got: %v", got)
	}

	s.Shutdown()
}

func TestPsize(t *testing.T) {

	expected := "1023"
	got := Psize(1023)
	if got != expected {
		t.Fatalf("Wrong human readable value. expected: %v, got: %v", expected, got)
	}

	expected = "1.0K"
	got = Psize(1024)
	if got != expected {
		t.Fatalf("Wrong human readable value. expected: %v, got: %v", expected, got)
	}

	expected = "1.0M"
	got = Psize(1024 * 1024)
	if got != expected {
		t.Fatalf("Wrong human readable value. expected: %v, got: %v", expected, got)
	}

	expected = "1.0G"
	got = Psize(1024 * 1024 * 1024)
	if got != expected {
		t.Fatalf("Wrong human readable value. expected: %v, got: %v", expected, got)
	}
}

func TestMonitorStats(t *testing.T) {
	engine := NewEngine("127.0.0.1", server.DEFAULT_HTTP_PORT, 10, 1)
	engine.SetupHTTP()
	s := runMonitorServer(server.DEFAULT_HTTP_PORT)
	defer s.Shutdown()

	go func() {
		err := engine.MonitorStats()
		if err != nil {
			t.Fatalf("Could not start info monitoring loop. expected no error, got: %v", err)
		}
	}()
	defer close(engine.ShutdownCh)

	select {
	case stats := <-engine.StatsCh:
		got := stats.Varz.Cores
		if got < 1 {
			t.Fatalf("Could not monitor number of cores. got: %v", got)
		}
		return
	case <-time.After(3 * time.Second):
		t.Fatalf("Timed out polling /varz via http")
	}
}

func TestMonitoringTLSConnectionUsingRootCA(t *testing.T) {
	srv, _ := gnatsd.RunServerWithConfig("./test/tls.conf")
	defer srv.Shutdown()

	engine := NewEngine("127.0.0.1", 8223, 10, 1)
	err := engine.SetupHTTPS("./test/ca.pem", "", "", false)
	if err != nil {
		t.Fatalf("Expected to be able to configure polling via HTTPS. Got: %s", err)
	}

	go func() {
		err := engine.MonitorStats()
		if err != nil {
			t.Fatalf("Could not start info monitoring loop. expected no error, got: %v", err)
		}
	}()
	defer close(engine.ShutdownCh)

	select {
	case stats := <-engine.StatsCh:
		got := stats.Varz.Cores
		if got < 1 {
			t.Fatalf("Could not monitor number of cores. got: %v", got)
		}
		return
	case <-time.After(3 * time.Second):
		t.Fatalf("Timed out polling /varz via https")
	}
}

func TestMonitoringTLSConnectionUsingRootCAWithCerts(t *testing.T) {
	srv, _ := gnatsd.RunServerWithConfig("./test/tls.conf")
	defer srv.Shutdown()

	engine := NewEngine("127.0.0.1", 8223, 10, 1)
	err := engine.SetupHTTPS("./test/ca.pem", "./test/client-cert.pem", "./test/client-key.pem", false)
	if err != nil {
		t.Fatalf("Expected to be able to configure polling via HTTPS. Got: %s", err)
	}

	go func() {
		err := engine.MonitorStats()
		if err != nil {
			t.Fatalf("Could not start info monitoring loop. expected no error, got: %v", err)
		}
	}()
	defer close(engine.ShutdownCh)

	select {
	case stats := <-engine.StatsCh:
		got := stats.Varz.Cores
		if got < 1 {
			t.Fatalf("Could not monitor number of cores. got: %v", got)
		}
		return
	case <-time.After(3 * time.Second):
		t.Fatalf("Timed out polling /varz via https")
	}
}

func TestMonitoringTLSConnectionUsingCertsAndInsecure(t *testing.T) {
	srv, _ := gnatsd.RunServerWithConfig("./test/tls.conf")
	defer srv.Shutdown()

	engine := NewEngine("127.0.0.1", 8223, 10, 1)
	err := engine.SetupHTTPS("", "./test/client-cert.pem", "./test/client-key.pem", true)
	if err != nil {
		t.Fatalf("Expected to be able to configure polling via HTTPS. Got: %s", err)
	}

	go func() {
		err := engine.MonitorStats()
		if err != nil {
			t.Fatalf("Could not start info monitoring loop. expected no error, got: %v", err)
		}
	}()
	defer close(engine.ShutdownCh)

	select {
	case stats := <-engine.StatsCh:
		got := stats.Varz.Cores
		if got < 1 {
			t.Fatalf("Could not monitor number of cores. got: %v", got)
		}
		return
	case <-time.After(3 * time.Second):
		t.Fatalf("Timed out polling /varz via https")
	}
}
