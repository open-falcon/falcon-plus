package jsonql

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// JSONQL - JSON Query Lang struct encapsulating the JSON data.
type JSONQL struct {
	Data interface{}
}

// NewStringQuery - creates a new &JSONQL from raw JSON string
func NewStringQuery(jsonString string) (*JSONQL, error) {
	var data = new(interface{})
	err := json.Unmarshal([]byte(jsonString), data)
	if err != nil {
		return nil, err
	}
	return &JSONQL{*data}, nil
}

// NewQuery - creates a new &JSONQL from an array of interface{} or a map of [string]interface{}
func NewQuery(jsonObject interface{}) *JSONQL {
	return &JSONQL{jsonObject}
}

// Query - queries against the JSON using the conditions specified in the where stirng.
func (thisJSONQL *JSONQL) Query(where string) (interface{}, error) {
	parser := &Parser{
		Operators: sqlOperators,
	}
	tokens := parser.Tokenize(where)
	rpn, err := parser.ParseRPN(tokens)
	if err != nil {
		return nil, err
	}
	switch v := thisJSONQL.Data.(type) {
	case []interface{}:
		ret := []interface{}{}
		for _, obj := range v {
			parser.SymbolTable = obj
			r, err := thisJSONQL.processObj(parser, *rpn)
			if err != nil {
				return nil, err
			}
			if r {
				ret = append(ret, obj)
			}
		}
		return ret, nil
	case map[string]interface{}:
		parser.SymbolTable = v
		r, err := thisJSONQL.processObj(parser, *rpn)
		if err != nil {
			return nil, err
		}
		if r {
			return v, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("Failed to parse input data.")
	}
}

func (thisJSONQL *JSONQL) processObj(parser *Parser, rpn Lifo) (bool, error) {
	result, err := parser.Evaluate(&rpn, true)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return strconv.ParseBool(result)
}
