package cmd

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	p := fmt.Println

	ll := listDirectory("/")
	p(ll)
}
