# pfDeploy
`pfDeploy` is a command-line tool capable of automatically configuring pf and pflog in a FreeBSD host machine. It is especially useful to quickly configure a newly spawned FreeBSD VM.

## Table of contents

<!-- vim-markdown-toc GFM -->

* [Installation](#installation)
	- [Dependencies](#dependencies)
	- [With go install](#with-go-install)
* [Usage](#usage)
	- [Deploy](#deploy)
	- [Check](#check)
* [Exemplary configuration files](#exemplary-configuration-files)

<!-- vim-markdown-toc -->

## Installation
The following installation steps will install an executable of **pfDeploy** in the path used by Go to store binaries (as of _Go 1.18_, you can check the installation path for binaries by running `go env` and looking for the `GOPATH` variable. The binaries will be installed in the `/bin` subfolder of `GOPATH`. Add your `GOPATH` to your shell's `PATH` variable in order to execute Go binaries without having to specify the whole `GOPATH`).

### Dependencies
* Go 1.18+

### With go install
1. Install the most recent version of the pfDeploy command-line utility to the `GOPATH` binaries subfolder with: 

```
go install github.com/erodrigufer/pfDeploy/cmd/pfDeploy@latest
```

## Usage
```
pfDeploy - Automatically setup pf in your new deployment.

USAGE:
   pfDeploy [global options] command [command options] [arguments...]

COMMANDS:
   deploy   Setup pf and pflog at boot and deploy a new pf rules file.
   check    Check the syntactical validity of a pf rules file.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### Deploy
The command `deploy` checks the syntax validity of a given pf ruleset file, if the file is syntactically correct, the script changes the file at `/etc/rc.conf` so that both `pf` and `pflog` are always initialized at boot. Then, the given pf ruleset file is used as the new pf ruleset. Finally, the host is rebooted (unless the `--no-reboot` flag is used).

```
pfDeploy deploy --file <FILE_RULESET> --no-reboot
```

To show more help for the deploy command run `pfDeplot deploy --help`

### Check
Check the syntax validity of a pf ruleset file without changing any system configuration with the command: 

```
pfDeploy check --file <FILE_RULESET>
```

## Exemplary configuration files
In the subfolder `/configFiles` are exemplary ruleset files for pf.

* `mongodb.conf` is a ruleset especially tailored for a VM hosting a remotely accessible MongoDB instance.
