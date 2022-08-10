package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erodrigufer/pfDeploy/internal/sysutils"
	"github.com/erodrigufer/pfDeploy/pfSetup"

	"github.com/urfave/cli/v2"
)

// setupApplication, configures all needed general parameters for the
// application.
func (app *application) setupApplication() {

	// Create a logger for INFO messages, the prefix "INFO" and a tab will be
	// displayed before each log message. The flags Ldate and Ltime provide the
	// local date and time.
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an ERROR messages logger, addiotionally use the Lshortfile flag to
	// display the file's name and line number for the error.
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}

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
	app.tui = &cli.App{
		Name:  "pfDeploy",
		Usage: "Automatically setup pf in your new deployment.",
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
						Name:    "reboot",
						Aliases: []string{"r"},
						Value:   true,
						Usage:   "Reboot host after deploying new pf configuration.",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if err := app.deploy(cCtx.String("file"), cCtx.Bool("reboot")); err != nil {
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

// deploy, runs the main application, by encapsulating the main application's
// methods, it makes the application more testable.
func (app *application) deploy(filePath string, rebootFlag bool) error {
	// Check the pf rules before enabling pf, if the rules have a problem return
	// before configuring the system any further.
	if err := app.checkRuleSet(filePath); err != nil {
		return fmt.Errorf("error while checking validity of pf ruleset file %s: %w", filePath, err)
	}

	if err := pfSetup.PFSetup(app.infoLog); err != nil {
		return fmt.Errorf("error while configuring rc: %w", err)
	}

	// Copy and configure file attributions of new pf rules file.
	if err := app.initializeRulesFile(filePath); err != nil {
		return fmt.Errorf("error while copying and configuring the pf ruleset file: %w", err)
	}

	if rebootFlag {
		app.infoLog.Print("Rebooting system to properly enable pf.")
		// A reboot is necessary after configuring pf for the first time.
		if err := sysutils.Reboot(); err != nil {
			return fmt.Errorf("error while rebooting: %w", err)
		}
	}
	return nil
}

// initializeRulesFile, copies the given pf rules file to its standard path,
// and sets the ownership and file attributes properly.
func (app *application) initializeRulesFile(localFile string) error {
	// Destination for pf rule set.
	des := "/etc/pf.conf"

	// Copy local pf rule set to /etc/pf.conf, this does not change the
	// ownership and file attributions of new file.
	if err := sysutils.CopyFile(localFile, des); err != nil {
		err = fmt.Errorf("error while copying local pf rule set to %s : %w", des, err)
		return err
	}
	app.infoLog.Printf("Succesfully copied new rule set to %s", des)

	// Change file mod and owners of new rule set file.
	if err := os.Chmod(des, 0644); err != nil {
		err = fmt.Errorf("error while changing file mod of %s : %w", des, err)
		return err
	}
	// File owned by root and in group wheel: uid=root(0); gid=wheel(0).
	if err := os.Chown(des, 0, 0); err != nil {
		err = fmt.Errorf("error while changing file owners of %s : %w", des, err)
		return err
	}

	return nil
}

// checkRuleSet, check the syntax validity of a pf ruleset file.
func (app *application) checkRuleSet(filePath string) error {
	outStr, err := pfSetup.CheckRuleSet(filePath)
	if err != nil {
		return fmt.Errorf("error while checking syntax validity of ruleset file %s: %w", filePath, err)
	}
	app.infoLog.Print("pf config file successfully passed syntax test.")
	if outStr != "" {
		app.infoLog.Print(outStr)
	}
	return nil

}
