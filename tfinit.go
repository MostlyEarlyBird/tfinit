package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Tag struct {
	Name  string
	Value string
}
type Var struct {
	Name        string
	Type        string
	Description string
}
type Mod struct {
	Vars []Var `yaml:"vars,flow"`
}
type Yaml struct {
	Region string         `yaml:"region,flow"`
	Tags   []Tag          `yaml:"tags,flow"`
	Mods   map[string]Mod `yaml:"modules,flow"`
}

func getConfig() (string, error) {
	// TODO: search with exact extentions (yml/yaml)
	file, err := filepath.Glob("config.*ml")
	fmt.Printf("file: %v\n", file)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	} else if len(file) == 0 || len(file) > 1 {
		return "", fmt.Errorf("Found: %v config files\n", len(file))
	} else {
		return file[0], nil
	}
}

func (yml *Yaml) readConf(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &yml)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatalf("err %v\n", err)
	}
	yml := new(Yaml)
	if err := yml.readConf(config); err != nil {
		log.Fatalf("err %v\n", err)
	}
	if err := yml.generateRoot(); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	for key, value := range yml.Mods {
		if err := value.CreateModule(key); err != nil {
			log.Fatalf("err: %v\n", err)
		}
	}
}
