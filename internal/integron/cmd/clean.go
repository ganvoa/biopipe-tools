package cmd

import (
	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_integron_clean = "integron:clean"

func IntegronClean() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_integron_clean,
		Args:  cobra.NoArgs,
		Short: "cleans integron finder results folder",
		Run:   runIntegronClean,
	}

	cmd.Flags().StringP("folder", "f", "output/Results_Integron_Finder", "integron finder result folder")
	cmd.Flags().Bool("debug", false, "Debug")
	cmd.MarkFlagRequired("folder")
	return cmd
}

func runIntegronClean(cmd *cobra.Command, args []string) {
	godotenv.Load()

	folderDir, _ := cmd.Flags().GetString("folder")
	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_integron_clean, debug)
	logger.Debug("started")

	logger.Debugf("result folder %s", folderDir)

	finder := integron.NewIntegronResultCleaner(logger)
	err := finder.Clean(folderDir)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("ending")
}
