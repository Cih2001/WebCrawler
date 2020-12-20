package controller

import (
	"Cih2001/WebCrawler/linkvalidator"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo"
	"golang.org/x/net/html"
)

// ReportResponse is a structure that contains all the necessary data
// to load report.html template.
type ReportResponse struct {
	Address       string
	Title         string
	Version       string
	Headings      [6]int
	InternalLinks int
	ExternalLinks int
	TotalLinks    int
	BrokenLinks   int
	LoginForm     bool
	Links         []linkvalidator.LinkInfo
}

// defaultLinkValidator will be initialized with the request url.
var defaultLinkValidator *linkvalidator.Validator

// IsLoginForm checks if a given form is for logging user in or not.
func IsLoginForm(node *html.Node) bool {
	// How can we check if a given form is meant for logging in? there is no fix method.
	// Login forms can be designed in various ways. we use some heuristics. This implies
	// no guarantee for the correctness, but it can detect typical situations.

	// we check into the attributes, texts, data, and names of any HTML element in forms.
	// if it contains some keywords, we assume that it is a login form.
	// keywords should be in lower case only
	keywords := []string{"log in", "login", "password"}

	// check attributes of the current node.
	for _, attr := range node.Attr {
		if attr.Key == "id" || attr.Key == "name" {
			for _, name := range keywords {
				if strings.Contains(strings.ToLower(attr.Val), name) {
					return true
				}
			}
		}
	}

	// check node data.
	for _, name := range keywords {
		if strings.Contains(strings.ToLower(node.Data), name) {
			return true
		}
	}

	// check other nodes.
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if IsLoginForm(c) {
			// if result is already true, we do not need to look into other nodes.
			return true
		}
	}

	// no indication of login in the current node or any other childern.
	return false
}

// examineLinks gets a slice of links, examines them and
// sets links related fields in Report. fields such as BrokenLinks,
// Links, InternalLinks and etc.
func (rr *ReportResponse) examineLinks(links []string) {
	var mutex sync.Mutex
	var wg sync.WaitGroup

	brokenLinksCount := 0
	externalLinksCount := 0

	getConcurrentLinkInfo := func(wg *sync.WaitGroup, link string) {
		// gets a link info and updates the report.
		defer wg.Done()

		linkInfo := defaultLinkValidator.GetLinkInfo(link)

		// accessing global variables within a thread can cause race condition.
		// these accessed should only be made with a use of mutex.
		mutex.Lock()
		defer mutex.Unlock()

		if linkInfo.IsBroken {
			brokenLinksCount++
		}
		if linkInfo.IsExternal {
			externalLinksCount++
		}

		rr.Links = append(rr.Links, linkInfo)
	}

	for _, l := range links {
		wg.Add(1)
		go getConcurrentLinkInfo(&wg, l) // isn't it of divine intellect?
	}

	wg.Wait() // for all links to be resolved.

	rr.BrokenLinks = brokenLinksCount
	rr.ExternalLinks = externalLinksCount
	rr.TotalLinks = len(rr.Links)
	rr.InternalLinks = rr.TotalLinks - rr.ExternalLinks
}

// prepareReport gets a response to a http request and makes a report based on that.
func (rr *ReportResponse) prepareReport(body io.ReadCloser) error {
	links := []string{}
	hasLoginForm := false

	doc, err := html.Parse(body)
	if err != nil {
		return err
	}

	var parseHTMLTree func(*html.Node)
	parseHTMLTree = func(n *html.Node) {
		// extract page title.
		if n.Type == html.ElementNode && n.Data == "title" {
			// there might be multiple html tags named title in a page
			// therefore we check the parent of title tag as well.
			if n.Parent.Data == "head" {
				rr.Title = n.FirstChild.Data
			}
		}

		// extract links.
		for _, attr := range n.Attr {
			// links are usually specified in href attribute.
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}

		// extract html version.
		if n.Type == html.DoctypeNode {
			if n.Data == "html" {
				// we found <!DOCTYPE html>. it should be html 5.
				rr.Version = "5"
			}
		}

		// count headings.
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1":
				rr.Headings[0]++
			case "h2":
				rr.Headings[1]++
			case "h3":
				rr.Headings[2]++
			case "h4":
				rr.Headings[3]++
			case "h5":
				rr.Headings[4]++
			case "h6":
				rr.Headings[5]++
			}
		}

		// check for login form
		if n.Type == html.ElementNode && n.Data == "form" {
			if !hasLoginForm {
				hasLoginForm = IsLoginForm(n)
			}

			// we probably do not need to parse the children of a form.
			// IsLoginForm does that for us. so it would make more sence to return here.
			// however, as it does not have a huge impact on performance, (forms are literally
			// small in size) we do not return here just to be safe. (for example catching
			// any potential link in the forms (relative links might be probable).
		}

		// parse other nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseHTMLTree(c)
		}
	}
	parseHTMLTree(doc)

	rr.examineLinks(links)
	rr.LoginForm = hasLoginForm
	return nil
}

// ReportHandler handles requests to /report?address=:addr
func ReportHandler(c echo.Context) error {
	// check if address in not empty
	webAddress := c.QueryParam("address")
	if webAddress == "" {
		return c.HTML(http.StatusBadRequest, "Bad request")
	}

	// download the content of the website.
	resp, err := http.Get(webAddress)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Error: "+err.Error())
	}

	responseData := ReportResponse{
		Address:  resp.Request.URL.String(),
		Version:  "4 or earlier", // Default html version
		Headings: [6]int{0, 0, 0, 0, 0, 0},
	}

	// initialize default link validator.
	defaultLinkValidator = &linkvalidator.Validator{
		BaseURL: resp.Request.URL,
	}

	// prepare report.
	err = responseData.prepareReport(resp.Body)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Error: "+err.Error())
	}

	// render report template
	return c.Render(http.StatusOK, "report.html", responseData)
}
