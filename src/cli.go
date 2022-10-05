package main

import (
	c "github.com/urfave/cli/v2"
)

func Before(ctx *c.Context) error {
	SetupDatabase()
	return nil
}

func Shutdown() {
	CloseDatabase()
}

func After(ctx *c.Context) error {
	Shutdown()
	return nil
}

func NewCli() *c.App {
	return &c.App{
		Name:                 "Go-API-CLI",
		Usage:                "Simple and Fast!",
		EnableBashCompletion: true,
		Commands: []*c.Command{
			{
				Name:    "runserver",
				Aliases: []string{"r", "run"},
				Usage:   "Start web server",
				Action: func(ctx *c.Context) error {
					return NewServer().Run(MustGetEnv("HOST", "127.0.0.1") + ":" + MustGetEnv("PORT", "3333"))
				},
				Before: Before,
				After:  After,
			},
			{
				Name:  "migration",
				Usage: "Database migrations",
				Subcommands: []*c.Command{
					{
						Name:      "make",
						Usage:     "Create blank migration file",
						ArgsUsage: "[file_name]",
						Action: func(ctx *c.Context) error {
							fileName := ctx.Args().First()
							if fileName == "" {
								return c.Exit("Migration file name not specified", 1)
							}
							MakeMigration(fileName)
							return nil
						},
					},
					{
						Name:  "migrate",
						Usage: "Run migrations",
						Action: func(ctx *c.Context) error {
							RunMigration()
							return nil
						},
					},
				},
			},
		},
	}
}
