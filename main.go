package main

import (
	"fmt"

	"ziggytwister.com/shadow-hunter/lexer"
)

func main() {
	input := `(){}[]`
	l := lexer.New(input)
	fmt.Println(l.NextToken())
}
