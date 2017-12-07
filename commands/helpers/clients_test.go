package helpers_test

import (
	"code.cloudfoundry.org/bbs/fake_bbs"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("ClientConfig", func() {
	var (
		fakeBBSClient *fake_bbs.FakeClient
	)

	BeforeEach(func() {
		fakeBBSClient = &fake_bbs.FakeClient{}
	})

})
