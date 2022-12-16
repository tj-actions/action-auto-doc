[![Codacy Badge](https://api.codacy.com/project/badge/Grade/bb5dc10c1d2645c0894fa6774300639b)](https://app.codacy.com/gh/tj-actions/auto-doc?utm_source=github.com\&utm_medium=referral\&utm_content=tj-actions/auto-doc\&utm_campaign=Badge_Grade_Settings)
![Coverage](https://img.shields.io/badge/Coverage-84.1%25-brightgreen)
[![Go Reference](https://pkg.go.dev/badge/github.com/tj-actions/auto-doc.svg)](https://pkg.go.dev/github.com/tj-actions/auto-doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/tj-actions/auto-doc)](https://goreportcard.com/report/github.com/tj-actions/auto-doc)
[![CI](https://github.com/tj-actions/auto-doc/workflows/CI/badge.svg)](https://github.com/tj-actions/auto-doc/actions?query=workflow%3ACI)
[![Update release version.](https://github.com/tj-actions/auto-doc/workflows/Update%20release%20version./badge.svg)](https://github.com/tj-actions/auto-doc/actions?query=workflow%3A%22Update+release+version.%22)
[![Public workflows that use this action.](https://img.shields.io/endpoint?url=https%3A%2F%2Fused-by.vercel.app%2Fapi%2Fgithub-actions%2Fused-by%3Faction%3Dtj-actions%2Fauto-doc%26badge%3Dtrue)](https://github.com/search?o=desc\&q=tj-actions+auto-doc+language%3AYAML\&s=\&type=Code)

[![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?logo=ubuntu\&logoColor=white)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Mac OS](https://img.shields.io/badge/mac%20os-000000?logo=macos\&logoColor=F0F0F0)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Windows](https://img.shields.io/badge/Windows-0078D6?logo=windows\&logoColor=white)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)

## auto-doc

Generate documentation for your actions.(yml|yaml).

## Table of Contents

*   [Usage](#usage)
*   [Inputs](#inputs)
*   [Examples](#examples)
*   [CLI](#cli)
    *   [Installation](#installation)
    *   [Synopsis](#synopsis)
    *   [Options](#options)
*   [Credits](#credits)
*   [Report Bugs](#report-bugs)

## Usage

Add the `Inputs` and/or `Outputs` [`H2` header](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#headers) to any markdown file.

```yaml
...
    steps:
      - uses: actions/checkout@v2
      - name: Run auto-doc
        uses: tj-actions/auto-doc@v1.4.1
```

## Inputs

<!-- AUTO-DOC-INPUT:START - Do not remove or modify this section -->

|     INPUT      |  TYPE  | REQUIRED |    DEFAULT     |                               DESCRIPTION                                |
|----------------|--------|----------|----------------|--------------------------------------------------------------------------|
|     action     | string |  false   | `"action.yml"` |                       Path to the action.yml file                        |
|    bin\_path    | string |  false   |                |                       Path to the auto-doc binary                        |
| col\_max\_width  | string |  false   |    `"1000"`    |                          Max width of a column                           |
| col\_max\_words  | string |  false   |     `"7"`      |               Max number of words per line in<br>a column                |
| input\_columns  | string |  false   |                | List of Input columns names to display,<br>default (display all columns) |
|     output     | string |  false   | `"README.md"`  |                         Path to the output file                          |
| output\_columns | string |  false   |                | List of Output column names to display,<br>default (display all columns) |

<!-- AUTO-DOC-INPUT:END -->

**👆 This is generated 👆 using :point\_right: [action.yml](./action.yml)**

## Examples

Create a pull request each time the action.yml inputs/outputs change

```yaml
name: Update README.md with the latest actions.yml

on:
  push:
    branches:
      - main

jobs:
  update-doc:
     runs-on: ubuntu-latest
     steps:
       - name: Checkout
         uses: actions/checkout@v2.4.0
         with:
           fetch-depth: 0  # otherwise, you will failed to push refs to dest repo

       - name: Run auto-doc
         uses: tj-actions/auto-doc@v1.4.1

       - name: Verify Changed files
         uses: tj-actions/verify-changed-files@v8.6
         id: verify-changed-files
         with:
           files: |
             README.md

       - name: Create Pull Request
         if: steps.verify-changed-files.outputs.files_changed == 'true'
         uses: peter-evans/create-pull-request@v3
         with:
           base: "main"
           title: "auto-doc: Updated README.md"
           branch: "chore/auto-doc-update-readme"
           commit-message: "auto-doc: Updated README.md"
           body: "auto-doc: Updated README.md"
```

## CLI

### Installation

Run

```shell script
go install github.com/tj-actions/auto-doc@latest
```

### Synopsis

Auto generate documentation for your github action.

    auto-doc [flags]

### Options

          --action string               action config file (default "action.yml")
          --colMaxWidth string          Max width of a column (default "1000")
          --colMaxWords string          Max number of words per line in a column (default "5")
      -h, --help                        help for auto-doc
          --inputColumns stringArray    list of input column names (default [Input,Type,Required,Default,Description])
          --output string               Output file (default "README.md")
          --outputColumns stringArray   list of output column names (default [Output,Type,Description])

*   Free software: [Apache License 2.0](LICENSE)

If you feel generous and want to show some extra appreciation:

[![Buy me a coffee][buymeacoffee-shield]][buymeacoffee]

[buymeacoffee]: https://www.buymeacoffee.com/jackton1

[buymeacoffee-shield]: https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png

## Credits

This package was created with [Cookiecutter](https://github.com/cookiecutter/cookiecutter) using [cookiecutter-action](https://github.com/tj-actions/cookiecutter-action)

*   [cobra](https://github.com/spf13/cobra)
*   [gobinaries](https://github.com/tj/gobinaries)
*   [goreleaser](https://github.com/goreleaser/goreleaser/)

## Report Bugs

Report bugs at https://github.com/tj-actions/auto-doc/issues.

If you are reporting a bug, please include:

*   Your operating system name and version.
*   Any details about your workflow that might be helpful in troubleshooting.
*   Detailed steps to reproduce the bug.
