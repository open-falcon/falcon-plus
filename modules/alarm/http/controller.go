package http

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/toolkits/file"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Version() {
	this.Ctx.WriteString(g.VERSION)
}

func (this *MainController) Health() {
	this.Ctx.WriteString("ok")
}

func (this *MainController) Workdir() {
	this.Ctx.WriteString(fmt.Sprintf("%s", file.SelfDir()))
}

func (this *MainController) ConfigReload() {
	remoteAddr := this.Ctx.Request.RemoteAddr
	if strings.HasPrefix(remoteAddr, "127.0.0.1") {
		g.ParseConfig(g.ConfigFile)
		this.Data["json"] = g.Config()
		this.ServeJSON()
	} else {
		this.Ctx.WriteString("no privilege")
	}
}

func (this *MainController) Index() {
	events := g.Events.Clone()

	defer func() {
		this.Data["Now"] = time.Now().Unix()
		this.TplName = "index.html"
	}()

	count := len(events)
	if count == 0 {
		this.Data["Events"] = []*g.EventDto{}
		return
	}

	// 按照持续时间排序
	beforeOrder := make([]*g.EventDto, count)
	i := 0
	for _, event := range events {
		beforeOrder[i] = event
		i++
	}

	sort.Sort(g.OrderedEvents(beforeOrder))
	this.Data["Events"] = beforeOrder
}

func (this *MainController) Solve() {
	ids := this.GetString("ids")
	if ids == "" {
		this.Ctx.WriteString("")
		return
	}

	idArr := strings.Split(ids, ",,")
	for i := 0; i < len(idArr); i++ {
		g.Events.Delete(idArr[i])
	}

	this.Ctx.WriteString("")
}
