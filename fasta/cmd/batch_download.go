package cmd

import (
	"log"
	"os"

	"github.com/ganvoa/biopipe-tools/fasta"
	"github.com/ganvoa/biopipe-tools/platform/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func FastaBatchDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fasta:batch-download",
		Short: "Downloads a fasta file from a pending list",
		Run:   runBatchDownload,
	}

	cmd.Flags().StringP("output", "o", "", "Output directory")
	cmd.MarkFlagRequired("output")

	cmd.Flags().StringP("database", "d", "ecoli", "Database name")
	cmd.Flags().StringP("index", "i", "enterobase", "Elasticsearch Index Name")

	return cmd
}

func runBatchDownload(cmd *cobra.Command, args []string) {

	godotenv.Load()
	sessionKey := os.Getenv("ENTEROBASE_SESSION")

	outputDir, _ := cmd.Flags().GetString("output")
	databaseName, _ := cmd.Flags().GetString("database")
	indexName, _ := cmd.Flags().GetString("index")

	client, err := database.NewClient(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	if err != nil {
		log.Fatal(err)
	}

	repository := fasta.NewRepository(indexName, client)
	downloader := fasta.NewDownloader(sessionKey, outputDir, databaseName)

	downloaderPersistentBatch := fasta.NewDownloaderPersistentBatch(downloader, repository)
	err = downloaderPersistentBatch.Download()
	if err != nil {
		log.Fatal(err)
	}
}
