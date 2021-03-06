package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/plin2k/goxhparser"
)

func main() {
	ParseByXMLFileInGoroutine()
	ParseByXMLFile()

}

func ParseByXMLFileInGoroutine() {
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex

	parser, err := goxhparser.NewParser("./example/golang_useful.xml")
	if err != nil {
		log.Fatalln(err)
	}

	err = parser.XMLToStruct()
	if err != nil {
		log.Fatalln(err)
	}
	for _, value := range parser.Service.Sources {
		wg.Add(1)
		go func(source goxhparser.Source) {
			content, err := parser.Parse(source)
			if err != nil {
				log.Fatalln(err)
			}

			mu.Lock()
			parser.Content = append(parser.Content, content...)
			mu.Unlock()

			wg.Done()
		}(value)
	}
	wg.Wait()

	for _, content := range parser.Content {
		fmt.Println(content)
	}
	fmt.Println(len(parser.Content))

	log.Println(time.Since(start))
}

func ParseByXMLFile() {
	start := time.Now()

	parser, err := goxhparser.NewParser("./example/golang_useful.xml")
	if err != nil {
		log.Fatalln(err)
	}

	err = parser.Exec(time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	for _, content := range parser.Content {
		fmt.Println(content, content.Date.Format("02 January 2006'"))
	}
	fmt.Println(len(parser.Content))

	log.Println(time.Since(start))
}
