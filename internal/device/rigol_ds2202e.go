package device

// Verified with physical DS2202E hardware. Other DS2000 models are intentionally
// not registered until their command syntax and payload have been checked.

var ds2202e = profile{
	model:              "DS2202E",
	family:             "Rigol DS2000E",
	screenshotCommand:  ":DISP:DATA?",
	maxScreenshotBytes: 16 * 1024 * 1024,
	evidence:           "Tested on local DS2202E hardware",
}
