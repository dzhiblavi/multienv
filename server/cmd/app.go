package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"multienv/common/util"
	"multienv/server/internal/doctor"
	"multienv/server/internal/host"
	"multienv/server/internal/httpserver"
	"multienv/server/internal/master"
	"multienv/server/internal/storage"
	redisstorage "multienv/server/internal/storage/redis"
)

type App struct {
	server  httpserver.Server
	storage io.Closer
}

func (a App) Close() {
	if a.storage != nil {
		a.storage.Close()
	}
}

func newRootCmd() *cobra.Command {
	var configPath string
	app := App{}

	cmd := &cobra.Command{
		Use:           "multienv-server",
		Short:         "coordinate sessions across connected workspace",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			defer app.Close()
			return app.server.Serve(cmd.Context())
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var config Config
			err := util.ReadYaml(ctx, configPath, &config)
			if err != nil {
				return fmt.Errorf("failed to parse config: %v", err)
			}

			data, err := redisstorage.New(ctx, config.Storage)
			if err != nil {
				return fmt.Errorf("create storage: %w", err)
			}

			app.storage = data
			hostDispatcher, err := host.New(config.Hosts)
			if err != nil {
				app.Close()
				return fmt.Errorf("configure hosts: %w", err)
			}

			diagnostics, err := doctor.New(hostDispatcher.Hostnames(), data, doctor.AgentRunner{})
			if err != nil {
				app.Close()
				return fmt.Errorf("create doctor: %w", err)
			}

			app.server = httpserver.New(config.Server)
			if err := app.server.Install("storage", storage.NewDispatcher(data)); err != nil {
				app.Close()
				return fmt.Errorf("install storage dispatcher: %w", err)
			}

			if err := app.server.Install("doctor", doctor.NewDispatcher(diagnostics)); err != nil {
				app.Close()
				return fmt.Errorf("install doctor: %w", err)
			}

			if err := app.server.Install("master", master.New()); err != nil {
				app.Close()
				return fmt.Errorf("install master dispatcher: %w", err)
			}

			if err := app.server.Install("host", hostDispatcher); err != nil {
				app.Close()
				return fmt.Errorf("install host dispatcher: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"~/.config/multienv/server.yaml",
		"path to configuration file",
	)

	return cmd
}
