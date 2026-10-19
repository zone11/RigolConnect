package device

import (
	"context"
	"fmt"
	"time"

	"github.com/zone11/RigolConnect/internal/scpi"
)

// Device keeps identity and capture on the same session. It supports one capture
// per session and sequential calls. Close it after use.
type Device struct {
	client   *scpi.Client
	identity Identity
}

// Open only issues *IDN?, including for unsupported models, so -info can explain
// why model-specific operations are unavailable.
func Open(ctx context.Context, address string, timeout time.Duration) (*Device, error) {
	client, err := scpi.Dial(ctx, address, timeout)
	if err != nil {
		return nil, err
	}
	raw, err := client.Identify()
	if err != nil {
		client.Close()
		return nil, err
	}
	id, err := ParseIdentity(raw)
	if err != nil {
		client.Close()
		return nil, err
	}
	return &Device{client: client, identity: id}, nil
}

func (d *Device) Close()                       { d.client.Close() }
func (d *Device) Identity() Identity           { return d.identity }
func (d *Device) Compatibility() Compatibility { return Check(d.identity) }

// Screenshot rejects unknown models before sending model-specific commands.
// The current profiles return BMP data for the screenshot conversion package.
func (d *Device) Screenshot() ([]byte, error) {
	p, ok := lookup(d.identity)
	if !ok {
		return nil, fmt.Errorf("screenshot not supported for %s %s: no matching model profile", d.identity.Manufacturer, d.identity.Model)
	}
	data, err := d.client.QueryBlock(p.screenshotCommand, p.maxScreenshotBytes)
	if err != nil {
		return nil, fmt.Errorf("read %s screenshot: %w", p.model, err)
	}
	return data, nil
}
