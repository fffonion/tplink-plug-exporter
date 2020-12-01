package kasa

import "strings"

type KasaClientSystemService struct {
	c   *KasaClient
	ctx *KasaRequestContext
}

type GetSysInfoRequest struct {
}

type SysInfoChildren struct {
	ID     string `mapstructure:"id"`
	State  int    `mapstructure:"state"`
	Alias  string `mapstructure:"alias"`
	OnTime int    `mapstructure:"on_time"`
}

type GetSysInfoResponse struct {
	MAC             string            `mapstructure:"mac"`
	Model           string            `mapstructure:"model"`
	Alias           string            `mapstructure:"alias"`
	Feature         string            `mapstructure:"feature"`
	RelayState      int               `mapstructure:"relay_state"`
	RSSI            int               `mapstructure:"rssi"`
	LEDOff          int               `mapstructure:"led_off"`
	OnTime          int               `mapstructure:"on_time"`
	DeviceID        string            `mapstructure:"deviceId"`
	SoftwareVersion string            `mapstructure:"sw_ver"`
	HardwareVersion string            `mapstructure:"hw_ver"`
	Children        []SysInfoChildren `mapstructure:"children"`
}

func (s *KasaClientSystemService) GetSysInfo() (*GetSysInfoResponse, error) {
	var response GetSysInfoResponse
	err := s.c.RPC("system", "get_sysinfo", s.ctx, GetSysInfoRequest{}, &response)

	return &response, err
}

func (s *KasaClientSystemService) EmeterSupported(r *GetSysInfoResponse) bool {
	return strings.Contains(r.Feature, "ENE")
}
