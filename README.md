# newtek-go

A Go library for interacting with NewTek TriCaster and 3Play devices via HTTP API and TCP connections.

## Developer

**FlowingSPDG** (Shugo Kawamura/河村 柊吾)

## Overview

This library provides a comprehensive Go interface for controlling and monitoring NewTek broadcast equipment, including TriCaster and 3Play devices. It supports both HTTP API communication and TCP-based real-time control, enabling developers to build custom automation and monitoring solutions.

## Features

### HTTP API Client (ClientV1)
- **Device Information**: Get product information, version details, and system status
- **Shortcut Control**: Execute device shortcuts via HTTP requests
- **Video Preview**: Capture video thumbnails and previews
- **Metadata Access**: Retrieve clip metadata and project information
- **State Monitoring**: Get current shortcut states and device status
- **WebSocket Notifications**: Real-time notifications for various events

### TCP Client (ClientV1TCP)
- **Direct TCP Communication**: Low-latency control via TCP port 5951
- **Shortcut Execution**: Send shortcut commands with parameters
- **Registration Management**: Register/unregister client connections
- **Binary Data Transmission**: Efficient data transfer for real-time control

### Supported Notification Types
- **Change Notifications**: Monitor system state changes
- **Shortcut Notifications**: Track shortcut executions
- **Timecode Notifications**: Receive timecode updates
- **VU Notifications**: Monitor audio levels
- **Video Notifications**: Track video state changes

## Installation

```bash
go get github.com/FlowingSPDG/newtek-go
```

## Dependencies

- Go 1.24+
- `github.com/coder/websocket` - WebSocket support
- `github.com/icholy/digest` - HTTP Digest authentication

## Quick Start

### HTTP API Client

```go
package main

import (
    "context"
    "fmt"
    "github.com/FlowingSPDG/newtek-go"
)

func main() {
    ctx := context.Background()
    
    // Create HTTP client
    client, err := newtek.NewClientV1("192.168.1.100", "admin", "admin")
    if err != nil {
        panic(err)
    }
    
    // Get device information
    product, err := client.GetProduct(ctx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Product: %s %s\n", product.ProductName, product.ProductVersion)
    fmt.Printf("Session: %s\n", product.SessionName)
    
    // Execute a shortcut
    err = client.ShortcutHTTP(ctx, "record", nil)
    if err != nil {
        panic(err)
    }
}
```

### TCP Client

```go
package main

import (
    "fmt"
    "github.com/FlowingSPDG/newtek-go"
)

func main() {
    // Create TCP client
    client, err := newtek.NewClientV1TCP("192.168.1.100")
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // Send shortcut command
    shortcuts := []newtek.Shortcut{
        {Name: "mode", Value: "2"},
    }
    
    req := newtek.Shortcuts{Shortcut: shortcuts}
    err = client.Shortcut(req)
    if err != nil {
        panic(err)
    }
}
```

### WebSocket Notifications

```go
package main

import (
    "context"
    "fmt"
    "github.com/FlowingSPDG/newtek-go"
)

func main() {
    ctx := context.Background()
    client, err := newtek.NewClientV1("192.168.1.100", "admin", "admin")
    if err != nil {
        panic(err)
    }
    
    // Listen for change notifications
    err = client.ChangeNotifications(ctx, func(msg string) {
        fmt.Printf("System change: %s\n", msg)
    })
    if err != nil {
        panic(err)
    }
    
    // Keep the program running
    select {}
}
```

## API Reference

### ClientV1 Interface

#### Device Information
- `GetLive(ctx context.Context) bool` - Check if device is live
- `GetProduct(ctx context.Context) (*ProductInformation, error)` - Get device information

#### Control
- `ShortcutHTTP(ctx context.Context, name string, kv map[string]string) error` - Execute shortcut via HTTP
- `Trigger(ctx context.Context, name string) error` - Trigger device action

#### Media
- `VideoPreview(ctx context.Context, name string, xres, yres, quality int) (image.Image, error)` - Get video preview
- `MetaData(ctx context.Context) (*Metadata, error)` - Get project metadata

#### State Monitoring
- `ShortcutStates(ctx context.Context) (*ShortcutStates, error)` - Get current shortcut states

#### Notifications
- `ChangeNotifications(ctx context.Context, cb func(msg string)) error` - Monitor system changes
- `ShortcutNotifications(ctx context.Context, cb func(msg string)) error` - Monitor shortcut executions
- `TimecodeNotifications(ctx context.Context, cb func(msg string)) error` - Monitor timecode
- `VUNotifications(ctx context.Context, name string, cb func(msg string)) error` - Monitor audio levels

### ClientV1TCP Interface

- `Register(name string) error` - Register client connection
- `Unregister(name string) error` - Unregister client connection
- `Shortcut(s Shortcuts) error` - Send shortcut command
- `Send(data string) error` - Send raw data
- `SendBytes(data []byte) error` - Send raw bytes
- `Close() error` - Close connection

## Data Models

### ProductInformation
Contains device information including model, version, serial number, and session details.

### ShortcutStates
Current state of all device shortcuts with their values and types.

### Metadata
Project metadata including clips, marks, and camera information.

### Tally
Tally information for video sources including program and preview states.

## Examples

The `example/` directory contains comprehensive examples demonstrating various features:

- `simple_shortcut/` - Basic shortcut execution
- `simple_shortcut_tcp/` - TCP-based shortcut control
- `change_notification/` - System change monitoring
- `shortcut_notifications/` - Shortcut execution monitoring
- `timecode/` - Timecode monitoring
- `vu_notifications/` - Audio level monitoring
- `thumbnail/` - Video preview capture
- `autoplay/` - Complex automation example with 3Play and TriCaster integration

## Use Cases

- **Broadcast Automation**: Automate recording, switching, and playback operations
- **Monitoring Systems**: Real-time monitoring of device status and audio levels
- **Integration**: Connect NewTek devices with external systems and workflows
- **Custom Control Panels**: Build custom control interfaces for specific workflows
- **Data Collection**: Gather device telemetry and operational data

## Network Requirements

- **HTTP API**: Port 80 (default) for HTTP communication
- **TCP Control**: Port 5951 for direct TCP control
- **WebSocket**: Port 80 for WebSocket notifications
- **Authentication**: HTTP Digest authentication required

## Error Handling

The library provides comprehensive error handling for network issues, authentication failures, and device communication errors. Always check returned errors and implement appropriate retry logic for production applications.

## Contributing

Contributions are welcome! Please ensure your code follows Go best practices and includes appropriate tests.

## License

This project is licensed under the terms specified in the LICENSE file.

## Support

For issues and questions, please create an issue in the GitHub repository.
