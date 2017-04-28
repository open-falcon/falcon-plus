//line yyparse.y:15
package falcon

import __yyfmt__ "fmt"

//line yyparse.y:16
import (
	"fmt"
	"os"
)

//line yyparse.y:25
type yySymType struct {
	yys  int
	num  int
	text string
	b    bool
}

const NUM = 57346
const TEXT = 57347
const IPA = 57348
const ON = 57349
const YES = 57350
const OFF = 57351
const NO = 57352
const INCLUDE = 57353
const ROOT = 57354
const PID_FILE = 57355
const LOG = 57356
const HOST = 57357
const DISABLED = 57358
const DEBUG = 57359
const CTRL = 57360
const AGENT = 57361
const LOADBALANCE = 57362
const BACKEND = 57363
const UPSTREAM = 57364
const METRIC = 57365
const MIGRATE = 57366

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUM",
	"TEXT",
	"IPA",
	"'{'",
	"'}'",
	"';'",
	"ON",
	"YES",
	"OFF",
	"NO",
	"INCLUDE",
	"ROOT",
	"PID_FILE",
	"LOG",
	"HOST",
	"DISABLED",
	"DEBUG",
	"CTRL",
	"AGENT",
	"LOADBALANCE",
	"BACKEND",
	"UPSTREAM",
	"METRIC",
	"MIGRATE",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yyparse.y:276

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 94
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 203

var yyAct = [...]int{

	139, 137, 156, 155, 154, 16, 19, 20, 21, 29,
	37, 47, 56, 58, 59, 60, 61, 153, 142, 18,
	17, 69, 122, 18, 17, 68, 158, 152, 75, 144,
	77, 80, 140, 150, 82, 145, 63, 85, 86, 89,
	90, 135, 132, 125, 93, 130, 95, 115, 97, 100,
	101, 79, 91, 103, 127, 105, 83, 107, 110, 88,
	128, 81, 67, 74, 92, 66, 78, 18, 17, 99,
	141, 136, 84, 102, 87, 65, 140, 62, 109, 18,
	17, 94, 18, 17, 98, 49, 70, 71, 72, 73,
	104, 57, 52, 108, 134, 54, 51, 53, 18, 17,
	126, 22, 114, 113, 55, 112, 111, 30, 24, 106,
	96, 27, 25, 26, 76, 64, 118, 119, 124, 28,
	50, 129, 143, 18, 17, 131, 31, 138, 123, 117,
	41, 32, 39, 33, 23, 11, 36, 34, 35, 10,
	146, 149, 9, 38, 8, 116, 2, 1, 0, 133,
	18, 17, 0, 40, 157, 0, 0, 0, 0, 48,
	43, 148, 0, 45, 42, 44, 151, 3, 0, 46,
	0, 0, 6, 7, 4, 5, 147, 0, 0, 12,
	13, 14, 15, 64, 18, 17, 0, 0, 0, 70,
	71, 72, 73, 18, 17, 0, 121, 0, 0, 0,
	0, 0, 120,
}
var yyPact = [...]int{

	-1000, 158, -1000, -1000, 74, 74, 74, 74, 93, 118,
	145, 77, 74, 74, 74, 74, 68, -1000, -1000, 111,
	66, 56, -1000, 53, 74, 76, 111, 74, 107, 179,
	74, -1000, 52, 74, 76, 111, 74, 179, 74, 74,
	-1000, 43, 76, 74, 111, 74, 103, 179, 74, -1000,
	41, 76, 74, 111, 74, 102, 179, 74, 99, 98,
	96, 95, -1000, 38, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 188, 14, 35, -1000,
	74, -1000, -1000, 36, 74, -1000, 33, 76, 87, 32,
	-1000, 64, -1000, -1000, -1000, -1000, -1000, 62, 10, 179,
	74, -1000, -1000, 24, 76, 20, 8, -5, -6, -7,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 18, -1000,
}
var yyPgo = [...]int{

	0, 21, 0, 36, 147, 146, 1, 145, 144, 142,
	139, 135, 134, 131, 130, 129, 128, 127, 122, 120,
	116, 100,
}
var yyR1 = [...]int{

	0, 4, 4, 1, 1, 1, 1, 1, 2, 2,
	3, 6, 6, 6, 6, 6, 7, 7, 7, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 8, 8,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 9, 9, 13, 13, 13, 13, 13, 13, 13,
	13, 13, 13, 13, 10, 10, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 15, 15, 16,
	16, 17, 17, 18, 18, 18, 11, 11, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 20,
	20, 21, 21, 21,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 1, 1, 1, 0, 1, 1,
	1, 0, 4, 4, 4, 4, 0, 2, 4, 1,
	3, 4, 3, 3, 2, 2, 2, 2, 3, 3,
	0, 2, 2, 1, 2, 2, 4, 2, 2, 2,
	2, 3, 3, 0, 2, 2, 1, 2, 2, 2,
	2, 2, 2, 2, 3, 3, 0, 2, 2, 1,
	2, 2, 4, 2, 2, 2, 2, 0, 3, 0,
	5, 0, 3, 0, 2, 4, 3, 3, 0, 2,
	2, 1, 2, 2, 4, 2, 2, 2, 2, 0,
	3, 0, 2, 4,
}
var yyChk = [...]int{

	-1000, -4, -5, 9, 16, 17, 14, 15, -8, -9,
	-10, -11, 21, 22, 23, 24, -2, 6, 5, -2,
	-2, -2, 8, -12, 15, 19, 20, 18, 26, -2,
	14, 8, -13, 15, 19, 20, 18, -2, 25, 14,
	8, -14, 19, 15, 20, 18, 24, -2, 14, 8,
	-19, 19, 15, 20, 18, 27, -2, 14, -2, -2,
	-2, -2, 9, -3, 4, 9, 9, 9, -2, -1,
	10, 11, 12, 13, -3, -2, 7, -2, -3, -1,
	-2, 9, -2, -1, -3, -2, -2, -3, -1, -2,
	-2, 9, -1, -2, -3, -2, 7, -2, -3, -1,
	-2, 9, -1, -2, -3, -2, 7, -2, -3, -1,
	-2, 7, 7, 7, 7, 9, -7, -15, -20, -2,
	14, 8, 8, -16, -2, 8, -21, 19, 25, -2,
	9, -2, 9, -1, 7, 9, 7, -6, -17, -2,
	14, 8, 8, -18, 19, 25, -2, -3, -1, -2,
	9, -1, 7, 9, 9, 9, 9, -6, 8,
}
var yyDef = [...]int{

	1, -2, 2, 19, 0, 0, 0, 0, 30, 43,
	56, 78, 0, 0, 0, 0, 0, 8, 9, 0,
	0, 0, 24, 0, 0, 7, 33, 0, 0, 7,
	0, 25, 0, 0, 7, 46, 0, 7, 0, 0,
	26, 0, 7, 0, 59, 0, 0, 7, 0, 27,
	0, 7, 0, 81, 0, 0, 7, 0, 0, 0,
	0, 0, 20, 0, 10, 22, 23, 29, 31, 32,
	3, 4, 5, 6, 34, 35, 16, 37, 38, 39,
	40, 42, 44, 45, 47, 48, 49, 50, 51, 52,
	53, 55, 57, 58, 60, 61, 67, 63, 64, 65,
	66, 77, 79, 80, 82, 83, 89, 85, 86, 87,
	88, 28, 41, 54, 76, 21, 0, 69, 91, 17,
	0, 36, 62, 0, 0, 84, 0, 7, 0, 0,
	68, 0, 90, 92, 11, 18, 71, 0, 73, 7,
	0, 93, 70, 0, 7, 0, 0, 0, 0, 0,
	72, 74, 11, 12, 13, 14, 15, 0, 75,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 9,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 8,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 10, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:50
		{
			yyVAL.b = true
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:51
		{
			yyVAL.b = true
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:52
		{
			yyVAL.b = false
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:53
		{
			yyVAL.b = false
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line yyparse.y:54
		{
			yyVAL.b = true
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:58
		{
			yyVAL.text = string(yy.t)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:59
		{
			yyVAL.text = exprText(yy.t)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:63
		{
			yyVAL.num = yy.i
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:67
		{
			yy_ss[yyDollar[2].text] = yyDollar[3].text
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:68
		{
			yy_ss[yyDollar[2].text] = fmt.Sprintf("%d", yyDollar[3].num)
		}
	case 14:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:69
		{
			yy_ss[yyDollar[2].text] = fmt.Sprintf("%v", yyDollar[3].b)
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:70
		{
			yy.include(yyDollar[3].text)
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:74
		{
			yy_as = append(yy_as, yyDollar[2].text)
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:75
		{
			yy.include(yyDollar[3].text)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:79
		{
			conf.pidFile = yyDollar[2].text
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:80
		{
			conf.log = yyDollar[2].text
			conf.logv = yyDollar[3].num
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:83
		{
			yy.include(yyDollar[2].text)
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:84
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:88
		{
			yy_ctrl.Ctrl.Set(APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)

			yy_ctrl.Name = fmt.Sprintf("ctrl_%s", yy_ctrl.Name)
			if yy_ctrl.Host == "" {
				yy_ctrl.Host, _ = os.Hostname()
			}

			if !yy_ctrl.Disabled || yy.debug {
				conf.conf = append(conf.conf, yy_ctrl)
			}
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:100
		{
			yy_agent.Configer.Set(APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)

			yy_agent.Name = fmt.Sprintf("agent_%s", yy_agent.Name)
			if yy_agent.Host == "" {
				yy_agent.Host, _ = os.Hostname()
			}

			if !yy_agent.Disabled || yy.debug {
				conf.conf = append(conf.conf, yy_agent)
			}
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:112
		{
			yy_loadbalance.Configer.Set(APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)

			yy_loadbalance.Name = fmt.Sprintf("loadbalance_%s", yy_loadbalance.Name)
			if yy_loadbalance.Host == "" {
				yy_loadbalance.Host, _ = os.Hostname()
			}
			if !yy_loadbalance.Disabled || yy.debug {
				conf.conf = append(conf.conf, yy_loadbalance)
			}
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:123
		{
			yy_backend.Configer.Set(APP_CONF_FILE, yy_ss2)
			yy_ss2 = make(map[string]string)

			yy_backend.Name = fmt.Sprintf("backend_%s", yy_backend.Name)
			if yy_backend.Host == "" {
				yy_backend.Host, _ = os.Hostname()
			}
			if !yy_backend.Disabled || yy.debug {
				conf.conf = append(conf.conf, yy_backend)
			}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:138
		{
			yy_ctrl = &ConfCtrl{}
			yy_ctrl.Ctrl.Set(APP_CONF_DEFAULT, ConfDefault["ctrl"])
			yy_ctrl.Name = yyDollar[2].text
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:146
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:150
		{
			yy_ctrl.Disabled = yyDollar[2].b
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:151
		{
			yy_ctrl.Debug = 1
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:152
		{
			yy_ctrl.Debug = yyDollar[2].num
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:153
		{
			yy_ctrl.Host = yyDollar[2].text
		}
	case 36:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:154
		{
			yy_ctrl.Metrics = yy_as
			yy_as = make([]string, 0)
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:157
		{
			yy_ss2[yyDollar[1].text] = yyDollar[2].text
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:158
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:159
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:160
		{
			yy.include(yyDollar[2].text)
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:164
		{
			yy_agent = &ConfAgent{}
			yy_agent.Configer.Set(APP_CONF_DEFAULT, ConfDefault["agent"])
			yy_agent.Name = yyDollar[2].text
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:172
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:176
		{
			yy_agent.Disabled = yyDollar[2].b
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:177
		{
			yy_agent.Debug = 1
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:178
		{
			yy_agent.Debug = yyDollar[2].num
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:179
		{
			yy_agent.Host = yyDollar[2].text
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:180
		{
			yy_ss2[yyDollar[1].text] = yyDollar[2].text
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:181
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:182
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 52:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:183
		{
			yy_ss2["upstream"] = yyDollar[2].text
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:184
		{
			yy.include(yyDollar[2].text)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:189
		{
			yy_loadbalance = &ConfLoadbalance{}
			yy_loadbalance.Configer.Set(APP_CONF_DEFAULT, ConfDefault["loadbalance"])
			yy_loadbalance.Name = yyDollar[2].text
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:197
		{
			yy_loadbalance.Disabled = yyDollar[2].b
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:198
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:202
		{
			yy_loadbalance.Debug = 1
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:203
		{
			yy_loadbalance.Debug = yyDollar[2].num
		}
	case 61:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:204
		{
			yy_loadbalance.Host = yyDollar[2].text
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:206
		{
			yy_ss2[yyDollar[1].text] = yyDollar[2].text
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:207
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:208
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:209
		{
			yy.include(yyDollar[2].text)
		}
	case 70:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line yyparse.y:218
		{
			yy_loadbalance_backend.Type = yyDollar[1].text
			yy_loadbalance_backend.Name = yyDollar[2].text
			if !yy_loadbalance_backend.Disabled || yy.debug {
				yy_loadbalance.Backend = append(yy_loadbalance.Backend, *yy_loadbalance_backend)
			}
			yy_loadbalance_backend = &LbBackend{}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:233
		{
			yy_loadbalance_backend.Disabled = yyDollar[2].b
		}
	case 75:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:234
		{
			yy_loadbalance_backend.Upstream = yy_ss
			yy_ss = make(map[string]string)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line yyparse.y:241
		{
			yy_backend = &ConfBackend{}
			yy_backend.Configer.Set(APP_CONF_DEFAULT, ConfDefault["backend"])
			yy_backend.Name = yyDollar[2].text
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:250
		{
			yy_backend.Disabled = yyDollar[2].b
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:251
		{
			if err := os.Chdir(yyDollar[2].text); err != nil {
				yy.Error(err.Error())
			}
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line yyparse.y:255
		{
			yy_backend.Debug = 1
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:256
		{
			yy_backend.Debug = yyDollar[2].num
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:257
		{
			yy_backend.Host = yyDollar[2].text
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:259
		{
			yy_ss2[yyDollar[1].text] = yyDollar[2].text
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:260
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%d", yyDollar[2].num)
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:261
		{
			yy_ss2[yyDollar[1].text] = fmt.Sprintf("%v", yyDollar[2].b)
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:262
		{
			yy.include(yyDollar[2].text)
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line yyparse.y:269
		{
			yy_backend.Migrate.Disabled = yyDollar[2].b
		}
	case 93:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line yyparse.y:270
		{
			yy_backend.Migrate.Upstream = yy_ss
			yy_ss = make(map[string]string)
		}
	}
	goto yystack /* stack new state and value */
}
