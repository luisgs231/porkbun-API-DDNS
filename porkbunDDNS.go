package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"gopkg.in/yaml.v3"
)

type Record struct {
	Domain     string `yaml:"domain"`
	Record     string `yaml:"record"`
	RecordType string `yaml:"recordType"`
	TTL        string `yaml:"ttl"`
}

type Config struct {
	APIKey         string   `yaml:"apikey"`
	SecretAPIKey   string   `yaml:"secretapikey"`
	IPSource       string   `yaml:"ipSource"`
	ReportEndpoint string   `yaml:"reportEndpoint"`
	Records        []Record `yaml:"records"`
}

func main() {
	fmt.Println("Reading configuration file...")
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("Error: cannot read config.yaml:", err)
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		fmt.Println("Error: cannot parse config.yaml:", err)
		os.Exit(1)
	}

	fmt.Println("Fetching current external IP...")
	resp, err := http.Get(cfg.IPSource)
	if err != nil {
		fmt.Println("Error: cannot fetch IP:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: cannot read IP response:", err)
		os.Exit(1)
	}
	currentIP := strings.TrimSpace(string(ipBytes))
	fmt.Println("Current external IP is:", currentIP)

	for _, r := range cfg.Records {
		fmt.Printf("\nProcessing record: %s.%s (%s)\n", r.Record, r.Domain, r.RecordType)

		// Check current DNS record
		checkURL := "https://api.porkbun.com/api/json/v3/dns/retrieveByNameType/" + r.Domain + "/" + r.RecordType + "/" + r.Record
		checkPayload := fmt.Sprintf(`{"apikey":"%s","secretapikey":"%s"}`, cfg.APIKey, cfg.SecretAPIKey)

		checkResp, err := http.Post(checkURL, "application/json", strings.NewReader(checkPayload))
		if err != nil {
			fmt.Println("Error: cannot check current record:", err)
			continue
		}
		defer checkResp.Body.Close()

		var checkData struct {
			Records []struct {
				Content string `yaml:"content" json:"content"`
			} `yaml:"records" json:"records"`
		}
		if err := yaml.NewDecoder(checkResp.Body).Decode(&checkData); err != nil {
			fmt.Println("Error: cannot parse record check:", err)
			continue
		}

		if len(checkData.Records) > 0 && strings.TrimSpace(checkData.Records[0].Content) == currentIP {
			fmt.Println("DNS record already up-to-date, no changes made.")
			reportStatus(cfg.ReportEndpoint, currentIP, "up-to-date")
			continue
		}

		// Update DNS record
		fmt.Println("Updating DNS record...")
		data := fmt.Sprintf(
			`{"apikey":"%s","secretapikey":"%s","content":"%s","ttl":"%s"}`,
			cfg.APIKey, cfg.SecretAPIKey, currentIP, r.TTL,
		)
		updateURL := "https://api.porkbun.com/api/json/v3/dns/editByNameType/" + r.Domain + "/" + r.RecordType + "/" + r.Record

		resp2, err := http.Post(updateURL, "application/json", strings.NewReader(data))
		if err != nil {
			fmt.Println("Error: cannot update record:", err)
			reportStatus(cfg.ReportEndpoint, currentIP, "update-failed")
			continue
		}
		defer resp2.Body.Close()

		if resp2.StatusCode >= 200 && resp2.StatusCode < 300 {
			fmt.Println("Record updated successfully.")
			reportStatus(cfg.ReportEndpoint, currentIP, "updated")
		} else {
			fmt.Println("Error: update failed with status", resp2.StatusCode)
			reportStatus(cfg.ReportEndpoint, currentIP, "update-failed")
		}
	}

	fmt.Println("\nAll records processed.")
}

// reportStatus calls push endpoint with the current IP
func reportStatus(endpoint, ip, msg string) {
	u, err := url.Parse(endpoint)
	if err != nil {
		fmt.Println("Error parsing report endpoint:", err)
		return
	}

	q := u.Query()
	q.Set("msg", fmt.Sprintf("%s", ip))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Error reporting status:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Reported status:", msg)
}
