package util

import "testing"

func TestFormatTime(t *testing.T) {
	t.Log(FormatTime(1000))
	t.Log(FormatTime(12345))
	t.Log(FormatTime(400))
	t.Log(FormatTime(123500))
}
