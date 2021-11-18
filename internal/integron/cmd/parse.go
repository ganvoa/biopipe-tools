package cmd

import (
	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_integron_parse = "integron:parse"

func IntegronParse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_integron_parse,
		Args:  cobra.MinimumNArgs(1),
		Short: "transforms a tsv file into a json file based on a integron finder result",
		Run:   runIntegronParse,
	}

	cmd.Flags().StringP("output", "o", "output/", "Output directory")
	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

func runIntegronParse(cmd *cobra.Command, args []string) {
	godotenv.Load()

	outputDir, _ := cmd.Flags().GetString("output")
	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_integron_parse, debug)
	logger.Debug("started")

	folder := args[0]

	logger.Debugf("output %s", outputDir)

	parser := integron.NewParser(logger)

	_, err := parser.Parse(folder)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("ending")
}
