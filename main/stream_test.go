package main

import (
	"context"
	"strings"
	"testing"
)

func TestNewUnCmdStream(t *testing.T) {
	r := strings.NewReader("hogehoge\n")
	in := NewUnCmdStream(NewCmdStreamFile(r))
	ctx := context.Background()
	line, err := in.ReadLine(&ctx)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}
	if line != "hogehoge" {
		t.Log("fail on normal case")
		t.Fail()
		return
	}
	in.UnreadLine("uhauha")
	line, err = in.ReadLine(&ctx)
	if err != nil || line != "uhauha" {
		t.Log("fail on unread case")
		t.Fail()
		return
	}
	return
}
