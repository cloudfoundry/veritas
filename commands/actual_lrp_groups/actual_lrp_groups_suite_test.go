package actual_lrp_groups_command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestActualLrpGroups(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ActualLrpGroups Suite")
}
