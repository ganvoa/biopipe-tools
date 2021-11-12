package cmd

import (
	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const command_integron_find = "integron:find"

func IntegronFind() *cobra.Command {
	cmd := &cobra.Command{
		Use:   command_integron_find,
		Args:  cobra.MinimumNArgs(1),
		Short: "runs integron finder on a specified fasta file",
		Run:   runIntegronFind,
	}

	cmd.Flags().StringP("output", "o", "output/", "Output directory")
	cmd.Flags().Bool("debug", false, "Debug")

	return cmd
}

/*
	si tiene sat2 se debe considerar como integron tipo 2 y guardar la cadena completa
	protein protein no sirve
	"qac" no siempre aparece aunque sea complete
	calin mantener si es que no tiene protein protein
	calin no trae el int1

	complete
	intI ANT_3pp_AadA1-NCBIFAM attC AAC_3_VIa-NCBIFAM attC 			-> ANT_3pp_AadA1-NCBIFAM|attC|AAC_3_VIa-NCBIFAM : integron
	intI AAC_3_VIa-NCBIFAM attC ANT_3pp_AadA1-NCBIFAM 				-> AAC_3_VIa-NCBIFAM|attC|ANT_3pp_AadA1-NCBIFAM : integron
	intI ANT_3pp_AadA1-NCBIFAM attC AAC_3_VIa-NCBIFAM attC-"qac" 	-> ANT_3pp_AadA1-NCBIFAM|attC|AAC_3_VIa-NCBIFAM

	calin
	ANT_3pp_AadA1-NCBIFAM-attC-proewasdas-attC -> propetxxx|attC|proewasdas
	propetxxx-attC-proewasdas -> propetxxx|attC|proewasdas

*/
func runIntegronFind(cmd *cobra.Command, args []string) {
	godotenv.Load()

	outputDir, _ := cmd.Flags().GetString("output")
	debug, _ := cmd.Flags().GetBool("debug")
	logger := platform.NewLogger(command_integron_find, debug)
	logger.Debug("started")

	fastaFilePath := args[0]

	logger.Debugf("output %s", outputDir)

	finder := integron.NewIntegronFinder(outputDir, logger)
	err := finder.Run(fastaFilePath)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("ending")
}
