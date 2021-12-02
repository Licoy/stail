package stail

const (
	OsWindows = "windows"
	OsLinux   = "linux"
	OsDarwin  = "darwin"
)

type Options struct {
	PowerShellPath string
	UnixTailPath   string
}
