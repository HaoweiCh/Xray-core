package udp

import (
	"github.com/haoweich/xray-core/common"
	"github.com/haoweich/xray-core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
