package fasta

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Downloader struct {
	sessionKey   string
	outputDir    string
	databaseName string
}

const downloadUrl = "http://enterobase.warwick.ac.uk/upload/download"

func NewDownloader(sessionKey string, outputDir string, databaseName string) *Downloader {
	d := new(Downloader)
	d.sessionKey = sessionKey
	d.outputDir = outputDir
	d.databaseName = databaseName
	return d
}

func (d Downloader) Download(assemblyId int) error {

	outputFilePath := d.outputDir + "/" + strconv.Itoa(assemblyId) + ".fasta"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	queryParams := url.Values{}
	queryParams.Add("assembly_id", strconv.Itoa(assemblyId))
	queryParams.Add("database", d.databaseName)
	req.URL.RawQuery = queryParams.Encode()

	sessionCookie := &http.Cookie{
		Name:  "session",
		Value: d.sessionKey,
	}
	req.AddCookie(sessionCookie)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("file not found")
	}

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
