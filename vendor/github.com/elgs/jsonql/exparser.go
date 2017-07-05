package jsonql

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Operator - encapsulates the Precedence and behavior logic of the operators.
type Operator struct {
	Precedence int
	Eval       func(symbolTable interface{}, left string, right string) (string, error)
}

// Parser - the main struct that contains operators, symbol table.
type Parser struct {
	Operators   map[string]*Operator
	SymbolTable interface{}
	maxOpLen    int
	initialized bool
}

// Init - init inspects the Operators and learns how long the longest operator string is
func (thisParser *Parser) Init() {
	for k := range thisParser.Operators {
		if len(k) > thisParser.maxOpLen {
			thisParser.maxOpLen = len(k)
		}
	}
}

// Calculate - gets the final result of the expression
func (thisParser *Parser) Calculate(expression string) (string, error) {
	tokens := thisParser.Tokenize(expression)
	//fmt.Println(expression, tokens)
	rpn, err := thisParser.ParseRPN(tokens)
	if err != nil {
		return "", err
	}
	return thisParser.Evaluate(rpn, true)
}

// Evaluate - evaluates the token stack until only one (as the final result) is left.
func (thisParser *Parser) Evaluate(ts *Lifo, postfix bool) (string, error) {
	newTs := &Lifo{}
	usefulWork := false
	for ti := ts.Pop(); ti != nil; ti = ts.Pop() {
		t := ti.(string)
		//		fmt.Println("t:", t)
		switch {
		case thisParser.Operators[t] != nil:
			// operators
			usefulWork = true
			if postfix {
				right := newTs.Pop()
				left := newTs.Pop()
				l := "0"
				r := "0"
				if left != nil {
					l = left.(string)
				}
				if right != nil {
					r = right.(string)
				}
				result, err := thisParser.Operators[t].Eval(thisParser.SymbolTable, l, r)
				newTs.Push(result)
				if err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", l, t, r) + " " + err.Error())
				}
			} else {
				right := ts.Pop()
				left := ts.Pop()
				l := ""
				r := ""
				if left != nil {
					l = left.(string)
				}
				if right != nil {
					r = right.(string)
				}
				result, err := thisParser.Operators[t].Eval(thisParser.SymbolTable, l, r)
				newTs.Push(result)
				if err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", l, t, r) + " " + err.Error())
				}
			}
		default:
			// operands
			newTs.Push(t)
		}
		//newTs.Print()
	}
	if !usefulWork {
		return "", errors.New("Failed to evaluate: no valid operator found.")
	}
	if newTs.Len() == 1 {
		return newTs.Pop().(string), nil
	}
	return thisParser.Evaluate(newTs, !postfix)
}

// false o1 in first, true o2 out first
func (thisParser *Parser) shunt(o1, o2 string) (bool, error) {
	op1 := thisParser.Operators[o1]
	op2 := thisParser.Operators[o2]
	if op1 == nil || op2 == nil {
		return false, errors.New(fmt.Sprint("Invalid operators:", o1, o2))
	}
	if op1.Precedence < op2.Precedence || (op1.Precedence <= op2.Precedence && op1.Precedence%2 == 1) {
		return true, nil
	}
	return false, nil
}

// ParseRPN - parses the RPN tokens
func (thisParser *Parser) ParseRPN(tokens []string) (output *Lifo, err error) {
	opStack := &Lifo{}
	outputQueue := []string{}
	for _, token := range tokens {
		switch {
		case thisParser.Operators[token] != nil:
			// operator
			for o2 := opStack.Peep(); o2 != nil; o2 = opStack.Peep() {
				stackToken := o2.(string)
				if thisParser.Operators[stackToken] == nil {
					break
				}
				o2First, err := thisParser.shunt(token, stackToken)
				if err != nil {
					return output, err
				}
				if o2First {
					outputQueue = append(outputQueue, opStack.Pop().(string))
				} else {
					break
				}
			}
			opStack.Push(token)
		case token == "(":
			opStack.Push(token)
		case token == ")":
			for o2 := opStack.Pop(); o2 != nil && o2.(string) != "("; o2 = opStack.Pop() {
				outputQueue = append(outputQueue, o2.(string))
			}
		default:
			outputQueue = append(outputQueue, token)
		}
	}
	for o2 := opStack.Pop(); o2 != nil; o2 = opStack.Pop() {
		outputQueue = append(outputQueue, o2.(string))
	}
	//fmt.Println(outputQueue)
	output = &Lifo{}
	for i := 0; i < len(outputQueue); i++ {
		(*output).Push(outputQueue[len(outputQueue)-i-1])
	}
	return
}

// Tokenize - splits the expression into tokens.
func (thisParser *Parser) Tokenize(exp string) (tokens []string) {
	if !thisParser.initialized {
		thisParser.Init()
	}
	sq, dq := false, false
	var tmp string
	expRunes := []rune(exp)
	for i := 0; i < len(expRunes); i++ {
		v := expRunes[i]
		s := string(v)
		switch {
		case unicode.IsSpace(v):
			if sq || dq {
				tmp += s
			} else if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
		case s == "'":
			tmp += s
			if !dq {
				sq = !sq
				if !sq {
					tokens = append(tokens, tmp)
					tmp = ""
				}
			}
		case s == "\"":
			tmp += s
			if !sq {
				dq = !dq
				if !dq {
					tokens = append(tokens, tmp)
					tmp = ""
				}
			}
		case s == "+" || s == "-" || s == "(" || s == ")":
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				lastToken := ""
				if len(tokens) > 0 {
					lastToken = tokens[len(tokens)-1]
				}
				if (s == "+" || s == "-") && (len(tokens) == 0 || lastToken == "(" || thisParser.Operators[lastToken] != nil) {
					// sign
					tmp += s
				} else {
					// operator
					tokens = append(tokens, s)
				}
			}
		default:
			if sq || dq {
				tmp += s
			} else {
				// until the max length of operators(n), check if next 1..n runes are operator, greedily
				opCandidateTmp := ""
				opCandidate := ""
				for j := 0; j < thisParser.maxOpLen && i < len(expRunes)-j-1; j++ {
					next := string(expRunes[i+j])
					opCandidateTmp += strings.ToUpper(next)
					if thisParser.Operators[opCandidateTmp] != nil {
						opCandidate = opCandidateTmp
					}
				}
				if len(opCandidate) > 0 {
					if len(tmp) > 0 {
						tokens = append(tokens, tmp)
						tmp = ""
					}
					tokens = append(tokens, opCandidate)
					i += len(opCandidate) - 1
				} else {
					tmp += s
				}
			}
		}
	}
	if len(tmp) > 0 {
		tokens = append(tokens, tmp)
		tmp = ""
	}
	return
}
