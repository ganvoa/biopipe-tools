package fasta

import "github.com/ganvoa/biopipe-tools/internal"

type downloaderPersistent struct {
	downloader *Downloader
	repository FastaRepository
	logger     internal.Logger
}

func NewDownloaderPersistent(downloader *Downloader, repository FastaRepository, logger internal.Logger) *downloaderPersistent {
	d := new(downloaderPersistent)
	d.downloader = downloader
	d.repository = repository
	d.logger = logger
	return d
}

func (d downloaderPersistent) Download(strainId int) error {

	d.logger.Debugf("retrieving strainId %d", strainId)
	strain, err := d.repository.GetByStrainId(strainId)
	if err != nil {
		return err
	}

	d.logger.Debugf("assemblyId %d", strain.AssemblyId)
	err = d.downloader.Download(strain.AssemblyId)
	if err != nil {
		return err
	}

	d.logger.Debugf("marking strainId %d as downloaded", strainId)
	err = d.repository.MarkAsDownloaded(strainId)
	if err != nil {
		return err
	}

	return nil
}
