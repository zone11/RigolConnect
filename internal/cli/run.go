// Package cli implements console input and output; device logic stays in the core.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zone11/RigolConnect/internal/buildinfo"
	"github.com/zone11/RigolConnect/internal/device"
	"github.com/zone11/RigolConnect/internal/screenshot"
)

// Run executes one CLI invocation and returns its process exit code.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fmt.Fprintf(stdout, "%s %s\n%s\n\n", buildinfo.Name, buildinfo.Version, buildinfo.Copyright)
	flags := flag.NewFlagSet("rigolconnect", flag.ContinueOnError)
	flags.SetOutput(stderr)
	host := flags.String("host", "", "Instrument IP address or hostname (required)")
	port := flags.Int("port", 5555, "TCP port")
	timeout := flags.Duration("timeout", 10*time.Second, "Timeout for connection and queries")
	version := flags.Bool("version", false, "Show application version and exit without connecting")
	info := flags.Bool("info", false, "Show instrument information and model compatibility")
	output := flags.String("screenshot", "", "Save screenshot to a new .png or .bmp file; otherwise query *IDN?")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if *version && flags.NArg() == 0 {
		return 0
	}
	if strings.TrimSpace(*host) == "" || *port < 1 || *port > 65535 || *timeout <= 0 || flags.NArg() != 0 {
		fmt.Fprintln(stderr, "Specify -host, a valid port and a positive timeout.")
		flags.PrintDefaults()
		return 2
	}
	address := net.JoinHostPort(*host, strconv.Itoa(*port))
	format := strings.TrimPrefix(strings.ToLower(filepath.Ext(*output)), ".")
	if *output != "" && format != "png" && format != "bmp" {
		fmt.Fprintln(stderr, "Screenshot filename must end with .png or .bmp.")
		return 2
	}
	if *output != "" {
		if _, err := os.Lstat(*output); err == nil {
			fmt.Fprintln(stderr, "Output file already exists:", *output)
			return 1
		} else if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, "Cannot access output file:", err)
			return 1
		}
	}
	instrument, err := device.Open(ctx, address, *timeout)
	if err != nil {
		fmt.Fprintln(stderr, "Instrument query failed:", err)
		return 1
	}
	defer instrument.Close()
	id := instrument.Identity()
	compatibility := instrument.Compatibility()
	if *info {
		printInfo(stdout, address, id, compatibility)
	}
	if *output == "" {
		if !*info {
			fmt.Fprintln(stdout, id.Raw)
		}
		return 0
	}
	data, err := instrument.Screenshot()
	if err == nil {
		data, err = screenshot.Encode(data, format)
	}
	if err == nil {
		err = saveNew(*output, data)
	}
	if err != nil {
		fmt.Fprintln(stderr, "Screenshot failed:", err)
		return 1
	}
	fmt.Fprintf(stdout, "Screenshot saved: %s (%d bytes)\n", *output, len(data))
	return 0
}

// O_EXCL also protects a file created after the initial existence check.
func saveNew(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		os.Remove(path) // Remove only the incomplete file created by this call.
		return err
	}
	return nil
}

func printInfo(w io.Writer, address string, id device.Identity, status device.Compatibility) {
	row := func(label, value string) { fmt.Fprintf(w, "  %-18s %s\n", label+":", value) }
	fmt.Fprintln(w, "Instrument")
	row("Manufacturer", id.Manufacturer)
	row("Model", id.Model)
	row("Serial number", id.SerialNumber)
	row("Firmware version", id.Firmware)
	if status.Family != "" {
		row("Instrument family", status.Family)
	}
	fmt.Fprintln(w, "\nConnection")
	row("Address", address)
	row("Protocol", "SCPI over TCP/IP")
	fmt.Fprintln(w, "\nSupport")
	row("Compatibility", status.Message)
	if status.Evidence != "" {
		row("Validation basis", status.Evidence)
	}
	if status.ModelSupported {
		row("Screenshot", "PNG / BMP")
	} else {
		row("Screenshot", "not available")
	}
	fmt.Fprintln(w, "\nSCPI identity")
	fmt.Fprintln(w, "  "+id.Raw)
	fmt.Fprintln(w)
}
