// Package cmd
/*
Copyright © 2021 Tonye Jack <jtonye@ymail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var InputsHeader = "## Inputs"
var OutputsHeader = "## Outputs"
var AutoDocStart = "<!-- AUTO-DOC-%s:START - Do not remove or modify this section -->"
var AutoDocEnd = "<!-- AUTO-DOC-%s:END -->"
var inputAutoDocStart = fmt.Sprintf(AutoDocStart, "INPUT")
var inputAutoDocEnd = fmt.Sprintf(AutoDocEnd, "INPUT")
var outputAutoDocStart = fmt.Sprintf(AutoDocStart, "OUTPUT")
var outputAutoDocEnd = fmt.Sprintf(AutoDocEnd, "OUTPUT")

var actionFileName string
var outputFileName string

type Input struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default,omitempty"`
}

type Output struct {
	Description string `yaml:"description"`
	Value       string `yaml:"default,omitempty"`
}

type Action struct {
	Inputs  map[string]Input  `yaml:"inputs,omitempty"`
	Outputs map[string]Output `yaml:"outputs,omitempty"`
}

func (a *Action) getAction() *Action {
	actionYaml, err := ioutil.ReadFile(actionFileName)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = yaml.Unmarshal(actionYaml, &a)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return a
}

func (a *Action) renderOutput() {
	var err error
	inputTableOutput := &strings.Builder{}

	if len(a.Inputs) > 0 {
		_, err = fmt.Fprintln(inputTableOutput, inputAutoDocStart)
		if err != nil {
			cobra.CheckErr(err)
		}

		inputTable := tablewriter.NewWriter(inputTableOutput)
		inputTable.SetHeader([]string{"Input", "Required", "Default", "Description"})
		inputTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		inputTable.SetCenterSeparator("|")

		keys := make([]string, 0, len(a.Inputs))
		for k := range a.Inputs {
			keys = append(keys, k)
		}

		for _, key := range keys {
			row := []string{key, strconv.FormatBool(a.Inputs[key].Required), a.Inputs[key].Default, a.Inputs[key].Description}
			inputTable.Append(row)
		}

		_, err = fmt.Fprintln(inputTableOutput)
		if err != nil {
			cobra.CheckErr(err)
		}

		inputTable.Render()

		_, err = fmt.Fprintln(inputTableOutput)
		if err != nil {
			cobra.CheckErr(err)
		}

		_, err = fmt.Fprint(inputTableOutput, inputAutoDocEnd)
		if err != nil {
			cobra.CheckErr(err)
		}
	}

	outputTableOutput := &strings.Builder{}

	if len(a.Outputs) > 0 {
		_, err = fmt.Fprintln(outputTableOutput, outputAutoDocStart)
		if err != nil {
			cobra.CheckErr(err)
		}

		outputTable := tablewriter.NewWriter(outputTableOutput)
		outputTable.SetHeader([]string{"Output", "Description", "Type"})
		outputTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		outputTable.SetCenterSeparator("|")

		keys := make([]string, 0, len(a.Outputs))
		for k := range a.Outputs {
			keys = append(keys, k)
		}

		for _, key := range keys {
			row := []string{key, a.Outputs[key].Description, "string"}
			outputTable.Append(row)
		}

		_, err = fmt.Fprintln(outputTableOutput)
		if err != nil {
			cobra.CheckErr(err)
		}

		outputTable.Render()

		_, err = fmt.Fprintln(outputTableOutput)
		if err != nil {
			cobra.CheckErr(err)
		}

		_, err = fmt.Fprint(outputTableOutput, outputAutoDocEnd)
		if err != nil {
			cobra.CheckErr(err)
		}
	}

	input, err := ioutil.ReadFile(outputFileName)

	if err != nil {
		cobra.CheckErr(err)
	}

	var output = []byte("")

	hasInputsData, inputStartIndex, inputEndIndex := HasBytesInBetween(
		input,
		[]byte(InputsHeader),
		[]byte(inputAutoDocEnd),
	)

	if hasInputsData {
		inputsStr := fmt.Sprintf("%s\n\n%v", InputsHeader, inputTableOutput.String())
		output = ReplaceBytesInBetween(input, inputStartIndex, inputEndIndex, []byte(inputsStr))
	} else {
		inputsStr := fmt.Sprintf("%s\n\n%v", InputsHeader, inputTableOutput.String())
		output = bytes.Replace(input, []byte(InputsHeader), []byte(inputsStr), -1)
	}

	hasOutputsData, outputStartIndex, outputEndIndex := HasBytesInBetween(
		output,
		[]byte(OutputsHeader),
		[]byte(outputAutoDocEnd),
	)

	if hasOutputsData {
		outputsStr := fmt.Sprintf("%s\n\n%v", OutputsHeader, outputTableOutput.String())
		output = ReplaceBytesInBetween(output, outputStartIndex, outputEndIndex, []byte(outputsStr))
	} else {
		outputsStr := fmt.Sprintf("%s\n\n%v", OutputsHeader, outputTableOutput.String())
		output = bytes.Replace(output, []byte(OutputsHeader), []byte(outputsStr), -1)
	}

	if len(output) > 0 {
		if err = ioutil.WriteFile(outputFileName, output, 0666); err != nil {
			cobra.CheckErr(err)
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "auto-doc",
	Short: "Auto doc generator for your github action",
	Long:  `Auto generate documentation for your github action.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			_, err := fmt.Fprintf(
				os.Stderr,
				"'%d' invalid arguments passed.\n",
				len(args),
			)
			if err != nil {
				cobra.CheckErr(err)
			}
			return
		}

		var action Action
		action.getAction()
		action.renderOutput()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Custom flags
	rootCmd.PersistentFlags().StringVar(
		&actionFileName,
		"action",
		"action.yml",
		"action config file",
	)
	rootCmd.PersistentFlags().StringVar(
		&outputFileName,
		"output",
		"README.md",
		"Output file",
	)
}

func HasBytesInBetween(value, start, end []byte) (found bool, startIndex int, endIndex int) {
	s := bytes.Index(value, start)

	if s == -1 {
		return false, -1, -1
	}

	e := bytes.Index(value, end)

	if e == -1 {
		return false, -1, -1
	}

	return true, s, e + len(end)
}

func ReplaceBytesInBetween(value []byte, startIndex int, endIndex int, new []byte) []byte {
	t := make([]byte, len(value)+len(new))
	w := 0

	w += copy(t[:startIndex], value[:startIndex])
	w += copy(t[w:w+len(new)], new)
	w += copy(t[w:], value[endIndex:])
	return t[0:w]
}
