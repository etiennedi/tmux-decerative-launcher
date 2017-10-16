package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func FatalOn(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	windowName := os.Args[1]
	config, err := GetConfig()
	FatalOn(err)

	err = insideTmux()
	FatalOn(err)

	window, err := config.GetWindow(windowName)
	FatalOn(err)

	log.Print(window)
}

type Split struct {
	Command string `yaml:"command"`
	Path    string `yaml:"path"`
}

type Window struct {
	Splits []Split `yaml:"splits"`
}
type Config struct {
	Windows map[string]*Window `yaml:"windows"`
}

func GetConfig() (*Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("could not load config.yml: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config.yml: %s", err)
	}

	return &config, nil
}

func insideTmux() error {
	if os.Getenv("TMUX") == "" {
		return errors.New("you cannot run tdl outside of tmux")
	}

	return nil
}

func (c *Config) GetWindow(windowName string) (*Window, error) {
	window, ok := c.Windows[windowName]
	if !ok {
		return nil, fmt.Errorf("window '%s' is not configured", windowName)
	}

	return window, nil
}
