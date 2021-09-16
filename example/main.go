package main

import (
	"fmt"
	"log"
	"plin2k.com/goxhparser"
	"sync"
)

func main() {
	ParseByXMLFileInGoroutine()
	ParseByXMLFile()

}

func ParseByXMLFileInGoroutine() {
	var wg sync.WaitGroup
	service, err := goxhparser.XMLToStruct("./golang_useful.xml")
	if err != nil {
		log.Fatalln(err)
	}
	for key, source := range service.Sources {
		wg.Add(1)
		go func(index int, source1 goxhparser.Source) {
			content, err := goxhparser.ParseSource(source1)
			if err != nil {
				log.Fatalln(err)
			}
			service.Sources[index].SourceContent = content
			wg.Done()
		}(key, source)
	}
	wg.Wait()
	fmt.Println(service.Sources)
}

func ParseByXMLFile() {
	sourcesContent, err := goxhparser.ParseByXMLFile("./golang_useful.xml")
	if err != nil {
		log.Fatalln(err)
	}
	for _, content := range sourcesContent {
		fmt.Println(content.SourceContent)
	}
}
