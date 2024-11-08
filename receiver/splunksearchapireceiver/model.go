package splunksearchapireceiver

// response structs for Splunk API calls
type CreateJobResponse struct {
	SID string `xml:"sid"`
}

type JobStatusResponse struct {
	Content struct {
		Type string `xml:"type,attr"`
		Dict Dict   `xml:"dict"`
	} `xml:"content"`
}

// Dict struct to represent <s:dict> elements
type Dict struct {
	Keys []Key `xml:"key"`
}

// Key struct to represent <s:key> elements
type Key struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
	Dict  *Dict  `xml:"dict,omitempty"`
	List  *List  `xml:"list,omitempty"`
}

type List struct {
	Items []struct {
		Value string `xml:",chardata"`
	} `xml:"item"`
}

type SearchResults struct {
	Results []struct {
		Raw  string `json:"_raw"`
		Time string `json:"_time"`
	} `json:"results"`
}
