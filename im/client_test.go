package im

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	Connect("127.0.0.1", 5060)
}

func TestBuildCrc1021Sign(t *testing.T) {
	v := []byte{76, 79, 71, 73, 78, 123, 34, 117, 115, 101, 114, 73, 100, 34, 58, 34, 49, 50, 51, 52, 53, 54, 55, 55, 55, 55, 55, 34, 44, 34, 110, 105, 99, 107, 78, 97, 109, 101, 34, 58, 34, 98, 97, 110, 99, 104, 117, 110, 34, 44, 34, 97, 118, 97, 116, 97, 114, 34, 58, 34, 34, 44, 34, 105, 109, 101, 105, 34, 58, 34, 34, 125}
	data := BuildCrc1021Sign(v)
	fmt.Println(data)
}
