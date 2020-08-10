package types

type Config struct {
	Global struct {
		KubeCfgPath     string
		TimeStartApp    int
		LevelLogs       int
		InterfaceBridge string
	}
}
