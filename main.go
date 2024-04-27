package main

import (
	"fmt"
	"github.com/MostlyEarlyBird/tfinit/terraformtree"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var temp = terraformtree.Temp{Temp: 3}

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

var config string

func init() {
	// TODO: search with exact extentions (yml/yaml)
	file, err := filepath.Glob("config.*ml")
	fmt.Printf("file: %v\n", file)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if len(file) == 0 || len(file) > 1 {
		fmt.Printf("Found: %v config files\n", len(file))
		return
	} else {
		config = file[0]
	}
}

func tfFiles() [3]string {
	return [3]string{"variables.tf", "main.tf", "outputs.tf"}
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
	yml := new(Yaml)
	if err := yml.readConf(config); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
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
