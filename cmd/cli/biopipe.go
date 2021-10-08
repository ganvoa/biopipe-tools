package main

import (
	fasta "github.com/ganvoa/biopipe-tools/internal/fasta/cmd"
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
	rootCmd.Execute()
}
