package cmd

import (
	"os"

	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_integron_parse_persist = "integron:parse-persist"

func IntegronParsePersist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_integron_parse_persist,
		Args:  cobra.MinimumNArgs(1),
		Short: "reads a tsv file and adds it to elasticsearch strain",
		Run:   run,
	}

	cmd.Flags().IntP("strain", "s", 0, "StrainId")
	cmd.Flags().Bool("debug", false, "Debug")
	cmd.Flags().StringP("index", "i", "enterobase", "Elasticsearch Index Name")

	cmd.MarkFlagRequired("strain")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	godotenv.Load()

	strainId, _ := cmd.Flags().GetInt("strain")
	debug, _ := cmd.Flags().GetBool("debug")
	indexName, _ := cmd.Flags().GetString("index")
	logger := platform.NewLogger(command_integron_parse, debug)
	logger.Debug("started")

	filePath := args[0]

	client, err := platform.NewElasticSearchConnection(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	repository := integron.NewRepository(indexName, client)
	persistent := integron.NewIntegronPersistentES(strainId, repository, logger)
	parser := integron.NewParser(filePath, persistent, logger)

	err = parser.Parse()
	if err != nil {
		logger.Fatal(err)
	}

	err = parser.Save()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("ending")
}
