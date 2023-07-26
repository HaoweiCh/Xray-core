package tcp

import (
	"github.com/haoweich/xray-core/common"
	"github.com/haoweich/xray-core/transport/internet"
)

const protocolName = "tcp"

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
