package desired_lrp_scheduling_infos_command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDesiredLrpSchedulingInfos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DesiredLrpSchedulingInfos Suite")
}
