package main

import "fmt"

// TokenList is list of tokens to parse, received from tokenizer.
type TokenList []Token

func NewTokenList(t []Token) TokenList {
	return TokenList(t)
}

func (l TokenList) PopAt(at int) (TreeNode, TokenList, error) {
	if len(l) <= at {
		return nil, nil, fmt.Errorf("Invalid index %d", at)
	}
	tkn := l[at]
	rest := NewTokenList(l[at+1 : len(l)])
	return AdaptTokenToNode(tkn), rest, nil
}

func (l TokenList) PopNext() (TreeNode, TokenList, error) {
	return l.PopAt(0)
}

func (l TokenList) LookAt(at int) (TreeNode, error) {
	if len(l) <= at {
		return nil, fmt.Errorf("Invalid index %d", at)
	}
	tkn := l[at]
	return AdaptTokenToNode(tkn), nil
}
