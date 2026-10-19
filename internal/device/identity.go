// Package device identifies Rigol oscilloscopes and selects verified protocols.
package device

import (
	"fmt"
	"strings"
)

type Identity struct {
	Manufacturer string
	Model        string
	SerialNumber string
	Firmware     string
	Raw          string
}

// ParseIdentity expects the four standard *IDN? fields. It never infers hardware
// revisions or installed options from a model name or a serial number.
func ParseIdentity(raw string) (Identity, error) {
	parts := strings.SplitN(strings.TrimSpace(raw), ",", 4)
	if len(parts) != 4 {
		return Identity{}, fmt.Errorf("invalid *IDN? response: expected manufacturer, model, serial number and firmware")
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
		if parts[i] == "" {
			return Identity{}, fmt.Errorf("empty field %d in *IDN? response", i+1)
		}
	}
	return Identity{Manufacturer: parts[0], Model: parts[1], SerialNumber: parts[2], Firmware: parts[3], Raw: strings.TrimSpace(raw)}, nil
}
