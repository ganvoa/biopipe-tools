package cmd

import (
	"os"
	"strconv"

	"github.com/ganvoa/biopipe-tools/internal/fasta"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_fasta_persistent_download = "fasta:persistent-download"

func FastaPersistentDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_fasta_persistent_download,
		Args:  cobra.MinimumNArgs(1),
		Short: "Downloads a fasta file from a list and marks it downloaded",
		Run:   runPersistentDownload,
	}

	cmd.Flags().StringP("output", "o", "", "Output directory")
	cmd.MarkFlagRequired("output")
	cmd.Flags().Bool("debug", false, "Debug")
	cmd.Flags().StringP("database", "d", "ecoli", "Database name")

	return cmd
}

func runPersistentDownload(cmd *cobra.Command, args []string) {

	godotenv.Load()

	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_fasta_persistent_download, debug)
	logger.Debug("started")

	sessionKey := os.Getenv("ENTEROBASE_SESSION")
	outputDir, _ := cmd.Flags().GetString("output")
	databaseName, _ := cmd.Flags().GetString("database")

	fastaId, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Fatal(err)
	}

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
	downloader := fasta.NewDownloader(sessionKey, outputDir, databaseName, logger)

	downloaderPersistent := fasta.NewDownloaderPersistent(downloader, repository, logger)
	err = downloaderPersistent.Download(fastaId)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("ending")
}
