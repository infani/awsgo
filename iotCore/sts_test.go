package iotCore

import "testing"

func Test_getEndpointAddress(t *testing.T) {
	got, err := getEndpointAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(got)
}
