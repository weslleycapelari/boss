//go:build !windows

package installer

import (
	"github.com/weslleycapelari/boss/pkg/models"
	"github.com/weslleycapelari/boss/pkg/msg"
)

func GlobalInstall(args []string, pkg *models.Package, lockedVersion bool, _ /* nosave */ bool) {
	EnsureDependency(pkg, args)
	DoInstall(pkg, lockedVersion)
	msg.Err("Cannot install global packages on this platform, only build and install local")
}
