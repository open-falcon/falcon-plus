/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

type Test struct {
	in  string
	out string
}

func init() {
	flag.Lookup("logtostderr").Value.Set("true")
	//flag.Lookup("v").Value.Set("5")
}

func TestExprText(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	goroot := os.Getenv("GOROOT")

	for i, test := range []Test{
		{"hello,{world}", "hello,{world}"},
		{"${GOPATH}", gopath},
		{"gopath:${GOPATH}", fmt.Sprintf("gopath:%s", gopath)},
		{"gopath:${GOPATH};goroot:${GOROOT}",
			fmt.Sprintf("gopath:%s;goroot:%s", gopath, goroot)},
	} {
		out := exprText([]byte(test.in))
		if out != test.out {
			t.Errorf(`#%d: exprText("%s")="%s"; want "%s"`,
				i, test.in, out, test.out)
		}
	}
}

func TestParse(t *testing.T) {
	conf := Parse("./etc/falcon.conf", true)
	fmt.Printf("conf:\n%s\n", conf)
}
