package parser_test

import (
	"log-transformer/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KernelLogParser", func() {
	var (
		kernelLogParser *parser.KernelLogParser
	)
	BeforeEach(func() {
		kernelLogParser = &parser.KernelLogParser{}

	})
	Describe("IsIPTablesLogData", func() {
		It("returns true if it contains OK_ or DENY_", func() {
			Expect(kernelLogParser.IsIPTablesLogData("stuff OK_ stuff")).To(BeTrue())
			Expect(kernelLogParser.IsIPTablesLogData("stuff DENY_ stuff")).To(BeTrue())
		})

		It("returns false", func() {
			Expect(kernelLogParser.IsIPTablesLogData("stuff stuff")).To(BeFalse())
		})
	})
	Describe("Parse", func() {
		It("returns the log line as parsed data", func() {
			Expect(kernelLogParser.Parse("Jun 28 18:21:24 localhost kernel: [100471.222018] OK_container-handle-1-longer IN=s-010255178004 OUT=eth0 MAC=aa:aa:0a:ff:b2:04:ee:ee:0a:ff:b2:04:08:00 SRC=10.255.0.1 DST=10.10.10.10 LEN=29 TOS=0x00 PREC=0x00 TTL=63 ID=2806 DF PROTO=UDP SPT=36556 DPT=11111 LEN=9 MARK=0x1")).To(Equal(map[string]interface{}{
				"IN":    "s-010255178004",
				"OUT":   "eth0",
				"MAC":   "aa:aa:0a:ff:b2:04:ee:ee:0a:ff:b2:04:08:00",
				"SRC":   "10.255.0.1",
				"DST":   "10.10.10.10",
				"LEN":   "29",
				"TOS":   "0x00",
				"PREC":  "0x00",
				"TTL":   "63",
				"ID":    "2806",
				"PROTO": "UDP",
				"SPT":   "36556",
				"DPT":   "11111",
				"MARK":  "0x1",
			}))
		})
		Context("when there is no parseable data", func() {
			It("returns an empty map", func() {
				Expect(kernelLogParser.Parse("stuff")).To(Equal(map[string]interface{}{}))
			})
		})
	})
})
