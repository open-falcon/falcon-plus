package funcs

import "testing"

func TestMatchProc(t *testing.T){
	data:=[]byte("1082 (bash) S 1081 1082 1082 8912896 -1 0 2704 2704 0 0 15 93 15 93 20 0 0 0 96121849 7487488 2521 345")
	procInfo:=matchProc(data,"bash")
	if procInfo==nil{
		t.Error("matchProc function incorrect")
	}
}
