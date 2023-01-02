package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	var dir string
	var debug bool
	var del bool

	flag.StringVar(&dir, "d", "", "directory to convert")
	flag.BoolVar(&del, "del", false, "delete original files")
	flag.BoolVar(&debug, "debug", false, "debug mode - outputs to stdout")
	flag.Parse()

	if dir == "" {
		flag.Usage()
		os.Exit(1)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	opt := &md.Options{
		HorizontalRule: "---",
		CodeBlockStyle: "fenced",
	}
	c := md.NewConverter("", true, opt)
	c.Use(plugin.GitHubFlavored())
	c.Use(plugin.EXPERIMENTALMoveFrontMatter())
	// c.Keep("#comment")

	// I could use the Keep() method, but I want to ensure the `more` comments are in the form Hugo expects
	commentRule := md.Rule{
		Filter: []string{"#comment"},
		Replacement: func(content string, selector *goquery.Selection, opt *md.Options) *string {
			c := selector.Nodes[0].Data
			if c == " more " {
				c = strings.Trim(c, " ")
			}
			return md.String("<!--" + c + "-->")
		},
	}
	c.AddRules(commentRule)

	for _, file := range files {
		html, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		markdown, err := c.ConvertString(string(html))
		if err != nil {
			log.Fatal(err)
		}
		if debug {
			fmt.Println(markdown)
		} else {
			err = os.WriteFile(strings.TrimSuffix(file, ".html")+".md", []byte(markdown), 0o644)
			if err != nil {
				log.Fatal(err)
			}
			if del {
				err = os.Remove(file)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
