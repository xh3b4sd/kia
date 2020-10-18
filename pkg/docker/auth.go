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
	Bytes []byte
}

type Auth struct {
	bytes []byte
}

func NewAuth(config AuthConfig) (*Auth, error) {
	if config.Bytes == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Bytes must not be empty", config)
	}

	a := &Auth{
		bytes: config.Bytes,
	}

	return a, nil
}

func (a *Auth) Encode() (string, error) {
	var ok bool

	m := map[string]string{}
	err := json.Unmarshal(a.bytes, &m)
	if err != nil {
		return "", tracer.Mask(err)
	}

	var password string
	{
		password, ok = m[Password]
		if !ok {
			return "", tracer.Maskf(executionFailedError, "no secret for %#q", Password)
		}
	}

	var registry string
	{
		registry, ok = m[Registry]
		if !ok {
			return "", tracer.Maskf(executionFailedError, "no secret for %#q", Registry)
		}
	}

	var username string
	{
		username, ok = m[Username]
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
