package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ComicData holds the extracted information for the JSON response.
type ComicData struct {
	Img      string `json:"img"`                // URL of the comic image for TRMNL markup
	ImageURL string `json:"imageUrl,omitempty"` // Backward-compatible field for existing clients
	Title    string `json:"title"`              // Caption/title of the comic
}

const targetURL = "https://www.thefarside.com/"

// fetchComicData fetches and parses the Far Side website to extract comic info.
func fetchComicData() (*ComicData, error) {
	// Make HTTP GET request with a browser-like User-Agent to avoid 403
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", targetURL, err)
	}
	defer res.Body.Close() // Ensure body is closed

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request failed with status code %d",
			res.StatusCode,
		)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find the first image from the daily comic container.
	img := doc.Find(".tfs-comic_image img").First()
	if img.Length() == 0 {
		// Fallback for legacy structure.
		img = doc.Find(".card-body img").First()
	}
	if img.Length() == 0 {
		return nil, fmt.Errorf("could not find comic image in '.tfs-comic_image img' or '.card-body img'")
	}

	// Prefer lazy-load URL in data-src, then fallback to src.
	imgSrc, exists := img.Attr("data-src")
	if !exists || imgSrc == "" {
		imgSrc, exists = img.Attr("src")
		if !exists || imgSrc == "" {
			return nil, fmt.Errorf("could not find img with 'data-src' or 'src' attribute")
		}
	}

	// Get caption text if present, then fallback to alt text.
	title := strings.TrimSpace(doc.Find(".tfs-comic_image .figure-caption, .card-body .figure-caption, .figure-caption").First().Text())
	if title == "" {
		if alt, ok := img.Attr("alt"); ok {
			title = strings.TrimSpace(alt)
		}
	}
	if title == "" {
		title = "The Far Side Daily Dose"
	}

	data := &ComicData{
		Img:      imgSrc,
		ImageURL: imgSrc,
		Title:    title,
	}

	return data, nil
}

// farsideHandler handles incoming web requests.
func farsideHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for %s from %s", r.URL.Path, r.RemoteAddr)

	// Fetch the latest comic data
	comicData, err := fetchComicData()
	if err != nil {
		log.Printf("Error fetching comic data: %v", err)
		http.Error(
			w,
			"Failed to retrieve Far Side comic data.",
			http.StatusInternalServerError,
		)
		return
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the data to JSON and write to response
	err = json.NewEncoder(w).Encode(comicData)
	if err != nil {
		// This error is less likely but possible (e.g., write error)
		log.Printf("Error encoding JSON response: %v", err)
		// Don't write another http.Error if headers are already sent
		return
	}
	log.Println("Successfully served comic data.")
}

func main() {
	// Register the handler function for the root path "/"
	http.HandleFunc("/", farsideHandler)

	port := "8123"
	log.Printf("Starting Far Side JSON server on port %s...", port)
	log.Printf("Access it at: http://localhost:%s", port)

	// Start the HTTP server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
