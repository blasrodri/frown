package lsof

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSockets(t *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
	}()
	l, err := net.Listen("tcp", ":3000")
	conn, err := l.Accept()
	if err != nil {
		return
	}

	defer conn.Close()
	// Do some stuff
	pid := os.Getpid()
	p, _ := newProcess(pid)
	files, err := ListOpenSockets(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(files))
}

func TestMonitorUserConnections(t *testing.T) {
	_, err := MonitorUserConnections()
	if err != nil {
		t.Fatal(err)
	}
}
