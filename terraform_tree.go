package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var files = []string{"variables.tf", "main.tf", "outputs.tf"}
var num int

func init() {
	flag.IntVar(&num, "n", 0, "number of modules")

}
func Usage() {
	fmt.Printf("Usage: %s argument ...\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
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
	trbody := terraform.Body().AppendNewBlock("required_providers", nil).Body()
	trbody.SetAttributeValue("aws", cty.ObjectVal(map[string]cty.Value{
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
    if _,err := os.Stat("main.tf") ; err == nil {
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
		for true {
			var name string
			fmt.Print("Enter a name: ")
			fmt.Scan(&name)
			dirpath := path.Join(wd,"modules" ,name)
			if list[dirpath] {
				log.Printf("%s already in list", name)
			} else if _, err := os.Stat(dirpath); err == nil {
				log.Printf("Directory by the name %s already exist", name)

			} else {
				list[dirpath] = true
				names = append(names, name)
				break
			}
		}
	}
	createDirs(list)
	generateMain(names)
}
