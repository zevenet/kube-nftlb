package defaults

import (
	"fmt"
	"os"

	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// key is not exportable, you must call SetNftlbKey to get a Header struct with the nftlb key.
var (
	key string
)

// Check if "app" has been executed with 1 argument only.
func init() {
	if len(os.Args) != 2 {
		err := fmt.Sprintf("Error: you must pass only 1 argument to be interpreted as the nftlb key, args passed: %d", len(os.Args)-1)
		panic(err)
	}
	key = os.Args[1]
}

// SetNftlbKey returns a Header with the KEY_NFTLB configured in build.sh.
func SetNftlbKey() *types.Header {
	return &types.Header{
		Key:   "Key",
		Value: key,
	}
}
