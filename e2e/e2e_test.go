package e2e

import (
	"flag"
	"testing"

	"github.com/alcideio/iskan/e2e/framework"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/spf13/viper"
)

func ViperizeFlags() {

	flag.Parse()

	//fmt.Println(flag.Args())
	//fmt.Println(pretty.Sprint(framework.GlobalConfig))

	// Part 2: Set Viper provided flags.
	// This must be done after common flags are registered, since Viper is a flag option.
	//viper.SetConfigName(TestContext.Viper)
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.Unmarshal(&framework.GlobalConfig)

	//AfterReadingAllFlags(&TestContext)
}

func init() {
	framework.RegisterFrameworkFlags()
}

func TestE2E(t *testing.T) {
	ViperizeFlags()
	gomega.RegisterFailHandler(ginkgo.Fail)
	RunE2ETests(t)
}
