package lsof

import (
	"encoding/hex"
	"fmt"
)

func getTCPConnections() {
}

func hexIPToDecimal(ipHex string) string {
	a, _ := hex.DecodeString(ipHex)
	s := fmt.Sprintf("%v.%v.%v.%v", a[3], a[2], a[1], a[0])
	return s
}

func hexPortToDecimal(portHex string) string {
	a, _ := hex.DecodeString(portHex)
	r := int(a[0])*256 + int(a[1])
	return fmt.Sprintf("%d", r)
}
