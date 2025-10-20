package main

import (
	"fmt"
	"github.com/alhaos-measurement/api/internal/config"
	"testing"
)

func TestDummy(t *testing.T) {

	conf, err := config.New(`C:\repo\alhaos-measurement\api\config\config.yml`)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", conf)

}
