package cmd

import (
	"os"

	"github.com/ganvoa/biopipe-tools/internal/fasta"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_fasta_batch_download = "fasta:batch-download"

func FastaBatchDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_fasta_batch_download,
		Short: "Downloads a fasta file from a pending list",
		Run:   runBatchDownload,
	}

	cmd.Flags().StringP("output", "o", "", "Output directory")
	cmd.MarkFlagRequired("output")
	cmd.Flags().Bool("debug", false, "Debug")
	cmd.Flags().StringP("database", "d", "ecoli", "Database name")
	cmd.Flags().StringP("index", "i", "enterobase", "Elasticsearch Index Name")

	return cmd
}

func runBatchDownload(cmd *cobra.Command, args []string) {

	godotenv.Load()

	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_fasta_batch_download, debug)
	logger.Debug("started")

	sessionKey := os.Getenv("ENTEROBASE_SESSION")

	outputDir, _ := cmd.Flags().GetString("output")
	databaseName, _ := cmd.Flags().GetString("database")
	indexName, _ := cmd.Flags().GetString("index")

	client, err := platform.NewClient(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	if err != nil {
		logger.Fatal(err)
	}

	repository := fasta.NewRepository(indexName, client)
	downloader := fasta.NewDownloader(sessionKey, outputDir, databaseName, logger)

	downloaderPersistentBatch := fasta.NewDownloaderPersistentBatch(downloader, repository, logger)
	err = downloaderPersistentBatch.Download()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("ending")
}
