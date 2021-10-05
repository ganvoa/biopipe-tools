package fasta

import (
	"fmt"
	"log"
)

type DownloaderPersistentBatch struct {
	downloader *Downloader
	repository *FastaRepositoryElasticSearch
}

func NewDownloaderPersistentBatch(downloader *Downloader, repository *FastaRepositoryElasticSearch) *DownloaderPersistentBatch {
	d := new(DownloaderPersistentBatch)
	d.downloader = downloader
	d.repository = repository
	return d
}

func (d DownloaderPersistentBatch) Download() error {

	strains, err := d.repository.NotDownloaded()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of strains: %v \n", len(strains))

	for _, strain := range strains {
		assemblyId := strain.AssemblyId
		strainId := strain.Id
		fmt.Printf("Downloading strainId: %v assemblyId: %v \n", strainId, assemblyId)

		err = d.downloader.Download(assemblyId)
		if err != nil {
			return err
		}

		err = d.repository.MarkAsDownloaded(strainId)
		if err != nil {
			return err
		}
	}

	return nil
}
