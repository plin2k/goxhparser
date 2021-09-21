package goxhparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXMLToStruct(t *testing.T) {
	service, err := XMLToStruct("./example/golang_useful.xml")
	assert.Nil(t,err, "Error from function")

	assert.NotNil(t, service.Title,"Title from xml file is empty")

	if len(service.Sources) <= 0 {
		t.Errorf("Sources from xml file not found")
	}

	assert.NotNil(t, service.Sources[0].Title,"Source Title from xml file is empty")
}