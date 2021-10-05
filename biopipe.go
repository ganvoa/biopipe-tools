package main

import (
	fasta "github.com/ganvoa/biopipe-tools/fasta/cmd"
	strain "github.com/ganvoa/biopipe-tools/strain/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "biopipe",
	}
	rootCmd.AddCommand(fasta.FastaDownloadCommand())
	rootCmd.AddCommand(fasta.FastaBatchDownloadCommand())
	rootCmd.AddCommand(strain.StrainBackupCommand())
	rootCmd.Execute()
}
