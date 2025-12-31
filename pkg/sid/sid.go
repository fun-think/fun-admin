package sid

import (
	"fmt"

	"github.com/sony/sonyflake"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func NewSid() (*Sid, error) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		return nil, fmt.Errorf("failed to create sonyflake instance")
	}
	return &Sid{sf}, nil
}
func (s Sid) GenString() (string, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return "", err
	}
	return IntToBase62(int(id)), nil
}
func (s Sid) GenUint64() (uint64, error) {
	return s.sf.NextID()
}
