package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/leyrem/thesis_crawler/crawler"
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
	header := []string{"Name", "Repository", "Category", "TagName", "TagLastPushed", "TagDigest", "TagSize", "StarCount", "PullCount", "DateDownloaded", "SubCategories"}
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
				strings.Join(image.SubCategories, "/"),
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

	// Define command-line flags
	executionCategory := flag.String("category", "", "the category of images to be collected, can be one of the following: official, verified, sponsored")
	inputImages := flag.String("path_images", "", "the path of the txt file containing the images to be collected")
	outputFile := flag.String("output_file", "", "the path of the CSV file where the collected images will be saved")

	flag.Parse()

	if *executionCategory == "verified" || *executionCategory == "sponsored" {
		if *inputImages == "" {
			panic(fmt.Errorf("error, you must specify the flag '-path_images=' to indicate the images to collect"))
		}
		if *outputFile == "" {
			panic(fmt.Errorf("error, you must specify the flag '-output_file=' to indicate the file path of the CSV file to save"))
		}
	}

	var err error
	switch *executionCategory {
	case "verified":
		err = collectVerified(*inputImages, *outputFile, "verified")
	case "official":
		err = collectOfficial()
	case "sponsored":
		err = collectVerified(*inputImages, *outputFile, "sponsored")
	default:
		panic(fmt.Errorf("error, you must specify the flag '-category=' with one of these values: official, verified or sponsored"))
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("DONE")
}
