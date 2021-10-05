package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/ganvoa/biopipe-tools/fasta"
	"github.com/ganvoa/biopipe-tools/platform/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func PersistentDownload() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fasta:persistent-download",
		Args:  cobra.MinimumNArgs(1),
		Short: "Downloads a fasta file from a list and marks it downloaded",
		Run:   runPersistentDownload,
	}

	cmd.Flags().StringP("output", "o", "", "Output directory")
	cmd.MarkFlagRequired("output")

	cmd.Flags().StringP("database", "d", "ecoli", "Database name")
	cmd.Flags().StringP("index", "i", "enterobase", "Elasticsearch Index Name")

	return cmd
}

func runPersistentDownload(cmd *cobra.Command, args []string) {

	godotenv.Load()
	sessionKey := os.Getenv("ENTEROBASE_SESSION")
	outputDir, _ := cmd.Flags().GetString("output")
	databaseName, _ := cmd.Flags().GetString("database")
	indexName, _ := cmd.Flags().GetString("index")

	fastaId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

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

	downloaderPersistent := fasta.NewDownloaderPersistent(downloader, repository)
	err = downloaderPersistent.Download(fastaId)
	if err != nil {
		log.Fatal(err)
	}
}
