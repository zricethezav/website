// main.go
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: go run main.go {title}")
		os.Exit(1)
	}

	title := strings.Replace(args[0], " ", "_", -1)

	// Create the new entry with boilerplate and store the date
	createEntry(title)

	// Update TIL.html to include the new entry
	updateTILPage(title)
}

func createEntry(title string) {
	now := time.Now()
	dateString := now.Format("2006-01-02") // Example: 2023-08-14

	header := `<!DOCTYPE html><link rel="icon" href="data:,"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Zach Rice</title>
<head>
      <link href="../style.css" rel="stylesheet" type="text/css" />
</head>
<i><a href="../index.html">Home</a> | <a href="../TIL.html">TIL</a>  | Feed | <a href="../resume.html">Resume</a></i>
<hr>
<!-- Date: ` + dateString + ` -->`

	filename := fmt.Sprintf("TIL/%s.html", title)

	os.WriteFile(filename, []byte(header), 0644)
}

func updateTILPage(title string) {
	// Read the current TIL.html content
	tilBytes, _ := os.ReadFile("TIL.html")
	tilContent := string(tilBytes)

	// Fetch all current links and their corresponding dates
	type LinkInfo struct {
		Date string
		Link string
	}
	var allLinks []LinkInfo

	lines := strings.Split(tilContent, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "<!-- Date:") {
			date := strings.TrimSpace(strings.Split(line, ":")[1])
			date = strings.TrimRight(date, " -->")
			allLinks = append(allLinks, LinkInfo{Date: date, Link: lines[i+1]})
			i++ // Skip the next line as we've already captured it
		}
	}

	// Add the new link
	now := time.Now()
	dateString := now.Format("2006-01-02")
	newLink := fmt.Sprintf("<a href=\"TIL/%s.html\">%s - %s</a><br/>", title, dateString, title)
	allLinks = append(allLinks, LinkInfo{Date: dateString, Link: newLink})

	// Sort the links based on the date
	sort.Slice(allLinks, func(i, j int) bool {
		return allLinks[i].Date < allLinks[j].Date
	})

	// Build the updated TIL.html content
	startIndex := strings.Index(tilContent, "<hr>") + len("<hr>")
	preContent := tilContent[:startIndex]
	postContent := "\n"
	for _, linkInfo := range allLinks {
		postContent += fmt.Sprintf("<!-- Date: %s -->\n%s\n", linkInfo.Date, linkInfo.Link)
	}

	newContent := preContent + postContent
	os.WriteFile("TIL.html", []byte(newContent), 0644)
}
