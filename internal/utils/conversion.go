package utils

import (
	"fmt"
	"math"
)

// FormatSize formats a number of bytes into a human readable string
func FormatSize(bytes uint) string {
	unitList := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	// Log(0) is -Inf, so special case
	if bytes == 0 {
		return "0B"
	}

	// number = bytes / 1024^{unit}
	// Compute both unit and number
	unit := int(math.Log(float64(bytes)) / math.Log(1024))
	// If we go over the unit list, we just use the last one
	unit = int(math.Min(float64(unit), float64(len(unitList)-1)))
	number := float64(bytes) / math.Pow(1024, float64(unit))

	return fmt.Sprintf("%.2f%s", number, unitList[unit])
}
