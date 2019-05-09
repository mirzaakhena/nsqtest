package config

import (
	"github.com/spf13/viper"
)

// IConfig is
type IConfig interface {
	//
	GetNSQdURL() string //

}

// RealtimeConfig is
type RealtimeConfig struct {
	NSQdURL string
}

// NewRealtimeConfig is
func NewRealtimeConfig(confiName, path string) *RealtimeConfig {

	viper.SetConfigName(confiName)

	// TODO move it to the same folder as executable apps stay
	viper.AddConfigPath(path)
	// viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	rc := RealtimeConfig{}

	// reload all the config
	rc.reload()

	return &rc

}

func (r *RealtimeConfig) reload() {

	r.NSQdURL = viper.GetString("messagebroker.nsqd_url")

}

// GetNSQdURL is
func (r *RealtimeConfig) GetNSQdURL() string {
	return r.NSQdURL
}
