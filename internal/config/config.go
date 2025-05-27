package config

type Config struct {
	LogLevel   LogLevel `name:"log-level" description:"Logging level for the application. One of debug, info, warn, or error" default:"info"`
	Scale      float64  `name:"scale" description:"Scale factor for the display." default:"2.0"`
	Fullscreen bool     `name:"fullscreen" description:"Enable fullscreen mode."`
	ROM        string   `name:"rom" description:"Path to the ROM file to load."`
	BIOS       string   `name:"bios" description:"Path to the BIOS file to load."`
}
