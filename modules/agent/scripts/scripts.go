package scripts

type Script struct {
	FilePath string
	MTime    int64
	Cycle    int
	Args     string
	ResultType string // text, json
}

var (
	Scripts              = make(map[string]*Script)
	ScriptsWithScheduler = make(map[string]*ScriptScheduler)
)

func DelNoUseScripts(newScripts map[string]*Script) {
	for currKey, currScript := range Scripts {
		newScript, ok := newScripts[currKey]
		if !ok || currScript.MTime != newScript.MTime {
			deleteScript(currKey)
		}
	}
}

func AddNewScripts(newScripts map[string]*Script) {
	for fpath, newScript := range newScripts {
		if _, ok := Scripts[fpath]; ok && newScript.MTime == Scripts[fpath].MTime {
			continue
		}

		Scripts[fpath] = newScript
		sch := NewScriptScheduler(newScript)
		ScriptsWithScheduler[fpath] = sch
		sch.Schedule()
	}
}

func ClearAllScripts() {
	for k := range Scripts {
		deleteScript(k)
	}
}

func deleteScript(key string) {
	v, ok := ScriptsWithScheduler[key]
	if ok {
		v.Stop()
		delete(ScriptsWithScheduler, key)
	}
	delete(Scripts, key)
}
