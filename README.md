# Codory

A cross-platform CLI application built with Bubble Tea for managing various development and system tools.

## Features

- 🎯 **Smart Category System**: Categories and actions automatically adapt to your platform
- 🔧 **Extensible Action System**: Easily add new actions without modifying existing code
- 🌍 **True Multi-Platform Support**:
  - 🐧 **Linux** (Debian/Ubuntu, Arch)
    - 📦 APT, Pacman, Yay (AUR), Homebrew support
  - 🍎 **macOS**
    - 🍺 Homebrew support (with auto-installation)
  - 🪟 **Windows**
    - 📦 Winget integration
- 🎨 **Beautiful TUI Interface**: Clean and intuitive terminal UI
- 🚀 **Package Manager Auto-Installation**: Automatically prompts to install required package managers
- 🔒 **Platform-Aware Actions**: Actions only appear when they're available for your platform

## Platform Support

| Feature | Linux (Debian) | Linux (Arch) | macOS | Windows |
|---------|----------------|--------------|-------|---------|
| Package Managers | APT | Pacman/Yay/Homebrew | Homebrew | Winget |
| Auto-install PM | ❌ | ✅ Yay | ✅ Brew | ❌* |
| Cross-Platform Tools | ✅ | ✅ | ✅ | ✅ |

*Winget must be installed manually via Microsoft Store

## Installation

### Homebrew (macOS/Linux)

```bash
brew install anthodev/tap/codory
