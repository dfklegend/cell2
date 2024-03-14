package proto

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Serializer implements the serialize.Serializer interface
type Serializer struct{}

var (
	protoSerializer = NewSerializer()
)

// NewSerializer returns a new Serializer.
func NewSerializer() *Serializer {
	return &Serializer{}
}

func GetDefaultSerializer() *Serializer {
	return protoSerializer
}

func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	if message, ok := v.(proto.Message); ok {
		bytes, err := proto.Marshal(message)
		if err != nil {
			return nil, err
		}

		return bytes, nil
	}
	return nil, fmt.Errorf("msg must be proto.Message")
}

func (s *Serializer) Unmarshal(bytes []byte, v interface{}) error {

	if message, ok := v.(proto.Message); ok {
		err := proto.Unmarshal(bytes, message)
		return err
	}

	return fmt.Errorf("msg must be proto.Message")
}

// GetName returns the name of the serializer.
func (s *Serializer) GetName() string {
	return "proto"
}
