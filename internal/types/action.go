//Package types contains all defined types
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
package types

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/tj-actions/auto-doc/v2/internal"
	"github.com/tj-actions/auto-doc/v2/internal/utils"
)

// ActionInput represents the input of the action.yml
type ActionInput struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default,omitempty"`
}

// ActionOutput represents the output of the action.yml
type ActionOutput struct {
	Description string `yaml:"description"`
	Value       string `yaml:"default,omitempty"`
}

// Action represents the action.yml
type Action struct {
	InputFileName  string
	OutputFileName string
	ColMaxWidth    string
	ColMaxWords    string
	InputColumns   []string
	OutputColumns  []string
	Inputs         map[string]ActionInput  `yaml:"inputs,omitempty"`
	Outputs        map[string]ActionOutput `yaml:"outputs,omitempty"`
}

// GetData parses the source yaml file
func (a *Action) GetData() error {
	actionYaml, err := os.ReadFile(a.InputFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(actionYaml, &a)
	return err
}

// WriteDocumentation write the table to the output file
func (a *Action) WriteDocumentation(inputTable, outputTable *strings.Builder) error {
	var err error
	input, err := os.ReadFile(a.OutputFileName)

	if err != nil {
		return err
	}

	var output []byte

	hasInputsData, indices := utils.HasBytesInBetween(
		input,
		[]byte(internal.InputAutoDocStart),
		[]byte(internal.InputAutoDocEnd),
	)

	output = input
	inputsStr := strings.TrimSpace(inputTable.String())

	if hasInputsData {
		output = utils.ReplaceBytesInBetween(output, indices, []byte(inputsStr))
	} else {
		re := regexp.MustCompile(fmt.Sprintf("(?m)^%s", internal.InputsHeader))
		output = re.ReplaceAllFunc(input, func(match []byte) []byte {
			if bytes.HasPrefix(match, []byte(internal.InputsHeader)) {
				return []byte(fmt.Sprintf("%s\n\n%v", internal.InputsHeader, inputsStr))
			}
			return match
		})
	}

	hasOutputsData, indices := utils.HasBytesInBetween(
		output,
		[]byte(internal.OutputAutoDocStart),
		[]byte(internal.OutputAutoDocEnd),
	)

	outputsStr := strings.TrimSpace(outputTable.String())

	if hasOutputsData {
		output = utils.ReplaceBytesInBetween(output, indices, []byte(outputsStr))
	} else {
		re := regexp.MustCompile(fmt.Sprintf("(?m)^%s", internal.OutputsHeader))
		output = re.ReplaceAllFunc(output, func(match []byte) []byte {
			if bytes.HasPrefix(match, []byte(internal.OutputsHeader)) {
				return []byte(fmt.Sprintf("%s\n\n%v", internal.OutputsHeader, outputsStr))
			}
			return match
		})
	}

	if err = os.WriteFile(a.OutputFileName, output, 0666); err != nil {
		cobra.CheckErr(err)
	}

	return nil
}

// RenderOutput renders the output and writes it to the given output
func (a *Action) RenderOutput() error {
	var err error
	maxWidth, err := strconv.Atoi(a.ColMaxWidth)
	if err != nil {
		return err
	}

	maxWords, err := strconv.Atoi(a.ColMaxWords)
	if err != nil {
		return err
	}

	inputTableOutput, err := renderActionInputTableOutput(a.Inputs, a.InputColumns, maxWidth, maxWords)
	if err != nil {
		return err
	}

	outputTableOutput, err := renderActionOutputTableOutput(a.Outputs, a.OutputColumns, maxWidth, maxWords)
	if err != nil {
		return err
	}

	err = a.WriteDocumentation(inputTableOutput, outputTableOutput)
	if err != nil {
		return err
	}

	return nil
}

// renderActionOutputTableOutput renders the action input table
func renderActionInputTableOutput(i map[string]ActionInput, inputColumns []string, maxWidth int, maxWords int) (*strings.Builder, error) {
	inputTableOutput := &strings.Builder{}

	if len(i) > 0 {
		_, err := fmt.Fprintln(inputTableOutput, internal.InputAutoDocStart)
		if err != nil {
			return inputTableOutput, err
		}

		inputTable := tablewriter.NewWriter(inputTableOutput)
		inputTable.SetHeader(inputColumns)
		inputTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		inputTable.SetCenterSeparator(internal.PipeSeparator)
		inputTable.SetAlignment(tablewriter.ALIGN_CENTER)

		keys := make([]string, 0, len(i))
		for k := range i {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		inputTable.SetColWidth(maxWidth)

		for _, key := range keys {
			var row []string

			for _, col := range inputColumns {
				switch col {
				case "Input":
					row = append(row, key)
				case "Type":
					row = append(row, "string")
				case "Required":
					row = append(row, strconv.FormatBool(i[key].Required))
				case "Default":
					row = append(row, utils.FormatValue(i[key].Default))
				case "Description":
					row = append(row, utils.WordWrap(i[key].Description, maxWords))
				default:
					return inputTableOutput, fmt.Errorf(
						"unknown input column: '%s'. Please specify any of the following columns: %s",
						col,
						strings.Join(internal.DefaultActionInputColumns, ", "),
					)
				}
			}
			inputTable.Append(row)
		}

		_, err = fmt.Fprintln(inputTableOutput)
		if err != nil {
			return inputTableOutput, err
		}

		inputTable.Render()

		_, err = fmt.Fprintln(inputTableOutput)
		if err != nil {
			return inputTableOutput, err
		}

		_, err = fmt.Fprint(inputTableOutput, internal.InputAutoDocEnd)
		if err != nil {
			return inputTableOutput, err
		}
	}
	return inputTableOutput, nil
}

// renderActionOutputTableOutput renders the action output table
func renderActionOutputTableOutput(o map[string]ActionOutput, outputColumns []string, maxWidth int, maxWords int) (*strings.Builder, error) {
	outputTableOutput := &strings.Builder{}

	if len(o) > 0 {
		_, err := fmt.Fprintln(outputTableOutput, internal.OutputAutoDocStart)
		if err != nil {
			return outputTableOutput, err
		}

		outputTable := tablewriter.NewWriter(outputTableOutput)
		outputTable.SetHeader(outputColumns)
		outputTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		outputTable.SetCenterSeparator(internal.PipeSeparator)
		outputTable.SetAlignment(tablewriter.ALIGN_CENTER)

		keys := make([]string, 0, len(o))
		for k := range o {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		outputTable.SetColWidth(maxWidth)
		for _, key := range keys {
			var row []string

			for _, col := range outputColumns {
				switch col {
				case "Output":
					row = append(row, key)
				case "Type":
					row = append(row, "string")
				case "Description":
					row = append(row, utils.WordWrap(o[key].Description, maxWords))
				default:
					return outputTableOutput, fmt.Errorf(
						"unknown output column: '%s'. Please specify any of the following columns: %s",
						col,
						strings.Join(internal.DefaultActionOutputColumns, ", "),
					)
				}
			}
			outputTable.Append(row)
		}

		_, err = fmt.Fprintln(outputTableOutput)
		if err != nil {
			return outputTableOutput, err
		}
		outputTable.Render()

		_, err = fmt.Fprintln(outputTableOutput)
		if err != nil {
			return outputTableOutput, err
		}

		_, err = fmt.Fprint(outputTableOutput, internal.OutputAutoDocEnd)
		if err != nil {
			return outputTableOutput, err
		}
	}
	return outputTableOutput, nil
}
