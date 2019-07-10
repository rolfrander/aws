package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// default config file
const CFGFILE = "config.yaml"

// Config directory on windows
const CFGDIR_WINDOWS = "Glacierbackup"

//Config directory on unix
const CFGDIR_UNIX = ".glacierbackup"

type config struct {
	GoogleClientFile string
	GoogleTokenFile  string
	SpreadsheetID    string
	GlacierVault     string
}

func configDir() string {
	if os.Getenv("APPDATA") != "" {
		return filepath.Join(os.Getenv("APPDATA"), CFGDIR_WINDOWS)
	}
	if os.Getenv("HOME") != "" {
		return filepath.Join(os.Getenv("HOME"), CFGDIR_UNIX)
	}
	// I have no idea, use current directory

	return "."
}

func defaultValue(v *string, def string) {
	if v != nil && *v == "" {
		*v = def
	}
}

func (c *config) readConfig() error {

	yamlData, err := ioutil.ReadFile(filepath.Join(configDir(), CFGFILE))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, c)
	if err != nil {
		return err
	}

	if c.SpreadsheetID == "" || c.GlacierVault == "" {
		return fmt.Errorf("Incomplete configuration, please run config")
	}

	defaultValue(&c.GoogleClientFile, filepath.Join(configDir(), "google-credentials.json"))
	defaultValue(&c.GoogleTokenFile, filepath.Join(configDir(), "google-token.json"))

	return nil
}

func (c *config) writeConfig() error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(configDir(), CFGFILE), yamlData, 0600)
	if err != nil {
		return err
	}
	return nil
}
