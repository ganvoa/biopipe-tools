package fasta

type FastaRepository interface {
	NotDownloaded() ([]Strain, error)
	MarkAsDownloaded(strainId int) error
}
