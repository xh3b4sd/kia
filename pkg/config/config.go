package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ghodss/yaml"
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
	v := struct {
		Kia *string `yaml:"kia"`
	}{}

	mustFromFile(&v)

	// At this point the kia base path was neither given via the process
	// environment nor found in the config file on the local file system.
	if v.Kia == nil {
		return ""
	}

	return *v.Kia
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
	v := struct {
		Sec *string `yaml:"sec"`
	}{}

	mustFromFile(&v)

	// At this point the sec base path was neither given via the process
	// environment nor found in the config file on the local file system.
	if v.Sec == nil {
		return ""
	}

	return *v.Sec
}

func SetKia(kia string) {
	// We read the whole config file and all its information in order to not
	// override other settings by applying the kia base path.
	v := map[string]interface{}{}

	mustFromFile(&v)

	v["kia"] = kia

	mustToFile(&v)
}

func SetSec(sec string) {
	// We read the whole config file and all its information in order to not
	// override other settings by applying the sec base path.
	v := map[string]interface{}{}

	mustFromFile(&v)

	v["sec"] = sec

	mustToFile(&v)
}

func mustFromFile(v interface{}) {
	b, err := ioutil.ReadFile(mustName())
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

func mustToFile(v interface{}) {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Dir(mustName()), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(mustName(), b, 0600)
	if err != nil {
		panic(err)
	}
}

// mustName returns the config file name as absolute path according to the
// current OS user known to the running process.
//
//     ~/.config/kia/config.yaml
//
func mustName() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	return filepath.Join(u.HomeDir, ".config/kia/config.yaml")
}
