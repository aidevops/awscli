package awscli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo/reporters"

	"os"
	"testing"
)

func TestAwscli(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv("TEAMCITY") == "true" {
		RunSpecsWithCustomReporters(t, "Awscli Test Suite", []Reporter{reporters.NewTeamCityReporter(os.Stdout)})
	} else {
		RunSpecs(t, "Awscli Test Suite")
	}
}
