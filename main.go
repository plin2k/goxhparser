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
	Features      string `xml:"features,attr"`
	FeaturesSlice []string
}

type Source struct {
	Link     string `xml:",chardata"`
	RuleName string `xml:"rule,attr"`
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
}

type Content struct {
	Title        string
	Link         string
	ShortContent string
	FullContent  string
	Author       string
	Rating       string
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
	parser.contentRules()
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
		for _, feature := range strings.Split(v.Features, ",") {
			parser.Service.ContentRule.Content[i].FeaturesSlice = append(parser.Service.ContentRule.Content[i].FeaturesSlice, feature)
		}
	}
	parser.Service.Rules = nil
}

func (parser *Parser) reverseContentSlice() {
	for i, j := 0, len(parser.Content)-1; i < j; i, j = i+1, j-1 {
		parser.Content[i], parser.Content[j] = parser.Content[j], parser.Content[i]
	}
}
