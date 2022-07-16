package main

import (
	"fmt"
)

// rcEnablePF, enable pf in /etc/rc.conf. So that it starts at boot.
func (app *application) rcEnablePF() (string, error) {
	// Check if pf has already been enabled, in that case, just return without
	// altering the /etc/rc.conf file again.
	// -n returns only the value of a variable in the rc.conf file.
	cmdOut, err := app.shCmd("sysrc", "-n", "pf_enable")
	if err != nil {
		// An error can be returned, when the variable does not exist yet in
		// the rc.conf file. So, do not inmediately return from the method,
		// try to still enable the variable.
		err = fmt.Errorf("error while checking value of pf_enable in rc.conf: %w", err)
		app.errorLog.Println(err)
	}
	// If no error is returned, check if the value is already 'YES', in that
	// case, return without further altering the rc.conf file.
	if err == nil {
		// cmdOut has a newline, to compare it properly, scan the result string
		// out of the whole string first into pfRCValue.
		var pfRCValue string
		fmt.Sscan(cmdOut, &pfRCValue)

		if pfRCValue == "YES" {
			app.infoLog.Println("pf is already enabled in /etc/rc.conf")
			return "", nil
		}
	}
	// If pf_enable has not been set to 'YES' yet, enable it.
	cmdOut, err = app.shCmd("sysrc", "pf_enable=YES")
	if err != nil {
		return "", fmt.Errorf("unable to enable pf in /etc/rc.conf: %w", err)
	}

	app.infoLog.Println("pf was successfully enabled in /etc/rc.conf")

	return cmdOut, nil

}

// setupRulesFile, establishes '/etc/pf.conf' as the rules file for pf.
func (app *application) setupRulesFile() (string, error) {
	cmdOut, err := app.shCmd("sysrc", "pf_rules=/etc/pf.conf")
	if err != nil {
		return "", fmt.Errorf("unable to setup /etc/pf.conf as the rules file for pf: %w", err)
	}

	app.infoLog.Println("pf_rules were successfully setup to be in /etc/pf.conf")

	return cmdOut, nil

}

// setupLogFile, establishes '/var/log/pflog' as the log file for pf.
func (app *application) setupLogFile() (string, error) {
	cmdOut, err := app.shCmd("sysrc", "pflog_logfile=/var/log/pflog")
	if err != nil {
		return "", fmt.Errorf("unable to setup /var/log/pflog as the log file for pflog: %w", err)
	}

	app.infoLog.Println("pflog file was successfully setup to be in /var/log/pflog")

	return cmdOut, nil

}

// enablePflog, enable pflog in /etc/rc.conf. So that it starts at boot.
func (app *application) enablePflog() (string, error) {
	// Check if pflog has already been enabled, in that case, just return without
	// altering the /etc/rc.conf file again.
	// -n returns only the value of a variable in the rc.conf file.
	cmdOut, err := app.shCmd("sysrc", "-n", "pflog_enable")
	if err != nil {
		// An error can be returned, when the variable does not exist yet in
		// the rc.conf file. So, do not inmediately return from the method,
		// try to still enable the variable.
		err = fmt.Errorf("error while checking value of pflog_enable in rc.conf: %w", err)
		app.errorLog.Println(err)
	}
	// If no error is returned, check if the value is already 'YES', in that
	// case, return without further altering the rc.conf file.
	if err == nil {
		// cmdOut has a newline, to compare it properly, scan the result string
		// out of the whole string first into pflogRCValue.
		var pflogRCValue string
		fmt.Sscan(cmdOut, &pflogRCValue)

		if pflogRCValue == "YES" {
			app.infoLog.Println("pflog is already enabled in /etc/rc.conf")
			return "", nil
		}
	}
	// If pflog_enable has not been set to 'YES' yet, enable it.
	cmdOut, err = app.shCmd("sysrc", "pflog_enable=YES")
	if err != nil {
		return "", fmt.Errorf("unable to enable pflog in /etc/rc.conf: %w", err)
	}

	app.infoLog.Println("pflog was successfully enabled in /etc/rc.conf")

	return cmdOut, nil

}

// checkRuleSet, runs `pfctl -n` to check the pf rules of a given file.
func (app *application) checkRuleSet(file string) (string, error) {
	// -n checks rules of -f file.
	cmdOut, err := app.shCmd("pfctl", "-nf", file)
	if err != nil {
		return "", fmt.Errorf("error checking the rules of pf file: %w", err)
	}
	// If the file does not exist, pfctl will throw an error, so it is unnecesa-
	// ry to check for the existance of the file.

	return cmdOut, nil

}

// activateRules, activates the rules in a file as the new pf rule set.
func (app *application) activateRules(file string) (string, error) {
	// Activate given file as new rule set.
	cmdOut, err := app.shCmd("pfctl", "-f", file)
	if err != nil {
		return "", fmt.Errorf("error activating new pf rule set from file (%s): %w", file, err)
	}
	// If the file does not exist, pfctl will throw an error, so it is unnecesa-
	// ry to check for the existance of the file.

	return cmdOut, nil

}

// rcConfiguration, does all the required configurations on /etc/rc.conf to have pf
// working after rebooting the system.
func (app *application) rcConfiguration() error {
	// Enable pf in rc.conf. After enabling pf, the default pf stance is to
	// accept all connections, so one will not be locked out of the SSH
	// connection with the server.
	outStr, err := app.rcEnablePF()
	if err != nil {
		return err
	}
	if outStr != "" {
		app.infoLog.Println(outStr)
	}

	// Let '/etc/pf.conf' be the rules file for pf.
	outStr, err = app.setupRulesFile()
	if err != nil {
		return err
	}
	if outStr != "" {
		app.infoLog.Println(outStr)
	}

	outStr, err = app.enablePflog()
	if err != nil {
		return err
	}
	if outStr != "" {
		app.infoLog.Println(outStr)
	}

	// Let '/var/log/pflog' be the log file for pflog.
	outStr, err = app.setupLogFile()
	if err != nil {
		return err
	}
	if outStr != "" {
		app.infoLog.Println(outStr)
	}

	return nil
}

// enablePF, enables PF. The firewall starts filtering packets, it is just as
// running 'pfctl -e'.
func (app *application) enablePF() (string, error) {
	// Check first if pfctl is already running, because, otherwise if it is
	// already running 'pfctl -e' returns an error.
	_, err := app.shCmd("pfctl", "-s Running")
	// 'pfctl -s Running' returns no error if pfctl is already running.
	if err == nil {
		return "pfctl is already running.", nil
	}

	outStr, err := app.shCmd("pfctl", "-e")
	if err != nil {
		return "", fmt.Errorf("error enabling pf: %w", err)
	}

	return outStr, nil

}
