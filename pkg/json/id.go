package json

import (
	"sync"
)

// ID is a hash table for backends, corresponding to [farm-name]backend-number.
var (
	ID = struct {
		sync.RWMutex
		farms map[string]int
	}{farms: make(map[string]int)}
)

// CreateFarmID initializes the key-value pair.
func CreateFarmID(farmName string) {
	ID.Lock()
	ID.farms[farmName] = 0
	ID.Unlock()
}

// IncreaseBackendID increases the backend ID given the farm name.
func IncreaseBackendID(farmName string) {
	ID.Lock()
	ID.farms[farmName]++
	ID.Unlock()
}

// DecreaseBackendID decreases the backend ID given the farm name.
func DecreaseBackendID(farmName string) {
	ID.Lock()
	ID.farms[farmName]--
	ID.Unlock()
}

// GetBackendID returns the backend ID given the farm name.
func GetBackendID(farmName string) int {
	ID.RLock()
	defer ID.RUnlock()
	return ID.farms[farmName]
}
