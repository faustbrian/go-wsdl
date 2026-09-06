package wsdl_test

import (
	"context"
	"fmt"

	wsdl "github.com/faustbrian/go-wsdl"
)

func Example() {
	source := []byte(`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"` +
		` name="Inventory" targetNamespace="urn:inventory"/>`)

	document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	definitions, ok := document.Definitions11()
	if !ok {
		fmt.Println("not a WSDL 1.1 document")
		return
	}

	fmt.Println(document.Version())
	fmt.Println(definitions.Name)
	fmt.Println(definitions.TargetNamespace)
	// Output:
	// 1.1
	// Inventory
	// urn:inventory
}
