package utils

import (
	"fmt"
	"strings"
)

// Exported values
const (
	KiloByte = 1024.0
	MegaByte = 1024 * KiloByte
	GigaByte = 1024 * MegaByte
	TeraByte = 1024 * GigaByte
)

// HumanBytes returns a human readable size
func HumanBytes(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TeraByte:
		unit = "TB"
		value = value / TeraByte
	case bytes >= GigaByte:
		unit = "GB"
		value = value / GigaByte
	case bytes >= MegaByte:
		unit = "MB"
		value = value / MegaByte
	case bytes >= KiloByte:
		unit = "KB"
		value = value / KiloByte
	case bytes >= 1.0:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	sv := fmt.Sprintf("%.1f", value)
	sv = strings.TrimSuffix(sv, ".0")
	return fmt.Sprintf("%s%s", sv, unit)
}
