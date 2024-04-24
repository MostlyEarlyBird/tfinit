package main

import (
	"fmt"
	"log"
	"os"

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
type Yml struct {
	Tags []Tag          `yaml:"tags,flow"`
	Mods map[string]Mod `yaml:"mods,flow"`
}

func readConf(filename string) (Yml, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("file")
	}
	c := Yml{}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		log.Fatalf("in file %q: %v", filename, err)
	}

	return c, err
}

func main() {
	c, err := readConf("conf.yml")
	if err != nil {
		log.Fatal(err)
	}
	for _, tag := range c.Tags {
		fmt.Printf("tag.Name: %v\n", tag.Name)
		fmt.Printf("tag.Value: %v\n", tag.Value)
	}
	for key, value := range c.Mods {
		fmt.Printf("key: %v\n", key)
		for _, v := range value.Vars {
			fmt.Printf("v.Name: %v\n", v.Name)
			fmt.Printf("v.Type: %v\n", v.Type)
			fmt.Printf("v.Description: %v\n", v.Description)
		}
	}
}
