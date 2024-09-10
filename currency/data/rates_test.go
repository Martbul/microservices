package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

//! when creating a tests in go the file must end at `_test.go` in order for the compiler to ignore the file and
//! every test function must start with `Test...`

func TestNewRates(t *testing.T) {
	testRates, err := NewRates(hclog.Default())
	if err != nil {
		//failing the test imediately
		t.Fatal(err)
	}


	fmt.Printf("Rates %#v", testRates.rates )
}