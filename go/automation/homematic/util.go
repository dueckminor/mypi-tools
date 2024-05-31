package homematic

import (
	"fmt"
	"math"
)

func makeFloat64(val interface{}) (float64, error) {
	switch i := val.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	}
	return math.NaN(), fmt.Errorf("Failed to convert type '%T' to float64", val)
}
