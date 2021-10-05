package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/ganvoa/biopipe-tools/fasta"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func FastaDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fasta:download",
		Args:  cobra.MinimumNArgs(1),
		Short: "Downloads a fasta file from enterobase",
		Run:   runDownload,
	}

	cmd.Flags().StringP("output", "o", "", "Output directory")
	cmd.MarkFlagRequired("output")

	cmd.Flags().StringP("database", "d", "ecoli", "Database name")

	return cmd
}

func runDownload(cmd *cobra.Command, args []string) {
	godotenv.Load()
	sessionKey := os.Getenv("ENTEROBASE_SESSION")
	outputDir, _ := cmd.Flags().GetString("output")
	databaseName, _ := cmd.Flags().GetString("database")

	assemblyId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}
	downloader := fasta.NewDownloader(sessionKey, outputDir, databaseName)
	err = downloader.Download(assemblyId)
	if err != nil {
		log.Fatal(err)
	}
}
