package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var logFile string

func logPreRun(cmd *cobra.Command, args []string) error {
	// log output to file
	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
			return fmt.Errorf("failed to create log dir %s %w", filepath.Base(logFile), err)
		}
		os.Remove(logFile)
		logFileWriter, err := os.OpenFile(
			logFile,
			os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			return fmt.Errorf("failed to create log file %s %w", logFile, err)
		}
		logger := slog.New(slog.NewTextHandler(logFileWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
		slog.Info("logging enabled")
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "advent-of-code-2024",
	Short:             "advent-of-code solutions for 2024",
	PersistentPreRunE: logPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		// Show usage
		cmd.Help()
		os.Exit(1)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		println(err)
		os.Exit(1)
	}
}

func init() {
	// all commands have debug mode
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "", "tmp/advent.log", "log file to send structured logs to")
}
