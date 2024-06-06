package crawler

import "time"

type ImageCollected struct {
	Name           string
	Repository     string
	Category       string
	Tags           []TagInfo
	StarCount      int64
	PullCount      int64
	DateDownloaded time.Time
	SubCategories  []string
}

type TagInfo struct {
	Name       string
	LastPushed time.Time
	Digest     string // Should this be collected when I acc execute the pull?
	Size       int64  // size dependson the architecture
}

type RepositoryResponse struct {
	// Define the fields according to the JSON response structure
	// Example: if the JSON response is {"key1": "value1", "key2": "value2"}
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Results `json:"results"`
}

type TagsResponse struct {
	// Define the fields according to the JSON response structure
	// Example: if the JSON response is {"key1": "value1", "key2": "value2"}
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []TagResponse `json:"results"`
}

type TagResponse struct {
	Creator int64 `json:"creator"`
	ID      int64 `json:"id"`
	//Images              []Image `json:"images"` do I need this??
	LastUpdated         string `json:"last_updated"`
	LastUpdater         int64  `json:"last_updater"`
	LastUpdaterUsername string `json:"last_updater_username"`
	Name                string `json:"name"`
	Repository          int64  `json:"repository"`
	FullSize            int64  `json:"full_size"`
	v2                  bool   `json:"v2"`
	TagStatus           string `json:"tag_status"`
	TagLastPulled       string `json:"tag_last_pulled"`
	TagLastPushed       string `json:"tag_last_pushed"`
	MediaType           string `json:"media_type"`
	ContentType         string `json:"content_type"`
	Digest              string `json:"digest"`
}

/*type Image struct { // this is inside the Tags request, for every tag there is multiple images available depending on the architecture
	Architecture string `json:"architecture"`
	Features     string `json:"features"`
	//Variant
	Digest     string `json:"digest"`
	OS         string `json:"os"`
	OSFeatures string `json:"os_features"`
	OSVersion  string `json:"os_version"`
	Size       int64  `json:"size"`
	Status     string `json:"status"`
	LastPulled string `json:"last_pulled"`
	LastPushed string `json:"last_pushed"`
}*/

type Results struct {
	Name              string     `json:"name"`
	Namespace         string     `json:"namespace"`
	RepositoryType    string     `json:"repository_type"`
	Status            int        `json:"status"`
	StatusDescription string     `json:"status_description"`
	Description       string     `json:"description"`
	IsPrivate         bool       `json:"is_private"`
	StarCount         int64      `json:"star_count"`
	PullCount         int64      `json:"pull_count"`
	LastUpdated       string     `json:"last_updated"`
	DateRegistered    string     `json:"date_registered"`
	Affiliation       string     `json:"affiliation"`
	MediaTypes        []string   `json:"media_types"`
	ContentTypes      []string   `json:"content_types"`
	Categories        []Category `json:"categories"`
}

type Category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type RepositoryImageResponse struct {
	User              string    `json:"user"`
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	RepositoryType    string    `json:"repository_type"`
	Status            int       `json:"status"`
	StatusDescription string    `json:"status_description"`
	Description       string    `json:"description"`
	IsPrivate         bool      `json:"is_private"`
	IsAutomated       bool      `json:"is_automated"`
	StarCount         int64     `json:"star_count"`
	PullCount         int64     `json:"pull_count"`
	LastUpdated       time.Time `json:"last_updated"`
	DateRegistered    time.Time `json:"date_registered"`
	CollaboratorCount int       `json:"collaborator_count"`
	//Affiliation string
	HubUser         string     `json:"hub_user"`
	HasStarred      bool       `json:"has_starred"`
	FullDescription string     `json:"full_description"`
	Categories      []Category `json:"categories"`
	//Permissions ....
	//MediaTypes ...
	//ContentTypes ...
}
