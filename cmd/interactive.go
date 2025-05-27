package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/USA-RedDragon/configulator"
	"github.com/USA-RedDragon/go-gb/internal/cartridge"
	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/cpu"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

func newInteractiveCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interactive",
		Version: fmt.Sprintf("%s - %s", version, commit),
		Annotations: map[string]string{
			"version": version,
			"commit":  commit,
		},
		RunE:              runInteractive,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}
	return cmd
}

func runInteractive(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	fmt.Printf("go-gb - %s (%s)\n", cmd.Annotations["version"], cmd.Annotations["commit"])

	c, err := configulator.FromContext[config.Config](ctx)
	if err != nil {
		return fmt.Errorf("failed to get config from context")
	}

	cfg, err := c.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var logger *slog.Logger
	switch cfg.LogLevel {
	case config.LogLevelDebug:
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	case config.LogLevelInfo:
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelInfo}))
	case config.LogLevelWarn:
		logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelWarn}))
	case config.LogLevelError:
		logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelError}))
	}
	slog.SetDefault(logger)

	cartridge, err := cartridge.NewCartridge(cfg.ROM)
	if err != nil {
		return fmt.Errorf("failed to load cartridge: %w", err)
	}

	cpu := cpu.NewSM83(cfg, cartridge)
	// Wait for the user to hit Enter, run the CPU step and repeat until control-C is pressed
	fmt.Println("Interactive mode started. Press Enter to step through the CPU instructions. Type exit or quit to exit.")
	for {
		var input string
		fmt.Scanln(&input) // Wait for user input
		if input == "exit" || input == "quit" {
			break
		}
		cpu.Step() // Execute one CPU instruction
		slog.Debug("CPU Step executed")
	}
	fmt.Println("Exiting interactive mode.")
	return nil
}
