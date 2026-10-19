package device

// The reference project explicitly names DS1054Z and DS1104Z (an upgraded
// DS1054Z). Both use :DISP:DATA? with a definite-length BMP block. The Rigol
// MSO1000Z/DS1000Z Programming Guide, section 2-63, confirms default BMP24.
// Sources: https://github.com/RoGeorge/DS1054Z_screen_capture
// https://www.rigol.com/dam/global/downloads/brochures/en/program-guide/oscilloscopes/DS1000Z_ProgrammingGuide_EN.pdf
// No local hardware tests are available for these models. No Python code was
// copied or ported; the existing Go transport/parser implements the protocol.

var ds1054z = ds1000zProfile("DS1054Z")
var ds1104z = ds1000zProfile("DS1104Z")

func ds1000zProfile(model string) profile {
	return profile{
		model:              model,
		family:             "Rigol DS1000Z",
		screenshotCommand:  ":DISP:DATA?",
		maxScreenshotBytes: 16 * 1024 * 1024,
		evidence:           "Reference project and Rigol manual; not tested on local hardware",
	}
}
