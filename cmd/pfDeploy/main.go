// pfDeploy is a command-line tool capable of automatically configuring pf
// and pflog in a FreeBSD host machine. It is especially useful to quickly
// configure a newly spawned FreeBSD VM.
// pfDeploy, configures pf on a system to sane defaults, checks a given rule set
// syntax, enables pf with the given rule set and finally reboots the system.
package main

func main() {
	app := new(application)
	app.setupApplication()

	app.runTUI()
}
