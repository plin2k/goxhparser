package goxhparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	parser := NewParser("./example/golang_useful.xml")
	err := parser.XMLToStruct()
	assert.Nil(t, err, "Error from function XMLToStruct")

	assert.NotNil(t, parser.Service.Title, "Title from xml file is empty")

	if len(parser.Service.Sources) <= 0 {
		t.Errorf("Sources from xml file not found")
	}

	assert.NotNil(t, parser.Service.Sources[0].Rule.Title, "Source Title Rule from xml file is empty")
}
