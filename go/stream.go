// Parses HTML in streaming
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

	"golang.org/x/net/html"
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

	fmt.Println("Body size:", len(body))

	startCPUTime := getCPUTime()

	// Stream parse HTML
	var modifiedHTML bytes.Buffer
	reader := strings.NewReader(string(body))
	tokenizer := html.NewTokenizer(reader)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			endCPUTime := getCPUTime()
			return endCPUTime - startCPUTime, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "title" {
				modifiedHTML.WriteString("<title>Modified Title</title>")
				// Skip the next token if it's the original title text
				tt = tokenizer.Next()
				if tt == html.TextToken {
					continue
				}
			} else {
				modifiedHTML.WriteString(token.String())
			}
		default:
			modifiedHTML.WriteString(tokenizer.Token().String())
		}
	}
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
