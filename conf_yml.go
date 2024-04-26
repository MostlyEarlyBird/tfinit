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
	Tags []Tag          `yaml:"tags,flow"`
	Mods map[string]Mod `yaml:"modules,flow"`
}

func (yml *Yaml) readConf(filename string) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("file")
	}
	err = yaml.Unmarshal(buf, &yml)
	if err != nil {
		log.Fatalf("in file %q: %v", filename, err)
	}
}

func main() {
	yml := new(Yaml)
	file, err := filepath.Glob("config.y?(a)ml")
	yml.readConf("conf.yml")
	fmt.Printf("yml: %v\n", yml)
	for _, tag := range yml.Tags {
		fmt.Printf("tag.Name: %v\n", tag.Name)
		fmt.Printf("tag.Value: %v\n", tag.Value)
	}
	for key, value := range yml.Mods {
		fmt.Printf("key: %v\n", key)
		for _, v := range value.Vars {
			fmt.Printf("v.Name: %v\n", v.Name)
			fmt.Printf("v.Type: %v\n", v.Type)
			fmt.Printf("v.Description: %v\n", v.Description)
		}
	}
}
