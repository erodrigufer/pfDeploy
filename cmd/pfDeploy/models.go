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
	// tui
	tui *cli.App
}
