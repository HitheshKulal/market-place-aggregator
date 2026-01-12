package dto

type ParsedFileResponse struct {
	DiscoveredColumns []string                 `json:"discoveredColumns"`
	SampleRows        []map[string]interface{} `json:"sampleRows"`
	RowCount          int                      `json:"rowCount"`
	FileName          string                   `json:"fileName"`
	FileType          string                   `json:"fileType"`
}
