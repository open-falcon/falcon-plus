/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/golang/glog"
)

const (
	eof           = 0
	MAX_CTX_LEVEL = 16
)

var (
	conf                   *FalconConfig
	yy                     *yyLex
	yy_ctrl                *ConfCtrl
	yy_agent               *ConfAgent
	yy_loadbalance         *ConfLoadbalance
	yy_backend             *ConfBackend
	yy_loadbalance_backend = &LbBackend{}
	yy_ss                  = make(map[string]string)
	yy_ss2                 = make(map[string]string)
	yy_as                  = make([]string, 0)

	f_ip      = regexp.MustCompile(`^[0-9]+\.[0-0]+\.[0-9]+\.[0-9]+[ \t\n;{}]{1}`)
	f_num     = regexp.MustCompile(`^0x[0-9a-fA-F]+|^[0-9]+[ \t\n;{}]{1}`)
	f_keyword = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]+[ \t\n;{}]{1}`)
	f_word    = regexp.MustCompile(`(^"[^"]+")|(^[^"\n \t;]+)`)
	f_env     = regexp.MustCompile(`\$\{[a-zA-Z][0-9a-zA-Z_]+\}`)

	keywords = map[string]int{
		//general
		"on":       ON,
		"yes":      YES,
		"off":      OFF,
		"no":       NO,
		"include":  INCLUDE,
		"root":     ROOT,
		"pidFile":  PID_FILE,
		"log":      LOG,
		"host":     HOST,
		"disabled": DISABLED,
		"debug":    DEBUG,

		// module name
		"ctrl":        CTRL,
		"agent":       AGENT,
		"loadbalance": LOADBALANCE,
		"backend":     BACKEND,

		// module
		"upstream": UPSTREAM,
		"metric":   METRIC,
		"migrate":  MIGRATE,
	}
)

type yyCtx struct {
	text []byte
	pos  int
	lino int
	file string
}

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type yyLex struct {
	ctxData [MAX_CTX_LEVEL]yyCtx
	ctxL    int
	ctx     *yyCtx
	t       []byte
	i       int
	debug   bool
}

func prefix(a, b []byte) bool {
	if len(a) < len(b) {
		return false
	}

	if len(a) == len(b) {
		return bytes.Equal(a, b)
	}

	return bytes.Equal(a[:len(b)], b)
}

func exprText(s []byte) (ret string) {
	var i int
	var es [][]int

	if es = f_env.FindAllIndex(s, -1); es == nil {
		return string(s)
	}

	for j := 0; j < len(es); j++ {
		if i < es[j][0] {
			ret += string(s[i:es[j][0]])
		}
		ret += os.Getenv(string(s[es[j][0]+2 : es[j][1]-1]))
		i = es[j][1]
	}
	return ret + string(s[i:])
}

func (p *yyLex) include(filename string) (err error) {
	p.ctxL++
	p.ctx = &p.ctxData[p.ctxL]
	p.ctx.lino = 1
	p.ctx.pos = 0
	p.ctx.file = filename
	if p.ctx.text, err = ioutil.ReadFile(filename); err != nil {
		dir, _ := os.Getwd()
		glog.Errorf(MODULE_NAME+"%s(curdir:%s)", err.Error(), dir)
		os.Exit(1)
	}
	glog.V(5).Infof(MODULE_NAME+"ctx level %d", p.ctxL)
	return nil
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (p *yyLex) Lex(yylval *yySymType) int {
	var f []byte
	var b bool

begin:
	text := p.ctx.text[p.ctx.pos:]
	for {
		if p.ctx.pos == len(p.ctx.text) {
			if p.ctxL > 0 {
				p.ctxL--
				p.ctx = &p.ctxData[p.ctxL]
				goto begin
			}
			return eof
		}

		for text[0] == ' ' || text[0] == '\t' || text[0] == '\n' {
			p.ctx.pos += 1
			if p.ctx.pos == len(p.ctx.text) {
				glog.V(5).Infof(MODULE_NAME+"ctx level %d", p.ctxL)
				if p.ctxL > 0 {
					p.ctxL--
					p.ctx = &p.ctxData[p.ctxL]
					goto begin
				}
				return eof
			}
			if text[0] == '\n' {
				p.ctx.lino++
			}
			text = p.ctx.text[p.ctx.pos:]
		}

		b = prefix(text, []byte("include"))
		if b {
			p.ctx.pos += len("include")
			return INCLUDE
		}

		f = f_ip.Find(text)
		if f != nil {
			s := f[:len(f)-1]
			p.ctx.pos += len(s)
			p.t = s[:]
			return IPA
		}

		f = f_num.Find(text)
		if f != nil {
			s := f[:len(f)-1]
			p.ctx.pos += len(s)
			p.t = s[:]
			i64, _ := strconv.ParseInt(string(s), 0, 0)
			p.i = int(i64)
			glog.V(5).Infof(MODULE_NAME+"return NUM %d\n", p.i)
			return NUM
		}

		// find keyword
		f = f_keyword.Find(text)
		if f != nil {
			s := f[:len(f)-1]
			if val, ok := keywords[string(s)]; ok {
				p.ctx.pos += len(s)
				glog.V(5).Infof(MODULE_NAME+"find %s return %d\n", string(s), val)
				return val
			}
		}

		if bytes.IndexByte([]byte(`={}:;,()+*/%<>~\[\]?!\|-`), text[0]) != -1 {
			if !prefix(text, []byte(`//`)) &&
				!prefix(text, []byte(`/*`)) {
				p.ctx.pos++
				glog.V(5).Infof(MODULE_NAME+"return '%c'\n", int(text[0]))
				return int(text[0])
			}
		}

		// comm
		if text[0] == '#' || prefix(text, []byte(`//`)) {
			for p.ctx.pos < len(p.ctx.text) {
				//glog.Infof(MODULE_NAME+"%c", p.ctx.text[p.ctx.pos])
				if p.ctx.text[p.ctx.pos] == '\n' {
					p.ctx.pos++
					p.ctx.lino++
					goto begin
				}
				p.ctx.pos++
			}
			return eof
		}

		// ccomm
		if prefix(text, []byte(`/*`)) {
			p.ctx.pos += 2
			for p.ctx.pos < len(p.ctx.text) {
				if p.ctx.text[p.ctx.pos] == '\n' {
					p.ctx.lino++
				}
				if p.ctx.text[p.ctx.pos] == '*' {
					if p.ctx.text[p.ctx.pos-1] == '/' {
						p.Error("Comment nesting not supported")
					}
					if p.ctx.text[p.ctx.pos+1] == '/' {
						p.ctx.pos += 2
						goto begin
					}
				}
				p.ctx.pos++
			}
		}

		// find text
		f = f_word.Find(text)
		if f != nil {
			p.ctx.pos += len(f)
			if f[0] == '"' {
				p.t = f[1 : len(f)-1]
			} else {
				p.t = f[:]
			}
			glog.V(5).Infof(MODULE_NAME+"return TEXT(%s)", string(p.t))
			return TEXT
		}
		p.Error(fmt.Sprintf("unknown character %c", text[0]))
	}
}

// The parser calls this method on a parse error.
func (p *yyLex) Error(s string) {
	bline := 3
	aline := 3
	p.ctx.pos--
	out := fmt.Sprintf("\x1B[31m%c\x1B[0m", p.ctx.text[p.ctx.pos])

	lino := p.ctx.lino
	for pos := p.ctx.pos - 1; pos > 0; pos-- {
		if p.ctx.text[pos] == '\n' {
			if p.ctx.lino-lino < bline {
				out = fmt.Sprintf("%3d%s", lino, out)
				lino--
			} else {
				out = fmt.Sprintf("%3d%s", lino, out)
				break
			}
		}
		out = fmt.Sprintf("%c%s", p.ctx.text[pos], out)
	}

	lino = p.ctx.lino
	for pos := p.ctx.pos + 1; pos < len(p.ctx.text); pos++ {
		out = fmt.Sprintf("%s%c", out, p.ctx.text[pos])
		if p.ctx.text[pos] == '\n' {
			lino++
			if lino-p.ctx.lino < aline {
				out = fmt.Sprintf("%s%3d", out, lino)
			} else {
				break
			}
		}
	}

	glog.Errorf(MODULE_NAME+"parse file(%s) error: %s\n%s",
		p.ctx.file, s, out)
	os.Exit(1)
}

func Parse(filename string, debug bool) *FalconConfig {
	var err error
	conf = &FalconConfig{ConfigFile: filename}
	yy = &yyLex{
		ctxL:  0,
		debug: debug,
	}
	yy.ctx = &yy.ctxData[0]
	yy.ctx.file = filename
	yy.ctx.lino = 1
	yy.ctx.pos = 0
	if yy.ctx.text, err = ioutil.ReadFile(filename); err != nil {
		glog.Errorf(MODULE_NAME+"parse file(%s) error: %s\n",
			filename, err.Error())
		os.Exit(1)
	}
	yyParse(yy)
	return conf
}
