package cmd

import (
	"context"
	"embed"
	"fmt"
	"os"
	"simple-api-go/cmd/root"
	appContext "simple-api-go/internal/app/context"
	"simple-api-go/internal/pkg/utils/atexit"
	"simple-api-go/variable"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   variable.AppName,
	Short: variable.AppDescShort,
	Long:  variable.AppDescLong,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		envString := root.FlagEnvValue.String()
		if err = root.FlagEnvValue.Set(envString); err != nil {
			return
		}

		root.App.SetEnvironment(envString)
		root.PreStart()
		root.Start()

		return
	},
}

func init() {
	RootCmd.Flags().VarP(
		&root.FlagEnvValue,
		root.FlagEnv,
		root.FlagEnvShort,
		fmt.Sprintf(`allowed values: %s`, root.AllowedEnvInfo),
	)
}

func Execute(
	ctx context.Context,
	cancel context.CancelFunc,
	eg *errgroup.Group,
	embedFS *embed.FS,
	timeZone string,
	timeFormat string,
	timePostgresFriendlyFormat string,
) {
	root.App = appContext.NewAppContext(
		ctx,
		cancel,
		eg,
		embedFS,
		timeZone,
		timeFormat,
		timePostgresFriendlyFormat,
	)

	atexit.Add(root.App.GetCtxCancel())

	if err := RootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
