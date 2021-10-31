# goxhparser

See example folder

# SIMPLE

```
sourcesContent, xmlStruct, err := goxhparser.ParseByXMLFile("./golang_useful.xml")
	if err != nil {
		log.Fatalln(err)
	}
	for _,content := range sourcesContent {
		fmt.Println(xmlStruct.EntityID,content.SourceContent)
	}
 ```
 
 
 
 # GOROUTINES
 
 ```
 var wg sync.WaitGroup
	service,err := goxhparser.XMLToStruct("./golang_useful.xml")
	if err != nil {
		log.Fatalln(err)
	}
	for key,source := range service.Sources {
		wg.Add(1)
		go func(index int, source1 goxhparser.Source){
			content,err := goxhparser.ParseSource(source1)
			if err != nil {
				log.Fatalln(err)
			}
			service.Sources[index].SourceContent = content
			wg.Done()
		}(key,source)
	}
	wg.Wait()
	fmt.Println(service.Sources)
 ```
