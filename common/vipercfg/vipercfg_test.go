package vipercfg

import (
	"os"
	"fmt"
	"testing"
	"github.com/spf13/pflag"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestViperLoaderSuite struct{}

var _ = Suite(&TestViperLoaderSuite{})

// Tests the parsing of command line
//
// WARNING: This test would be conflict with TestBuildFacadeConfig()
func (suite *TestViperLoaderSuite) TestMustParseCmd(c *C) {
	viper := NewOwlConfigLoader().MustParseCmd()

	pflag.Set("version", "true")
	pflag.Set("config", "sample-1.json")

	c.Assert(viper.GetBool("version"), Equals, true)
	c.Assert(viper.GetString("config"), Equals, "sample-1.json")
}

// Tests the loading of config file
func (suite *TestViperLoaderSuite) TestMustLoadConfigFile(c *C) {
	testedLoader := NewOwlConfigLoader()
	testedLoader.MustParseCmd()

	pwd, _ := os.Getwd()
	pflag.Set("config", fmt.Sprintf("%s/sample-config.json", pwd))

	testedViper := testedLoader.MustLoadConfigFile()

	c.Assert(testedViper.GetString("db.username"), Equals, "hello")
	c.Assert(testedViper.GetString("db.password"), Equals, "nice")
}

// Tests the loading of combined configurations
//
// WARNING: This test would be conflict with TestMustParseCmd()
func (suite *TestViperLoaderSuite) TestBuildFacadeConfig(c *C) {
	testedLoader := NewOwlConfigLoader()
	testedLoader.FlagDefiner = func() {
		OwlDefaultPflagDefiner()
		pflag.String("target_name", "cc-tg", "new target")
	}
	testedLoader.MustParseCmd()

	pwd, _ := os.Getwd()
	pflag.Set("config", fmt.Sprintf("%s/sample-config.json", pwd))
	pflag.Set("target_name", "true-tg-01")

	testedViper, err := testedLoader.BuildFacadeConfig()
	c.Assert(err, IsNil)

	c.Assert(testedViper.GetString("db.username"), Equals, "hello")
	c.Assert(testedViper.GetString("db.password"), Equals, "nice")
	c.Assert(testedViper.GetString("target_name"), Equals, "true-tg-01")
}
