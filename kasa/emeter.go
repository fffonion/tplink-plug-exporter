package kasa

type KasaClientEmeterService struct {
	c *KasaClient
}

type GetRealtimeRequest struct {
}

type GetRealtimeResponse struct {
	Current float64 `mapstructure:"current"`
	Voltage float64 `mapstructure:"voltage"`
	Power   float64 `mapstructure:"power"`
	Total   float64 `mapstructure:"total"`
}

func (s *KasaClientEmeterService) GetRealtime() (*GetRealtimeResponse, error) {
	var response GetRealtimeResponse
	err := s.c.RPC("emeter", "get_realtime", GetRealtimeRequest{}, &response)

	return &response, err
}
