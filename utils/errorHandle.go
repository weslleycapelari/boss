package utils

import "github.com/hashload/boss/pkg/msg"

func HandleError(err error) {
	if err != nil {
		msg.Err(err.Error())
	}
}

func HandleErrorFatal(err error) {
	if err != nil {
		msg.Die(err.Error())
	}
}
