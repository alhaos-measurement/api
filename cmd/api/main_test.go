package main

import (
	"testing"
	"time"
)

func TestDummy(t *testing.T) {

	t.Log(time.Now().UTC().Format(time.RFC3339))

	ti, err := time.Parse(time.RFC3339, "2025-10-13T07:40:34Z")
	if err != nil {
		t.Error(err)
	}

	t.Log(ti.Local().Format(time.RFC3339))

}
