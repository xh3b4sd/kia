package project

var (
	description = "Command line tool for generating github workflows."
	gitSHA      = "n/a"
	name        = "workflow"
	source      = "https://github.com/xh3b4sd/kia"
	version     = "n/a"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
