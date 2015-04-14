package conf

import (
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := NewConfig(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	_ = c

}
