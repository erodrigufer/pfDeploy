package main

import (
	"fmt"
	"os"

	"github.com/erodrigufer/pfDeploy/internal/version"
	"github.com/urfave/cli/v2"
)

// runTUI, run the TUI (Terminal User Interface), handled by the CLI package.
func (app *application) runTUI() {
	app.setupCLI()

	if err := app.tui.Run(os.Args); err != nil {
		app.errorLog.Fatal(err)
	}
}

// setupCLI, configure and initialize all commands, flags and options of
// the TUI.
func (app *application) setupCLI() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "revision",
		Aliases: []string{"r"},
		Usage:   "Print the VCS revision when the binary was built (if the binary is a modified version of a commit the suffix 'dirty' will be added to the revision hash).",
	}
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf("revision=%s\n", version.GetRevision())
	}

	app.tui = &cli.App{
		Name:    "pfDeploy",
		Version: version.GetRevision(),
		Usage:   "Automatically setup pf in your new deployment.",
		// This options enables short flag abbreviations to be merged into a
		// single flag with a '-' prefix.
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "deploy",
				Usage: "Setup pf and pflog at boot and deploy a new pf rules file.",
				// Flags for deploy command.
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "`PATH` to the file used as the new pf rule set.",
					},
					&cli.BoolFlag{
						Name:  "no-reboot",
						Usage: "Do not reboot the host after deploying new pf configuration.",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if err := app.deploy(cCtx.String("file"), cCtx.Bool("no-reboot")); err != nil {
						err = fmt.Errorf("error while executing 'deploy' command: %w", err)
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "check",
				Usage: "Check the syntactical validity of a pf rules file.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "`PATH` to the file used as the pf rule set.",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Read input for -file flag.
					filePath := cCtx.String("file")
					if err := app.checkRuleSet(filePath); err != nil {
						err = fmt.Errorf("error while executing 'check' command: %w", err)
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
		},
	}

}
