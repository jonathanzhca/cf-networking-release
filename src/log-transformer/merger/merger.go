package merger

import (
	"code.cloudfoundry.org/lager"
)

type IPTablesLogData struct {
	Message string
	Data    lager.Data
}

type Merger struct {
}
