package utils

import (
	"math"
)

func ToGridX(x float32, width int) int {
	return ToGridV(x, width)
}

func GridToX(x int, width int) float32 {
	return GridToV(x, width)
}

func ToGridZ(z float32, height int) int {
	return ToGridV(z, height)
}

func GridToZ(z int, height int) float32 {
	return GridToV(z, height)
}

func ToGridV(x float32, size int) int {
	halfSize := size / 2
	return int(math.Round(float64(x))) + halfSize
}

func GridToV(z int, size int) float32 {
	halfSize := size / 2
	return float32(z - halfSize)
}
