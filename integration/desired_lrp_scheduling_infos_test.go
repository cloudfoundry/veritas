package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/bbs/models"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("desired-lrp-scheduling-infos", func() {
	var sess *gexec.Session

	BeforeEach(func() {
		bbsServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/v1/desired_lrp_scheduling_infos/list"),
				ghttp.RespondWithProto(200, &models.DesiredLRPSchedulingInfosResponse{
					DesiredLrpSchedulingInfos: []*models.DesiredLRPSchedulingInfo{
						{
							Instances: 1,
						},
					},
				}),
			),
		)
	})

	JustBeforeEach(func() {
		cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "desired-lrp-scheduling-infos")

		var err error
		sess, err = gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(sess.Exited).Should(BeClosed())
	})

	It("exits with status code of 0", func() {
		Expect(sess.ExitCode()).To(Equal(0))
	})

	It("returns the json encoding of the desired lrp scheduling info", func() {
		Expect(sess.Out).To(gbytes.Say(`"Instances":1`))
	})
})
