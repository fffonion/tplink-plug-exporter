package kasa

type KasaClientEmeterService struct {
	c *KasaClient
}

type GetRealtimeRequest struct {
}

type GetRealtimeResponse struct {
	Current float64 `mapstructure:"current_ma"`
	Voltage float64 `mapstructure:"voltage_mv"`
	Power   float64 `mapstructure:"power_mw"`
	Total   float64 `mapstructure:"total_wh"`
}

func (s *KasaClientEmeterService) GetRealtime() (*GetRealtimeResponse, error) {
	var response GetRealtimeResponse
	err := s.c.RPC("emeter", "get_realtime", GetRealtimeRequest{}, &response)

	return &response, err
}
