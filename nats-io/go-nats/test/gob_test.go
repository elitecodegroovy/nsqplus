// Copyright 2012-2017 Apcera Inc. All rights reserved.

package test

import (
	"reflect"
	"testing"

	"nats-io/go-nats"
)

func NewGobEncodedConn(tl TestLogger) *nats.EncodedConn {
	ec, err := nats.NewEncodedConn(NewConnection(tl, TEST_PORT), nats.GOB_ENCODER)
	if err != nil {
		tl.Fatalf("Failed to create an encoded connection: %v\n", err)
	}
	return ec
}

func TestEncBuiltinGobMarshalString(t *testing.T) {
	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	ec := NewGobEncodedConn(t)
	defer ec.Close()
	ch := make(chan bool)

	testString := "Hello World!"

	ec.Subscribe("gob_string", func(s string) {
		if s != testString {
			t.Fatalf("Received test string of '%s', wanted '%s'\n", s, testString)
		}
		ch <- true
	})
	ec.Publish("gob_string", testString)
	if e := Wait(ch); e != nil {
		t.Fatal("Did not receive the message")
	}
}

func TestEncBuiltinGobMarshalInt(t *testing.T) {
	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	ec := NewGobEncodedConn(t)
	defer ec.Close()
	ch := make(chan bool)

	testN := 22

	ec.Subscribe("gob_int", func(n int) {
		if n != testN {
			t.Fatalf("Received test int of '%d', wanted '%d'\n", n, testN)
		}
		ch <- true
	})
	ec.Publish("gob_int", testN)
	if e := Wait(ch); e != nil {
		t.Fatal("Did not receive the message")
	}
}

func TestEncBuiltinGobMarshalStruct(t *testing.T) {
	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	ec := NewGobEncodedConn(t)
	defer ec.Close()
	ch := make(chan bool)

	me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery St"}
	me.Children = make(map[string]*person)

	me.Children["sam"] = &person{Name: "sam", Age: 19, Address: "140 New Montgomery St"}
	me.Children["meg"] = &person{Name: "meg", Age: 17, Address: "140 New Montgomery St"}

	me.Assets = make(map[string]uint)
	me.Assets["house"] = 1000
	me.Assets["car"] = 100

	ec.Subscribe("gob_struct", func(p *person) {
		if !reflect.DeepEqual(p, me) {
			t.Fatalf("Did not receive the correct struct response")
		}
		ch <- true
	})

	ec.Publish("gob_struct", me)
	if e := Wait(ch); e != nil {
		t.Fatal("Did not receive the message")
	}
}

func BenchmarkPublishGobStruct(b *testing.B) {
	// stop benchmark for set-up
	b.StopTimer()

	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	ec := NewGobEncodedConn(b)
	defer ec.Close()
	ch := make(chan bool)

	me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery St"}
	me.Children = make(map[string]*person)

	me.Children["sam"] = &person{Name: "sam", Age: 19, Address: "140 New Montgomery St"}
	me.Children["meg"] = &person{Name: "meg", Age: 17, Address: "140 New Montgomery St"}

	ec.Subscribe("gob_struct", func(p *person) {
		if !reflect.DeepEqual(p, me) {
			b.Fatalf("Did not receive the correct struct response")
		}
		ch <- true
	})

	// resume benchmark
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		ec.Publish("gob_struct", me)
		if e := Wait(ch); e != nil {
			b.Fatal("Did not receive the message")
		}
	}
}
