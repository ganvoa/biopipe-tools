package cmd

import (
	"path/filepath"

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
		Run:   runBackup,
	}

	cmd.Flags().StringP("output", "o", "output/", "Output directory")
	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

func runBackup(cmd *cobra.Command, args []string) {
	godotenv.Load()

	outputDir, _ := cmd.Flags().GetString("output")
	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_integron_parse, debug)
	logger.Debug("started")

	filePath := args[0]
	_, file := filepath.Split(filePath)

	logger.Debugf("output %s", outputDir)

	outputPath := outputDir + "/" + file + ".json"

	persistent := integron.NewIntegronPersistentFile(outputPath, logger)
	parser := integron.NewParser(filePath, persistent, logger)

	err := parser.Parse()
	if err != nil {
		logger.Fatal(err)
	}

	err = parser.Save()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("ending")
}