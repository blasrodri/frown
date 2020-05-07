package lsof

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHexIpToDecimal(t *testing.T){
	z := "0A010100"
	result := hexIpToDecimal(z)
	assert.Equal(t, result, "10.1.1.0")
}

func TestHexPortToDecimal(t *testing.T){
	z := "5050"
	result := hexPortToDecimal(z)
	assert.Equal(t, "8080", result)
}
