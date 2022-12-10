package config

import (
	"flag"
	"time"

	"example.com/webrtc-game/pkg/emulator/libretro/image"
	"example.com/webrtc-game/pkg/monitoring"
	"github.com/spf13/pflag"
)

const DefaultSTUNTURN = `[{"urls":"stun:stun-turn.webgame2d.com:3478"},{"urls":"turn:stun-turn.webgame2d.com:3478","username":"root","credential":"root"}]`
const CODEC_VP8 = "VP8"
const CODEC_H264 = "H264"

const AUDIO_RATE = 48000
const AUDIO_CHANNELS = 2
const AUDIO_MS = 20
const AUDIO_FRAME = AUDIO_RATE * AUDIO_MS / 1000 * AUDIO_CHANNELS

var FrontendSTUNTURN = flag.String("stunturn", DefaultSTUNTURN, "Frontend STUN TURN servers")
var Mode = flag.String("mode", "dev", "Environment")
var StunTurnTemplate = `[{"urls":"stun:stun.l.google.com:19302"},{"urls":"stun:%s:3478"},{"urls":"turn:%s:3478","username":"root","credential":"root"}]`
var HttpPort = flag.String("httpPort", "8000", "User agent port of the app")
var HttpsPort = flag.Int("httpsPort", 443, "Https Port")
var HttpsKey = flag.String("httpsKey", "", "Https Key")
var HttpsChain = flag.String("httpsChain", "", "Https Chain")

var WSWait = 20 * time.Second
var MatchWorkerRandom = false
var ProdEnv = "prod"
var StagingEnv = "staging"

const NumKeys = 10

var FileTypeToEmulator = map[string]string{
	"gba": "gba",
	"gbc": "gba",
	"cue": "pcsx",
	"zip": "mame",
	"nes": "nes",
	"smc": "snes",
	"sfc": "snes",
	"swc": "snes",
	"fig": "snes",
	"bs":  "snes",
	"n64": "n64",
	"v64": "n64",
	"z64": "n64",
}

// There is no good way to determine main width and height of the emulator.
// When game run, frame width and height can scale abnormally.
type EmulatorMeta struct {
	Path            string
	Config          string
	Width           int
	Height          int
	AudioSampleRate int
	Fps             float64
	BaseWidth       int
	BaseHeight      int
	Ratio           float64
	Rotation        image.Rotate
	IsGlAllowed     bool
	UsesLibCo       bool
	HasMultitap     bool
}

// 게임 경로를 변경처리함.
var EmulatorConfig = map[string]EmulatorMeta{
	"gba": {
		Path:   "/vagrant/sample/game/assets/emulator/libretro/cores/mgba_libretro",
		Width:  240,
		Height: 160,
	},
	"pcsx": {
		Path:   "/vagrant/sample/game/assets/emulator/libretro/cores/pcsx_rearmed_libretro",
		Width:  350,
		Height: 240,
	},
	"nes": {
		//Path:   "assets/emulator/libretro/cores/nestopia_libretro",
		//Path:   "/usr/local/share/cloud-game/assets/emulator/libretro/cores/nestopia_libretro",
		Path: "/vagrant/sample/game/assets/emulator/libretro/cores/nestopia_libretro",
		Width:  256,
		Height: 240,
	},
	"snes": {
		Path:   "/vagrant/sample/game/assets/emulator/libretro/cores/snes9x_libretro",
		Width:  256,
		Height: 224,
		HasMultitap: true,
	},
	"mame": {
		Path:   "/vagrant/sample/game/assets/emulator/libretro/cores/fbneo_libretro",
		Width:  240,
		Height: 160,
	},
	"n64": {
		Path:   "/vagrant/sample/game/assets/emulator/libretro/cores/mupen64plus_next_libretro",
		Config:   "/vagrant/sample/game/assets/emulator/libretro/cores/mupen64plus_next_libretro.cfg",
		Width:  320,
		Height: 240,
		IsGlAllowed: true,
		UsesLibCo: true,
	},
}

var EmulatorExtension = []string{".so", ".armv7-neon-hf.so", ".dylib", ".dll"}


type Config struct {
	//Port               int
	//CoordinatorAddress string
	//HttpPort           int

	// video
	Scale             int
	EnableAspectRatio bool
	Width             int
	Height            int
	Zone              string

	// WithoutGame to launch encoding with Game
	WithoutGame bool

	MonitoringConfig monitoring.ServerMonitoringConfig
}

func NewDefaultConfig() Config {
	return Config{

		Scale:              1,
		EnableAspectRatio:  false,
		Width:              320,
		Height:             240,
		WithoutGame:        false,
		Zone:               "",

		MonitoringConfig: monitoring.ServerMonitoringConfig{
			Port:          6601,
			URLPrefix:     "/worker",
			MetricEnabled: true,
		},
	}
}

func (c *Config) AddFlags(fs *pflag.FlagSet) *Config {
	/*
	fs.IntVarP(&c.Port, "port", "", 8800, "Worker server port")
	fs.StringVarP(&c.CoordinatorAddress, "coordinatorhost", "", c.CoordinatorAddress, "Worker URL to connect")
	fs.IntVarP(&c.HttpPort, "httpPort", "", c.HttpPort, "Set external HTTP port")
	fs.StringVarP(&c.Zone, "zone", "z", c.Zone, "Zone of the worker")

	fs.IntVarP(&c.Scale, "scale", "s", c.Scale, "Set output viewport scale factor")
	fs.BoolVarP(&c.EnableAspectRatio, "ar", "", c.EnableAspectRatio, "Enable Aspect Ratio")
	fs.IntVarP(&c.Width, "width", "w", c.Width, "Set custom viewport width")
	fs.IntVarP(&c.Height, "height", "h", c.Height, "Set custom viewport height")
	fs.BoolVarP(&c.WithoutGame, "wogame", "", c.WithoutGame, "launch worker with game")
	*/
	fs.BoolVarP(&c.MonitoringConfig.MetricEnabled, "monitoring.metric", "m", c.MonitoringConfig.MetricEnabled, "Enable prometheus metric for server")
	fs.BoolVarP(&c.MonitoringConfig.ProfilingEnabled, "monitoring.pprof", "p", c.MonitoringConfig.ProfilingEnabled, "Enable golang pprof for server")
	fs.IntVarP(&c.MonitoringConfig.Port, "monitoring.port", "", c.MonitoringConfig.Port, "Monitoring server port")
	fs.StringVarP(&c.MonitoringConfig.URLPrefix, "monitoring.prefix", "", c.MonitoringConfig.URLPrefix, "Monitoring server url prefix")

	return c
}
