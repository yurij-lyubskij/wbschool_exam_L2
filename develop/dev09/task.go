package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func findLinks(n *html.Node, s *strings.Builder) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Fprintln(s, a.Val)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findLinks(c, s)
	}
}

var maxLevel = 2
var dir string

func init() {
	dir, _ = os.Getwd()
}

func recursiveDownload(link string, level int, name string) {
	baseLink, err := url.Parse(link)
	if err != nil {
		log.Fatal("parse ", link, err.Error())
	}
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal("get ", link, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(name)
	file, err := os.Create(name)
	if err != nil {
		log.Fatal("create, name = ", name, "err= ", err.Error(), " link =", link)
	}
	defer file.Close()

	if level == maxLevel {
		file.Write(body)
		return
	}

	reader := bytes.NewReader(body)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal("htmlparse", err)
	}
	s := strings.Builder{}
	findLinks(doc, &s)
	links := strings.Split(s.String(), "\n")
	for _, oneLink := range links {
		if oneLink == "/" {
			continue
		}
		parsedLink, err := url.Parse(oneLink)
		if err != nil {
			log.Fatal("parseOne ", oneLink, err.Error())
		}
		parts := strings.Split(oneLink, "/")
		nextName := uuid.NewString() + parts[len(parts)-1]
		nextURL := baseLink.ResolveReference(parsedLink)
		replaceLink := "file://" + dir + nextName
		fmt.Println(nextURL.String())
		body = bytes.ReplaceAll(body, []byte(oneLink), []byte(replaceLink))
		recursiveDownload(nextURL.String(), level+1, nextName)
	}
	file.Write(body)
}

func main() {
	link := "https://yandex.ru/"
	fname := "index.html"
	recursiveDownload(link, 1, fname)
}
