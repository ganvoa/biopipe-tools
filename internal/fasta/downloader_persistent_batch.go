package fasta

import (
	"github.com/ganvoa/biopipe-tools/internal"
)

type downloaderPersistentBatch struct {
	downloader *Downloader
	repository FastaRepository
	logger     internal.Logger
}

func NewDownloaderPersistentBatch(downloader *Downloader, repository FastaRepository, logger internal.Logger) *downloaderPersistentBatch {
	d := new(downloaderPersistentBatch)
	d.downloader = downloader
	d.repository = repository
	d.logger = logger
	return d
}

func (d downloaderPersistentBatch) Download() error {

	d.logger.Debug("retrieving not downloaded strains")
	strains, err := d.repository.NotDownloaded()
	if err != nil {
		return err
	}
	d.logger.Infof("found %d strains", len(strains))

	for _, strain := range strains {
		assemblyId := strain.AssemblyId
		strainId := strain.Id
		d.logger.Infof("downloading strainId: %d assemblyId: %d", strainId, assemblyId)

		err = d.downloader.Download(assemblyId)
		if err != nil {
			return err
		}

		d.logger.Debugf("marking strainId %d as downloaded", strainId)
		err = d.repository.MarkAsDownloaded(strainId)
		if err != nil {
			return err
		}
	}

	return nil
}
