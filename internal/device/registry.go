package device

import "strings"

// Profiles are exact model matches backed by protocol documentation or hardware
// tests. Evidence is explicit; firmware is informational and never gates support.
type profile struct {
	model              string
	family             string
	screenshotCommand  string
	maxScreenshotBytes int
	evidence           string
}

var profiles = []profile{ds2202e, ds1054z, ds1104z}

type Compatibility struct {
	ModelSupported bool
	Family         string
	Message        string
	Evidence       string
}

func lookup(id Identity) (profile, bool) {
	manufacturer := strings.ToUpper(strings.TrimSpace(id.Manufacturer))
	if manufacturer != "RIGOL TECHNOLOGIES" && manufacturer != "RIGOL" {
		return profile{}, false
	}
	for _, p := range profiles {
		if strings.EqualFold(p.model, strings.TrimSpace(id.Model)) {
			return p, true
		}
	}
	return profile{}, false
}

// Check selects by manufacturer and model only. Firmware remains visible in
// Identity but is deliberately not compared with any version or allowlist.
func Check(id Identity) Compatibility {
	p, ok := lookup(id)
	if !ok {
		return Compatibility{Message: "not supported: no matching Rigol oscilloscope profile; screenshot disabled"}
	}
	return Compatibility{ModelSupported: true, Family: p.family, Message: "Model supported", Evidence: p.evidence}
}
