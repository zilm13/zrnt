package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type ContstantsPresetData struct {
	Name    string
	Entries []string
}

type ConstantsPreset struct {
	BuildConstraint string
	Data *ContstantsPresetData
}

// hex input should not be prefixed with 0x
func hexStrToLiteralStr(hex string) string {
	byteCount := len(hex)
	formattedValue := fmt.Sprintf("[%d]byte{", byteCount)
	for i := 0; i < len(hex); i += 2 {
		formattedValue += "0x" + hex[i:i+2]
		if i+2 < len(hex) {
			formattedValue += ", "
		}
	}
	formattedValue += "}"
	return formattedValue
}

func buildPreset(path string) (*ContstantsPresetData, error) {
	presetName := filepath.Base(path)
	presetName = presetName[:len(presetName)-len(".yaml")]

	fmt.Println("processing preset", presetName)

	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rawPreset := yaml.MapSlice{}
	if err := yaml.Unmarshal(yamlBytes, &rawPreset); err != nil {
		return nil, err
	}

	preset := ContstantsPresetData{
		Name:    presetName,
		Entries: make([]string, 0, len(rawPreset)),
	}
	for _, item := range rawPreset {
		k := item.Key.(string)
		v := item.Value
		intFormat := "%d"
		if strings.HasPrefix(k, "DOMAIN_") {
			intFormat = "0x%08x"
		}
		formattedValue := ""
		formattedStart := "const " + k + " = "
		if strV, ok := v.(string); ok {
			if intV, err := strconv.ParseInt(strV, 0, 64); err == nil {
				formattedValue = fmt.Sprintf("%d", intV)
			} else if strings.HasPrefix(strV, "0x") {
				strNibbles := strV[2:]
				if len(strNibbles)%2 != 0 {
					return nil, errors.New(fmt.Sprintf("invalid constant, %s has value %s", k, strV))
				}
				formattedValue = hexStrToLiteralStr(strNibbles)
				// arrays cannot be constants
				formattedStart = "var " + k + " = "
			} else {
				return nil, errors.New(fmt.Sprintf("could not convert string formatted value in %s, key: %s, value: %s", presetName, k, strV))
			}
		} else if uintV, ok := v.(uint64); ok {
			formattedValue = fmt.Sprintf(intFormat, uintV)
		} else if uintV, ok := v.(uint32); ok {
			formattedValue = fmt.Sprintf(intFormat, uintV)
		} else if intV, ok := v.(int); ok {
			formattedValue = fmt.Sprintf(intFormat, intV)
		} else {
			formattedStart = "// " + formattedStart
			formattedValue = fmt.Sprintf("(unrecognized type) %v", v)
		}

		preset.Entries = append(preset.Entries, formattedStart+formattedValue)
	}
	return &preset, nil
}

func main() {
	var presetsDirPath, outputDirPath, defaultPreset string
	flag.StringVar(&presetsDirPath, "presets-dir", "", "The file path to the directory containing yaml constant presets")
	flag.StringVar(&outputDirPath, "output-dir", "", "The file path to the directory to output generated Go code to")
	flag.StringVar(&defaultPreset, "default-preset", "mainnet", "The name of the preset to duplicate, with build constraints to use it if no other preset is used")
	flag.Parse()

	templ := template.Must(template.New("constants_file").Parse(constantsFileTemplate))

	var presetNames []string
	if err := filepath.Walk(presetsDirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fmt.Println("found preset file", path)

			extension := filepath.Ext(path)
			if extension != ".yaml" {
				return nil
			}

			presetName := filepath.Base(path)
			presetName = presetName[:len(presetName)-len(".yaml")]

			presetNames = append(presetNames, presetName)
			return nil
		}); err != nil {
		panic(err)
	}

	if err := filepath.Walk(presetsDirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fmt.Println("processing preset file", path)

			extension := filepath.Ext(path)
			if extension != ".yaml" {
				return nil
			}

			presetData, err := buildPreset(path)
			if err != nil {
				return err
			}

			outPath := filepath.Join(outputDirPath, presetData.Name+".go")
			fmt.Printf("writing constants preset %s to %s\n", presetData.Name, outPath)
			f, err := os.Create(outPath)
			if err != nil {
				return err
			}
			preset := ConstantsPreset{
				BuildConstraint: "preset_"+presetData.Name,
				Data: presetData,
			}
			if err := templ.Execute(f, preset); err != nil {
				return err
			}

			if presetData.Name == defaultPreset {
				outPath := filepath.Join(outputDirPath, "defaults.go")
				fmt.Printf("writing default preset, alias %s to %s\n", presetData.Name, outPath)
				f, err := os.Create(outPath)
				if err != nil {
					return err
				}
				preset := ConstantsPreset{
					// when all other presets are not active, then use the default
					BuildConstraint: "!preset_"+strings.Join(presetNames, ",!preset_"),
					Data: presetData,
				}
				if err := templ.Execute(f, preset); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
		panic(err)
	}

}

var constantsFileTemplate = `// +build {{.BuildConstraint}}

package generated

const PRESET_NAME string = "{{.Data.Name}}"
{{ range .Data.Entries }}
{{.}}
{{ end }}`
