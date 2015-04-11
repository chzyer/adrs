package conf

import (
	"strings"
	"testing"

	"gopkg.in/logex.v1"
)

func TestConfig(t *testing.T) {
	c, err := NewConfig(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	logex.Struct(c)
}
