package lsof

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetUserPids(t *testing.T) {
	// Check that this pid is found as part of the user's pids
	pid := os.Getpid()
	result, err := getUserPids()
	if err != nil {
		t.Fatal(err)
	}
	assertionResult := false
	for _, v := range result {
		if v == pid {
			assertionResult = true
		}

	}
	assert.True(t, assertionResult)
}
