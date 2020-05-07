package lsof

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func getTcpConnections() {
}

func hexIpToDecimal(ipHex string) string {
	a, _ := hex.DecodeString(ipHex)
	s := fmt.Sprintf("%v.%v.%v.%v", a[0], a[1], a[2], a[3])
	return s
}


func hexPortToDecimal(portHex string) string {
	a, _ := hex.DecodeString(portHex)
	s := fmt.Sprintf("%v%v", a[0], a[1])
	s = strings.TrimPrefix(s, "0")
	return s
}
