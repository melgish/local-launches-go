package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Launch struct {
	Date        string
	Description string
	Mission     string
	Site        string
	Time        string
	LastUpdated string
}

// Site HTML looks like this:
//
// <div class="datename">
//    <div class="launchdate">{Date}</div>
//    <div class="mission">{Mission}</div>
// </div>
// <div class="missiondata">
//    <span class="strong">Launch time:</span>
//    {Time}
//    <br>
//    <span  class="strong">Launch site:</span>
//    {Site}
// </div>
// <div class="missdescrip">
//   <p>{Description}</p>
//   <p style="">Updated {Last Updated}</p>
// </div>
//

const (
	DataXPath        = "following-sibling::div[contains(@class, 'missiondata')]"
	DateXPath        = "*[contains(@class, 'launchdate')]"
	DescriptionXPath = "following-sibling::div[contains(@class,'missdescrip')]/p[1]"
	LastUpdatedXPath = "following-sibling::div[contains(@class,'missdescrip')]/p[2]"
	LaunchXPath      = "//div[contains(@class, 'datename')]"
	MissionXPath     = "*[contains(@class, 'mission')]"
	SiteURL          = "https://spaceflightnow.com/launch-schedule/"
	SiteXPath        = "span[contains(text(), 'Launch site:')]/following-sibling::text()[1]"
	TimeXPath        = "span[contains(text(), 'Launch time:')]/following-sibling::text()[1]"
)

var (
	whiteSpace    = regexp.MustCompile(`\s+`)
	siteFilter    = regexp.MustCompile(`(?i)(canaveral|kennedy|patrick)`)
	launchesCache []*Launch
	lastRefresh   time.Time
)

// findOneAndTrim finds a node by XPath and return its trimmed inner text
//
// Returns inner text of the node if found, otherwise returns an empty string.
func findOneAndTrim(node *html.Node, xpath string) string {
	child := htmlquery.FindOne(node, xpath)
	if child != nil {
		text := htmlquery.InnerText(child)
		text = whiteSpace.ReplaceAllString(text, " ") // Replace multiple spaces with a single space
		return strings.TrimSpace(text)
	}
	return ""
}

// newLaunch creates a new Launch struct from the provided HTML node.
//
// Returns a new launch instance.
func newLaunch(node *html.Node) *Launch {
	// Create a new Launch struct from the HTML node
	launch := Launch{Site: "", Time: ""}

	// Extract the date, mission, site, time and description
	// from the node using XPath queries
	launch.Date = findOneAndTrim(node, DateXPath)
	launch.Mission = findOneAndTrim(node, MissionXPath)
	launch.Description = findOneAndTrim(node, DescriptionXPath)
	launch.LastUpdated = findOneAndTrim(node, LastUpdatedXPath)

	data := htmlquery.FindOne(node, DataXPath)
	if data != nil {
		launch.Site = findOneAndTrim(data, SiteXPath)
		launch.Time = findOneAndTrim(data, TimeXPath)
	}

	return &launch
}

// getLaunches retrieves a list of launches from the SpaceFlightNow HTML page.
//
// Returns a slice of pointers to Launch structs and an error if any occurs.
func getLaunches() ([]*Launch, error) {
	// Load the SpaceFlightNow HTML page
	// doc, err := htmlquery.LoadDoc("./www/sample.html")
	doc, err := htmlquery.LoadURL(SiteURL)
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.QueryAll(doc, LaunchXPath)
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d launches\n", len(nodes))
	var launches []*Launch
	for _, node := range nodes {
		// Filter out launches where the Site does not match the regex
		launch := newLaunch(node)
		if siteFilter.MatchString(launch.Site) {
			launches = append(launches, launch)
		}
	}

	return launches, nil
}

// updateLaunchesPeriodically calls getLaunches every 4 hours and updates the cache
func updateLaunchesPeriodically(refreshInterval time.Duration, dataReady chan struct{}) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	once := true

	for {
		log.Println("Fetching launches...")
		lastRefresh = time.Now()
		launches, err := getLaunches()
		if err != nil {
			log.Printf("Failed to fetch launches: %v. Keeping previous data.", err)
		} else {
			// Update the cache with the new data.
			launchesCache = launches
			log.Println("Launches updated successfully.")
			if once {
				// Signal that the initial data fetch is complete
				close(dataReady)
				once = false
			}
		}

		// Wait for the next tick
		<-ticker.C
	}
}

func timeUntilNextRefresh(refreshInterval time.Duration) string {
	remaining := time.Until(lastRefresh.Add(refreshInterval))
	// Format the remaining time as "HH:MM:SS"
	if remaining < 0 {
		return "00:00:00"
	}
	// Calculate hours, minutes, and seconds
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	seconds := int(remaining.Seconds()) % 60

	// Format as HH:MM:SS
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
