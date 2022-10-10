package g

var (
	BinaryName string
	Version    string
	GitCommit  string
)

func VersionMsg() string {
	return Version + "@" + GitCommit
}
