package s3_test

import (
	"log"
	"os"

	"github.com/kubernetes-csi/csi-test/v4/pkg/sanity"
	"github.com/majst01/csi-driver-s3/pkg/s3"
)

var _ = Describe("S3Driver", func() {

	Context("s3fs", func() {
		socket := "/tmp/csi-driver-s3fs.sock"
		csiEndpoint := "unix://" + socket
		if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
			Expect(err).NotTo(HaveOccurred())
		}
		driver, err := s3.New("test-node", csiEndpoint)
		if err != nil {
			log.Fatal(err)
		}
		go driver.Run()

		Describe("CSI sanity", func() {
			sanityCfg := &sanity.TestConfig{
				TargetPath:  os.TempDir() + "/s3fs-target",
				StagingPath: os.TempDir() + "/s3fs-staging",
				Address:     csiEndpoint,
				SecretsFile: "../../test/secret-minio.yaml",
				TestVolumeParameters: map[string]string{
					"mounter": "s3fs",
				},
			}
			sanity.GinkgoTest(sanityCfg)
		})
	})

})
