package parser

import (
	"strings"
)

type ParsedData struct {
	Source      string
	Destination string
}

type KernelLogParser struct {
}

func (k *KernelLogParser) IsIPTablesLogData(line string) bool {
	return strings.Contains(line, "OK_") || strings.Contains(line, "DENY_")

}
func (k *KernelLogParser) Parse(line string) map[string]interface{} {
	data := map[string]interface{}{}
	words := strings.Fields(line)
	for _, word := range words {
		if equalSignIndex := strings.Index(word, "="); equalSignIndex > -1 {
			key := word[:equalSignIndex]
			value := word[equalSignIndex+1:]
			if _, ok := data[key]; !ok {
				data[key] = value
			}
		}
	}
	return data
}
