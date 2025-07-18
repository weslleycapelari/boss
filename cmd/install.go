package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weslleycapelari/boss/pkg/installer"
)

func installCmdRegister(root *cobra.Command) {
	var noSaveInstall bool

	var installCmd = &cobra.Command{
		Use:     "install",
		Short:   "Install a new dependency",
		Long:    `This command install a new dependency on your project`,
		Aliases: []string{"i", "add"},
		Example: `  Add a new dependency:
  boss install <pkg>

  Add a new version-specific dependency:
  boss install <pkg>@<version>

  Install a dependency without add it from the boss.json file:
  boss install <pkg> --no-save`,
		Run: func(_ *cobra.Command, args []string) {
			installer.InstallModules(args, true, noSaveInstall)
		},
	}

	root.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&noSaveInstall, "no-save", false, "prevents saving to dependencies")
}
