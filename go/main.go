// Parses HTML entirely
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sys/unix"
)

func getCPUTime() time.Duration {
	var ru unix.Rusage
	unix.Getrusage(unix.RUSAGE_SELF, &ru)
	return time.Duration(ru.Utime.Sec)*time.Second + time.Duration(ru.Utime.Usec)*time.Microsecond
}

func fetchAndProcessHTML(url string) (time.Duration, error) {
	var body []byte
	var err error

	if strings.HasPrefix(url, "https://") {
		// Fetch the HTML
		resp, err := http.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("HTTP error: %s", resp.Status)
		}
	} else {
		// Read the file content
		body, err = os.ReadFile(url)
		if err != nil {
			return 0, err
		}
	}

	startCPUTime := getCPUTime()

	// Read the entire body into memory
	reader := strings.NewReader(string(body))
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return 0, err
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return 0, err
	}

	// Modify the title
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		s.SetText("Modified Title")
	})

	// Serialize back to HTML
	_, err = doc.Html()
	if err != nil {
		return 0, err
	}

	endCPUTime := getCPUTime()
	return endCPUTime - startCPUTime, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <URL>")
	}

	url := os.Args[1]

	cpuTime, err := fetchAndProcessHTML(url)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Output result
	fmt.Printf("CPU Time: %v\n", cpuTime)
	fmt.Printf("Num Goroutines: %d\n", runtime.NumGoroutine())
}
