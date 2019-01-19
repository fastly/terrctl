package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/jedisct1/dlog"
)

// TerrariumUploadRequest - Query data for a new upload
type TerrariumUploadRequest struct {
	Language       string   `json:"lang"`
	Options        []string `json:"options"`
	EncodedTarFile string   `json:"tar"`
}

// TerrariumUploadResponse - Response data for a new upload
type TerrariumUploadResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TerrariumInstance - A Terrarium instance
type TerrariumInstance struct {
	ID string
}

// TerrariumStatusResponse - Rseponse data for a status query
type TerrariumStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Done    bool   `json:"done"`
}

// URL - Return the URL of an instance
func (instance *TerrariumInstance) URL() (*url.URL, error) {
	return url.Parse("https://" + url.PathEscape(instance.ID) + "." + URLInstance)
}

// Status - Return the current status of an instance
func (instance *TerrariumInstance) Status() (*TerrariumStatusResponse, error) {
	httpClient := http.Client{Timeout: Config().HTTPClientTimeout}
	url, err := url.Parse(URLStatus + "/" + url.PathEscape(instance.ID))
	if err != nil {
		return nil, err
	}
	header := map[string][]string{"User-Agent": {UserAgent}, "Accept": {"application/json"}}
	req := &http.Request{
		Method: "GET",
		URL:    url,
		Header: header,
		Close:  false,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	bin, err := ioutil.ReadAll(io.LimitReader(resp.Body, Config().MaxResponseBodySize))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var statusResponse TerrariumStatusResponse
	err = json.Unmarshal(bin, &statusResponse)
	if err != nil {
		return nil, err
	}
	return &statusResponse, nil
}

// WaitForDeployment - Wait for an instance to be deployed
func (instance *TerrariumInstance) WaitForDeployment() error {
	var previousStatus TerrariumStatusResponse
	tsStart := time.Now()
	for {
		status, err := instance.Status()
		if err != nil {
			return err
		}
		if *status != previousStatus {
			previousStatus = *status
			if !status.Success {
				return errors.New(status.Message)
			}
			dlog.Info(status.Message)
			if status.Done {
				break
			}
		}
		time.Sleep(1 * time.Second)
		now := time.Now()
		if now.Before(tsStart) {
			tsStart = now
		}
		if now.Sub(tsStart) > config.DeployTimeout {
			return errors.New("Timeout")
		}
	}
	return nil
}

// WaitForHealth - Wait for the HTTP service of an instance to be ready
func (instance *TerrariumInstance) WaitForHealth() error {
	tsStart := time.Now()
	for {
		healthy, err := instance.IsHealthy()
		if err != nil {
			return err
		}
		if healthy {
			break
		}
		dlog.Debug("Instance is not reachable over HTTPS yet")
		time.Sleep(1 * time.Second)
		now := time.Now()
		if now.Before(tsStart) {
			tsStart = now
		}
		if now.Sub(tsStart) > config.HealthTimeout {
			return errors.New("Timeout")
		}
	}
	return nil
}

// NewTerrariumInstance - Create a new instance
func NewTerrariumInstance(root string, language string) (*TerrariumInstance, error) {
	dlog.Infof("Preparing upload of directory [%v]", root)
	files, err := FileWalk(root)
	if err != nil {
		return nil, err
	}
	if language == "auto" || language == "" {
		language, err = GuessLanguage(files)
		if err != nil {
			dlog.Fatal(err)
		}
		dlog.Infof("Guessed programming language: %v", language)
	}
	tarFile, err := CreateTarFile(files)
	encodedTarFile := base64.RawStdEncoding.EncodeToString(tarFile)
	uploadRequest := TerrariumUploadRequest{
		Language:       language,
		EncodedTarFile: encodedTarFile,
	}
	uploadRequestJSON, err := json.MarshalIndent(uploadRequest, "", " ")
	if err != nil {
		return nil, err
	}
	body := uploadRequestJSON
	httpClient := http.Client{Timeout: Config().HTTPClientTimeout}
	url, err := url.Parse(URLDeploy)
	if err != nil {
		return nil, err
	}
	header := map[string][]string{"User-Agent": {UserAgent}, "Accept": {"application/json"}, "Content-Type": {"application/json"}}
	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: header,
		Close:  false,
	}
	req.ContentLength = int64(len(body))
	bc := ioutil.NopCloser(bytes.NewReader(body))
	req.Body = bc
	dlog.Notice("Upload in progress...")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	bin, err := ioutil.ReadAll(io.LimitReader(resp.Body, Config().MaxResponseBodySize))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var uploadResponse TerrariumUploadResponse
	err = json.Unmarshal(bin, &uploadResponse)
	if err != nil {
		return nil, err
	}
	dlog.Notice("Upload done, compilation in progress...")
	terrariumInstance := TerrariumInstance{ID: uploadResponse.ID}

	return &terrariumInstance, nil
}

// IsHealthy - Check if the HTTP service is healthy
func (instance *TerrariumInstance) IsHealthy() (bool, error) {
	httpClient := http.Client{Timeout: Config().HTTPClientTimeout}
	url, err := instance.URL()
	if err != nil {
		return false, err
	}
	url.Path = HealthPath
	header := map[string][]string{"User-Agent": {UserAgent}}
	req := &http.Request{
		Method: "GET",
		URL:    url,
		Header: header,
		Close:  false,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 || resp.ContentLength == 0 {
		dlog.Debug("Unexpected response from health check endpoint")
		return false, nil
	}
	return true, nil
}
