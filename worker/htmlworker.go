// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"github.com/Matir/gobuster/logging"
	"github.com/Matir/gobuster/util"
	"github.com/Matir/gobuster/workqueue"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTMLWorker struct {
	// Function to add future work
	adder workqueue.QueueAddFunc
}

func NewHTMLWorker(adder workqueue.QueueAddFunc) *HTMLWorker {
	return &HTMLWorker{adder: adder}
}

func (w *HTMLWorker) Handle(URL *url.URL, body io.Reader) {
	links := w.GetLinks(body)
	foundURLs := make([]*url.URL, 0, len(links))
	for _, l := range links {
		u, err := url.Parse(l)
		if err != nil {
			logging.Logf(logging.LogInfo, "Error parsing URL (%s): %s", l, err.Error())
			continue
		}
		foundURLs = append(foundURLs, URL.ResolveReference(u))
	}
	w.adder(foundURLs...)
}

func (*HTMLWorker) Eligible(resp *http.Response) bool {
	ct := resp.Header.Get("Content-type")
	if strings.ToLower(ct) != "text/html" {
		return false
	}
	return resp.ContentLength > 0 && resp.ContentLength < 1024*1024
}

func (*HTMLWorker) GetLinks(body io.Reader) []string {
	tree, err := html.Parse(body)
	if err != nil {
		logging.Logf(logging.LogInfo, "Unable to parse HTML document: %s", err.Error())
		return nil
	}
	links := make([]string, 0)
	var handleNode func(*html.Node)
	handleNode = func(node *html.Node) {
		if node.Type == html.ElementNode {
			if strings.ToLower(node.Data) == "a" {
				for _, a := range node.Attr {
					if strings.ToLower(a.Key) == "href" {
						links = append(links, a.Val)
						break
					}
				}
			}
		}
		// Handle children
		for n := node.FirstChild; n != nil; n = n.NextSibling {
			handleNode(n)
		}
	}
	handleNode(tree)
	return util.DedupeStrings(links)
}
