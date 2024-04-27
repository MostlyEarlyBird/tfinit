package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func init() {
	fmt.Println("create")
	if _, err := os.Stat("main.tf"); err == nil {
		log.Fatal("Root main.tf already exists")
	}
}

func tfFiles() [3]string {
	return [3]string{"variables.tf", "main.tf", "outputs.tf"}
}

func (mod *Mod) CreateModule(name string) error {
	dirpath := path.Join("modules", name)
	if _, err := os.Stat(dirpath); err == nil {
		return fmt.Errorf("Directory by the name %s already exist\n", dirpath)
	}
	fmt.Printf("Creating %s ...\n", name)
	if err := os.MkdirAll(dirpath, 0777); err != nil {
		return err
	}
	for _, file := range tfFiles() {
		out, err := os.Create(path.Join(dirpath, file))
		defer out.Close()
		if err != nil {
			return err
		}
		if file == "variables.tf" {
			mod.writeVars(out)
		}
	}
	return nil
}

func (mod *Mod) writeVars(out *os.File) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	rootBody.AppendNewline()
	for _, v := range mod.Vars {
		variable := rootBody.AppendNewBlock("variable", []string{v.Name}).Body()
		variable.SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: v.Type}})
		variable.SetAttributeValue("description", cty.StringVal(v.Description))
	}
	if _, err := out.Write(f.Bytes()); err != nil {
		return err
	}
	return nil
}

func (yml *Yaml) generateRoot() error {
	mainFile := hclwrite.NewEmptyFile()
	tfVarsFile := hclwrite.NewEmptyFile()
	tfVarsBody := tfVarsFile.Body()
	mainBody := mainFile.Body()
	terraform := mainBody.AppendNewBlock("terraform", nil)
	trbody := terraform.Body()
	reqbody := trbody.AppendNewBlock("required_providers", nil).Body()
	reqbody.SetAttributeValue("aws", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("hashicorp/aws"),
		"version": cty.StringVal("~> 5.0"),
	}))
	mainBody.AppendNewline()
	provider := mainBody.AppendNewBlock("provider", []string{"aws"}).Body()
	provider.SetAttributeValue("region", cty.StringVal(yml.Region))
	if len(yml.Tags) > 0 {
		tagsBlock := provider.AppendNewBlock("default_tags", nil).Body()
		tags := map[string]cty.Value{}
		for _, t := range yml.Tags {
			tags[t.Name] = cty.StringVal(t.Value)
		}
		tagsBlock.SetAttributeValue("tags", cty.ObjectVal(tags))
	}
	mainBody.AppendNewline()

	modOut, err := os.Create("variables.tf")
	defer modOut.Close()
	if err != nil {
		return err
	}
	for mod_name, mod := range yml.Mods {
		modBlock := mainBody.AppendNewBlock("module", []string{mod_name}).Body()
		modBlock.SetAttributeValue("source", cty.StringVal("./modules/"+mod_name))
		for _, v := range mod.Vars {
			modBlock.SetAttributeTraversal(v.Name, hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: v.Name},
			})
			tfVarsBody.SetAttributeValue(v.Name, cty.StringVal("placeholder"))
		}
		mod.writeVars(modOut)
		mainBody.AppendNewline()
		tfVarsBody.AppendNewline()
	}

	out, err := os.Create("main.tf")
	defer out.Close()
	if err != nil {
		return err
	}
	_, err = out.Write(mainFile.Bytes())
	if err != nil {
		return err
	}
	varsOut, err := os.Create("terraform.tfvars")
	defer varsOut.Close()
	if err != nil {
		return err
	}
	_, err = varsOut.Write(tfVarsFile.Bytes())
	if err != nil {
		return err
	}
	return nil
}
