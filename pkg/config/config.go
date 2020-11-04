package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/file"
)

func GetKia(kia string) string {
	// Here we want to prefer the kia base path given via e.g. the process
	// environment.
	if kia != "" {
		return kia
	}

	// At this point the kia base path is not provided with the process
	// environment. We read the config file information from the local file
	// system.
	c := Config{}

	mustUnmarshal(&c)

	// At this point the kia base path was neither given via the process
	// environment nor found in the config file on the local file system.
	return c.Kia
}

func GetSec(sec string) string {
	// Here we want to prefer the sec base path given via e.g. the process
	// environment.
	if sec != "" {
		return sec
	}

	// At this point the sec base path is not provided with the process
	// environment. We read the config file information from the local file
	// system.
	c := Config{}

	mustUnmarshal(&c)

	// At this point we can iterate over the result of the user's config file.
	// Once we find the currently selected sec base path we return it.
	for _, i := range c.Org.List {
		if i.Org == c.Org.Selected {
			return i.Sec
		}
	}

	// At this point the sec base path was neither given via the process
	// environment nor found in the config file on the local file system.
	return ""
}

func Select(org string) Config {
	// We read the whole config file and all its information in order to not
	// override other settings when applying the kia base path.
	c := Config{}

	mustUnmarshal(&c)

	c.Org.Selected = org

	return c
}

func Validate(c Config) error {
	{
		abs, err := filepath.Abs(strings.TrimPrefix(c.Kia, "~/"))
		if err != nil {
			return tracer.Mask(err)
		}

		if !file.Exists(abs) {
			return tracer.Maskf(invalidConfigError, "config.kia is %#q but there does no such file exist on the file system", c.Kia)
		}
	}

	{
		for n, i := range c.Org.List {
			abs, err := filepath.Abs(strings.TrimPrefix(i.Sec, "~/"))
			if err != nil {
				return tracer.Mask(err)
			}

			if !file.Exists(abs) {
				return tracer.Maskf(invalidConfigError, "config.org.list[%d].sec is %#q but there does no such file exist on the file system", n, i.Sec)
			}
		}
	}

	{
		var ok bool

		for _, i := range c.Org.List {
			if i.Org == c.Org.Selected {
				ok = true
			}
		}

		if !ok {
			return tracer.Maskf(invalidConfigError, "config.org.selected is %#q but there is no org configured with this name", c.Org.Selected)
		}
	}

	return nil
}

func Write(c Config) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return tracer.Mask(err)
	}

	err = os.MkdirAll(filepath.Dir(mustPath()), os.ModePerm)
	if err != nil {
		return tracer.Mask(err)
	}

	err = ioutil.WriteFile(mustPath(), b, 0600)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func mustUnmarshal(v interface{}) {
	b, err := ioutil.ReadFile(mustPath())
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, v)
	if err != nil {
		panic(err)
	}
}

// mustPath returns the config file name as absolute path according to the
// current OS user known to the running process.
//
//     ~/.config/kia/config.yaml
//
func mustPath() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	return filepath.Join(u.HomeDir, ".config/kia/config.yaml")
}
