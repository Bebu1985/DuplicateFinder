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
	flag.Parse()

	if *hostFile == "" {
		log.Fatal("no hostfile given")
	}

	dataInfos := make(map[string]advancedFileInfo)

	hostNames, err := getHostNames(*hostFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = filepath.Walk(*directory,
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

	for _, fileInfo := range dataInfos {
		for _, hostName := range hostNames {
			if strings.Contains(fileInfo.FileInfo.Name(), hostName) {
				fmt.Printf("Full Path: %s, Name: %s; Size: %d; Modified: %v;\n", fileInfo.FilePath, fileInfo.FileInfo.Name(), fileInfo.FileInfo.Size(), fileInfo.FileInfo.ModTime())
			}
		}
	}

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
