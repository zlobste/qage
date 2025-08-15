package version

var (
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)

func String() string {
	if Commit != "" {
		return Version + " (" + Commit + ")"
	}
	return Version
}
