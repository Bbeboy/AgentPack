package version

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

func Value() string {
	if Version == "" {
		return "dev"
	}
	return Version
}
