package docker

type DockerConfigJSON struct {
	Auths map[string]DockerConfigJSONAuth `json:"auths"`
}

type DockerConfigJSONAuth struct {
	Auth string `json:"auth"`
	Pass string `json:"password"`
	User string `json:"username"`
}
