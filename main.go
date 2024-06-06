package main

import (
	"crawlproject/crawler"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func saveToCSV(filePath string, data []crawler.ImageCollected) error {

	// Open file for writing
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Name", "Repository", "Category", "TagName", "TagLastPushed", "TagDigest", "TagSize", "StarCount", "PullCount", "DateDownloaded"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Failed to write header to file:", err)
		panic(err)
	}
	writer.Write(header)

	// Loop through data and convert to string slices
	for _, image := range data {

		for _, tag := range image.Tags {
			record := []string{
				image.Name,
				image.Repository,
				image.Category,
				tag.Name,
				tag.LastPushed.Format(time.RFC3339),
				tag.Digest,
				strconv.FormatInt(tag.Size, 10),
				strconv.FormatInt(image.StarCount, 10),
				strconv.FormatInt(image.PullCount, 10),
				image.DateDownloaded.Format(time.RFC3339),
			}

			if err := writer.Write(record); err != nil {
				fmt.Println("Failed to write record to file:", err)
				panic(err)
			}
		}

	}

	return writer.Error()
}

func main() {
	err := collectVerified()
	if err != nil {
		panic(err)
	}
	fmt.Println("DONE")
}
