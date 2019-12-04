package kasa

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/mitchellh/mapstructure"
)

type KasaClient struct {
	addr  string
	model string
}

type KasaClientConfig struct {
	Host string
}

type RPCResponse struct {
	ErrCode int `json:"err_code"`
}

func New(c *KasaClientConfig) *KasaClient {
	return &KasaClient{
		addr: c.Host + ":9999",
	}
}

const key = 171

func packInt(in int32) []byte {
	out := make([]byte, 4)
	for i := 3; i > 0; i-- {
		out[i] = byte(in & 0xFF)
		in >>= 8
	}
	return out
}

func unpackInt(in []byte) int32 {
	length := 0
	for i := 0; i < 4; i++ {
		length <<= 8
		length += int(in[i])
	}
	return int32(length)
}

func encrypt(in []byte) []byte {
	length := len(in)
	out := make([]byte, length)

	key := key
	for i, r := range in {
		key = key ^ int(r)
		out[i] = byte(key)
	}
	return out
}

func decrypt(in []byte) []byte {
	length := len(in)
	out := make([]byte, length)
	key := key
	for i := 0; i < length; i++ {
		b := int(in[i])
		out[i] = byte(key ^ b)
		key = b
	}
	return out
}

func (c *KasaClient) Request(payload interface{}) ([]byte, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}
	// the tcp server on smart plug can only handle one connection
	// at a time, make sure we don't wait too long to block others
	conn.SetWriteDeadline(time.Now().Add(time.Second))

	jpayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	_, err = conn.Write(packInt(int32(len(jpayload))))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(encrypt(jpayload))
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(time.Second))

	var buf []byte
	tmp := make([]byte, 1024)

	read := -1
	length := 0
	for read < length {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		read += n
		if length == 0 {
			length = int(unpackInt(tmp[:4]))
			tmp = tmp[4:]
			n -= 4
		}
		buf = append(buf, tmp[:n]...)

	}
	err = conn.Close()
	if err != nil {
		return nil, err
	}
	buf = decrypt(buf)

	return buf, nil
}

func (c *KasaClient) RPC(service string, cmd string, payload interface{}, out interface{}) error {
	payload = map[string]interface{}{
		service: map[string]interface{}{
			cmd: payload,
		},
	}

	response, err := c.Request(payload)
	if err != nil {
		return err
	}

	var outMarshal map[string]map[string]map[string]interface{}

	err = json.Unmarshal(response, &outMarshal)
	if err != nil {
		return err
	}

	if outMarshal[service] == nil || outMarshal[service][cmd] == nil {
		return fmt.Errorf("malformed response: %v", outMarshal)
	}
	var r RPCResponse
	mapstructure.Decode(outMarshal[service][cmd], &r)
	if r.ErrCode != 0 {
		return fmt.Errorf("rpc error: %v", outMarshal)
	}

	mapstructure.Decode(outMarshal[service][cmd], &out)

	return nil

}

func (c *KasaClient) SystemService() *KasaClientSystemService {
	return &KasaClientSystemService{
		c: c,
	}
}

func (c *KasaClient) EmeterService() *KasaClientEmeterService {
	return &KasaClientEmeterService{
		c: c,
	}
}
