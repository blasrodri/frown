package lsof

import (
	"testing"
	"fmt"
)

func TestLoadSockets(t *testing.T){
	p := NewProcess(2459)
	files, err := listOpenSockets(p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s", files)
}
