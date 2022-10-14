// version gets version information embedded during the build process in the
// executable, e.g. revision and revision status from the VCS.
//
// This codebase has been modified from its original form, taken from
// autostrada.dev aka. Alex Edwards (2022).
package version

import (
	"fmt"
	"runtime/debug"
)

// GetRevision, returns a string with the commit hash used to build the
// executable. If the executable was built with files with some modifications
// related to the last commit, the suffix '-dirty' will be added to the revision
// returned. If no revision information could be retrieved, the the method
// returns 'unavailable'.
func GetRevision() string {
	var revision string
	var modified bool

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if revision == "" {
		return "unavailable"
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
