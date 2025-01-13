package config

import (
	"testing"

	"github.com/kr/pretty"
)

func TestLua(t *testing.T) {
	_, err := LoadStair("housedExample.lua")
	if err != nil {
		pretty.Println(err)
		t.Fail()
	}
}
