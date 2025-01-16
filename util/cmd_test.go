package util

import (
	"fmt"
	"testing"
)

func TestListFSCmd(t *testing.T) {
	m, _ := ListFSCmd("nfs")
	fmt.Println(m)
	m, _ = ListFSCmd("cifs")
	fmt.Println(m)
}

func TestExecCmd(t *testing.T) {
	m, _ := RunBash("dmidecode -s system-serial-number | sed '/^#/ d'")
	fmt.Println(m)
	m, _ = RunBash("dmidecode -s system-uuid | sed '/^#/ d'")
	fmt.Println(m)
}
