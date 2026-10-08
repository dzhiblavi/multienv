package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/common/util"
	"multienv/server/internal/codex"
	"multienv/server/internal/httpserver"
	"multienv/server/internal/storage"
	redisstorage "multienv/server/internal/storage/redis"
	"multienv/server/internal/tmux"
)

type App struct {
	server       httpserver.Server
	tmuxStorage  storage.Storage
	codexStorage storage.Storage
}

func (a App) Close() {
	if a.tmuxStorage != nil {
		a.tmuxStorage.Close()
	}
	if a.codexStorage != nil {
		a.codexStorage.Close()
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

			tmuxStorage, err := redisstorage.New(ctx, config.Storage, "@tmux")
			if err != nil {
				return fmt.Errorf("create tmux storage: %w", err)
			}
			app.tmuxStorage = tmuxStorage
			if err := tmuxStorage.Init(ctx); err != nil {
				app.Close()
				return fmt.Errorf("initialize tmux storage: %w", err)
			}
			tmuxDispatcher := tmux.New(config.Tmux, tmuxStorage)

			codexStorage, err := redisstorage.New(ctx, config.Storage, "@codex")
			if err != nil {
				app.Close()
				return fmt.Errorf("create codex storage: %w", err)
			}
			app.codexStorage = codexStorage
			if err := codexStorage.Init(ctx); err != nil {
				app.Close()
				return fmt.Errorf("initialize codex storage: %w", err)
			}
			codexDispatcher := codex.New(config.Codex, codexStorage)

			app.server = httpserver.New(config.Server)

			if err := app.server.Install("tmux", tmuxDispatcher); err != nil {
				app.Close()
				return fmt.Errorf("install tmux dispatcher: %w", err)
			}

			if err := app.server.Install("codex", codexDispatcher); err != nil {
				app.Close()
				return fmt.Errorf("install codex dispatcher: %w", err)
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
