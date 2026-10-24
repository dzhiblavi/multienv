package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	doctorservice "multienv/agent/doctor/service"
	"multienv/common/log"
	"multienv/tui/internal/app"
	"multienv/tui/internal/backend"
	"multienv/tui/internal/config"
	"multienv/tui/internal/doctorview"
	"multienv/tui/internal/navigation"
)

type options struct {
	configPath string
	noPopup    bool
	popupChild bool
	resultFile string
	resumeID   string
	doctor     bool
	mode       app.Mode
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := newCommand().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	opts := options{mode: app.Tmux}
	cmd := &cobra.Command{Use: "multienv-tui [tmux|codex|doctor]", Short: "manage workspace sessions",
		SilenceErrors: true, SilenceUsage: true, Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.mode = app.Mode(args[0])
			}

			if opts.mode != app.Tmux && opts.mode != app.Codex && opts.mode != app.Doctor {
				return fmt.Errorf("unknown view %q", opts.mode)
			}

			return run(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", config.DefaultPath, "path to TUI configuration")
	cmd.Flags().BoolVar(&opts.noPopup, "no-popup", false, "use the current pane")
	cmd.Flags().BoolVar(&opts.doctor, "doctor", false, "print global diagnostics and TUI configuration status")
	cmd.Flags().BoolVar(&opts.popupChild, "popup-child", false, "")
	cmd.Flags().StringVar(&opts.resultFile, "result-file", "", "")
	cmd.Flags().StringVar(&opts.resumeID, "resume-when-visible", "", "")
	for _, name := range []string{"popup-child", "result-file", "resume-when-visible"} {
		_ = cmd.Flags().MarkHidden(name)
	}

	return cmd
}

func run(ctx context.Context, opts options) error {
	cfg, configErr := config.Load(ctx, opts.configPath)
	if opts.doctor {
		return printDoctor(ctx, configErr)
	}

	if configErr != nil {
		return configErr
	}

	if os.Getenv("TMUX") == "" {
		return fmt.Errorf("run multienv-tui inside tmux")
	}

	return runInteractive(ctx, cfg, opts)
}

func printDoctor(ctx context.Context, configErr error) error {
	report := doctorview.WithConfig(doctorservice.Global(ctx, config.AgentPath), configErr)
	fmt.Print(doctorview.Plain(report))
	if report.Failed() {
		return fmt.Errorf("doctor found failed checks")
	}

	return nil
}

func runInteractive(ctx context.Context, cfg config.Config, opts options) (runErr error) {
	if err := os.MkdirAll(filepath.Dir(cfg.LogFile), 0700); err != nil {
		return err
	}

	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()
	log.Setup(file)
	ctx = log.InitContext(ctx)
	defer func() {
		if runErr != nil {
			log.Error(ctx, "TUI failed", "error", runErr)
		}
	}()

	if opts.resumeID != "" {
		return navigation.ResumeWhenVisible(ctx, opts.resumeID)
	}

	if cfg.Popup.Enabled && !opts.noPopup && !opts.popupChild {
		action, err := navigation.Popup(ctx, cfg.Popup, []string{"--config", opts.configPath, string(opts.mode)})
		if err != nil {
			return err
		}

		return navigation.Execute(ctx, cfg, action)
	}

	return runUI(ctx, cfg, opts)
}

func runUI(ctx context.Context, cfg config.Config, opts options) error {
	services, err := backend.New(ctx)
	if err != nil {
		return err
	}

	options := app.Options{Mode: opts.mode, ResumeInPlace: opts.noPopup && !opts.popupChild}
	action, err := app.Run(ctx, services, cfg, options)
	if err != nil {
		return err
	}

	if opts.resultFile != "" {
		return navigation.WriteResult(opts.resultFile, action)
	}

	return navigation.Execute(ctx, cfg, action)
}
