package defaults

import (
	"github.com/zevenet/kube-nftlb/pkg/env"
	"github.com/zevenet/kube-nftlb/pkg/types"
)

var (
	protocol = env.GetString("NFTLB_PROTOCOL")
	host     = env.GetString("NFTLB_HOST")
	port     = env.GetInt("NFTLB_PORT")
)

// GetURL returns an already configured URL, except for the path.
func GetURL() *types.URL {
	return &types.URL{
		Protocol: protocol,
		IP:       host,
		Port:     port,
	}
}
