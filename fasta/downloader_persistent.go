package fasta

type DownloaderPersistent struct {
	downloader *Downloader
	repository *FastaRepositoryElasticSearch
}

func NewDownloaderPersistent(downloader *Downloader, repository *FastaRepositoryElasticSearch) *DownloaderPersistent {
	d := new(DownloaderPersistent)
	d.downloader = downloader
	d.repository = repository
	return d
}

func (d DownloaderPersistent) Download(assemblyId int) error {

	err := d.downloader.Download(assemblyId)
	if err != nil {
		return err
	}

	err = d.repository.MarkAsDownloaded(assemblyId)
	if err != nil {
		return err
	}

	return nil
}
