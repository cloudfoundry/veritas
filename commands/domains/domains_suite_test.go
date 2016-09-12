package domains_command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDomains(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domains Suite")
}
