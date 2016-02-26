package awscli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAwscli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Awscli Suite")
}
