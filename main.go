package main

import (
	"fmt"
	p "pattern/pattern"
)

func main() {
	p.TestFacade()
	fmt.Println()
	p.TestBuilder()
	fmt.Println()
	p.TestVisitor()
}
