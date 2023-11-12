package main

import (
	"fmt"
	"os"

	"github.com/lvjp/s3impl/internal/app"
	"github.com/spf13/cobra"
)

var configPath string

var cmd = &cobra.Command{
	Use:                   "s3impl",
	Short:                 "s3impl is a simple S3 implementation",
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	SilenceUsage:          true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.Execute(configPath)
	},
}

func init() {
	cmd.Flags().StringVar(&configPath, "config", "examples/config.yaml", "Path to configuration file")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
