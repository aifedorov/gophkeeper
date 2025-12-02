package version

import (
	"fmt"
	"runtime"
)

// Build information. Populated at build-time via ldflags.
var (
	Version   = "dev"
	BuildDate = "unknown"
)

// Info returns formatted version information.
func Info() string {
	return fmt.Sprintf(
		"GophKeeper Client\nVersion:    %s\nBuild Date: %s\nGo Version: %s\nOS/Arch:    %s/%s",
		Version,
		BuildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
