package main

import (
	"flag"
	"fmt"
	"log"
	"log-transformer/config"
	"os"
	"path/filepath"

	"github.com/hpcloud/tail"

	"code.cloudfoundry.org/lager"
)

var (
	logPrefix = "cfnetworking"
)

func main() {
	configFilePath := flag.String("config-file", "", "path to config file")
	flag.Parse()
	conf, err := config.New(*configFilePath)
	if err != nil {
		log.Fatalf("%s.log-transformer: reading config: %s", logPrefix, err)
	}

	logger := lager.NewLogger(fmt.Sprintf("%s.log-transformer", logPrefix))
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.DEBUG)
	logger.RegisterSink(sink)

	sink.SetMinLevel(lager.DEBUG)

	logger.Info("starting")

	file, err := os.Create(filepath.Join(conf.OutputDirectory, "iptables.log"))
	if err != nil {
		logger.Fatal("create-output-file", err)
	}

	t, err := tail.TailFile(conf.InputFile, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
		MustExist: true,
		Follow:    true,
		Poll:      true,
	})
	if err != nil {
		logger.Fatal("tail-input", err)
	}

	go func() {
		fmt.Println("*** hello ***")
		for {
			select {
			case line := <-t.Lines:
				file.Write([]byte(line.Text))
			}
		}
	}()

	done := make(chan struct{})
	<-done
}
