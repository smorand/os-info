# OS Info - AI Development Documentation

## Project Overview

This is a Go-based GUI application that displays comprehensive system information in a fullscreen overlay window. Built using Fyne toolkit for cross-platform GUI and gopsutil for system information gathering. Features lazy loading for external IP and country detection, custom theming, color-coded sections, and click-anywhere-to-close functionality.

The project follows Go standard project layout with proper separation of concerns.

## Architecture

### Project Structure

```
os-info/
├── cmd/
│   └── os-info/
│       └── main.go              # Application entry point (minimal)
├── internal/
│   ├── sysinfo/                 # System information gathering
│   │   ├── sysinfo.go          # Core Info struct and collection orchestration
│   │   ├── battery.go          # Battery information collection
│   │   ├── disk.go             # Disk information collection
│   │   └── network.go          # Network information collection
│   └── ui/                      # User interface components
│       ├── theme.go            # Custom Fyne theme (1.5x font size)
│       ├── widgets.go          # Custom widgets (TappableContainer)
│       └── display.go          # Display creation and rendering
├── bin/                         # Build output (gitignored)
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Makefile                     # Build automation
├── README.md                    # User documentation
└── CLAUDE.md                    # AI development documentation
```

### Core Packages

#### cmd/os-info

**main.go** - Minimal entry point that wires everything together:
- Creates Fyne application with custom theme
- Initializes system information collection
- Sets up UI display
- Configures window behavior (fullscreen, click-to-close)

#### internal/sysinfo

**sysinfo.go** - Core system information orchestration:
- `Info` struct: Main data structure holding all system information
- `New()` constructor: Initializes and collects all system data
- `UpdateExternalNetworkInfo()`: Asynchronous external network data collection
- Date/time collection with ordinal suffixes (1st, 2nd, 3rd, etc.)
- OS type and version detection (macOS/Linux)

**battery.go** - Battery information collection:
- Platform-specific battery reading (Linux via `/sys/class/power_supply`, macOS via pmset)
- Battery percentage, status, temperature monitoring
- AC adapter detection

**disk.go** - Disk information collection:
- `DiskInfo` struct: Per-partition information
- Physical disk filtering (excludes virtual filesystems, loop devices, snap mounts)
- Excludes `/boot` partitions
- Formatted table output with mount point, total, used, free, and usage percentage

**network.go** - Network information collection:
- `NetworkInfo` struct: Per-interface information including external IP and country
- Active interface detection via `/proc/net/route`
- WiFi ESSID detection using `iwgetid` or `iw` commands
- Default gateway and DNS server detection
- External IP detection via `api.ipify.org`
- Country detection via `ip-api.com` GeoIP API
- Async HTTP requests with 30-second timeouts

#### internal/ui

**theme.go** - Custom Fyne theme:
- `CustomTheme` struct implementing `fyne.Theme` interface
- 1.5x font size multiplier for all text
- Delegates colors, fonts, and icons to default theme

**widgets.go** - Custom widgets:
- `TappableContainer`: Widget that executes callback on tap
- `tappableRenderer`: Renderer implementation for tappable container
- Used for click-anywhere-to-close functionality

**display.go** - UI display creation:
- `CreateInfoDisplay()`: Main display builder
- Color-coded section creators for different information types
- Data binding setup for async network updates
- Specialized rendering for date/time, system, battery, disk, and network sections
- Bold formatting for emphasis (date, battery percentage, system name)

### Data Flow

```
main()
  → app.New() with CustomTheme
  → sysinfo.New()
    → collectDateTimeInfo() - date with ordinal, uptime
    → collectOSInfo() - OS type, distribution, kernel
    → collectDiskInfo() - physical disks only
    → collectBatteryInfo() - battery %, status, temp, adapter
    → collectNetworkInfo() - active interface, ExternalIP/Country = "searching..."
  → ui.CreateInfoDisplay(sysInfo, window)
    → Create color-coded sections with icons
    → Set up data binding for network section
    → Launch UpdateExternalNetworkInfo() goroutine
      → getExternalIP() - async HTTP to api.ipify.org
      → getCountry() - async HTTP to ip-api.com
      → Update binding (triggers UI refresh)
  → ui.NewTappableContainer() - wrap display
  → w.SetFullScreen(true)
  → w.ShowAndRun()
```

## Code Organization Principles

### File Element Order

Each Go file follows this structure:
1. Package declaration
2. Import statements (grouped: stdlib, external, internal)
3. Constants (if any)
4. Type/Struct definitions (alphabetically by field name)
5. Constructor functions (`New*`)
6. Public methods (alphabetically)
7. Private methods/helper functions (alphabetically)

### Naming Conventions

- **Packages**: Lowercase, single word (sysinfo, ui)
- **Structs**: PascalCase, exported (Info, DiskInfo, NetworkInfo)
- **Functions**: camelCase for private, PascalCase for exported
- **Methods**: Receiver name is single letter or abbreviation (i for Info, d for DiskInfo)
- **Constants**: PascalCase for exported, camelCase for private

### Error Handling

- Graceful degradation: Failed reads return "N/A" or "Unknown" instead of panicking
- Explicit error ignoring: `_, _ = func()` pattern for intentionally ignored errors
- Deferred cleanup: `defer func() { _ = resp.Body.Close() }()` for resource cleanup

## Platform-Specific Implementation

### Linux

- **Battery**: Reads from `/sys/class/power_supply/BAT[0-1]/`
  - Capacity: `capacity` file
  - Status: `status` file
  - Temperature: `temp` file (tenths of degrees Celsius)
  - Adapter: Checks `AC*/online` files
- **Network gateway**: Parses `/proc/net/route` for default route
- **DNS servers**: Parses `/etc/resolv.conf`
- **WiFi ESSID**: Executes `iwgetid -r` or `iw dev [interface] link`
- **Active interface**: Determined from default route in `/proc/net/route`
- **Disk filtering**: Filters virtual filesystems and devices

### macOS

- **Battery**: Placeholder (would use `pmset` command)
- **Other features**: Use gopsutil cross-platform APIs

### Cross-Platform

- **OS info**: Uses gopsutil's `host.Info()` and Go's `runtime.GOOS`
- **Disk info**: Uses gopsutil's `disk.Partitions()` and `disk.Usage()`
- **Network interfaces**: Uses gopsutil's `net.Interfaces()`
- **Date/Time**: Standard Go `time` package
- **External IP/Country**: HTTP APIs work on all platforms

## Dependencies

### Go Version

- **Required**: Go 1.21 or later

### External Libraries

```go
// GUI Framework
"fyne.io/fyne/v2"                     // Core GUI toolkit
"fyne.io/fyne/v2/app"                 // Application creation
"fyne.io/fyne/v2/canvas"              // Canvas primitives
"fyne.io/fyne/v2/container"           // Layout containers
"fyne.io/fyne/v2/data/binding"        // Data binding for reactive UI
"fyne.io/fyne/v2/theme"               // Theming system
"fyne.io/fyne/v2/widget"              // Standard widgets

// System Information
"github.com/shirou/gopsutil/v3/disk"  // Disk information
"github.com/shirou/gopsutil/v3/host"  // Host/OS information
"github.com/shirou/gopsutil/v3/net"   // Network information
```

### GUI Build Dependencies

#### Linux (Debian/Ubuntu)
```bash
sudo apt-get install libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

#### Linux (Fedora/RHEL/CentOS)
```bash
sudo dnf install libX11-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel mesa-libGL-devel libXxf86vm-devel
```

## Build and Development

### Build Commands

```bash
make build     # Build the application
make run       # Build and run
make clean     # Remove build artifacts
make install   # Install dependencies
make test      # Run tests
make fmt       # Format code
make lint      # Run linter
make help      # Show available targets
```

### Build Output

- Binary location: `bin/os-info`
- Build command: `go build -o bin/os-info ./cmd/os-info`

## Key Features Implementation

### Lazy Loading Network Information

External IP and country detection happen asynchronously to avoid blocking the UI:

1. Initial display shows "searching..." for external IP and country
2. `UpdateExternalNetworkInfo()` launches goroutine
3. Goroutine fetches external IP, then country
4. Updates `Info.Networks[0]` fields
5. Invokes callback to update data binding
6. UI automatically refreshes via binding

### Custom Theme

1.5x font size multiplier applied to all text elements via custom theme implementation.

### Click-Anywhere-to-Close

Custom `TappableContainer` widget wraps entire display and triggers window close on any tap event.

### Color-Coded Sections

Each information category has a distinct background color:
- Date/Time: Cornflower blue (RGB 100, 149, 237)
- System: Medium sea green (RGB 60, 179, 113)
- Disk: Dark orange (RGB 255, 140, 0)
- Battery: Crimson (RGB 220, 20, 60)
- Network: Medium purple (RGB 147, 112, 219)

## Testing Considerations

When testing:
- Verify all system info displays correctly on your platform
- Test with missing battery (desktop systems)
- Test with multiple network interfaces
- Test with various disk configurations
- Verify fullscreen and always-on-top behavior
- Confirm click-to-close and key-press-to-close functionality
- Test external IP/country async loading with slow network

## Common Issues

1. **Build fails with X11 errors**: Install X11 development libraries
2. **Battery shows N/A**: Normal on desktops or systems without battery
3. **External IP/Country show "searching..." for long time**: Network latency, 30s timeout
4. **WiFi ESSID shows N/A**: Requires `iwgetid` or `iw` command
5. **Module import errors**: Run `make install` to download dependencies

## Extending the Application

### Adding New System Information

1. Add field to `Info` struct in `internal/sysinfo/sysinfo.go`
2. Create `collect*Info()` method (or add to existing collector)
3. Call from `New()` constructor
4. Create formatting method (e.g., `Get*InfoTable()`)
5. Add display section in `internal/ui/display.go`

### Adding New UI Section

1. Create section builder function in `internal/ui/display.go`
2. Choose background color
3. Use `createColoredSectionMultiLineMonospaceWithIcon()` or create custom
4. Add to `CreateInfoDisplay()` VBox container
5. Add separator if needed with `widget.NewSeparator()`

### Modifying Layout

- Edit `CreateInfoDisplay()` in `internal/ui/display.go`
- Use Fyne container types: `container.NewVBox()`, `container.NewHBox()`, `container.NewGrid()`
- Update color scheme in section creation functions

## Future Enhancements

Potential improvements:
- Memory usage information (total, used, free, cached)
- CPU usage/temperature monitoring
- Process list with top consumers
- System load averages (1m, 5m, 15m)
- Alternative external IP services as fallbacks
- VPN status detection
- Battery health metrics (cycle count, design capacity vs current)
- Network bandwidth usage statistics
- Configuration file for customization (colors, font size, sections)
- Refresh button to update information on demand
- macOS battery reading implementation

## Code Style

- Following Go standard conventions and golang skill recommendations
- Package-oriented design with clear separation of concerns
- No `/src` directory (forbidden in Go standards)
- One struct per responsibility
- Private functions for internal logic, public methods for API
- Early returns with default "N/A" values on errors
- Descriptive variable names
- Error handling: Explicit ignore with `_, _ = func()` when safe
- Deferred cleanup: `defer func() { _ = obj.Close() }()` pattern
- Alphabetical field ordering in structs for consistency
- Comments explain "why" not "what" (code is self-documenting)
