package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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

// run, runs the main application, by encapsulating the main applications
// methods, it makes the application more testable.
func (app *application) run() {
	// Check the pf rules before enabling pf, if the rules have a problem return
	// before configuring the system any further.
	outStr, err := app.checkRuleSet("./pf.conf")
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	app.infoLog.Println("pf config file successfully passed syntax test.")
	if outStr != "" {
		app.infoLog.Println(outStr)
	}

	if err := app.rcConfiguration(); err != nil {
		app.errorLog.Fatalln(err)
	}

	// Copy local pf rule set to /etc/pf.conf
	if err := copyFile("./pf.conf", "/etc/pf.conf"); err != nil {
		err = fmt.Errorf("error while copying local pf rule set to /etc/pf.conf: %w", err)
		app.errorLog.Fatalln(err)
	}
	app.infoLog.Println("Copied rule set to /etc/pf.conf.")

	app.infoLog.Println("Rebooting system to properly enable pf.")
	// A reboot is necessary after configuring pf for the first time.
	if err := app.reboot(); err != nil {
		app.errorLog.Fatalln(err)
	}
}

// parseFlags, parses any flags if they are present.
func (app *application) parseFlags() {
	flag.BoolVar(&app.configurations.debugMode, "debugMode", false, "Debug mode activates the debug logger.")
	flag.Parse()
}

// TODO: move reboot and copyFile to internal

// reboot, reboots the system. Required after activating pf for the first time.
func (app *application) reboot() error {
	_, err := app.shCmd("reboot")
	if err != nil {
		err = fmt.Errorf("reboot attempt failed: %w", err)
		return err
	}
	// The program should never come this far, since it would reboot before.
	return nil
}

// copyFile the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error while trying to open file ('%s'): %w", src, err)
	}
	defer in.Close()

	// The 'dst' file will be created, or truncated if it already exists
	// (overwritten). 'dst' file has file mode 0666.
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error while trying to create file ('%s'): %w", dst, err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("error while trying to copy data from file '%s' to file '%s': %w", src, dst, err)
	}

	return nil
}
