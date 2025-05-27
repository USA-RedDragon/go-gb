package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/USA-RedDragon/configulator"
	"github.com/USA-RedDragon/go-gb/internal/cartridge"
	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/emulator"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

func NewCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "go-gb",
		Version: fmt.Sprintf("%s - %s", version, commit),
		Annotations: map[string]string{
			"version": version,
			"commit":  commit,
		},
		RunE:              runRoot,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(newInteractiveCommand(version, commit))
	cmd.AddCommand(newCPUCommand(version, commit))
	return cmd
}

func runRoot(cmd *cobra.Command, _ []string) error {
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

	emu := emulator.New(cfg, cartridge)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			fmt.Println("Exiting")
			emu.Stop()
		}
	}()

	ebiten.SetWindowSize(int(cfg.Scale*160), int(cfg.Scale*144))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(cfg.Fullscreen)
	ebiten.SetScreenClearedEveryFrame(true)
	if cartridge.Title != "" {
		ebiten.SetWindowTitle(fmt.Sprintf("%s (v%d) - %s [%s] | go-gb", cartridge.Title, cartridge.Version, cartridge.Publisher, cartridge.CartridgeType))
	} else {
		ebiten.SetWindowTitle(fmt.Sprintf("%s [%s] | go-gb", filepath.Base(cfg.ROM), cartridge.CartridgeType))
	}

	return ebiten.RunGame(emu)
}
