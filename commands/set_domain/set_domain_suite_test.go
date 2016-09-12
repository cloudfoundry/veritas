package set_domain_command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSetDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SetDomain Suite")
}
