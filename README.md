# RigolConnect

RigolConnect 0.1 is a Go command-line application that captures whatever is displayed on the screen of a Rigol oscilloscope and saves it as a WYSIWYG (What You See Is What You Get) PNG or BMP image.

It supports the DS2202E, DS1054Z and DS1104Z. The DS2202E has been tested on physical hardware; DS1054Z and DS1104Z support is based on the original project and Rigol documentation.

SCPI (Standard Commands for Programmable Instruments) commands are sent over a direct TCP connection on port 5555. The computer and oscilloscope must be connected by LAN. No USB, VISA, IVI or Rigol drivers are required.

## Build

With Go 1.27.2 or later installed:

```sh
git clone https://github.com/zone11/RigolConnect.git
cd RigolConnect
go build -o rigolconnect ./cmd/rigolconnect
```

On Windows, use `-o rigolconnect.exe`.

## Example

Save the current screen as PNG:

```sh
./rigolconnect -host 192.168.44.174 -screenshot capture.png
```

Use a `.bmp` filename to save as BMP. Existing files are not overwritten.

Show device information:

```sh
./rigolconnect -host 192.168.44.174 -info
```

Run `./rigolconnect -h` for all options or `./rigolconnect -version` for the application version.

![Rigol DS2202E screen capture](https://github.com/user-attachments/assets/2555dc3e-3d3d-4ff0-a676-22f28e935562)

Based on the original [DS2202E-ScreenCapture Python project](https://github.com/zone11/DS2202E-ScreenCapture), which builds on work by [RoGeorge](https://github.com/RoGeorge/DS1054Z_screen_capture) and Python 3 contributions by [doug-a-brunner](https://github.com/doug-a-brunner).
