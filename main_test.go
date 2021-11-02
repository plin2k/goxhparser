package goxhparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	parser, err := NewParser("./example/golang_useful.xml")
	assert.Nil(t, err, "Error from function NewParser")
	err = parser.XMLToStruct()
	assert.Nil(t, err, "Error from function XMLToStruct")

	assert.NotNil(t, parser.Service.Title, "Title from xml file is empty")

	if len(parser.Service.Sources) <= 0 {
		t.Errorf("Sources from xml file not found")
	}

	assert.NotNil(t, parser.Service.Sources[0].Rule.Title, "Source Title Rule from xml file is empty")
}
