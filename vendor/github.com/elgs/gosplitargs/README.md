# gosplitargs
Splits strings into tokens by given separator except treating quoted part as a single token.

#Installation
`go get -u github.com/elgs/gosplitargs`

#Usage
```go
package main

import (
	"fmt"
	"github.com/elgs/gosplitargs"
)

func main() {
	i1 := "I said 'I am sorry.', and he said \"it doesn't matter.\""
	o1, _ := gosplitargs.SplitArgs(i1, "\\s")
	for _, s := range o1 {
		fmt.Println(s)
	}
	/*
	   [ 'I',
	     'said',
	     'I am sorry.,',
	     'and',
	     'he',
	     'said',
	     'it doesn\'t matter.' ]
	*/

	i2 := "I said \"I am sorry.\", and he said \"it doesn't matter.\""
	o2, _ := gosplitargs.SplitArgs(i2, "\\s")
	for _, s := range o2 {
		fmt.Println(s)
	}
	/*
	   [ 'I',
	     'said',
	     'I am sorry.,',
	     'and',
	     'he',
	     'said',
	     'it doesn\'t matter.' ]
	*/

	i3 := `I said "I am sorry.", and he said "it doesn't matter."`
	o3, _ := gosplitargs.SplitArgs(i3, "\\s")
	for _, s := range o3 {
		fmt.Println(s)
	}
	/*
	   [ 'I',
	     'said',
	     'I am sorry.,',
	     'and',
	     'he',
	     'said',
	     'it doesn\'t matter.' ]
	*/

	i4 := `I said 'I am sorry.', and he said "it doesn't matter."`
	o4, _ := gosplitargs.SplitArgs(i4, "\\s")
	for _, s := range o4 {
		fmt.Println(s)
	}
	/*
	   [ 'I',
	     'said',
	     'I am sorry.,',
	     'and',
	     'he',
	     'said',
	     'it doesn\'t matter.' ]
	*/
}

```