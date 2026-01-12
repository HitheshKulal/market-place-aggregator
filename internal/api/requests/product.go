package requests

type UploadProductsResponse struct {
	FileName          string                   `json:"fileName"`
	TotalRows         int                      `json:"totalRows"`
	SuccessCount      int                      `json:"successCount"`
	FailedCount       int                      `json:"failedCount"`
	DiscoveredColumns []string                 `json:"discoveredColumns"`
	SampleProducts    []map[string]interface{} `json:"sampleProducts"` // First 5 created products
	Errors            []ProductError           `json:"errors,omitempty"`
}

type ProductError struct {
	Row     int    `json:"row"`
	SKU     string `json:"sku,omitempty"`
	Message string `json:"message"`
}
