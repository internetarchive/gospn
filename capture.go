package spn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Capture execute a capture via https://web.archive.org/save
// and return the response. Options for the capture can be specified
// when calling the method
func (c Connector) Capture(URL string, options CaptureOptions) (captureResponse CaptureResponse, err error) {
	c.GetAvailableCaptureSlot()

	// Build request
	urlValues := options.Encode()
	urlValues.Set("url", URL)
	req, err := http.NewRequest("POST", "https://web.archive.org/save", strings.NewReader(urlValues.Encode()))
	if err != nil {
		return captureResponse, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	logger.Debug("Executing capture request", "payload", urlValues.Encode())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return captureResponse, err
	}

	body, _ := io.ReadAll(resp.Body)
	logger.Debug("Capture response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode != 200 {
		return captureResponse, fmt.Errorf("SPN Capture failed with status code %d, response: %s", resp.StatusCode, body)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		return captureResponse, fmt.Errorf("SPN response is not JSON, Content-Type: %s, response: %s", resp.Header.Get("Content-Type"), body)
	}

	if err := json.Unmarshal(body, &captureResponse); err != nil {
		return captureResponse, fmt.Errorf("Failed to unmarshal JSON: %s", err)
	}

	return captureResponse, nil
}
