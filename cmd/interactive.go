package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

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

	var cart *cartridge.Cartridge
	if cfg.ROM != "" {
		cart, err = cartridge.NewCartridge(cfg.ROM)
		if err != nil {
			return fmt.Errorf("failed to load cartridge: %w", err)
		}
	}

	cpu := cpu.NewSM83(cfg, cart)
	// Wait for the user to hit Enter, run the CPU step and repeat until control-C is pressed
	fmt.Println("Interactive mode started. Press Enter to step through the CPU instructions. Type exit or quit to exit.")
	for {
		var input string
		_, err := fmt.Scanln(&input) // Wait for user input
		if err != nil {
			if err.Error() == "expected newline" {
				// If the input is empty, just continue to the next step
				continue
			}
			return fmt.Errorf("failed to read input: %w", err)
		}
		if input == "exit" || input == "quit" {
			break
		}
		if strings.HasPrefix(input, "0x") {
			// Convert hex address to integer
			addr, err := strconv.ParseUint(input[2:], 16, 16)
			if err != nil {
				fmt.Printf("Invalid address: %s\n", input)
				continue
			}
			// Run Step until the PC reaches the specified address
			for cpu.GetPC() != uint16(addr) {
				cpu.Step() // Execute one CPU instruction
			}
		}
		if steps, err := strconv.Atoi(input); err == nil {
			// If input is a number, run Step that many times
			for range steps {
				cpu.Step() // Execute one CPU instruction
				slog.Debug("CPU Step executed")
			}
			continue
		}
		cpu.Step() // Execute one CPU instruction
		slog.Debug("CPU Step executed")
	}
	fmt.Println("Exiting interactive mode.")
	return nil
}
