# OS Info - System Information Display

A lightweight GUI application written in Go that displays comprehensive system information in a floating window.

## Features

The application displays the following system information:

- **Date & Time**: Current date, time, and system uptime
- **System Information**: OS type, distribution, and kernel version
- **Disk Information**: Mount points with total, used, and free space
- **Battery Status**: Battery percentage, charging/discharging status, and temperature
- **Network Information**:
  - Network interfaces (WiFi/Ethernet)
  - MAC addresses
  - Local IP addresses
  - WiFi ESSID (network name)
  - Default gateway
  - DNS servers
  - External IP address (loaded asynchronously)
  - Country detection via GeoIP (loaded asynchronously)

## Window Behavior

- Window appears **fullscreen** with no borders
- **Closes automatically** when you click anywhere or press any key
- Larger font size (1.5x) for better readability
- Color-coded sections with icons for easy navigation
- Lazy loading for external IP and country information

## Prerequisites

### Go Installation

You need Go 1.21 or later installed on your system.

#### Install Go on Linux (Debian/Ubuntu)

```bash
# Option 1: Using apt (may have older version)
sudo apt-get update
sudo apt-get install golang-go

# Option 2: Install latest version from official source
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Install Go on Linux (Fedora/RHEL/CentOS)

```bash
# Option 1: Using dnf/yum
sudo dnf install golang

# Option 2: Install latest version from official source
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Install Go on macOS

```bash
# Using Homebrew (recommended)
brew install go

# Or download from https://go.dev/dl/
```

### GUI Dependencies

#### Linux (Debian/Ubuntu)

Before building, install the required X11 development libraries:

```bash
sudo apt-get update
sudo apt-get install libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

#### Linux (Fedora/RHEL/CentOS)

```bash
sudo dnf install libX11-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel mesa-libGL-devel libXxf86vm-devel
```

#### Linux (Arch Linux)

```bash
sudo pacman -S libx11 libxrandr libxcursor libxinerama libxi mesa libxxf86vm
```

#### macOS

No additional GUI dependencies required beyond Go.

## Installation

1. Clone or download this repository
2. Install dependencies (see Prerequisites above)
3. Build the application:

```bash
make build
```

Or manually:

```bash
go build -o bin/os-info ./cmd/os-info
```

## Usage

### Run the application:

```bash
make run
```

Or directly:

```bash
./bin/os-info
```

The application will display a fullscreen window with all system information. Click anywhere or press any key to close it.

**Note**: External IP and country information will appear as "searching..." initially and update automatically once fetched (may take 10-30 seconds depending on network speed).

## Building

### Available Make targets:

- `make build` - Build the application
- `make run` - Build and run the application
- `make clean` - Remove build artifacts
- `make install` - Install/update dependencies
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Run linter (golangci-lint)
- `make help` - Show available targets

## Project Structure

```
os-info/
├── cmd/
│   └── os-info/
│       └── main.go              # Application entry point (minimal)
├── internal/
│   ├── sysinfo/                 # System information gathering
│   │   ├── sysinfo.go          # Core Info struct and orchestration
│   │   ├── battery.go          # Battery information collection
│   │   ├── disk.go             # Disk information collection
│   │   └── network.go          # Network information collection
│   └── ui/                      # User interface components
│       ├── theme.go            # Custom Fyne theme (1.5x font)
│       ├── widgets.go          # Custom widgets (TappableContainer)
│       └── display.go          # Display creation and rendering
├── bin/                         # Compiled binaries (gitignored)
│   └── os-info
├── Makefile                     # Build automation
├── README.md                    # User documentation (this file)
├── CLAUDE.md                    # AI development documentation
├── go.mod                       # Go module definition
└── go.sum                       # Go dependencies checksums
```

## Dependencies

- **[Fyne v2](https://fyne.io/)** - Cross-platform GUI toolkit
- **[gopsutil v3](https://github.com/shirou/gopsutil)** - Cross-platform system and process utilities

## Platform Support

- **Linux**: Full support (tested on Debian-based distributions)
- **macOS**: Full support
- **Windows**: Partial support (battery temperature may not be available)

## Technical Details

### System Information Gathering

- **Date/Time**: Uses Go's `time` package with custom ordinal formatting
- **Uptime**: Retrieved via gopsutil's `host.Info()`
- **OS Info**: Retrieved via gopsutil and Go's `runtime` package
- **Disk Info**: Uses gopsutil's `disk.Partitions()` and `disk.Usage()`, filters virtual filesystems
- **Battery**: Reads from `/sys/class/power_supply/BAT*` on Linux
- **Network**: Uses gopsutil's `net.Interfaces()` and parses `/proc/net/route` for gateway
- **WiFi ESSID**: Uses `iwgetid` command with fallback to `iw dev`
- **External IP**: HTTP request to `api.ipify.org` (loaded asynchronously)
- **Country**: HTTP request to `ip-api.com` JSON API (loaded asynchronously)

### UI Features

- **Custom Theme**: 1.5x font size multiplier for better readability
- **Color Coding**: Each section has a distinct color (blue, green, orange, red, purple)
- **Icons**: Material design icons for each section
- **Lazy Loading**: External network calls run in background goroutines with Fyne data binding
- **Click-to-Close**: Custom tappable container widget for anywhere-click closing
- **Fullscreen**: Borderless fullscreen mode for overlay display

## License

This project is provided as-is for personal use.

## Author

Sebastien MORAND (seb.morand@gmail.com)
