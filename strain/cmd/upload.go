package cmd

import (
	"log"
	"os"

	"github.com/ganvoa/biopipe-tools/platform/database"
	"github.com/ganvoa/biopipe-tools/strain"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func StrainBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strain:backup [enterobase.json]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Backups a json file from enterobase to elasticsearch",
		Run:   runBackup,
	}

	cmd.Flags().StringP("index", "i", "enterobase", "Elasticsearch Index Name")

	return cmd
}

func runBackup(cmd *cobra.Command, args []string) {
	godotenv.Load()
	indexName, _ := cmd.Flags().GetString("index")
	filePath := args[0]
	client, err := database.NewClient(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}
	repository := strain.NewRepository(indexName, client)
	parser := strain.NewStrainParser(filePath)
	backuper := strain.NewStrainBackuper(repository, parser)
	err = backuper.Backup()
	if err != nil {
		log.Fatal(err)
	}
}
