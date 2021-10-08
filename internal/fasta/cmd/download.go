package cmd

import (
	"os"
	"strconv"

	"github.com/ganvoa/biopipe-tools/internal/fasta"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_fasta_download = "fasta:download"

func FastaDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_fasta_download,
		Args:  cobra.MinimumNArgs(1),
		Short: "Downloads a fasta file from enterobase",
		Run:   runDownload,
	}

	cmd.Flags().StringP("output", "o", "output", "Output directory")
	cmd.MarkFlagRequired("output")
	cmd.Flags().Bool("debug", false, "Debug")
	cmd.Flags().StringP("database", "d", "ecoli", "Database name")

	return cmd
}

func runDownload(cmd *cobra.Command, args []string) {
	godotenv.Load()

	sessionKey := os.Getenv("ENTEROBASE_SESSION")
	outputDir, _ := cmd.Flags().GetString("output")
	databaseName, _ := cmd.Flags().GetString("database")
	debug, _ := cmd.Flags().GetBool("debug")

	logger := platform.NewLogger(command_fasta_download, debug)
	logger.Debug("started")

	assemblyId, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Fatal(err)
	}

	downloader := fasta.NewDownloader(sessionKey, outputDir, databaseName, logger)
	err = downloader.Download(assemblyId)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("ending")
}
