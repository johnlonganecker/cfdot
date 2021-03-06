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

var _ = Describe("desired-lrps", func() {
	var sess *gexec.Session

	itValidatesBBSFlags("desired-lrps")
	itHasNoArgs("desired-lrps", false)

	Context("when no filters are passed", func() {
		BeforeEach(func() {
			bbsServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/desired_lrps/list.r2"),
					ghttp.RespondWithProto(200, &models.DesiredLRPsResponse{
						DesiredLrps: []*models.DesiredLRP{
							{
								Instances: 1,
							},
						},
					}),
				),
			)
		})

		It("returns the json encoding of the desired lrp scheduling info", func() {
			cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "desired-lrps")

			var err error
			sess, err = gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess).Should(gexec.Exit(0))
			Expect(sess.Out).To(gbytes.Say(`"instances":1`))
		})
	})

	Context("when passing filters", func() {
		BeforeEach(func() {
			bbsServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/desired_lrps/list.r2"),
					ghttp.VerifyProtoRepresenting(&models.DesiredLRPsRequest{
						Domain: "cf-apps",
					}),
					ghttp.RespondWithProto(200, &models.DesiredLRPsResponse{
						DesiredLrps: []*models.DesiredLRP{
							{
								Instances: 1,
							},
						},
					}),
				),
			)
		})

		Context("when -d is used as a filter flag", func() {
			It("exits with a status code of 0", func() {
				cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "desired-lrps", "-d", "cf-apps")
				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))
			})
		})

		Context("when --domain is used as a filter flag", func() {
			It("exits with a status code of 0", func() {
				cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "desired-lrps", "--domain", "cf-apps")
				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))
			})
		})

		Context("when --domain and -d are supplied as filter flags", func() {
			It("exits with a status code of 3", func() {
				cfdotCmd := exec.Command(cfdotPath, "--bbsURL", bbsServer.URL(), "desired-lrps", "--domain", "cf-apps", "-d", "cf-apps")
				sess, err := gexec.Start(cfdotCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(3))
			})
		})
	})
})
