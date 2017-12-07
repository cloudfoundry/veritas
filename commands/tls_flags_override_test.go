package commands_test

// import (
// 	"code.cloudfoundry.org/cfdot/commands"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	"github.com/onsi/gomega/gbytes"
// 	"github.com/spf13/cobra"
// )

// var _ = Describe("TLSFlagsOverride", func() {
// 	var validTLSFlags map[string]string
// 	var dummyCmd *cobra.Command
// 	var err error
// 	var output *gbytes.Buffer

// 	BeforeEach(func() {
// 		dummyCmd = &cobra.Command{
// 			Use: "dummy",
// 			Run: func(cmd *cobra.Command, args []string) {},
// 		}
// 	})

// 	JustBeforeEach(func() {
// 		commands.AddTLSFlags(dummyCmd)
// 		err = dummyCmd.PreRunE(dummyCmd, dummyCmd.Flags().Args())
// 	})

// 	Describe("BBSFlags override", func() {
// 		BeforeEach(func() {
// 			commands.AddBBSFlags(dummyCmd)
// 			output = gbytes.NewBuffer()
// 			dummyCmd.SetOutput(output)

// 			validTLSFlags = map[string]string{
// 				"--bbsSkipCertVerify": "false",
// 				"--bbsURL":            "https://example.com",
// 				"--bbsCACertFile":     "fixtures/bbsCACert.crt",
// 				"--bbsCertFile":       "fixtures/bbsClient.crt",
// 				"--bbsKeyFile":        "fixtures/bbsClient.key",
// 			}
// 		})

// 		It("prefers clientCertFile over bbsCertFile", func() {
// 			tlsFlags := map[string]string{
// 				"--clientCertFile": "fixtures/clientClient.crt",
// 			}
// 			parseFlagsErr := dummyCmd.ParseFlags(mergeFlags(validTLSFlags, tlsFlags))
// 			Expect(parseFlagsErr).NotTo(HaveOccurred())
// 			Expect(commands.BBSClientConfig.CertFile).To(Equal("fixtures/clientClient.crt"))
// 		})

// 	})

// })
