package main

import (
	"log"

	"github.com/urfave/cli/v2"
)

// application, sets the types/objects which are needed application-wide.
type application struct {
	// errorLog, error log handler.
	errorLog *log.Logger
	// infoLog, info log handler.
	infoLog *log.Logger
	// debugLog, debug log handler.
	debugLog *log.Logger
	// tui
	tui *cli.App
	// userConfigurations is the struct that stores all the user-defined
	// configuration values.
	configurations userConfigurations
}

type userConfigurations struct {
	// debugMode, if true, run the debug logger for more explicit logging.
	debugMode bool
}
