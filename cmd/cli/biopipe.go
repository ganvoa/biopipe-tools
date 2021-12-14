package main

import (
	fasta "github.com/ganvoa/biopipe-tools/internal/fasta/cmd"
	integron "github.com/ganvoa/biopipe-tools/internal/integron/cmd"
	strain "github.com/ganvoa/biopipe-tools/internal/strain/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "biopipe",
	}
	rootCmd.AddCommand(fasta.FastaDownloadCommand())
	rootCmd.AddCommand(fasta.FastaBatchDownloadCommand())
	rootCmd.AddCommand(fasta.FastaPersistentDownloadCommand())
	rootCmd.AddCommand(strain.StrainBackupCommand())
	rootCmd.AddCommand(integron.IntegronParse())
	rootCmd.AddCommand(integron.IntegronFind())
	rootCmd.AddCommand(integron.IntegronClean())
	rootCmd.AddCommand(integron.GetIntegrons())
	rootCmd.AddCommand(integron.NormalizeIntegrons())
	rootCmd.Execute()
}
