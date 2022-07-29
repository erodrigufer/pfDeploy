package pfSetup

import (
	"fmt"
	"log"

	"github.com/erodrigufer/pfDeploy/internal/sysutils"
)

// RCEnablePF, enable pf in /etc/rc.conf. So that it starts at boot.
func RCEnablePF() (string, error) {
	// Check if pf has already been enabled, in that case, just return without
	// altering the /etc/rc.conf file again.
	// -n returns only the value of a variable in the rc.conf file.
	cmdOut, err := sysutils.ShCmd("sysrc", "-n", "pf_enable")
	// An error can be returned, when the variable does not exist yet in
	// the rc.conf file. So, do not inmediately return from the method,
	// try to still enable the variable.
	// If no error is returned, check if the value is already 'YES', in that
	// case, return without further altering the rc.conf file.
	if err == nil {
		// cmdOut has a newline, to compare it properly, scan the result string
		// out of the whole string first into pfRCValue.
		var pfRCValue string
		fmt.Sscan(cmdOut, &pfRCValue)

		if pfRCValue == "YES" {
			return "pf is already enabled in /etc/rc.conf", nil
		}
	}
	// If pf_enable has not been set to 'YES' yet, enable it.
	cmdOut, err = sysutils.ShCmd("sysrc", "pf_enable=YES")
	if err != nil {
		return "", fmt.Errorf("unable to enable pf in /etc/rc.conf: %w", err)
	}

	cmdOut = fmt.Sprintf("pf was successfully enabled in /etc/rc.conf: %s", cmdOut)

	return cmdOut, nil

}

// setupRulesFile, establishes '/etc/pf.conf' as the rules file for pf.
func setupRulesFile() (string, error) {
	cmdOut, err := sysutils.ShCmd("sysrc", "pf_rules=/etc/pf.conf")
	if err != nil {
		return "", fmt.Errorf("unable to setup /etc/pf.conf as the rules file for pf: %w", err)
	}

	cmdOut = fmt.Sprintf("pf_rules were successfully setup to be in /etc/pf.conf: %s", cmdOut)

	return cmdOut, nil

}

// setupLogFile, establishes '/var/log/pflog' as the log file for pf.
func setupLogFile() (string, error) {
	cmdOut, err := sysutils.ShCmd("sysrc", "pflog_logfile=/var/log/pflog")
	if err != nil {
		return "", fmt.Errorf("unable to setup /var/log/pflog as the log file for pflog: %w", err)
	}

	cmdOut = fmt.Sprintf("pflog file was successfully setup to be in /var/log/pflog: %s", cmdOut)
	return cmdOut, nil

}

// RCEnablePflog, enable pflog in /etc/rc.conf. So that it starts at boot.
func RCEnablePflog() (string, error) {
	// Check if pflog has already been enabled, in that case, just return without
	// altering the /etc/rc.conf file again.
	// -n returns only the value of a variable in the rc.conf file.
	cmdOut, err := sysutils.ShCmd("sysrc", "-n", "pflog_enable")
	// An error can be returned, when the variable does not exist yet in
	// the rc.conf file. So, do not inmediately return from the method,
	// try to still enable the variable.
	// If no error is returned, check if the value is already 'YES', in that
	// case, return without further altering the rc.conf file.
	if err == nil {
		// cmdOut has a newline, to compare it properly, scan the result string
		// out of the whole string first into pflogRCValue.
		var pflogRCValue string
		fmt.Sscan(cmdOut, &pflogRCValue)

		if pflogRCValue == "YES" {
			return "pflog is already enabled in /etc/rc.conf", nil
		}
	}
	// If pflog_enable has not been set to 'YES' yet, enable it.
	cmdOut, err = sysutils.ShCmd("sysrc", "pflog_enable=YES")
	if err != nil {
		return "", fmt.Errorf("unable to enable pflog in /etc/rc.conf: %w", err)
	}

	cmdOut = fmt.Sprintf("pflog was successfully enabled in /etc/rc.conf: %s", cmdOut)
	return cmdOut, nil

}

// CheckRuleSet, runs `pfctl -n` to check the syntax validity of the pf rules
// of a given file.
func CheckRuleSet(file string) (string, error) {
	// -n checks rules of -f file.
	cmdOut, err := sysutils.ShCmd("pfctl", "-nf", file)
	if err != nil {
		return "", fmt.Errorf("error checking the rules of pf file: %w", err)
	}
	// If the file does not exist, pfctl will throw an error, so it is unnecesa-
	// ry to check for the existance of the file.

	return cmdOut, nil

}

// activateRules, activates the rules in a file as the new pf rule set.
// func activateRules(file string) (string, error) {
// 	// Activate given file as new rule set.
// 	cmdOut, err := sysutils.ShCmd("pfctl", "-f", file)
// 	if err != nil {
// 		return "", fmt.Errorf("error activating new pf rule set from file (%s): %w", file, err)
// 	}
// 	// If the file does not exist, pfctl will throw an error, so it is unnecesa-
// 	// ry to check for the existance of the file.

// 	return cmdOut, nil

// }

// PFSetup, does all the required configurations on /etc/rc.conf to
// have pf working after rebooting the system.
// This function accepts a *log.Logger informational logger to print out
// information after running commands.
func PFSetup(infoLog *log.Logger) error {
	// Enable pf in rc.conf. After enabling pf, the default pf stance is to
	// accept all connections, so one will not be locked out of the SSH
	// connection with the server.
	outStr, err := RCEnablePF()
	if err != nil {
		return err
	}
	if outStr != "" {
		infoLog.Println(outStr)
	}

	// Let '/etc/pf.conf' be the rules file for pf.
	outStr, err = setupRulesFile()
	if err != nil {
		return err
	}
	if outStr != "" {
		infoLog.Println(outStr)
	}

	outStr, err = RCEnablePflog()
	if err != nil {
		return err
	}
	if outStr != "" {
		infoLog.Println(outStr)
	}

	// Let '/var/log/pflog' be the log file for pflog.
	outStr, err = setupLogFile()
	if err != nil {
		return err
	}
	if outStr != "" {
		infoLog.Println(outStr)
	}

	return nil
}

// EnablePF, enables PF. The firewall starts filtering packets, it is just as
// running 'pfctl -e'.
func EnablePF() (string, error) {
	// Check first if pfctl is already running, because, otherwise if it is
	// already running 'pfctl -e' returns an error.
	_, err := sysutils.ShCmd("pfctl", "-s Running")
	// 'pfctl -s Running' returns no error if pfctl is already running.
	if err == nil {
		return "pfctl is already running.", nil
	}

	outStr, err := sysutils.ShCmd("pfctl", "-e")
	if err != nil {
		return "", fmt.Errorf("error enabling pf: %w", err)
	}

	return outStr, nil

}
