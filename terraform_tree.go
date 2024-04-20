package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/erikgeiser/promptkit/textinput"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var files = []string{"variables.tf", "main.tf", "outputs.tf"}
var num int
var defaultTags bool

func init() {
	flag.IntVar(&num, "n", 0, "number of modules")
	flag.BoolVar(&defaultTags, "t", false, "add default tags")

}
func Usage() {
	fmt.Printf("Usage: %s [-t] [-n int] \n", path.Base(os.Args[0]))
	flag.PrintDefaults()
	os.Exit(0)
}

func createDirs(dirs map[string]bool) {
	for dir := range dirs {
		log.Printf("%s", dir)
		os.MkdirAll(dir, 0777)
		for _, file := range files {
			out, err := os.Create(path.Join(dir, file))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()
		}
	}
}

func generateMain(dirs []string) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	terraform := rootBody.AppendNewBlock("terraform", nil)
	trbody := terraform.Body()
	reqbody := trbody.AppendNewBlock("required_providers", nil).Body()
	reqbody.SetAttributeValue("aws", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("hashicorp/aws"),
		"version": cty.StringVal("~> 4.16"),
	}))
	trbody.AppendNewline()
	trbody.SetAttributeValue("required_version", cty.StringVal(">= 1.2.0"))
	rootBody.AppendNewline()
	provider := rootBody.AppendNewBlock("provider", []string{"aws"}).Body()
	provider.SetAttributeValue("region", cty.StringVal("ap-south-1"))
	rootBody.AppendNewline()

	for _, dir := range dirs {
		mod_name := dir[strings.LastIndex(dir, "/")+1:]
		mod := rootBody.AppendNewBlock("module", []string{mod_name}).Body()
		mod.SetAttributeValue("source", cty.StringVal("./modules/"+mod_name))
		rootBody.AppendNewline()
	}

	out, err := os.Create("main.tf")
	if err != nil {
		return
	}
	defer out.Close()
	_, err = out.Write(f.Bytes())
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	if _, err := os.Stat("main.tf"); err == nil {
		log.Fatal("main.tf already exists")

	}
	flag.Usage = Usage
	flag.Parse()
	if num == 0 {
		Usage()
	}
	wd, _ := os.Getwd()
	names := []string{}
	var list = make(map[string]bool)

	for i := 0; i < num; i++ {
		for {
			// tagName := textinput.New("Enter a name: ")
			// tagName.AutoComplete = textinput.AutoCompleteFromSlice([]string{
			// "Owner",
			// "bootcamp",
			// "expiration_date",
			// })
			moduleImpot := textinput.New("Enter a name")
			moduleImpot.Validate = func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("can't be empty")
				}

				return nil
			}
			moduleName, err := moduleImpot.RunPrompt()
			if err != nil {
				log.Fatalf("Error %v\n", err)
			}
			moduleName = strings.TrimSpace(moduleName)
			dirpath := path.Join(wd, "modules", moduleName)
			if list[dirpath] {
				log.Printf("%s already in list", moduleName)
			} else if _, err := os.Stat(dirpath); err == nil {
				log.Printf("Directory by the name %s already exist", moduleName)

			} else {
				list[dirpath] = true
				names = append(names, moduleName)
				break
			}
		}
	}
	createDirs(list)
	generateMain(names)
}
