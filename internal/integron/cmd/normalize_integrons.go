package cmd

import (
	"os"

	"github.com/ganvoa/biopipe-tools/internal/fasta"
	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_normalize_integron = "integron:normalize"

func NormalizeIntegrons() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_normalize_integron,
		Short: "normalize integrons found",
		Run:   runNormalizeIntegrons,
	}

	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

func runNormalizeIntegrons(cmd *cobra.Command, args []string) {
	godotenv.Load()

	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_normalize_integron, debug)
	logger.Debug("started")
	indexName := os.Getenv("ELASTICSEARCH_INDEX")

	client, err := platform.NewElasticSearchConnection(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	if err != nil {
		logger.Fatal(err)
	}

	repository := fasta.NewRepository(indexName, client)
	cleaner := integron.NewIntegronResultCleaner(logger)
	parser := integron.NewParser(logger)

	normalizer := integron.NewIntegronNormalizer(cleaner, parser, repository, logger)
	err = normalizer.Run()
	if err != nil {
		logger.Fatal(err)
	}
}
