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

	totalDownloaded := 0
	totalFailed := 0
	currentFrom := 10000000000
	for {

		strains, err := d.repository.NotDownloaded(currentFrom)
		if err != nil {
			return err
		}

		d.logger.Infof("downloading %d strains", len(strains))
		if len(strains) == 0 {
			break
		}

		for _, strain := range strains {

			assemblyId := strain.AssemblyId
			strainId := strain.Id
			currentFrom = strainId

			d.logger.Debugf("downloading strainId %d", strainId)
			d.logger.Debugf("assemblyId %d", assemblyId)

			err = d.downloader.Download(assemblyId)
			if err != nil {
				totalFailed = totalFailed + 1
				d.logger.Warnf("couldnt donwload strainId %d", strainId)
				continue
			}

			d.logger.Debugf("marking strainId %d as downloaded", strainId)
			err = d.repository.MarkAsDownloaded(strainId)
			if err != nil {
				totalFailed = totalFailed + 1
				d.logger.Warnf("couldnt mark as downloaded strainId %d", strainId)
				continue
			}

			totalDownloaded = totalDownloaded + 1
		}

		d.logger.Infof("downloaded until now %d", totalDownloaded)
		d.logger.Infof("failed until now %d", totalFailed)
	}

	d.logger.Infof("total downloaded %d", totalDownloaded)
	d.logger.Infof("total failed %d", totalFailed)

	return nil
}
