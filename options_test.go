package spn

import (
	"fmt"
	"testing"
)

func TestOptionsEmpty(t *testing.T) {
	o := CaptureOptions{}
	query := o.Encode().Encode()
	if query != "" {
		t.Errorf("Expected empty query, got %s", query)
	}
}

func TestOptionsMany(t *testing.T) {
	o := CaptureOptions{
		IfNotArchivedWithin: "3d",
		CaptureCookie:       "balabalabala",
		CaptureAll:          true,
	}
	want := "capture_all=1&capture_cookie=balabalabala&if_not_archived_within=3d"
	query := o.Encode().Encode()
	if query != want {
		t.Errorf("Expected %s, got %s", want, query)
	}
	fmt.Println(query)
}
