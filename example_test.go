package gltf

import (
	"fmt"
)

func ExampleOpen() {
	doc, err := Open("fake")
	if err != nil {
		panic(err)
	}
	fmt.Print(doc.Asset)
}
