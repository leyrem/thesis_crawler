package main

import (
	"fmt"
	"strings"

	"crawlproject/crawler"
)

func collectOfficial() error {

	uri1 := "https://hub.docker.com/v2/repositories/library/?page_size=100"

	resp := &crawler.RepositoryResponse{}
	resp2 := &crawler.RepositoryResponse{}

	err := crawler.ExecuteRequest(uri1, resp, nil, nil)
	if err != nil {
		return err
	} else {
		uri2 := resp.Next
		err = crawler.ExecuteRequest(uri2, resp2, nil, nil)
		if err != nil {
			return err
		}
		if resp2.Next != "" {
			return fmt.Errorf("there are more results than additional pages, execute request again!!")
		}
	}
	combinedResults := append(resp.Results, resp2.Results...)
	fmt.Printf("There are %d Official images collected \n", len(combinedResults))

	deprecatedCounter := 0
	inactiveCounter := 0
	var finalResults []crawler.ImageCollected
	c := 1

	for _, image := range combinedResults {
		if strings.Contains(image.Description, "DEPRECATED") {
			deprecatedCounter++
		}
		if image.StatusDescription != "active" {
			inactiveCounter++
		}
		if !strings.Contains(image.Description, "DEPRECATED") && image.StatusDescription == "active" {
			fmt.Println("--> EXECUTING request: ", c)
			c++

			// Add the category types
			var subCat []string
			for _, cat := range image.Categories {
				subCat = append(subCat, cat.Slug)
			}

			f := crawler.ImageCollected{
				Name:          image.Name,
				Repository:    "library",
				Category:      "Official",
				StarCount:     image.StarCount,
				PullCount:     image.PullCount,
				SubCategories: subCat,
				//DateDownloaded: ,
			}

			uri := "https://hub.docker.com/v2/repositories/library/" + f.Name + "/tags/?page_size=100"
			tagsInfo, err := crawler.ExecuteTagRequest(uri, f.Name)
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

	if err := saveToCSV("imagesOfficial.csv", finalResults); err != nil {
		return err
	}
	return nil
}
