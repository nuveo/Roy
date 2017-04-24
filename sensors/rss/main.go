package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/lytics/multibayes"
)

/*
https://blog.golang.org/feed.atom
https://news.ycombinator.com/rss
https://www.reddit.com/r/golang.rss
https://www.goinggo.net/index.xml
*/

var classifier *multibayes.Classifier

func init() {
	classifier = multibayes.NewClassifier()
	classifier.MinClassSize = 0
}

func getFeed(uri string) *rss.Feed {

	f := func(url string) (resp *http.Response, err error) {
		client := http.DefaultClient
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", "news bot")

		return client.Do(req)
	}

	for {
		feed, err := rss.FetchByFunc(f, uri)
		if err == nil {
			return feed
		}
		log.Printf("Error fetching %s: %s", uri, err)

		<-time.After(time.Duration(50 * time.Second))
	}
}

func getClass(p map[string]float64) (c1, c2 string, v1, v2 float64) {
	r := 0.0

	for k, v := range p {
		if v > r {
			c1 = k
			v1 = v
			r = v
		}
	}

	r = 0.0
	for k, v := range p {
		if v > r {
			if k == c1 {
				continue
			}
			c2 = k
			v2 = v
			r = v
		}
	}

	return
}

func main() {
	fmt.Println("stanting news bot")

	documents := []struct {
		Text    string
		Classes []string
	}{
		{
			Text:    "A Mini Guide to Google’s Golang and Why It’s Perfect for DevOps • r/devops",
			Classes: []string{"DevOps"},
		},
		{
			Text:    "Deep Learning from Scratch in Go Equations Are Graphs",
			Classes: []string{"Machine Learning", "AI"},
		},
		{
			Text:    "YAML to GO: Convert YAML to Struct",
			Classes: []string{"Parser"},
		},
	}

	// train the classifier
	for _, document := range documents {
		classifier.Add(document.Text, document.Classes)
	}

	for {

		feed := getFeed("https://www.reddit.com/r/golang/new.rss")

		for i := len(feed.Items) - 1; i > -1; i-- {
			fmt.Println(i, "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

			title := feed.Items[i].Title
			fmt.Println(title)
			fmt.Println(feed.Items[i].Link)

			probs := classifier.Posterior(title)

			//for k, v := range probs {
			//	fmt.Println(k, v)
			//}
			c1, c2, v1, v2 := getClass(probs)

			fmt.Println(c1, v1, c2, v2)

		}
		<-time.After(time.Duration(10 * time.Minute))

	}

}
