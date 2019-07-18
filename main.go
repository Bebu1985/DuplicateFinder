package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type advancedFileInfo struct {
	FileInfo os.FileInfo
	FilePath string
}

func main() {
	startTime := time.Now()

	directory := flag.String("d", ".", "the directory which is crawled")
	hostFile := flag.String("i", "", "textfile which contains the hostnames")
	outputDirectory := flag.String("o", "", "directory for output csv file. The directory must exit!")
	flag.Parse()

	if *hostFile == "" {
		log.Fatal("no hostfile given")
	}

	hostNames, err := getHostNames(*hostFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	dataInfos := getFileInformation(*directory)
	if err != nil {
		log.Println(err)
	}

	fileRows := buildFileRows(dataInfos, hostNames)

	outputFileName := buildFileName(*directory)

	writeOutputFile(*outputDirectory, outputFileName, fileRows)

	log.Printf("runtime: %s\n", time.Since(startTime))
}

func getHostNames(hostsFilePath string) ([]string, error) {
	file, err := os.Open(hostsFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var hostNames []string
	for scanner.Scan() {
		hostNames = append(hostNames, scanner.Text())
	}

	return hostNames, nil
}

func getFileInformation(directory string) map[string]advancedFileInfo {
	dataInfos := make(map[string]advancedFileInfo)
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				advancedFileInfo := advancedFileInfo{FileInfo: info, FilePath: path}
				dataInfos[info.Name()] = advancedFileInfo
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return dataInfos
}

func buildFileRows(fileInfos map[string]advancedFileInfo, hostNames []string) []string {
	fileRows := []string{}
	for _, fileInfo := range fileInfos {
		for _, hostName := range hostNames {
			if strings.Contains(fileInfo.FileInfo.Name(), hostName) {
				row := fmt.Sprintf("%s,%s,%s,%s\n", fileInfo.FilePath, fileInfo.FileInfo.Name(), byteCountDecimal(fileInfo.FileInfo.Size()), fileInfo.FileInfo.ModTime().Format("2006-01-02 15:04:05"))
				fileRows = append(fileRows, row)
			}
		}
	}
	return fileRows
}

func buildFileName(directory string) string {
	baseDirectoryName := filepath.Base(directory)
	timestamp := time.Now()
	stringTimeStamp := fmt.Sprintf("%02d-%02d-%02d--%02d-%02d-%02d", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second())
	outputFileName := baseDirectoryName + "_" + stringTimeStamp + ".csv"
	outputFileName = strings.Replace(outputFileName, " ", "", -1)

	return outputFileName
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func writeOutputFile(outputDirectory string, outputFileName string, fileRows []string) {
	outputFile, err := os.Create(filepath.Join(outputDirectory, outputFileName))
	if err != nil {
		outputFile.Close()
		log.Fatalf("error creating file: %s\n", err.Error())
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, "Voller Pfad,Dateiname,Groesse,Letzte Aenderung")
	for _, fileRow := range fileRows {
		fmt.Fprint(outputFile, fileRow)
	}
}
