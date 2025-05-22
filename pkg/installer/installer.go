package installer

import (
	"os"

	"github.com/hashload/boss/pkg/env"
	"github.com/hashload/boss/pkg/models"
	"github.com/hashload/boss/pkg/msg"
)

func InstallModules(args []string, lockedVersion bool, noSave bool) {
	pkg, err := models.LoadPackage(env.GetGlobal())
	if err != nil {
		if os.IsNotExist(err) {
			msg.Die("boss.json not exists in " + env.GetCurrentDir())
		} else {
			msg.Die("Fail on open dependencies file: %s", err)
		}
	}

	// Sincroniza boss.json global com o local
	_ = SyncGlobalBossFileWithLocal(pkg)

	if env.GetGlobal() {
		GlobalInstall(args, pkg, lockedVersion, noSave)
	} else {
		LocalInstall(args, pkg, lockedVersion, noSave)
	}
}

func UninstallModules(args []string, noSave bool) {
	pkg, err := models.LoadPackage(false)
	if err != nil && !os.IsNotExist(err) {
		msg.Die("Fail on open dependencies file: %s", err)
	}

	if pkg == nil {
		return
	}

	for _, arg := range args {
		dependencyRepository := ParseDependency(arg)
		pkg.UninstallDependency(dependencyRepository)
	}

	pkg.Save()

	// TODO implement remove without reinstall process
	InstallModules([]string{}, false, noSave)
}

// Sincroniza os pacotes do boss.json local com o global, sem remover os de outros projetos
type SyncResult int

const (
	SyncNone SyncResult = iota
	SyncUpdated
)

func SyncGlobalBossFileWithLocal(localPkg *models.Package) SyncResult {
	globalPath := env.GetGlobalBossFile()
	globalPkg, err := models.LoadPackageOther(globalPath)
	if err != nil && !os.IsNotExist(err) {
		msg.Err("Erro ao carregar boss.json global", err)
		return SyncNone
	}

	updated := false
	for dep, ver := range localPkg.Dependencies {
		if globalPkg.Dependencies[dep] != ver {
			globalPkg.Dependencies[dep] = ver
			updated = true
		}
	}

	if updated {
		globalPkg.Save()
		return SyncUpdated
	}
	return SyncNone
}