package integration_test

import (
	"io"
	"io/ioutil"
	"lib/datastore"
	"lib/filelock"
	"lib/serial"
	"log-transformer/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	outputDir  string
	outputFile string
)

var _ = Describe("Integration", func() {
	var (
		session               *gexec.Session
		conf                  config.LogTransformer
		kernelLogFile         *os.File
		containerMetadataFile *os.File
		store                 *datastore.Store
	)

	BeforeEach(func() {
		kernelLogFile, _ = ioutil.TempFile("", "")
		containerMetadataFile, _ = ioutil.TempFile("", "")
		outputDir, _ := ioutil.TempDir("", "")
		conf = config.LogTransformer{
			KernelLogFile:         kernelLogFile.Name(),
			ContainerMetadataFile: containerMetadataFile.Name(),
			OutputDirectory:       outputDir,
		}
		configFilePath := WriteConfigFile(conf)

		var err error
		logTransformerCmd := exec.Command(binaryPath, "-config-file", configFilePath)
		session, err = gexec.Start(logTransformerCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		outputFile = filepath.Join(outputDir, "iptables.log")

		store = &datastore.Store{
			Serializer: &serial.Serial{},
			Locker:     filelock.NewLocker(containerMetadataFile.Name()),
		}
		AddToContainerMetadata(store, "container-handle-1-longer-than-29-chars", "10.255.0.1", map[string]interface{}{
			"org_id":          "organization_id_1",
			"space_id":        "space_id_1",
			"app_id":          "app_id_1",
			"policy_group_id": "policy_group_id_1",
		})
		AddToContainerMetadata(store, "container-handle-2-longer-than-29-chars", "10.255.0.2", map[string]interface{}{
			"org_id":          "organization_id_1",
			"space_id":        "space_id_2",
			"app_id":          "app_id_2",
			"policy_group_id": "policy_group_id_2",
		})
	})

	AfterEach(func() {
		session.Interrupt()
		Eventually(session, DEFAULT_TIMEOUT).Should(gexec.Exit())
	})

	It("should log when starting", func() {
		Eventually(session.Out).Should(gbytes.Say("cfnetworking.log-transformer.*starting"))
	})

	It("should run as a daemon", func() {
		Consistently(session, DEFAULT_TIMEOUT).ShouldNot(gexec.Exit())
	})

	FIt("logs data about packets", func() {
		By("logging successful egress packets")
		go AddToKernelLog("Jun 28 18:21:24 localhost kernel: [100471.222018] OK_container-handle-1-longer IN=s-010255178004 OUT=eth0 MAC=aa:aa:0a:ff:b2:04:ee:ee:0a:ff:b2:04:08:00 SRC=10.255.0.1 DST=10.10.10.10 LEN=29 TOS=0x00 PREC=0x00 TTL=63 ID=2806 DF PROTO=UDP SPT=36556 DPT=11111 LEN=9 MARK=0x1\n", kernelLogFile)
		Eventually(outputFile).Should(BeAnExistingFile())
		Eventually(ReadLines, "5s").Should(ContainElement(MatchJSON(`{
			"source": "log-transformer",
			"message": "egress-allowed",
			"log_level": "1",
			"data": {
				"source": {
					"container_id": "container-handle-1-longer-than-29-chars",
					"app_guid": "app_id_1",
					"space_guid": "space_id_1",
					"organization_guid": "organization_id_1"
				},
				"packet": {
					"src_ip": "10.255.0.1",
					"src_port": 36556,
					"dst_ip": "10.10.10.10",
					"dst_port": 11111,
					"protocol": "udp",
					"mark": "0x1"
				}
			}
		}`)))

		By("logging denied egress packets")
		// go AddToKernelLog("Jun 30 16:07:06 localhost kernel: [265213.303412] DENY_container-handle-1-long IN=s-010255095010 OUT=eth0 MAC=aa:aa:0a:ff:5f:0a:ee:ee:0a:ff:5f:0a:08:00 SRC=10.255.0.1 DST=10.10.10.10 LEN=30 TOS=0x00 PREC=0x00 TTL=63 ID=2535 DF PROTO=UDP SPT=45564 DPT=25555 LEN=10 MARK=0x1", kernelLogFile)

	})
})

func AddToContainerMetadata(store *datastore.Store, containerID, containerIP string, metadata map[string]interface{}) {
	err := store.Add(containerID, containerIP, metadata)
	Expect(err).NotTo(HaveOccurred())
}
func AddToKernelLog(line string, w io.Writer) {
	defer GinkgoRecover()

	time.Sleep(200 * time.Millisecond)
	_, err := w.Write([]byte(line))
	// Jun 28 18:21:24 localhost kernel: [100471.222018] OK_fc2901d5-f631-40f0-6222-8 IN=s-010255178004 OUT=eth0 MAC=aa:aa:0a:ff:b2:04:ee:ee:0a:ff:b2:04:08:00 SRC=10.255.178.4 DST=10.10.10.10 LEN=29 TOS=0x00 PREC=0x00 TTL=63 ID=2806 DF PROTO=UDP SPT=36556 DPT=11111 LEN=9
	// Jun 28 18:21:30 localhost kernel: [100477.151949] DENY_fc2901d5-f631-40f0-6222 IN=s-010255178004 OUT=eth0 MAC=aa:aa:0a:ff:b2:04:ee:ee:0a:ff:b2:04:08:00 SRC=10.255.178.4 DST=10.10.10.10 LEN=29 TOS=0x00 PREC=0x00 TTL=63 ID=3568 DF PROTO=UDP SPT=54137 DPT=22222 LEN=9
	Expect(err).NotTo(HaveOccurred())
}

func ReadLines() []string {
	return strings.Split(ReadOutput(), "\n")
}
func ReadOutput() string {
	bytes, err := ioutil.ReadFile(outputFile)
	Expect(err).NotTo(HaveOccurred())
	return string(bytes)
}
