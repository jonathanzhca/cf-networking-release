package runner

import (
	"lib/datastore"
	"log-transformer/merger"
	"log-transformer/parser"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/hpcloud/tail"
)

//go:generate counterfeiter -o fakes/logMerger.go --fake-name LogMerger . logMerger
type logMerger interface {
	Merge(parser.ParsedData, datastore.Container, datastore.Container) (merger.IPTablesLogData, error)
}

//go:generate counterfeiter -o fakes/kernel_log_parser.go --fake-name KernelLogParser . kernelLogParser
type kernelLogParser interface {
	IsIPTablesLogData(line string) bool
	Parse(line string) parser.ParsedData
}

type Runner struct {
	Lines          chan *tail.Line
	Parser         kernelLogParser
	Merger         logMerger
	Store          datastore.Datastore
	Logger         lager.Logger
	IPTablesLogger lager.Logger
}

func (r *Runner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)
	r.Logger.Info("started", lager.Data{})
	for {
		select {
		case <-signals:
			r.Logger.Info("exited", lager.Data{})
			return nil
		case line := <-r.Lines:
			if line.Err != nil {
				r.Logger.Error("tail-kernel-logs", line.Err)
				continue
			}
			if r.Parser.IsIPTablesLogData(line.Text) {
				parsed := r.Parser.Parse(line.Text)
				containers, err := r.Store.ReadAll()
				if err != nil {
					panic(err)
				}
				src := containers[parsed.Source]
				dst := containers[parsed.Destination]
				merged, err := r.Merger.Merge(parsed, src, dst)
				if err != nil {
					r.Logger.Error("merge-kernel-logs", err)
					continue
				}
				r.IPTablesLogger.Info(merged.Message, merged.Data)
			}
		}
	}

	return nil
}
