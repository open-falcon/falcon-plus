# gojq
JSON query in Golang.

## Install
`go get -u github.com/elgs/gojq`

This library serves three purposes:

* makes parsing JSON configuration file much easier
* enables JSON expression evaluation
* reduces the pain of type assertion parsing JSON 


## Query from JSON Object
```go
package main

import (
	"fmt"

	"github.com/elgs/gojq"
)

var jsonObj = `
{
  "name": "sam",
  "gender": "m",
  "pet": null,
  "skills": [
    "Eating",
    "Sleeping",
    "Crawling"
  ],
  "hello.world":true
}
`

func main() {
	parser, err := gojq.NewStringQuery(jsonObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(parser.Query("name"))          // sam <nil>
	fmt.Println(parser.Query("gender"))        // m <nil>
	fmt.Println(parser.Query("skills.[1]"))    // Sleeping <nil>
	fmt.Println(parser.Query("hello"))         // <nil> hello does not exist.
	fmt.Println(parser.Query("pet"))           // <nil> <nil>
	fmt.Println(parser.Query("."))             // map[name:sam gender:m pet:<nil> skills:[Eating Sleeping Crawling] hello.world:true] <nil>
	fmt.Println(parser.Query("'hello.world'")) // true <nil>
}

```

## Query from JSON Array
```go
package main

import (
	"fmt"
	"github.com/elgs/gojq"
)

var jsonArray = `
[
  {
    "name": "elgs",
    "gender": "m",
    "skills": [
      "Golang",
      "Java",
      "C"
    ]
  },
  {
    "name": "enny",
    "gender": "f",
    "skills": [
      "IC",
      "Electric design",
      "Verification"
    ]
  },
  {
    "name": "sam",
    "gender": "m",
	"pet": null,
    "skills": [
      "Eating",
      "Sleeping",
      "Crawling"
    ]
  }
]
`

func main() {
	parser, err := gojq.NewStringQuery(jsonArray)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(parser.Query("[0].name"))       // elgs <nil>
	fmt.Println(parser.Query("[1].gender"))     // f <nil>
	fmt.Println(parser.Query("[2].skills.[1]")) // Sleeping <nil>
	fmt.Println(parser.Query("[2].hello"))      // <nil> hello does not exist.
	fmt.Println(parser.Query("[2].pet"))        // <nil> <nil>
}
```
## Netsted Query
```go
package main

import (
	"fmt"
	"github.com/elgs/gojq"
)

var jsonArray = `
[
  {
    "name": "elgs",
    "gender": "m",
    "skills": [
      "Golang",
      "Java",
      "C"
    ]
  },
  {
    "name": "enny",
    "gender": "f",
    "skills": [
      "IC",
      "Electric design",
      "Verification"
    ]
  },
  {
    "name": "sam",
    "gender": "m",
	"pet": null,
    "skills": [
      "Eating",
      "Sleeping",
      "Crawling"
    ]
  }
]
`

func main() {
	parser, err := gojq.NewStringQuery(jsonArray)
	if err != nil {
		fmt.Println(err)
		return
	}
	samSkills, err := parser.Query("[2].skills")
	fmt.Println(samSkills, err) //[Eating Sleeping Crawling] <nil>
	samSkillParser := gojq.NewQuery(samSkills)
	fmt.Println(samSkillParser.Query("[1]")) //Sleeping <nil>
}
```