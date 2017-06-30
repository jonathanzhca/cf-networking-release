package integration_test

import (
	"io"
	"io/ioutil"
	"log-transformer/config"
	"os"
	"os/exec"
	"path/filepath"
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
		session   *gexec.Session
		conf      config.LogTransformer
		inputFile *os.File
	)

	BeforeEach(func() {
		inputFile, _ = ioutil.TempFile("", "")
		outputDir, _ := ioutil.TempDir("", "")
		conf = config.LogTransformer{
			InputFile:       inputFile.Name(),
			OutputDirectory: outputDir,
		}
		configFilePath := WriteConfigFile(conf)

		var err error
		logTransformerCmd := exec.Command(binaryPath, "-config-file", configFilePath)
		session, err = gexec.Start(logTransformerCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		outputFile = filepath.Join(outputDir, "iptables.log")
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

	It("writes the input data to the output file", func() {
		go WriteLines(5, inputFile)

		Eventually(outputFile).Should(BeAnExistingFile())
		Eventually(ReadOutput, "5s").Should(Equal("12345"))
	})
})

func WriteLines(n int, w io.Writer) {
	defer GinkgoRecover()

	for i := 1; i <= n; i++ {
		time.Sleep(200 * time.Millisecond)
		_, err := w.Write([]byte("hello"))
		Expect(err).NotTo(HaveOccurred())
	}
}

func ReadOutput() string {
	bytes, err := ioutil.ReadFile(outputFile)
	Expect(err).NotTo(HaveOccurred())
	return string(bytes)
}
