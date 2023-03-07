package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

// Private
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func findKeyEditValue(node interface{}, key, value string) (interface{}, bool) {
	m, ok := node.(map[interface{}]interface{})
	if !ok {
		return nil, false
	}

	for k, v := range m {
		if k == key {
			m[key] = value
		}

		if result, found := findKeyEditValue(v, key, value); found {
			return result, true
		}
	}

	return nil, false
}

// Public
func CreateIfNotExists(name string) {
	file, _ := filepath.Abs(name)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(name)
		check(err)
		defer f.Close()
	}
}

func ToAbsPath(ans interface{}) interface{} {
	p := ans.(string)
	if len(p) == 0 {
		return nil
	}
	abs_path, _ := filepath.Abs(p)
	return abs_path
}

func Contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func ConvertBoolToStr(s string) bool {
	if s != "" {
		parsed, err := strconv.ParseBool(s)
		if err != nil {
			log.Fatalf(
				`error parsing value: %s
Available values is: 
	"1", "t", "T", "true", "TRUE", "True"
	"0", "f", "F", "false", "FALSE", "False"`, s)
		}
		return parsed
	}
	return true

}

func CheckFileExt(name string) {
	fileExt := filepath.Ext(name)
	allowedExt := []string{
		".yaml",
		".yml",
		".tf",
		".ini",
		".env",
	}
	if !Contains(allowedExt, fileExt) {
		log.Fatalf("file extension is not supported: %s\nsupported extenstions: %s", fileExt, strings.Join(allowedExt, " "))
	}
}

func RunEditor(outFile string) {
	editor := "vim"
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}
	if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}
	f, err := filepath.Abs(outFile)
	check(err)

	cmd := exec.Command(editor, f)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Run()
}

func EditYaml(filepath, key, value string) {
	file, err := os.ReadFile(filepath)
	check(err)

	info, err := os.Stat(filepath)
	check(err)

	mode := info.Mode()

	y := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &y)
	check(err)

	updatedNode, found := findKeyEditValue(y, key, value)

	if !found {
		y[key] = value
		updatedNode = y
	}

	yamlFile, err := yaml.Marshal(updatedNode)
	check(err)

	err = os.WriteFile(filepath, yamlFile, mode)
	check(err)
}

func EditHCL(filepath, key, value string) {
	file, err := os.ReadFile(filepath)
	check(err)

	info, err := os.Stat(filepath)
	check(err)

	mode := info.Mode()

	hcf, diags := hclwrite.ParseConfig(file, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		fmt.Printf("errors: %s", diags)
		return
	}
	hcfb := hcf.Body()
	hcfb.SetAttributeValue(key, cty.StringVal(value))
	err = os.WriteFile(filepath, hcf.Bytes(), mode)
	check(err)
}

func EditIni(filepath, key, value string) {
	cfg, err := ini.Load(filepath)
	check(err)
	re, _ := regexp.Compile(`^.+\_\_`)
	section := re.Find([]byte(key))
	if section != nil {
		key := strings.Replace(key, string(section), "", 1)
		section := strings.Replace(string(section), "__", "", 1)
		cfg.Section(section).Key(key).SetValue(value)
	} else {
		cfg.Section("").Key(key).SetValue(value)
	}
	cfg.SaveTo(filepath)
}

func EditDotEnv(filepath, key, value string) {
	ini.PrettyFormat = false
	ini.PrettyEqual = false
	_, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	check(err)

	cfg, err := ini.Load(filepath)
	check(err)

	cfg.Section("").Key(strings.ToUpper(key)).SetValue(value)
	cfg.SaveTo(filepath)
}

func WriteToFile(file string, answers map[string]interface{}) {
	fileExt := filepath.Ext(file)
	switch fileExt {
	case ".yaml", ".yml":
		for k, v := range answers {
			EditYaml(file, k, v.(string))
		}
	case ".env":
		for k, v := range answers {
			EditDotEnv(file, k, v.(string))
		}
	case ".ini":
		for k, v := range answers {
			EditIni(file, k, v.(string))
		}
	case ".tf":
		for k, v := range answers {
			EditHCL(file, k, v.(string))
		}
	}
}
