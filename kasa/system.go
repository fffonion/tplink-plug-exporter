package kasa

import "strings"

type KasaClientSystemService struct {
	c *KasaClient
}

type GetSysInfoRequest struct {
}

type GetSysInfoResponse struct {
	MAC        string `mapstructure:"mac"`
	Model      string `mapstructure:"model"`
	Alias      string `mapstructure:"alias"`
	Feature    string `mapstructure:"feature"`
	RelayState int    `mapstructure:"relay_state"`
	RSSI       int    `mapstructure:"rssi"`
	LEDOff     int    `mapstructure:"led_off"`
	OnTime     int    `mapstructure:"on_time"`
}

func (s *KasaClientSystemService) GetSysInfo() (*GetSysInfoResponse, error) {
	var response GetSysInfoResponse
	err := s.c.RPC("system", "get_sysinfo", GetSysInfoRequest{}, &response)

	return &response, err
}

func (s *KasaClientSystemService) EmeterSupported(r *GetSysInfoResponse) bool {
	return strings.Contains(r.Feature, "ENE")
}
