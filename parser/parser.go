package parser

import (
	"ziggytwister.com/shadow-hunter/ast"
	"ziggytwister.com/shadow-hunter/lexer"
	"ziggytwister.com/shadow-hunter/token"
)

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
	case token.KEYWORD:
		return p.parseKeywordElememt()
	case token.LCURLY:
		return p.parseMapElement()
	case token.LSQBRACKET:
		return p.parseVectorElement()
	default:
		return nil
	}
}

func (p *Parser) parseKeywordElememt() *ast.KeywordElement {
	return &ast.KeywordElement{Token: p.curToken, Value: p.curToken.Literal}
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

	return vector
}

func (p *Parser) parseMapElement() *ast.MapElement {
	element := &ast.MapElement{Token: p.curToken}

	var i = 0
	for !p.curTokenIs(token.RCURLY) {
		if i/2 == 0 {
			element.Keys = append(element.Keys, p.parseElement())
		} else {
			element.Values = append(element.Values, p.parseElement())
		}
		i += 1
		p.nextToken()
	}

	p.nextToken()

	return element
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
