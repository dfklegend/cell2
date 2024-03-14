package serialize

type (
	Serializer interface {
		Marshal(interface{}) ([]byte, error)
		Unmarshal([]byte, interface{}) error
		GetName() string
	}
)
