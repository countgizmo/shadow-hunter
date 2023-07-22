package parser

import (
	"fmt"
	"strconv"

	"ziggytwister.com/shadow-hunter/ast"
	"ziggytwister.com/shadow-hunter/lexer"
	"ziggytwister.com/shadow-hunter/token"
)

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens to set cur and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseEDN() *ast.EDN {
	edn := &ast.EDN{}
	edn.Elements = []ast.Element{}

	for p.curToken.Type != token.EOF {
		element := p.parseElement()
		if element != nil {
			edn.Elements = append(edn.Elements, element)
		}
		p.nextToken()
	}

	return edn
}

func (p *Parser) parseElement() ast.Element {
	switch p.curToken.Type {
	case token.INT:
		return p.parseIntegerLiteral()
	case token.KEYWORD:
		return p.parseKeywordElememt()
	case token.LCURLY:
		return p.parseMapElement()
	case token.LSQBRACKET:
		return p.parseVectorElement()
	case token.STRING:
		return p.parseStringElement()
	default:
		return nil
	}
}

func (p *Parser) parseKeywordElememt() *ast.KeywordElement {
	return &ast.KeywordElement{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStringElement() *ast.StringElement {
	return &ast.StringElement{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseVectorElement() *ast.VectorElement {
	vector := &ast.VectorElement{Token: p.curToken}
	p.nextToken()

	for p.curToken.Type != token.RSQBRACKET && p.curToken.Type != token.EOF {
		element := p.parseElement()
		if element != nil {
			vector.Elements = append(vector.Elements, element)
		}
		p.nextToken()
	}

	p.nextToken()

	return vector
}

func (p *Parser) parseMapElement() *ast.MapElement {
	mapElement := &ast.MapElement{Token: p.curToken}
	p.nextToken()

	var i = 0
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.RCURLY) {
		if i%2 == 0 {
			mapElement.Keys = append(mapElement.Keys, p.parseElement())
		} else {
			mapElement.Values = append(mapElement.Values, p.parseElement())
		}
		i += 1
		p.nextToken()
	}

	p.nextToken()

	return mapElement
}

func (p *Parser) parseIntegerLiteral() ast.Element {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		fmt.Println(msg)
		return nil
	}

	lit.Value = value
	return lit
}
