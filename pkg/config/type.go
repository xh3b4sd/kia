package config

type Config struct {
	Kia string    `json:"kia" yaml:"kia"`
	Org ConfigOrg `json:"org" yaml:"org"`
}

type ConfigOrg struct {
	List     []ConfigOrgItem `json:"list" yaml:"list"`
	Selected string          `json:"selected" yaml:"selected"`
}

type ConfigOrgItem struct {
	Org string `json:"org" yaml:"org"`
	Sec string `json:"sec" yaml:"sec"`
}
