// linkvalidator package provides some helper methods for examination of the links
package linkvalidator

import (
	"net/http"
	"net/url"
	"strings"
)

type Validator struct {
	BaseURL *url.URL
}

// LinkInfo contains information about links
type LinkInfo struct {
	IsExternal bool
	IsBroken   bool
	URL        string
	FullPath   string
}

// getFullPath converts an internal link to an external one.
func (v *Validator) getFullPath(link string) string {
	link = strings.ReplaceAll(link, " ", "")
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		// it is alread a full path
		return link
	}

	if linkURL, err := url.Parse(link); err == nil {
		result := v.BaseURL.ResolveReference(linkURL).String()
		result = strings.ReplaceAll(result, " ", "") // just trimming
		return result
	}

	return ""
}

// isBrokenLink checks if a link is accessible or not
func (v *Validator) isBrokenLink(url string) bool {
	fullPath := v.getFullPath(url)
	resp, err := http.Get(fullPath)
	if err != nil {
		// fmt.Println(fullPath, err.Error())
		return true
	}

	// usually, an status code of 400 or above means a broken link
	// this is not accurate though.
	if resp.StatusCode >= 400 {
		// fmt.Println(fullPath, resp.StatusCode)
		return true
	}

	return false
}

// isExternalLink checks if a link is external or not
func (v *Validator) isExternalLink(link string) bool {
	if linkURL, err := url.Parse(link); err == nil {
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			// it is alread a full path
			return linkURL.Hostname() != v.BaseURL.Hostname()
		}
		// link is not a full path, first we compute the full path.
		fp := v.BaseURL.ResolveReference(linkURL)
		return fp.Hostname() != v.BaseURL.Hostname()
	}

	// an error happend, we just assume it is a local link
	// TODO: fix this
	return false
}

func (v *Validator) GetLinkInfo(url string) LinkInfo {
	// LinkInfo is a relativly small struct. therefore we return the whole datastuct
	// instead of a pointer to it. There will be some copy and construction around
	// but the overhead is negligable.

	result := LinkInfo{
		IsExternal: v.isExternalLink(url),
		IsBroken:   v.isBrokenLink(url),
		URL:        url,
		FullPath:   v.getFullPath(url),
	}
	return result
}
