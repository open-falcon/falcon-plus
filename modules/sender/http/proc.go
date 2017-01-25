package http

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/sender/proc"
	"net/http"
)

func configProcRoutes() {

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("sms:%v, mail:%v", proc.GetSmsCount(), proc.GetMailCount())))
	})

}
