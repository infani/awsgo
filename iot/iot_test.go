package iot

import (
	"context"
	"testing"

	"github.com/gookit/goutil/dump"
)

func TestListThings(t *testing.T) {
	thingNames, err := ListThings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	dump.P(thingNames)
}
