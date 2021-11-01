# goxhparser

See example folder

# SIMPLE

```go
    parser := goxhparser.NewParser("./example/golang_useful.xml")
    err := parser.Exec()
    if err != nil {
        log.Fatalln(err)
    }

    for _, content := range parser.Content {
        fmt.Println(content)
    }	
 ```
 
 
 
 # GOROUTINES
 
 ```go
    var wg sync.WaitGroup
    var mu sync.Mutex
    parser := goxhparser.NewParser("./example/golang_useful.xml")
    err := parser.XMLToStruct()
    if err != nil {
        log.Fatalln(err)
    }
    for _, value := range parser.Service.Sources {
        wg.Add(1)
        go func(source goxhparser.Source) {
            content,err := parser.Parse(source)
            if err != nil {
                log.Fatalln(err)
            }

            mu.Lock()
            parser.Content = append(parser.Content,content...)
            mu.Unlock()

            wg.Done()
        }(value)
    }
    wg.Wait()

    for _, content := range parser.Content {
        fmt.Println(content)
    }
 ```
