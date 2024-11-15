package spn

import (
	"encoding/json"
	"net/http"
	"time"
)

// UserStatus represent the data returned by the /save/status/user endpoint
type UserStatus struct {
	DailyCaptures      int `json:"daily_captures"`
	DailyCapturesLimit int `json:"daily_captures_limit"`
	Available          int `json:"available"`
	Processing         int `json:"processing"`
}

func (to *UserStatus) Update(from UserStatus) {
	to.DailyCaptures = from.DailyCaptures
	to.DailyCapturesLimit = from.DailyCapturesLimit
	to.Available = from.Available
	to.Processing = from.Processing
}

// CaptureStatus represent the date returned by the /save/status/{job_id} endpoint
type CaptureStatus struct {
	Timestamp   string   `json:"timestamp"`
	DurationSec float64  `json:"duration_sec"`
	OriginalURL string   `json:"original_url"`
	Status      string   `json:"status"`
	StatusExt   string   `json:"status_ext"`
	JobID       string   `json:"job_id"`
	Outlinks    []string `json:"outlinks"`
	Resources   []string `json:"resources"`
	Exception   string   `json:"exception"`
	Message     string   `json:"message"`
}

// GetCaptureStatus retrieve the informations about a SPN job
func (c Connector) GetCaptureStatus(jobID string) (captureStatus CaptureStatus, err error) {
	// Build request
	req, err := http.NewRequest("GET", "https://web.archive.org/save/status/"+jobID, nil)
	if err != nil {
		return captureStatus, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return captureStatus, err
	}

	json.NewDecoder(resp.Body).Decode(&captureStatus)

	if captureStatus.Outlinks == nil {
		captureStatus.Outlinks = []string{""}
	}

	if captureStatus.Resources == nil {
		captureStatus.Resources = []string{""}
	}

	return captureStatus, nil
}

// GetUserStatus retrieve the user status for a given SPN account
func (c Connector) GetUserStatus() (userStatus UserStatus, err error) {

	// Build request
	req, err := http.NewRequest("GET", "https://web.archive.org/save/status/user", nil)
	if err != nil {
		return userStatus, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return userStatus, err
	}

	json.NewDecoder(resp.Body).Decode(&userStatus)

	return userStatus, nil
}

// Refresh the cached user status in the background
func (c *Connector) cachedUserStatusFetcher() {
	logger.Debug("Starting cachedUserStatusFetcher")
	defer logger.Debug("cachedUserStatusFetcher exited")
	lastFetch := time.Unix(0, 0)
	for {
		select {
		case <-c.cachedStatusFetcherIntr:
			return
		default:
		}

		if time.Since(lastFetch) < time.Second*10 && c.cachedStatus.Available > 2 {
			time.Sleep(time.Second)
			continue
		}

		// Make sure we don't fetch the status too often (< 2s)
		wait := time.Second*2 - time.Since(lastFetch)
		if wait > 0 {
			logger.Info("Waiting before fetching user status", "wait", wait)
			time.Sleep(wait)
		}

		logger.Debug("Fetching user status")

		lastFetch = time.Now()
		userStatus, err := c.GetUserStatus()
		if err != nil {
			logger.Error("Failed to fetch user status", "error", err)
			continue
		}

		logger.Debug("User status fetched", "status", userStatus)
		c.cachedStatus.Update(userStatus)
		logger.Debug("cachedStatus updated", "cachedStatus", c.cachedStatus)
	}
}

// Wait until a capture slot is available
func (c Connector) GetAvailableCaptureSlot() (err error) {
	for {
		if c.cachedStatus.Available > 0 {
			c.cachedStatus.Available--
			c.cachedStatus.Processing++
			logger.Debug("AwaitAvailableSlot return", "cachedStatus", c.cachedStatus)
			return nil
		}

		logger.Debug("AwaitAvailableSlot waiting")
		time.Sleep(time.Second)
	}
}
