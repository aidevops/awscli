package logger_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo/reporters"

	"github.com/aidevops/awscli/logger"
)

var suite = "Logger Test Suite"

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv("TEAMCITY") == "true" {
		RunSpecsWithCustomReporters(t, suite, []Reporter{reporters.NewTeamCityReporter(os.Stdout)})
	} else {
		RunSpecs(t, suite)
	}
}

var _ = Describe(suite, func() {

	var (
		log      *logger.Logger
		logLevel string
		logFile  string
		context  string
		format   string
	)

	BeforeEach(func() {
		// mock logging options
		logLevel = "error"
		logFile = "/dev/null"
		context = "test"
		format = "text"
		log = logger.NewLogger(logLevel, logFile, context, format)
	})

	Describe("Logger", func() {
		Context("Creation of a new instance", func() {
			It("Can create and populate new Logger instance", func() {
				Expect(log).NotTo(BeNil())
			})
			It("Can add a context", func() {
				log.AddContext("foo")
				Expect(len(log.Context)).To(Equal(2))
			})
			It("Can remove a context", func() {
				log.AddContext("foo")
				removed := log.RemoveContext()
				Expect(removed).To(Equal("foo"))
				Expect(len(log.Context)).To(Equal(1))
			})
			It("Can Debugf a context", func() {
				log.SetLevel("debug")
				log.Debugf("debug test %s=%d\n", "zero", 0)
			})
			It("Can Infof a context", func() {
				log.SetLevel("info")
				log.Infof("info test %s=%d\n", "zero", 0)
			})
			It("Can Warnf a context", func() {
				log.SetLevel("warn")
				log.Warnf("warn test %s=%d\n", "zero", 0)
			})
			It("Can Errorf a context", func() {
				log.SetLevel("error")
				log.Errorf("error test %s=%d\n", "zero", 0)
			})
		})
	})
})
