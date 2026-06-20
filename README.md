# Codory

A cross-platform CLI application built with Bubble Tea for managing various development and system tools.

## Features

- 🎯 **Smart Category System**: Categories and actions automatically adapt to your platform
- 🔧 **Extensible Action System**: Easily add new actions without modifying existing code
- 🌍 **True Multi-Platform Support**:
  - 🐧 **Linux**
    - 📦 APT (Debian-based), Pacman (Arch-based), Yay (AUR) (Arch-based), Homebrew support (with auto-installation)
  - 🍎 **macOS**
    - 🍺 Homebrew support (with auto-installation)
  - 🪟 **Windows**
    - 📦 Winget integration
- 🎨 **Beautiful TUI Interface**: Clean and intuitive terminal UI
- 🚀 **Package Manager Auto-Installation**: Automatically prompts to install required package managers
- 🔒 **Platform-Aware Actions**: Actions only appear when they're available for your platform

## Platform Support

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| Package Managers | APT (Debian-based)/Pacman (Arch-based)/Yay (Arch-based)/Homebrew | Homebrew | Winget |
| Auto-install PM | ✅ Yay/Brew | ✅ Brew | ❌* |
| Cross-Platform Tools | ✅ | ✅ | ✅ |

*Winget must be installed manually via Microsoft Store

## Installation

### Homebrew (macOS/Linux)

```bash
brew install anthodev/tap/codory
```

### From source

```bash
git clone https://github.com/anthodev/codory.git
cd codory
go build -o codory .
```

## Usage

```bash
./codory
```
