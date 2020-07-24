package defaults

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"os"

	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// key is not exportable, you must call SetNftlbKey to get a Header struct with the nftlb key.
var (
	key string
)

// Check if "app" has been executed with 1 argument only.
func Init() *types.Config {
	if len(os.Args) != 3 {
		err := fmt.Sprintf("Error: kube-nftlb expectes the nftlb key and the configuration file, args passed: %d", len(os.Args)-1)
		panic(err)
	}
	key = os.Args[1]
	fmt.Println("Key set")

	var cfg types.Config
	err := gcfg.ReadFileInto(&cfg, os.Args[2])
	if err != nil {
		panic(err)
	}

	return &cfg
}

// SetNftlbKey returns a Header with the KEY_NFTLB configured in build.sh.
func SetNftlbKey() *types.Header {
	return &types.Header{
		Key:   "Key",
		Value: key,
	}
}
