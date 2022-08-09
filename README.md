# pfDeploy
`pfDeploy` is a command-line tool capable of automatically configuring pf and pflog in a FreeBSD host machine. It is especially useful to quickly configure a newly spawned FreeBSD VM.

## Usage
```bash
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

## Exemplary configuration files
In the subfolder `configFiles` are exemplary ruleset files for pf.
* `mongodb.conf` is a ruleset especially tailored for a VM hosting a remotely accessible MongoDB instance.
