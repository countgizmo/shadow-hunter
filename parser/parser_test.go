package parser

import (
	"fmt"
	"testing"

	"ziggytwister.com/shadow-hunter/ast"
	"ziggytwister.com/shadow-hunter/lexer"
)

func TestKeywordElement(t *testing.T) {
	input := ":name"

	l := lexer.New(input)
	p := New(l)

	edn := p.ParseEDN()

	if edn == nil {
		t.Fatal("ParseEDN() returned nil")
	}

	if len(edn.Elements) != 1 {
		t.Fatalf("EDN expected to contain 1 element go %d",
			len(edn.Elements))
	}

	if edn.Elements[0].TokenLiteral() != ":name" {
		t.Fatalf("Expected first element to be :name got %v",
			edn.Elements[0].TokenLiteral())
	}
}

func TestVectorElement(t *testing.T) {
	input := "[:name :age 12]"

	l := lexer.New(input)
	p := New(l)

	edn := p.ParseEDN()

	if edn == nil {
		t.Fatal("ParseEDN() returned nil")
	}

	if len(edn.Elements) != 1 {
		t.Fatalf("EDN expected to contain 1 element go %d",
			len(edn.Elements))
	}

	if edn.Elements[0].TokenLiteral() != "vector" {
		t.Fatalf("Expected first element to be vector got %v",
			edn.Elements[0].TokenLiteral())
	}

	if v := edn.Elements[0].String(); v != "[:name :age 12]" {
		t.Fatalf("Vector should look like [:name :age 12] got %s", v)
	}

	vectorElement, ok := edn.Elements[0].(*ast.VectorElement)

	if !ok {
		t.Fatalf("Element not *ast.VectorElement. got=%T", edn.Elements[0])
	}

	if ve := vectorElement.Elements[0].TokenLiteral(); ve != ":name" {
		t.Fatalf("Expected first element of vector to be :name got %v", ve)
	}

	if ve := vectorElement.Elements[1].TokenLiteral(); ve != ":age" {
		t.Fatalf("Expected second element of vector to be :age got %v", ve)
	}

	if ve := vectorElement.Elements[2].TokenLiteral(); ve != "12" {
		t.Fatalf("Expected third element of vector to be '12' got %v", ve)
	}

	num, ok := vectorElement.Elements[2].(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Element not *ast.IntegerLiteral. got=%T", vectorElement.Elements[2])
	}

	if num.Value != 12 {
		t.Fatalf("Expect third element to be an integer 12 got %d", num.Value)
	}
}

func TestEmptyVectorElement(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)

	edn := p.ParseEDN()

	if edn == nil {
		t.Fatal("ParseEDN() returned nil")
	}

	if len(edn.Elements) != 1 {
		t.Fatalf("EDN expected to contain 0 element go %d",
			len(edn.Elements))
	}

	if v := edn.Elements[0].String(); v != "[]" {
		t.Fatalf("Empty vector should look like [] got %s", v)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	edn := p.ParseEDN()

	if len(edn.Elements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(edn.Elements))
	}

	element, ok := edn.Elements[0].(ast.Element)
	if !ok {
		t.Fatalf("edn.Elements[0] is not ast.Element got=%T", edn.Elements[0])
	}

	literal, ok := element.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", element)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestMapElement(t *testing.T) {
	input := `
	{:name "Jack"
	:age 39}
	`

	l := lexer.New(input)
	p := New(l)

	edn := p.ParseEDN()

	fmt.Println(edn)
	if edn == nil {
		t.Fatal("ParseEDN() returned nil")
	}
	if len(edn.Elements) != 1 {
		t.Fatalf("EDN expected to contain 1 element go %d",
			len(edn.Elements))
	}

	tests := []struct {
		expectedKey   string
		expectedValue string
	}{
		{":name", "Jack"},
		{":age", "39"},
	}

	mapElement := edn.Elements[0].(*ast.MapElement)

	for i, tt := range tests {

		if actualKey := mapElement.Keys[i]; actualKey.String() != tt.expectedKey {
			t.Fatalf("Expected key %s got %s", tt.expectedKey, actualKey)
		}

		if actualValue := mapElement.Values[i]; actualValue.String() != tt.expectedValue {
			t.Fatalf("Expected value %s got %s", tt.expectedValue, actualValue)
		}
	}
}
