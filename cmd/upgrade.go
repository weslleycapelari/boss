package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weslleycapelari/boss/internal/upgrade"
)

func upgradeCmdRegister(root *cobra.Command) {
	var preRelease bool

	var upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the client version",
		Example: `  Upgrade boss:
  boss upgrade

  Upgrade boss with pre-release:
  boss upgrade --dev`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return upgrade.BossUpgrade(preRelease)
		},
	}

	root.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVar(&preRelease, "dev", false, "pre-release")
}
