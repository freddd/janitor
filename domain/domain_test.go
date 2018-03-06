package domain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestCert(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CertSuite")
}

var _ = Describe("CertSuite", func() {
	It("", func() {
		cv := DomainVerifier{}
	})
})