package cron

import (
    "github.com/open-falcon/falcon-plus/modules/agent/g"
    "github.com/open-falcon/falcon-plus/modules/agent/scripts"
    "log"
    "time"
)

func SyncMineScripts() {
    if !g.Config().Script.Enabled {
        return
    }

    go syncMineScripts()
}

func syncMineScripts() {
    duration := time.Duration(g.Config().Script.ReloadCfgFilesSec) * time.Second
    for {
        time.Sleep(duration)

        desiredAll := make(map[string]*scripts.Script)
        ps := scripts.ListScripts()
        for k, p := range ps {
            desiredAll[k] = p
        }

        scripts.DelNoUseScripts(desiredAll)
        scripts.AddNewScripts(desiredAll)

        if g.Config().Debug {
            log.Printf("current scripts:%v\n", scripts.Scripts)
        }
    }
}
