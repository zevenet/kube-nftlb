package defaults

import (
	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// SetNftlbURL sets the default URL to communicate with nftlb
func SetNftlbURL(path string) *types.URL {
	return &types.URL{
		Protocol:  types.HTTP,
		IPversion: types.IPv4,
		IP:        types.LocalHostIPv4,
		Port:      types.NFTLBport,
		Path:      path,
	}
}
