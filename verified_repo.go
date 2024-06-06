package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"crawlproject/crawler"
)

func ReadFile(filePath string) ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, fmt.Errorf("error opening file: %s", err.Error())
	}
	defer file.Close()

	// Create a new scanner for the file
	scanner := bufio.NewScanner(file)

	var lines []string
	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("error reading file: %s", err.Error())
	}
	return lines, nil
}

func collectVerified(inputFile string, outputFile string, categoryType string) error {

	images, err := ReadFile(inputFile)
	if err != nil {
		return err
	}

	deprecatedCounter := 0
	inactiveCounter := 0
	var finalResults []crawler.ImageCollected
	c := 1

	for i := 0; i < len(images); {
		img := images[i]
		uri := "https://hub.docker.com/v2/repositories/" + img + "/"
		repoImgResp := &crawler.RepositoryImageResponse{}

		fmt.Printf("About to get info for image: %s\n", img)

		err := crawler.ExecuteRequest(uri, nil, nil, repoImgResp)
		if err != nil {
			if err.Error() == "Non-OK HTTP status: 429" {
				fmt.Println("------- WE'VE BEEN RATE LIMITED, SLEEP 5 SECONDS ------")
				time.Sleep(5 * time.Second)
				continue
			} else {
				return err
			}
		}
		i++

		if strings.Contains(repoImgResp.Description, "DEPRECATED") {
			deprecatedCounter++
		}
		if repoImgResp.StatusDescription != "active" {
			inactiveCounter++
		}

		//fmt.Printf("Got info for image %s: %+v\n", img, repoImgResp)
		if !strings.Contains(repoImgResp.Description, "DEPRECATED") && repoImgResp.StatusDescription == "active" {
			fmt.Println("--> EXECUTING request: ", c)
			c++

			var subCat []string
			for _, cat := range repoImgResp.Categories {
				subCat = append(subCat, cat.Slug)
			}

			f := crawler.ImageCollected{
				Name:          repoImgResp.Name,
				Repository:    repoImgResp.Namespace,
				Category:      categoryType,
				StarCount:     repoImgResp.StarCount,
				PullCount:     repoImgResp.PullCount,
				SubCategories: subCat,
				//DateDownloaded: ,
			}

			uriTags := "https://hub.docker.com/v2/repositories/" + f.Repository + "/" + f.Name + "/tags/?page_size=100"
			tagsInfo, err := crawler.ExecuteTagRequest(uriTags, f.Name)
			if err != nil {
				return err
			} else if len(tagsInfo) == 0 {
				continue
			}
			f.Tags = append(f.Tags, tagsInfo...)
			finalResults = append(finalResults, f)
		}

	}
	fmt.Printf("There are %d deprecated images\n", deprecatedCounter)
	fmt.Printf("There are %d inactive images\n", inactiveCounter)

	if err := saveToCSV(outputFile, finalResults); err != nil {
		return err
	}
	return nil
}
