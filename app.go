package main

import (
	"github.com/weslleycapelari/boss/cmd"
	"github.com/weslleycapelari/boss/pkg/msg"
)

func main() {
	if err := cmd.Execute(); err != nil {
		msg.Die(err.Error())
	}
}
