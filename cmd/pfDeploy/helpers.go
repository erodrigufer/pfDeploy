package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erodrigufer/pfDeploy/internal/sysutils"
	"github.com/erodrigufer/pfDeploy/pfSetup"

	"github.com/urfave/cli/v2"
)

// setupApplication, it configures all needed general parameters for the
// application.
func (app *application) setupApplication() {

	// Create a logger for INFO messages, the prefix "INFO" and a tab will be
	// displayed before each log message. The flags Ldate and Ltime provide the
	// local date and time.
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an ERROR messages logger, addiotionally use the Lshortfile flag to
	// display the file's name and line number for the error.
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// if app.configurations.debugMode {
	// 	app.debugLog = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	// } else {
	// 	// Discard the output of the debug logger.
	// 	app.debugLog = log.New(io.Discard, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	// }
}

func (app *application) runTUI() {
	app.setupCLI()

	if err := app.tui.Run(os.Args); err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *application) setupCLI() {
	app.tui = &cli.App{
		Name:  "pfDeploy",
		Usage: "Automatically setup pf in your new deployment.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "./pf.conf",
				Usage:   "`PATH` to the file used as the new pf rule set.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			app.run()
			return nil
		},
	}

}

// run, runs the main application, by encapsulating the main application's
// methods, it makes the application more testable.
func (app *application) run() {
	// Check the pf rules before enabling pf, if the rules have a problem return
	// before configuring the system any further.
	outStr, err := pfSetup.CheckRuleSet("./pf.conf")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.infoLog.Print("pf config file successfully passed syntax test.")
	if outStr != "" {
		app.infoLog.Print(outStr)
	}

	if err := pfSetup.PFSetup(app.infoLog); err != nil {
		app.errorLog.Fatal(err)
	}

	// Copy and configure file attributions of new pf rules file.
	if err := app.initializeRulesFile(); err != nil {
		app.errorLog.Fatal(err)
	}

	app.infoLog.Print("Rebooting system to properly enable pf.")
	// A reboot is necessary after configuring pf for the first time.
	if err := sysutils.Reboot(); err != nil {
		app.errorLog.Fatal(err)
	}
}

// parseFlags, parses any flags if they are present.
// func (app *application) parseFlags() {
// 	flag.BoolVar(&app.configurations.debugMode, "debugMode", false, "Debug mode activates the debug logger.")
// 	flag.Parse()
// }

// initializeRulesFile, copies the given pf rules file to its standard path,
// and sets the ownership and file attributes properly.
func (app *application) initializeRulesFile() error {
	// Destination for pf rule set.
	des := "/etc/pf.conf"

	// Copy local pf rule set to /etc/pf.conf, this does not change the
	// ownership and file attributions of new file.
	if err := sysutils.CopyFile("./pf.conf", des); err != nil {
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
