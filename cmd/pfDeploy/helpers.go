package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/erodrigufer/pfDeploy/internal/sysutils"
	"github.com/erodrigufer/pfDeploy/pfSetup"
)

// setupApplication, it configures all needed general parameters for the
// application.
func (app *application) setupApplication() {
	app.parseFlags()

	// Create a logger for INFO messages, the prefix "INFO" and a tab will be
	// displayed before each log message. The flags Ldate and Ltime provide the
	// local date and time.
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an ERROR messages logger, addiotionally use the Lshortfile flag to
	// display the file's name and line number for the error.
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if app.configurations.debugMode {
		app.debugLog = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// Discard the output of the debug logger.
		app.debugLog = log.New(io.Discard, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// run, runs the main application, by encapsulating the main application's
// methods, it makes the application more testable.
func (app *application) run() {
	// Check the pf rules before enabling pf, if the rules have a problem return
	// before configuring the system any further.
	outStr, err := pfSetup.CheckRuleSet("./pf.conf")
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	app.infoLog.Println("pf config file successfully passed syntax test.")
	if outStr != "" {
		app.infoLog.Println(outStr)
	}

	if err := pfSetup.PFSetup(app.infoLog); err != nil {
		app.errorLog.Fatalln(err)
	}

	// Destination for pf rule set.
	des := "/etc/pf.conf"
	// Copy local pf rule set to /etc/pf.conf
	if err := sysutils.CopyFile("./pf.conf", des); err != nil {
		err = fmt.Errorf("error while copying local pf rule set to %s : %w", des, err)
		app.errorLog.Fatalln(err)
	}
	app.infoLog.Println("Copied rule set to /etc/pf.conf.")

	// Change file mod and owners of new rule set file.
	if err := os.Chmod(des, 0644); err != nil {
		err = fmt.Errorf("error while changing file mod of %s : %w", des, err)
		app.errorLog.Fatalln(err)
	}
	// uid=root(0); gid=wheel(0)
	if err := os.Chown(des, 0, 0); err != nil {
		err = fmt.Errorf("error while changing file owners of %s : %w", des, err)
		app.errorLog.Fatalln(err)
	}
	app.infoLog.Println("Rebooting system to properly enable pf.")
	// A reboot is necessary after configuring pf for the first time.
	if err := sysutils.Reboot(); err != nil {
		app.errorLog.Fatalln(err)
	}
}

// parseFlags, parses any flags if they are present.
func (app *application) parseFlags() {
	flag.BoolVar(&app.configurations.debugMode, "debugMode", false, "Debug mode activates the debug logger.")
	flag.Parse()
}
