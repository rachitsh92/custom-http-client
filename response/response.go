package res

import (
	"encoding/json"
	"log"
)

type Response  map[string]interface{}

func FormattedResponse(result map[string]interface{}) []byte {
	formattedResponse, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        log.Fatalf("Error formatting response body: %v", err)
    }
	return formattedResponse
}