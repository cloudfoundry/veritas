package commands_test

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/bbs/fake_bbs"
	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/cfdot/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/spf13/cobra"
)

var _ = FDescribe("CreateDesiredLRP", func() {
	var (
		fakeBBSClient      *fake_bbs.FakeClient
		returnedError      error
		stdout, stderr     *gbytes.Buffer
		expectedDesiredLRP *models.DesiredLRP
		spec               []byte
		cmd                *cobra.Command
	)

	BeforeEach(func() {
		cmd = &cobra.Command{}
		fakeBBSClient = &fake_bbs.FakeClient{}
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()

		fakeBBSClient.DesireLRPReturns(returnedError)
		expectedDesiredLRP = &models.DesiredLRP{
			ProcessGuid: "some-desired-lrp",
		}
		var err error
		spec, err = json.Marshal(expectedDesiredLRP)
		Expect(err).NotTo(HaveOccurred())
	})

	It("creates the desired lrp", func() {
		err := commands.CreateDesiredLRP(cmd, stdout, stderr, fakeBBSClient, string(spec))
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeBBSClient.DesireLRPCallCount()).To(Equal(1))
		_, lrp := fakeBBSClient.DesireLRPArgsForCall(0)
		Expect(lrp).To(Equal(expectedDesiredLRP))
	})

	Context("when a file is passed as an argument", func() {
		var filename string

		BeforeEach(func() {
			f, err := ioutil.TempFile(os.TempDir(), "spec_file")
			Expect(err).NotTo(HaveOccurred())
			defer f.Close()
			Expect(f.Write(spec)).To(Succeed())
			filename = f.Name()
		})

		It("creates the desired lrp", func() {
			err := commands.CreateDesiredLRP(cmd, stdout, stderr, fakeBBSClient, "@"+filename)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeBBSClient.DesireLRPCallCount()).To(Equal(1))
			_, lrp := fakeBBSClient.DesireLRPArgsForCall(0)
			Expect(lrp).To(Equal(expectedDesiredLRP))
		})
	})

	Context("when the bbs errors", func() {
		BeforeEach(func() {
			fakeBBSClient.DesireLRPReturns(models.ErrUnknownError)
		})

		It("fails with a relevant error", func() {
			err := commands.CreateDesiredLRP(cmd, stdout, stderr, fakeBBSClient, "{}")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(models.ErrUnknownError))
		})
	})
})
