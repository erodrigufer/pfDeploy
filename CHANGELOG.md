# v0.3.2
* Merge PR with dependancies update (`urfave/cli`).

# v0.3.1
* [BUG] Fix go version configuration in GitHub Action for GoReleaser.
	* The Go version being used was too old, so it did not support the features required to get the revision from an executable.
* Use Go v1.19 in go.mod file.

# v0.3.0
* Add the `--revision`/`-r` flags to show the hash of the commit used to build the executable.
	* If the executable was built with files with uncommitted changes the '-dirty' suffix will be added to the revision.

# v0.2.3
* Improve installation and usage documentation.

# v0.2.2
* Add more information to main package description.
* Add INSTALLATION guide to README.
* Add Examples to README.

# v0.2.1
* Add LICENSE

# v0.2.0
* Fix bug, reboot flag is now called `--no-reboot`.

# v0.1.2
* First tests with automatic releasing using goreleaser.

# v0.1.1
* Add exemplary pf conf files (configuration for MongoDB).

# v0.1.0
* Use `urfave/cli` to create the TUI.
* Commands `deploy` and `check` are working properly.
