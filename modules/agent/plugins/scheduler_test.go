package plugins

import (
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"io/ioutil"
	"os"
	"testing"
)

func TestProcessOutput(t *testing.T) {
	cfgFile, err := ioutil.TempFile(os.TempDir(), "falcon-plus.test.")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(cfgFile.Name())

	cfgFile.WriteString(`{ "debug": true}`)
	cfgFile.Close()

	g.ParseConfig(cfgFile.Name())

	plugin := Plugin{
		FilePath: "/tmp/test-plugin-file",
		MTime:    0,
		Cycle:    60,
	}

	{
		dataJson := `[{
	"metric": "cpu.idle",
	"tags": "_project=sre,number=1,core=2,test=测试",
	"value": 301,
	"timestamp": 1498712372,
	"judgeType": "G"
}, {
	"metric": "cpu.idle",
	"tags": "_project=sre,number=1,core=2,test=测试",
	"value": 302,
	"timestamp": 1498712372,
	"judgeType": "G"
}]`
		metrics := processOutput(&plugin, []byte(dataJson))
		if len(metrics) != 2 {
			t.Fatalf("processOutput failed")
		}
	}

	{
		dataLines := `
metric 102
metric tag1=v1,tag2=v2 104
metric tag1=v1,tag2=v2 105 1498712372
metric 106 1498712372
metric invalid-value 1498712372
metric 108 invalid-timestamp
{"metric": "cpu.idle", "value": 201, "timestamp": 1498712372, "judgeType": "G"}`
		metrics := processOutput(&plugin, []byte(dataLines))
		if len(metrics) != 5 {
			t.Log(metrics)
			t.Fatalf("processOutput failed")
		}
	}
}
