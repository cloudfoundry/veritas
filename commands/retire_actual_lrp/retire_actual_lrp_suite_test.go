package retire_actual_lrp_command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRetireActualLrp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RetireActualLrp Suite")
}
