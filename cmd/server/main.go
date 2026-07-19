package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/ngothanhtung/go-tutorials/internal/app"
	"github.com/ngothanhtung/go-tutorials/internal/config"
	"github.com/ngothanhtung/go-tutorials/internal/features/seed"
)

func main() {
	var envFile string
	root := &cobra.Command{
		Use:   "server",
		Short: "go-tutorials RESTful service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(envFile)
			if err != nil {
				return err
			}
			application, err := app.New(cfg)
			if err != nil {
				return err
			}
			// seed dev admin
			if cfg.App.Env != "production" {
				_ = seed.Run(context.Background(), application.DB,
					"admin@go-tutorials.local", "Admin@123456")
			}
			return application.Run()
		},
	}
	root.Flags().StringVar(&envFile, "env", "configs/.env", "path to .env file")
	root.AddCommand(migrateCmd())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "run database migrations (up)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return exec.Command("bash", "scripts/migrate.sh", "up", "configs/.env").Run()
		},
	}
}
