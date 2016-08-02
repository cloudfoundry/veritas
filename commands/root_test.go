package commands_test

import (
	"code.cloudfoundry.org/cfdot/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Root", func() {
	var dummyCmd *cobra.Command
	var validFlags []string
	var err error

	BeforeEach(func() {
		dummyCmd = &cobra.Command{
			Use: "dummy",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		commands.AddBBSFlags(dummyCmd)

		validFlags = []string{"--bbsURL", "http://example.com"}
	})

	JustBeforeEach(func() {
		err = dummyCmd.PreRunE(dummyCmd, dummyCmd.Flags().Args())
	})

	Context("when the --bbsURL isn't given", func() {
		BeforeEach(func() {
			dummyCmd.ParseFlags(removeFlag(validFlags, "--bbsURL"))
		})

		It("returns an error message", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("BBS URL not set"))
		})

		// It("returns an exit code of 3", func() {
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err.Error()).Should(ContainSubstring("BBS URL not set"))
		// })
	})
})

func removeFlag(flags []string, toRemove string) []string {
	var flagIdx int
	for idx := range flags {
		if flags[idx] == toRemove {
			flagIdx = idx
			break
		}
	}

	return append(flags[:flagIdx], flags[flagIdx+2:]...)
}
