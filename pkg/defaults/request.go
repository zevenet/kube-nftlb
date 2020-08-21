package defaults

import (
	"github.com/zevenet/kube-nftlb/pkg/env"
	"github.com/zevenet/kube-nftlb/pkg/types"
)

var nftlbKey = env.GetString("NFTLB_KEY")

// GetHeader returns the header for nftlb requests.
func GetHeader() *types.Header {
	return &types.Header{
		Key:   "Key",
		Value: nftlbKey,
	}
}
