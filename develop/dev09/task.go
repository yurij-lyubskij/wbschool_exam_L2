package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	resp, err := http.Get("http://yandex.ru")
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
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
