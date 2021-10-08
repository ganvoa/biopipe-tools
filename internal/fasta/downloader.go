package fasta

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ganvoa/biopipe-tools/internal"
)

type Downloader struct {
	sessionKey   string
	outputDir    string
	databaseName string
	logger       internal.Logger
}

const downloadUrl = "http://enterobase.warwick.ac.uk/upload/download"

func NewDownloader(sessionKey string, outputDir string, databaseName string, logger internal.Logger) *Downloader {
	d := new(Downloader)
	d.sessionKey = sessionKey
	d.outputDir = outputDir
	d.databaseName = databaseName
	d.logger = logger
	return d
}

func (d Downloader) Download(assemblyId int) error {

	d.logger.Infof("downloading assemblyId %d", assemblyId)

	outputFilePath := d.outputDir + "/" + strconv.Itoa(assemblyId) + ".fasta"
	d.logger.Debugf("output file: %s", outputFilePath)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	d.logger.Debugf("new request to url %s", downloadUrl)
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	queryParams := url.Values{}
	queryParams.Add("assembly_id", strconv.Itoa(assemblyId))
	queryParams.Add("database", d.databaseName)
	req.URL.RawQuery = queryParams.Encode()

	d.logger.Debugf("assemblyId: %d", assemblyId)
	d.logger.Debugf("databaseName: %s", d.databaseName)

	sessionCookie := &http.Cookie{
		Name:  "session",
		Value: d.sessionKey,
	}
	req.AddCookie(sessionCookie)

	d.logger.Infof("get request to %s", downloadUrl)
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("file not found")
	}

	d.logger.Infof("saving to file to %s", outputFilePath)
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}

	return nil
}
