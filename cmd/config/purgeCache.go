package config

import (
	"github.com/spf13/cobra"
	"github.com/weslleycapelari/boss/pkg/gc"
)

func RegisterCmd(cmd *cobra.Command) {
	purgeCacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Configure cache",
	}

	rmCacheCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove cache",
		RunE: func(_ *cobra.Command, _ []string) error {
			return gc.RunGC(true)
		},
	}

	purgeCacheCmd.AddCommand(rmCacheCmd)

	cmd.AddCommand(purgeCacheCmd)
}
