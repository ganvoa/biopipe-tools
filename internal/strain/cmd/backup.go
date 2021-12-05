package cmd

import (
	"os"

	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/ganvoa/biopipe-tools/internal/strain"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_strain_backup = "strain:backup"

func StrainBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_strain_backup,
		Args:  cobra.MinimumNArgs(1),
		Short: "Backups a json file from enterobase to elasticsearch",
		Run:   run,
	}

	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	godotenv.Load()

	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_strain_backup, debug)
	logger.Debug("started")

	filePath := args[0]

	indexName := os.Getenv("ELASTICSEARCH_INDEX")
	client, err := platform.NewElasticSearchConnection(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	if err != nil {
		logger.Fatal(err)
	}

	repository := strain.NewRepository(indexName, client, logger)
	parser := strain.NewStrainParser(filePath)
	backuper := strain.NewStrainBackuper(repository, parser, logger)
	err = backuper.Backup()
	if err != nil {
		logger.Fatal(err)
	}
}
