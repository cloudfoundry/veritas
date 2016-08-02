package main_test

import (
	"os"
	"os/exec"

	"code.cloudfoundry.org/bbs/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("domains", func() {
	Context("when the server responds with domains", func() {
		BeforeEach(func() {
			bbsServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/domains/list"),
					ghttp.RespondWithProto(200, &models.DomainsResponse{
						Error:   nil,
						Domains: []string{"domain-1", "domain-2"},
					}),
				),
			)
		})

		It("prints a json stream of all the domains", func() {
			cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "domains")

			sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-sess.Exited
			Expect(sess.ExitCode()).To(Equal(0))

			Expect(sess.Out).To(gbytes.Say(`"domain-1"\n"domain-2"\n`))
		})
	})

	Context("when the server doesn't respond", func() {
		It("fails with a relevant error message", func() {
			cfdotCmd := exec.Command(cfdotPath, "--bbsURL", "http://127.1.1.1:1", "domains")

			sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-sess.Exited
			Expect(sess.ExitCode()).To(Equal(1))

			Expect(sess.Err).To(gbytes.Say("(error)|(connection refused)"))
		})
	})

	Context("when connecting to a non TLS server", func() {
		Describe("flag parsing", func() {
			BeforeEach(func() {
				bbsServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v1/domains/list"),
						ghttp.RespondWithProto(200, &models.DomainsResponse{}),
					),
				)
			})

			Context("when the URL is HTTP", func() {
				It("ignores bbs TLS related flags", func() {
					cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "domains", "--bbsCACertFile", "invalid/path/to/bbs/ca/cert/file", "--bbsSkipCertVerify", "--bbsCertFile", "invalid/path/to/bbs/cert/file", "--bbsKeyFile", "invalid/path/to/key/file")

					sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					<-sess.Exited
					Expect(sess.ExitCode()).To(Equal(0))
				})

				It("ignores bbs TLS related environment variables", func() {
					os.Setenv("BBS_URL", bbsServer.URL())
					os.Setenv("BBS_CA_CERT_FILE", "invalid/path/to/bbs/ca/cert/file")
					os.Setenv("BBS_SKIP_CERT_VERIFY", "true")
					os.Setenv("BBS_CERT_FILE", "invalid/path/to/bbs/cert/file")
					os.Setenv("BBS_KEY_FILE", "invalid/path/to/key/file")
					defer os.Unsetenv("BBS_URL")
					defer os.Unsetenv("BBS_CA_CERT_FILE")
					defer os.Unsetenv("BBS_SKIP_CERT_VERIFY")
					defer os.Unsetenv("BBS_CERT_FILE")
					defer os.Unsetenv("BBS_KEY_FILE")

					cfdotCmd := exec.Command(cfdotPath, "domains")

					sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					<-sess.Exited
					Expect(sess.ExitCode()).To(Equal(0))
				})
			})

			It("fails when the URL is not HTTP or HTTPS", func() {
				cfdotCmd := exec.Command(cfdotPath, "--bbsURL", "nohttp.com", "domains")

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(3))

				Expect(sess.Err).To(gbytes.Say(
					"The URL 'nohttp.com' does not have an 'http' or 'https' scheme. Please specify one with the '--bbsURL' flag or the 'BBS_URL' environment variable.",
				))
				Expect(sess.Err).To(gbytes.Say("List fresh domains"))
				Expect(sess.Err).To(gbytes.Say("Usage:"))
			})

			It("fails when specifying a non-empty, invalid bbsURL", func() {
				cfdotCmd := exec.Command(cfdotPath, "domains", "--bbsURL", ":")

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(3))

				Expect(sess.Err).To(gbytes.Say(
					"The value ':' is not a valid BBS URL. Please specify one with the '--bbsURL' flag or the 'BBS_URL' environment variable.",
				))
				Expect(sess.Err).To(gbytes.Say("List fresh domains"))
				Expect(sess.Err).To(gbytes.Say("Usage:"))
			})

			It("fails when not specifying a bbs URL", func() {
				cfdotCmd := exec.Command(cfdotPath, "domains")

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(3))

				Expect(sess.Err).To(gbytes.Say(
					"BBS URL not set. Please specify one with the '--bbsURL' flag or the 'BBS_URL' environment variable.",
				))
				Expect(sess.Err).To(gbytes.Say("List fresh domains"))
				Expect(sess.Err).To(gbytes.Say("Usage:"))
			})

			It("works with a --bbsURL flag specified before domains", func() {
				cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "domains")

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(0))
			})

			It("works with a --bbsURL flag specified after domains", func() {
				cfdotCmd := exec.Command(cfdotPath, "domains", "--bbsURL", bbsServer.URL())

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(0))
			})

			It("works with a BBS_URL environment variable", func() {
				os.Setenv("BBS_URL", bbsServer.URL())
				defer os.Unsetenv("BBS_URL")

				cfdotCmd := exec.Command(cfdotPath, "domains")

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(0))
			})

			Context("prefers flags over environment", func() {
				It("for bbsURL flag --bbsURL flag over the environment variable", func() {
					os.Setenv("BBS_URL", "broken url")
					defer os.Unsetenv("BBS_URL")

					cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "domains")

					sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					<-sess.Exited
					Expect(sess.ExitCode()).To(Equal(0))
				})
			})
		})
	})

	Context("when connecting to a TLS server", func() {
		Describe("flag parsing", func() {
			BeforeEach(func() {
				bbsTLSServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v1/domains/list"),
						ghttp.RespondWithProto(200, &models.DomainsResponse{}),
					),
				)
			})

			Context("when validating TLS related flags", func() {
				It("verifies each FILE flag is associated with a readable file", func() {

				})
				It("verifies if one of BBS client cert or client key files is set, the other must be as well", func() {
				})
				It("verifies if the skip-cert-verfy option is not enabled, the CA cert file must be provided", func() {
				})
			})

			It("prefers flags over environment", func() {
				os.Setenv("BBS_URL", "broken url")
				os.Setenv("BBS_CA_CERT_FILE", "invalid/path/to/bbs/ca/cert/file")
				os.Setenv("BBS_SKIP_CERT_VERIFY", "false")
				os.Setenv("BBS_CERT_FILE", "invalid/path/to/bbs/cert/file")
				os.Setenv("BBS_KEY_FILE", "invalid/path/to/key/file")

				defer os.Unsetenv("BBS_URL")
				defer os.Unsetenv("BBS_CA_CERT_FILE")
				defer os.Unsetenv("BBS_SKIP_CERT_VERIFY")
				defer os.Unsetenv("BBS_CERT_FILE")
				defer os.Unsetenv("BBS_KEY_FILE")

				cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsTLSServer.URL(), "domains", "--bbsCACertFile")

				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(0))
			})
		})
	})
})
