package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
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

func recursiveDownload(prev []byte, links []string, level int) {

}

func main() {
	resp, err := http.Get("https://yandex.ru")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	name := "index.html"
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	file.Write(body)
	reader := bytes.NewReader(body)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node, *strings.Builder)
	s := strings.Builder{}
	f = findLinks
	f(doc, &s)
	links := strings.Split(s.String(), "\n")
	for i := range links {
		fmt.Println(links[i])
	}

}
