package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/logger"
)

var output string
var configs []string
var rootCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate file from template",
	Long:  `A helper tool for generate file from template`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		t, err := template.ParseFiles(src)
		if err != nil {
			logger.Event("can-not-parse-file").Field("error", err).Fatal()
		}

		fin, err := os.Stat(src)
		if err != nil {
			logger.Event("invalid-input-file").Field("error", err).Fatal()
		}

		if output == "" {
			if strings.HasSuffix(src, ".template") {
				output = src[:len(src)-9]
			} else {
				output = src + ".out"
			}
		}

		fout, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fin.Mode())
		if err != nil {
			logger.Event("failed-to-open-output").Field("error", err).Fatal()
		}

		for i := range configs {
			if err := config.Add(configs[i]); err != nil {
				logger.Event("add-config-error").Field("error", err).Fatal()
			}
		}

		if err := t.Execute(fout, config.ToMap()); err != nil {
			logger.Event("failed-to-execute-template").Field("error", err).Fatal()
		}

		fmt.Println("Generate file successfully")
	},
}

func main() {
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Specifiy the output path")
	rootCmd.Flags().StringArrayVarP(&configs, "config", "c", nil, "Specify overridden config files")
	if err := rootCmd.Execute(); err != nil {
		logger.Event("template-command-error").Field("error", err).Fatal()
	}
}
