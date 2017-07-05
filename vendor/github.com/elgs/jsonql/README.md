# jsonql
JSON query expression library in Golang.

This library enables query against JSON. Currently supported operators are: (precedences from low to high)

```
||
&&
= != > < >= <= ~= !~=
+ -
* / %
^
( )
```

## Changes
Previously I was hoping to make the query as similar to SQL `WHERE` clause as possible. Later I found a problem parsing the term `PRIORITY>5`. The tokenizer split it as `PRI`, `OR`, `ITY`, `>`, `5`, since `OR` was then an operator, which is terribly wrong. At that point, I thought of two choices:

1. to force the query expression to contain at least one white space between tokens, thus `PRIORITY>5` should be written as `PRIORITY > 5`;
2. to replace operators as follows:
	* `AND` to `&&`
	* `OR`	to `||`
	* `RLIKE` to `~=`
	* `NOT RLIKE` to `!~=`

I adopted the second choice as the new operators I believe are still quite intuitive to most programmers, rather than forcing to put white spaces between tokens.

## Install
`go get -u github.com/elgs/jsonql`

## TODD
* Implement `IS NULL` and `IS NOT NULL` (Patches welcome)

## Example
```go
package main

import (
	"fmt"
	"github.com/elgs/jsonql"
)

var jsonString = `
[
  {
    "name": "elgs",
    "gender": "m",
	"age": 35,
    "skills": [
      "Golang",
      "Java",
      "C"
    ]
  },
  {
    "name": "enny",
    "gender": "f",
    "age": 36,
	"skills": [
      "IC",
      "Electric design",
      "Verification"
    ]
  },
  {
    "name": "sam",
    "gender": "m",
	"age": 1,
    "skills": [
      "Eating",
      "Sleeping",
      "Crawling"
    ]
  }
]
`

func main() {
	parser, err := jsonql.NewStringQuery(jsonString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(parser.Query("name='elgs'"))
	//[map[skills:[Golang Java C] name:elgs gender:m age:35]] <nil>

	fmt.Println(parser.Query("name='elgs' && gender='f'"))
	//[] <nil>

	fmt.Println(parser.Query("age<10 || (name='enny' && gender='f')"))
	//[map[gender:f age:36 skills:[IC Electric design Verification] name:enny] map[age:1 skills:[Eating Sleeping Crawling] name:sam gender:m]] <nil>

	fmt.Println(parser.Query("age<10"))
	//[map[name:sam gender:m age:1 skills:[Eating Sleeping Crawling]]] <nil>

	fmt.Println(parser.Query("1=0"))
	//[] <nil>

	fmt.Println(parser.Query("age=(2*3)^2"))
	//[map[skills:[IC Electric design Verification] name:enny gender:f age:36]] <nil>
	
	fmt.Println(parser.Query("name ~= 'e.*'"))
	//[map[age:35 skills:[Golang Java C] name:elgs gender:m] map[skills:[IC Electric design Verification] name:enny gender:f age:36]] <nil>
	
	fmt.Println(parser.Query("name='el'+'gs'"))
	fmt.Println(parser.Query("age=30+5.0"))
	fmt.Println(parser.Query("age=40.0-5"))
	fmt.Println(parser.Query("age=70-5*7"))
	fmt.Println(parser.Query("age=70.0/2.0"))
	fmt.Println(parser.Query("age=71%36"))
	// [map[name:elgs gender:m age:35 skills:[Golang Java C]]] <nil>
}
```