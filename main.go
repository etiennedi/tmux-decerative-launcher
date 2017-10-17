package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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

	err = exec.Command("tmux", "new-window").Run()
	FatalOn(err)

	err = exec.Command("tmux", "rename-window", windowName).Run()
	FatalOn(err)

	err = window.CreatePanes()
	FatalOn(err)
}

type Pane struct {
	Command string `yaml:"command"`
	Path    string `yaml:"path"`
	Height  *int   `yaml:"height"`
}

func (s Pane) NavigateToPath() error {
	if s.Path == "" {
		return errors.New("path in pane must be set")
	}

	shellCommand := fmt.Sprintf("cd %s", s.Path)
	err := exec.Command("tmux", "send-keys", shellCommand, "Enter").Run()
	if err != nil {
		return fmt.Errorf("could not navigate to path '%s': %s", s.Path, err)
	}

	return nil
}

func (s Pane) SetHeight() error {
	if s.Height == nil {
		return nil
	}

	err := exec.Command("tmux", "resize-pane", "-y", fmt.Sprintf("%d", *s.Height)).Run()
	if err != nil {
		return fmt.Errorf("could not set pane height to '%d': %s", s.Height, err)
	}

	return nil
}

func (s Pane) RunCommand() error {
	if s.Command == "" {
		return errors.New("command in pane must be set")
	}

	err := exec.Command("tmux", "send-keys", s.Command, "Enter").Run()
	if err != nil {
		return fmt.Errorf("could not navigate to path '%s': %s", s.Path, err)
	}

	return nil
}

type Window struct {
	Panes []Pane `yaml:"panes"`
}

func (w *Window) CreatePanes() error {
	if len(w.Panes) == 0 {
		return errors.New("No panes configured for this window")
	}

	for i, pane := range w.Panes {
		if i != 0 {
			err := exec.Command("tmux", "split-window", "-v").Run()
			if err != nil {
				return err
			}
		}

		err := pane.SetHeight()
		if err != nil {
			return err
		}

		err = pane.NavigateToPath()
		if err != nil {
			return err
		}

		err = pane.RunCommand()
		if err != nil {
			return err
		}
	}

	return nil
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
