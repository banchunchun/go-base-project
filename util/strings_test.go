package util

import (
	"strings"
	"testing"
)

func TestTrim(t *testing.T) {
	s := "    "
	s = strings.TrimSpace(s)
	t.Log(s)
	s = "  \r \t \n aa  "
	s = strings.TrimSpace(s)
	t.Log(s)
}

func TestGetCallerFrame(t *testing.T) {
	frame, _ := GetCallerFrame(0)
	t.Logf("%s:%d", frame.File, frame.Line)
	t.Logf("%+v", frame)
}

func TestFormatMessage(t *testing.T) {
	t.Log(FormatMessage("%+v %d", t, 1))
	t.Log(FormatMessage("hello"))
	t.Log(FormatMessage("", t))
	t.Log(FormatMessage("", t, 4))
}

func TestSplitString(t *testing.T) {
	t.Log(SplitString("  , asdf, asdfd, ", ","))
}

func TestSplitMultiLineString(t *testing.T) {
	t.Log(SplitMultiLineString("wyg\r\njyw\r\n", false))
	t.Log(SplitMultiLineString("wyg\r\njyw\r\n", true))
}
