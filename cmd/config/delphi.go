package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/weslleycapelari/boss/pkg/env"
	"github.com/weslleycapelari/boss/pkg/msg"
	"github.com/weslleycapelari/boss/utils/dcc32"
)

func delphiCmd(root *cobra.Command) {
	delphiCmd := &cobra.Command{
		Use:   "delphi",
		Short: "Configure Delphi version",
		Long:  `Configure Delphi version to compile modules`,
		Run: func(cmd *cobra.Command, _ []string) {
			msg.Info("Running in path %s", env.GlobalConfiguration().DelphiPath)
			_ = cmd.Usage()
		},
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List Delphi versions",
		Long:  `List Delphi versions to compile modules`,
		Run: func(_ *cobra.Command, _ []string) {
			paths := dcc32.GetDcc32DirByCmd()
			if len(paths) == 0 {
				msg.Warn("Installations not found in $PATH")
				return
			}

			msg.Warn("Installations found:")
			for index, path := range paths {
				msg.Info("  [%d] %s", index, path)
			}
		},
	}

	use := &cobra.Command{
		Use:   "use [path]",
		Short: "Use Delphi version",
		Long:  `Use Delphi version to compile modules`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			if _, err := strconv.Atoi(args[0]); err != nil {
				if _, err = os.Stat(args[0]); os.IsNotExist(err) {
					return errors.New("invalid path")
				}
			}

			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			var path = args[0]
			config := env.GlobalConfiguration()
			if index, err := strconv.Atoi(path); err == nil {
				delphiPaths := dcc32.GetDcc32DirByCmd()
				config.DelphiPath = delphiPaths[index]
			} else {
				config.DelphiPath = args[0]
			}

			config.SaveConfiguration()
			msg.Info("Successful!")
		},
	}

	root.AddCommand(delphiCmd)

	delphiCmd.AddCommand(list)
	delphiCmd.AddCommand(use)
}
