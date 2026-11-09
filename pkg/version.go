package version

var (
	Version   = "0.0.3"
	BuildTime = ""
	GitCommit = ""
)

func GetVersion() string {
	return Version
}

func GetBuildInfo() string {
	return "Version: " + Version + ", BuildTime: " + BuildTime + ", GitCommit: " + GitCommit
}
