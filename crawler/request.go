package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"time"
)

const TAGS_COLLECTED_PER_IMAGE = 5

// ExecuteRequest issues a HTTP request to the given uri and collects response in either one
// of the data structures: RepositoryResponse or TagsResponse
func ExecuteRequest(uri string, repoResponse *RepositoryResponse, tagsResponse *TagsResponse, repoImageResponse *RepositoryImageResponse) error {
	// Make the HTTP GET request
	resp, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non-OK HTTP status: %d", resp.StatusCode)
	}

	if repoResponse == nil && repoImageResponse == nil {
		err = json.Unmarshal(body, &tagsResponse)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal JSON response: %v", err)
		}
		//fmt.Printf("Response: %+v\n", response)
	} else if repoResponse == nil && tagsResponse == nil {
		err = json.Unmarshal(body, &repoImageResponse)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal JSON response: %v", err)
		}
	} else {
		err = json.Unmarshal(body, &repoResponse)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal JSON response: %v", err)
		}
		//fmt.Printf("Response: %+v\n", response)
	}
	return nil
}

// ExecuteRequestWithPagination issues HTTP requests to get the tags for a given image
// and handles pagination.
func ExecuteTagRequest(uriTags, imageName string) ([]TagInfo, error) {
	fmt.Printf("Fetching tags for image: %s\n", imageName)

	tagsResp := &TagsResponse{}

	err := ExecuteRequest(uriTags, nil, tagsResp, nil)
	if err != nil {
		if err.Error() == "Non-OK HTTP status: 429" {
			fmt.Println("-------- WE'VE BEEN RATE LIMITED, SLEEPING 10 SECONDS [at start] ---------")
			time.Sleep(10 * time.Second)
			err = ExecuteRequest(uriTags, nil, tagsResp, nil)
			if err != nil {
				return []TagInfo{}, fmt.Errorf("Error, rate limit occured at start of ExecuteTagRequest for the 2nd time")
			}
		} else {
			return []TagInfo{}, err
		}
	}

	count := tagsResp.Count
	combinedResults := tagsResp.Results

	// Collect all results from other pages too.
	for tagsResp.Next != "" {
		next := tagsResp.Next
		tagsResp = &TagsResponse{}
		err := ExecuteRequest(next, nil, tagsResp, nil)
		if err != nil {
			if err.Error() == "Non-OK HTTP status: 429" {
				fmt.Println("-------- WE'VE BEEN RATE LIMITED, SLEEPING 10 SECONDS [tag request] ---------")
				time.Sleep(10 * time.Second)
				tagsResp.Next = next
				continue
			} else {
				return []TagInfo{}, err
			}
		}
		combinedResults = append(tagsResp.Results, combinedResults...)
	}

	if count != len(combinedResults) {
		return []TagInfo{}, fmt.Errorf("The number of combined results for the tags for image %s and the count are not the same; count = %d, results = %d", imageName, count, len(combinedResults))
	}

	if count == 0 {
		return []TagInfo{}, nil
	}

	var tagsInfo []TagInfo
	for _, tag := range combinedResults {
		if tag.Name == "latest" || tag.TagStatus != "active" {
			continue
		}
		timestamp, err := time.Parse(time.RFC3339Nano, tag.TagLastPushed)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			return []TagInfo{}, err
		}
		tagInfo := TagInfo{
			Name:       tag.Name,
			LastPushed: timestamp,
			Digest:     tag.Digest,   // every tag has multiple images for different architectures, so its better to save this field when i pull the actual image from docker hub
			Size:       tag.FullSize, // size also depends on the architecture
		}
		tagsInfo = append(tagsInfo, tagInfo)
	}

	// Sort by LastPushed (most recent first)
	sort.Slice(tagsInfo, func(i, j int) bool {
		return !tagsInfo[i].LastPushed.Before(tagsInfo[j].LastPushed)
	})
	// If there are less tags than the number we want to collect, ensure we take the min.
	numTake := int(math.Min(float64(len(tagsInfo)), TAGS_COLLECTED_PER_IMAGE))

	fmt.Println("tagsInfo:")
	for i := 0; i < numTake; i++ {
		tag := tagsInfo[i]
		fmt.Printf("  - Name: %s, LastPushed: %s, Digest: %s\n", tag.Name, tag.LastPushed.Format(time.RFC3339), tag.Digest)
	}
	return tagsInfo[:numTake], nil
}
