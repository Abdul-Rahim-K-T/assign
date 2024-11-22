package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Experience struct {
	Company  string `json:"company"`
	Position string `json:"position"`
	Duration string `json:"duration"`
}
type Education struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	Year        string `json:"year"`
}

type ParsedResume struct {
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Phone      string       `json:"phone"`
	Skills     []string     `json:"skills"`
	Education  []Education  `json:"education"`
	Experience []Experience `json:"experience"`
}

func ParseResumeWithAPILayer(resumeData []byte) (*ParsedResume, error) {
	apiURL := "https://api.apilayer.com/resume_parser/upload" // Replace with the actual APILayer endpoint
	apiKey := os.Getenv("RESUME_API_KEY")
	log.Println(apiKey) // Replace with your actual API key

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(resumeData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream") // Adjust the content type based on your resume file format
	req.Header.Set("apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("API Error: %v - Response :%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("received non-200 response code: %v", resp.StatusCode)
	}

	log.Println("TRACK TEST")
	var parsedResume ParsedResume
	if err := json.NewDecoder(resp.Body).Decode(&parsedResume); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &parsedResume, nil
}

// ParseResume uploads a resume file to an external API and parses the response to extract relevant details

// func ParseResume(file []byte) (map[string]string, error) {
// 	// Prepare the POST request
// 	req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parseser/upload", bytes.NewReader(file))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Set API key for authentication
// 	req.Header.Set("Authorization", "Bearer YOUR_API_KEY")

// 	// Iitialize HTTP client
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Read the response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if the response is successful (status code 200)
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errors.New("failed to parse resume: " + string(body))
// 	}

// 	// Parse the response body to extract relevant information
// 	// You can expand this part to extract different fields based on the API response
// 	parsedData := map[string]string{
// 		"skills":    string(body), // Customize this to extract specific data
// 		"education": string(body), // Same here, depending on the response format
// 	}

// 	return parsedData, nil

// }
// https://chatgpt.com/c/673d8290-1db0-800f-9db5-c02863fd8d38
// https://chatgpt.com/c/673d8290-1db0-800f-9db5-c02863fd8d38
