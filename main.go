package main

import (
	"encoding/xml"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)
 type loc struct {
 	Value string `xml:"loc"`
 }

 type urlset struct {
 	Urls []loc   `xml:"url"`
 }


func main() {
		urlFlag := flag.String("url","https://gophercises.com/","the url you want to build the site map")
		maxDepth := flag.Int("depth",1000,"the maximum number of links depp to traverse")

		flag.Parse()

	    pages := bfs(*urlFlag,*maxDepth)

		var toXml urlset
	    for _, page := range pages {
	    	toXml.Urls = append(toXml.Urls,loc{page})
		}
		enc := xml.NewEncoder(os.Stdout)
		if err := enc.Encode(toXml); err != nil {
			panic(err)
		}
}
type empty struct{}

func bfs (urlStr string , maxDepth int) []string {
	seen := make(map[string]empty)
	var q map[string]empty
	nq := map[string]empty {
		urlStr:empty{},
	}
	for i:= 0 ;i<= maxDepth ;i++ {
		q,nq = nq,make(map[string]empty)
		if len(q)==0 {
			break
		}
		for urlNode, _ := range q {
			if _,ok := seen[urlNode]; ok {
				continue
			}
			seen[urlNode]= empty{}

			for _,link := range get(urlNode) {
				nq[link] = empty{}
			}
		}
	}
	 ret := make([]string,0,len(seen))

	for urlNode,_ := range seen {
		ret = append(ret,urlNode)
	}
	return ret
 }


func get(urlString string)[]string {

	resp, err := http.Get(urlString)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	requestUrl := resp.Request.URL

	baseUrl := &url.URL{
		Scheme: requestUrl.Scheme,
		Host:   requestUrl.Host,
	}

	base := baseUrl.String()

	hrefs := filter(href(resp.Body, base),hasPrefixWithBase(base))

	return hrefs
}

func href(resp io.Reader, base string) []string {

	link, err := Parse(resp)

	if err != nil {
		panic(err)
	}
	var hrefs []string

	for _, l := range link {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)

		}

	}
	return hrefs
}
func filter (links []string, keepFun func(string)bool) []string {

	var newlinks []string

	for _,link := range links {
		if keepFun(link) == true {
			newlinks = append(newlinks, link)
		}
	}
	return newlinks
}
func hasPrefixWithBase(base string) func(string)bool{
	return func(link string)bool {
			if strings.HasPrefix(link,base){
				return true
			}
			return false
	}
}