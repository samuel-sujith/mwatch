package types

import (
	"github.com/go-kit/kit/log"
)

//Cfg is the configuration for the listener
type Cfg struct {
	Listenaddress string
	DesiredMetric string
}

//Interimconfig takes all config in one place for usage
type Interimconfig struct {
	Configuration       Cfg
	Logger              log.Logger
	Cert                string
	Key                 string
	SkipServerCertCheck bool
}
