//go:build !windows
// +build !windows

package librarypath

import (
	"github.com/weslleycapelari/boss/pkg/msg"
)

func updateGlobalLibraryPath() {
	msg.Warn("updateGlobalLibraryPath not implemented on this platform")
}

func updateGlobalBrowsingByProject(_ string, _ bool) {
	msg.Warn("updateGlobalBrowsingByProject not implemented on this platform")
}
