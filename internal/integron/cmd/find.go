package cmd

import (
	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_integron_find = "integron:find"

func IntegronFind() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_integron_find,
		Args:  cobra.MinimumNArgs(2),
		Short: "runs integron finder on a specified fasta file",
		Run:   runIntegronFind,
	}

	cmd.Flags().StringP("output", "o", "output/", "Output directory")
	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

func runIntegronFind(cmd *cobra.Command, args []string) {
	godotenv.Load()

	outputDir, _ := cmd.Flags().GetString("output")
	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_integron_find, debug)
	logger.Debug("started")

	fastaFileDir := args[0]
	fastaFileName := args[2]

	logger.Debugf("output %s", outputDir)

	finder := integron.NewIntegronFinder(outputDir, logger)
	resultsFolder, err := finder.Run(fastaFileDir, fastaFileName)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debugf("result on %s", resultsFolder)
}
