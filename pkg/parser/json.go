package parser

import (
	"encoding/json"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/types"
)

// NftlbAsJSON parses a given Nftlb struct and returns a JSON string that can be interpreted by nftlb.
func NftlbAsJSON(data *types.Nftlb) (string, error) {
	indentedJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}
	// Add spaces before every :, it's required by nftlb
	return strings.ReplaceAll(string(indentedJSON), "\":", "\" :"), nil
}
