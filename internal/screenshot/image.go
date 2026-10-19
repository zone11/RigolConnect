// Package screenshot validates instrument bitmap data and converts it for export.
package screenshot

import (
	"bytes"
	"fmt"
	"image/png"

	"golang.org/x/image/bmp"
)

// Encode validates the complete BMP image and encodes the decoded pixels.
// Re-encoding BMP repairs the incorrect file-size field sent by DS2202E firmware.
func Encode(data []byte, format string) ([]byte, error) {
	if format != "png" && format != "bmp" {
		return nil, fmt.Errorf("unsupported output format %q", format)
	}
	config, err := bmp.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("read BMP configuration: %w", err)
	}
	const maxPixels = 16 * 1024 * 1024
	if config.Width <= 0 || config.Height <= 0 || config.Width > maxPixels/config.Height {
		return nil, fmt.Errorf("invalid or excessive image dimensions: %dx%d", config.Width, config.Height)
	}
	img, err := bmp.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode BMP: %w", err)
	}
	var result bytes.Buffer
	if format == "bmp" {
		err = bmp.Encode(&result, img)
	} else {
		err = png.Encode(&result, img)
	}
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", format, err)
	}
	return result.Bytes(), nil
}
