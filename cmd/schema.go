/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/scs/sbom-utility/schema"
	"github.com/spf13/cobra"
)

var SCHEMA_LIST_TITLES = []string{"Format", "Version", "Variant", "File", "Source"}

func NewCommandSchema() *cobra.Command {
	var command = new(cobra.Command)
	command.Use = "schema"
	command.Short = "view built-in SBOM schemas."
	command.Long = fmt.Sprintf("view built-in SBOM schemas supported by the utility. The default command produces a list based upon `%s`", DEFAULT_SCHEMA_CONFIG)
	command.RunE = schemaCmdImpl
	initCommandSchema(command)
	return command
}

func initCommandSchema(command *cobra.Command) {
	getLogger().Enter()
	defer getLogger().Exit()

	command.Flags().Bool("list", true, "List all configured schemas by format")
	//rootCmd.AddCommand(command)
}

func formatSchemaVariant(variant string) (formattedVariant string) {
	var variantName string = schema.SCHEMA_VARIANT_LATEST
	if variant != "" {
		variantName = variant
	}
	formattedVariant = "(" + variantName + ")"
	return
}

func schemaCmdImpl(cmd *cobra.Command, args []string) error {
	getLogger().Enter()
	defer getLogger().Exit()

	// initialize tabwriter
	w := new(tabwriter.Writer)

	// min-width, tab-width, padding, pad-char, flags
	w.Init(os.Stdout, 8, 2, 2, ' ', 0)

	defer w.Flush()

	if len(schema.SupportedFormatConfig.Formats) > 0 {
		var formatName = ""

		// Create title row and add tabs between column titles for the tabWRiter
		titles, underlines := createTitleRows(SCHEMA_LIST_TITLES, nil)
		fmt.Fprintf(w, "%s\n", strings.Join(titles, "\t"))
		fmt.Fprintf(w, "%s\n", strings.Join(underlines, "\t"))

		for _, format := range (schema.SupportedFormatConfig).Formats {
			formatName = format.CanonicalName

			if len(format.Schemas) > 0 {
				for _, currentSchema := range format.Schemas {
					fmt.Fprintf(w, "%v\t%s\t%s\t%s\t%s\n",
						formatName,
						currentSchema.Version,
						formatSchemaVariant(currentSchema.Variant),
						currentSchema.File,
						currentSchema.Url)
				}
			} else {
				getLogger().Warningf("No supported schemas for format `%s`.\n", formatName)
			}
		}
	} else {
		getLogger().Warningf("No supported built-in formats found in `%s`.\n", DEFAULT_SCHEMA_CONFIG)
	}

	fmt.Fprintln(w, "")
	return nil
}
