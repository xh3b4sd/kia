package docker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/xh3b4sd/tracer"
)

const (
	Password = "docker.password"
	Registry = "docker.registry"
	Username = "docker.username"
)

type AuthConfig struct {
	Secrets map[string]string
}

type Auth struct {
	secrets map[string]string
}

func NewAuth(config AuthConfig) (*Auth, error) {
	if config.Secrets == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Secrets must not be empty", config)
	}

	a := &Auth{
		secrets: config.Secrets,
	}

	return a, nil
}

func (a *Auth) Encode() (string, error) {
	var ok bool

	var password string
	{
		password, ok = a.secrets[Password]
		if !ok {
			return "", tracer.Maskf(executionFailedError, "no secret for %#q", Password)
		}
	}

	var registry string
	{
		registry, ok = a.secrets[Registry]
		if !ok {
			return "", tracer.Maskf(executionFailedError, "no secret for %#q", Registry)
		}
	}

	var username string
	{
		username, ok = a.secrets[Username]
		if !ok {
			return "", tracer.Maskf(executionFailedError, "no secret for %#q", Username)
		}
	}

	var auth string
	{
		auth = base64.URLEncoding.EncodeToString([]byte(username + ":" + password))
	}

	v := DockerConfigJSON{
		Auths: map[string]DockerConfigJSONAuth{
			registry: {
				Auth: auth,
				User: username,
				Pass: password,
			},
		},
	}

	var enc string
	{
		b, err := json.Marshal(v)
		if err != nil {
			return "", tracer.Mask(err)
		}

		enc = fmt.Sprintf("%s\n", base64.URLEncoding.EncodeToString(b))
	}

	return enc, nil
}
