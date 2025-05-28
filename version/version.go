package version

import (
	"fmt"
	"time"
)

var (
	// Version number that is being run at the moment.
	Version = "0.4.0"

	// The git commit that was compiled. This will be filled in via Makefile.
	GitCommit string

	// BuildDate is RFC3339 formatted time of the git commit used to build
	// the program. This will be filled in via Makefile.
	BuildDate string

	// The name of the binary. This will be filled in via Makefile.
	BinaryName string
)

// ShortVersion returns only the version.
func ShortVersion() string {
	return BinaryName + " v" + Version
}

// FullVersion returns the version and the commit hash.
func FullVersion() string {
	version := ShortVersion()

	ts, _ := time.Parse(time.RFC3339, BuildDate)
	if !ts.IsZero() {
		version += "\n" + fmt.Sprintf("Built on %s", ts.Format(time.DateTime))
	}

	version += "\n" + fmt.Sprintf("Revision %s", GitCommit)

	return version
}
