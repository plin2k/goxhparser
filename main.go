package goxhparser

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Parser struct {
	File    []byte
	Service Service
	Content []Content
}

type FeaturesSlice map[string][]string

type Service struct {
	XMLName     xml.Name `xml:"xml"`
	Title       string   `xml:"title"`
	Name        string   `xml:"name"`
	EntityID    string   `xml:"entity_id"`
	EntityType  string   `xml:"entity_type"`
	Sources     []Source `xml:"source"`
	Rules       []Rule   `xml:"rule"`
	ContentRule struct {
		Content []ContentRuleField `xml:"content"`
	} `xml:"content_rule"`
}

type ContentRuleField struct {
	Field         string `xml:",chardata"`
	Prefix        string `xml:"prefix,attr"`
	Postfix       string `xml:"postfix,attr"`
	BottomPadding int    `xml:"bottom_padding,attr"`
	TopPadding    int    `xml:"top_padding,attr"`
	Bold          bool   `xml:"bold,attr"`
	Italic        bool   `xml:"italic,attr"`
	Features      string `xml:"features,attr"`
	FeaturesSlice
}

type Source struct {
	Link     string `xml:",chardata"`
	RuleName string `xml:"rule,attr"`
	TagName  string `xml:"tag,attr"`
	Rule     Rule
}

type Rule struct {
	Name  string `xml:"name,attr"`
	Block string `xml:"block"`
	Title string `xml:"title"`
	Link  struct {
		Href   string `xml:",chardata"`
		Prefix string `xml:"prefix,attr"`
	} `xml:"link"`
	ShortContent string `xml:"short_content"`
	FullContent  string `xml:"full_content"`
	Author       string `xml:"author"`
	Rating       string `xml:"rating"`
	Date         struct {
		Time      string `xml:",chardata"`
		Layout    string `xml:"layout,attr"`
		Attribute string `xml:"attribute,attr"`
	} `xml:"date"`
}

type Content struct {
	Title        string
	Link         string
	ShortContent string
	FullContent  string
	Author       string
	Rating       string
	Date         time.Time

	SourceTagName string
}

func NewParser(filename string) (*Parser, error) {
	var parser Parser

	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	parser.File, err = ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}
	return &parser, nil
}

func (parser *Parser) XMLToStruct() error {

	err := xml.Unmarshal(parser.File, &parser.Service)
	if err != nil {
		return err
	}
	parser.ruleToSource()
	parser.contentRules()
	return nil
}

func (parser *Parser) Exec(delay time.Duration) error {
	err := parser.XMLToStruct()
	if err != nil {
		return err
	}
	for _, source := range parser.Service.Sources {
		content, err := parser.Parse(source)
		if err != nil {
			return err
		}
		parser.Content = append(parser.Content, content...)
		time.Sleep(delay)
	}
	parser.reverseContentSlice()
	return nil
}

func (parser *Parser) ruleToSource() {
	for _, rule := range parser.Service.Rules {
		for skey, source := range parser.Service.Sources {
			if source.RuleName == rule.Name {
				parser.Service.Sources[skey].Rule = rule
			}
		}
	}
	parser.Service.Rules = nil
}

func (parser *Parser) contentRules() {
	for i, v := range parser.Service.ContentRule.Content {
		parser.contentRuleFeatures(i, v)
	}
	parser.Service.Rules = nil
}

func (parser *Parser) contentRuleFeatures(i int, v ContentRuleField) {
	var featureSlice = make(FeaturesSlice)
	for _, feature := range strings.Split(v.Features, ";") {
		if featureArr := strings.Split(feature, ":"); len(featureArr) > 1 {
			featureSlice[featureArr[0]] = strings.Split(featureArr[1], ",")
		} else {
			featureSlice[featureArr[0]] = []string{}
		}
	}
	parser.Service.ContentRule.Content[i].FeaturesSlice = featureSlice
}

func (parser *Parser) reverseContentSlice() {
	for i, j := 0, len(parser.Content)-1; i < j; i, j = i+1, j-1 {
		parser.Content[i], parser.Content[j] = parser.Content[j], parser.Content[i]
	}

	for i, j := 0, len(parser.Content)-1; i < j; i, j = i+1, j-1 {
		parser.Content[i], parser.Content[j] = parser.Content[j], parser.Content[i]
	}

}
