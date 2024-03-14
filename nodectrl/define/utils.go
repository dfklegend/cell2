package define

import (
	"fmt"
)

func GetStateName(state NodeState) string {
	switch state {
	case Init:
		return "init"
	case Working:
		return "working"
	case Retiring:
		return "retiring"
	case Retired:
		return "retired"
	}
	return fmt.Sprintf("%v", state)
}
