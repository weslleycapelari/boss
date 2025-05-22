package env

import (
	//nolint:gosec // We are not using this for security purposes
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashload/boss/pkg/consts"
	"github.com/hashload/boss/pkg/msg"
	"github.com/hashload/boss/utils/dcc32"
	"github.com/mitchellh/go-homedir"
)

//nolint:gochecknoglobals //TODO: Refactor this
var (
	global                 bool
	internal               = false
	globalConfiguration, _ = LoadConfiguration(GetBossHome())
)

func SetGlobal(b bool) {
	global = b
}

func SetInternal(b bool) {
	internal = b
}

func GetInternal() bool {
	return internal
}

func GetGlobal() bool {
	return global
}

func GlobalConfiguration() *Configuration {
	return globalConfiguration
}

func HashDelphiPath() string {
	//nolint:gosec // We are not using this for security purposes
	hasher := md5.New()
	hasher.Write([]byte(strings.ToLower(GlobalConfiguration().DelphiPath)))
	hashString := hex.EncodeToString(hasher.Sum(nil))
	if internal {
		hashString = consts.BossInternalDir + hashString
	}
	return hashString
}

func GetInternalGlobalDir() string {
	internalOld := internal
	internal = true
	result := filepath.Join(GetBossHome(), consts.FolderDependencies, HashDelphiPath())
	internal = internalOld
	return result
}

func getwd() string {
	if global {
		return filepath.Join(GetBossHome(), consts.FolderDependencies, HashDelphiPath())
	}

	dir, err := os.Getwd()
	if err != nil {
		msg.Err("Error to get paths", err)
		return ""
	}

	return dir
}

func GetCacheDir() string {
	return filepath.Join(GetBossHome(), "cache")
}

func GetBossHome() string {
	homeDir := os.Getenv("FIORILLI_BOSS_HOME")

	if homeDir == "" {
		systemHome, err := homedir.Dir()
		homeDir = systemHome
		if err != nil {
			msg.Err("Error to get cache paths", err)
		}

		homeDir = filepath.FromSlash(homeDir)
	}
	return filepath.Join(GetProjectFolder(consts.FolderProject), consts.FolderRefactor, consts.FolderBossHome)
	// return filepath.Join(homeDir, consts.FolderBossHome)
}

func GetBossFile() string {
	return filepath.Join(GetCurrentDir(), consts.FilePackage)
}

func GetGlobalBossFile() string {
	return filepath.Join(GetProjectFolder(consts.FolderProject), consts.FolderRefactor, consts.FilePackage)
}

func GetModulesDir() string {
	return filepath.Join(GetProjectFolder(consts.FolderProject), consts.FolderRefactor, consts.FolderDependencies)
	//return filepath.Join(GetCurrentDir(), consts.FolderDependencies)
}

func GetCurrentDir() string {
	return getwd()
}

func GetGlobalEnvBpl() string {
	return filepath.Join(GetBossHome(), consts.FolderEnvBpl)
}
func GetGlobalEnvDcp() string {
	return filepath.Join(GetBossHome(), consts.FolderEnvDcp)
}
func GetGlobalEnvDcu() string {
	return filepath.Join(GetBossHome(), consts.FolderEnvDcu)
}

func GetGlobalBinPath() string {
	return filepath.Join(GetBossHome(), consts.FolderDependencies, consts.BinFolder)
}

func GetDcc32Dir() string {
	if GlobalConfiguration().DelphiPath != "" {
		return GlobalConfiguration().DelphiPath
	}

	byCmd := dcc32.GetDcc32DirByCmd()
	if len(byCmd) > 0 {
		return byCmd[0]
	}

	return ""
}

func GetProjectFolder(target string) string {
	dir, _ := os.Getwd()
	for {
		// Caminho potencial da pasta alvo
		potential := filepath.Join(dir, target)

		if stat, err := os.Stat(potential); err == nil && stat.IsDir() {
			return potential
		}

		// Sobe um diretório
		parent := filepath.Dir(dir)
		if parent == dir {
			// Chegou na raiz e não encontrou
			break
		}
		dir = parent
	}
	return ""
}

