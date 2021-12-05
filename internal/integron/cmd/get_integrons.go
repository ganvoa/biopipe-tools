package cmd

import (
	"os"

	"github.com/ganvoa/biopipe-tools/internal/fasta"
	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_integron_get = "integron:get"

func GetIntegrons() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_integron_get,
		Args:  cobra.MinimumNArgs(1),
		Short: "runs integron finder on a all downloaded fasta files",
		Run:   runGetIntegrons,
	}

	cmd.Flags().StringP("output", "o", "output/", "Output directory")
	cmd.Flags().StringP("index", "i", "enterobase", "Elasticsearch Index Name")
	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

func runGetIntegrons(cmd *cobra.Command, args []string) {
	godotenv.Load()

	outputDir, _ := cmd.Flags().GetString("output")
	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_integron_get, debug)
	logger.Debug("started")

	downloadDir := args[0]

	logger.Debugf("output %s", outputDir)
	logger.Debugf("download %s", downloadDir)

	indexName, _ := cmd.Flags().GetString("index")

	client, err := platform.NewElasticSearchConnection(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	if err != nil {
		logger.Fatal(err)
	}

	repository := fasta.NewRepository(indexName, client)
	finder := integron.NewIntegronFinder(outputDir, logger)
	cleaner := integron.NewIntegronResultCleaner(logger)
	parser := integron.NewParser(logger)

	getIntegrons := integron.NewGetIntegron(finder, cleaner, parser, repository, logger)
	err = getIntegrons.Run(downloadDir, outputDir)
	if err != nil {
		logger.Fatal(err)
	}
}
