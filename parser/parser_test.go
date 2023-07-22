package parser

import (
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
	input := "[:name :age]"

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

	if v := edn.Elements[0].String(); v != "[:name :age]" {
		t.Fatalf("Vector should look like [:name :age] got %s", v)
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

//func TestMapElement(t *testing.T) {
//	input := `
//{:name "Jack"}
//`
//
//	l := lexer.New(input)
//	p := New(l)
//
//	edn := p.ParseEDN()
//	if edn == nil {
//		t.Fatal("ParseEDN() returned nil")
//	}
//	if len(edn.Elements) != 1 {
//		t.Fatalf("EDN expected to contain 1 element go %d",
//			len(edn.Elements))
//	}
//
//	tests := []struct {
//		expectedKey   string
//		expectedValue string
//	}{
//		{":name", "Jack"},
//	}
//
//	for i, tt := range tests {
//		element := edn.Elements[i]
//
//		if !testMapElement(t, element, tt.expectedKey, tt.expectedValue) {
//			return
//		}
//	}
//}
//
//func testMapElement(t *testing.T, e ast.Element, key string, value string) bool {
//	if e.TokenLiteral() != "{" {
//		t.Errorf("I don't know what I'm doing yet")
//		return false
//	}
//
//	return true
//}
