package actions

import (
	"context"
	"slices"
)

type ActionType string

const (
	ActionTypeCommand  ActionType = "command"
	ActionTypeFunction ActionType = "function"
)

// Platform represents a specific platform
type Platform string

const (
	PlatformDebian  Platform = "debian"
	PlatformArch    Platform = "arch"
	PlatformMacOS   Platform = "macos"
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
	PlatformAny     Platform = "any"
)

// PackageSource indicates the source of the package
type PackageSource string

const (
	PackageSourceOfficial PackageSource = "official"
	PackageSourceAUR      PackageSource = "aur"
	PackageSourceBrew     PackageSource = "brew"
	PackageSourceWinget   PackageSource = "winget"
	PackageSourceAny      PackageSource = "any"
)

// PlatformCommand represents a command for a specific platform
type PlatformCommand struct {
	Command       string
	PackageSource PackageSource
	CheckCommand  string
	Interactive   bool // Whether this command requires user interaction
}

type ActionArgument struct {
	Name        string
	Description string
	Required    bool
}

// Action represents an executable action
type Action struct {
	ID             string
	Name           string
	Description    string
	Type           ActionType
	Handler        ActionHandler
	Arguments      []ActionArgument
	SuccessMessage string

	// For system commands (new structure)
	PlatformCommands map[Platform]PlatformCommand

	VisibleOnPlatforms []Platform
	HiddenOnPlatforms  []Platform
}

// ActionHandler is a function that executes an action
type ActionHandler func(ctx context.Context) (string, error)

// Category represents a category of actions
type Category struct {
	ID            string
	Name          string
	Description   string
	SubCategories []*Category
	Actions       []*Action

	// Platforms where the category is hidden
	HiddenOnPlatforms []Platform
}

// IsLeaf returns true if the category has no sub-categories
func (c *Category) IsLeaf() bool {
	return len(c.SubCategories) == 0
}

// HasActions returns true if the category has actions
func (c *Category) HasActions() bool {
	return len(c.Actions) > 0
}

// IsVisibleOnPlatform checks if the category is visible on a platform
func (c *Category) IsVisibleOnPlatform(platform Platform) bool {
	return !slices.Contains(c.HiddenOnPlatforms, platform)
}

// IsVisibleOnPlatform checks if the action is visible on a platform
func (a *Action) IsVisibleOnPlatform(platform Platform) bool {
	return !slices.Contains(a.HiddenOnPlatforms, platform)
}

// IsHiddenOnPlatform checks if the action is hidden on a platform
func (a *Action) IsHiddenOnPlatform(platform Platform) bool {
	return slices.Contains(a.HiddenOnPlatforms, platform)
}

// IsHiddenOnPlatform checks if the category is hidden on a platform
func (c *Category) IsHiddenOnPlatform(platform Platform) bool {
	return slices.Contains(c.HiddenOnPlatforms, platform)
}

// GetPlatformCommand returns the command for a platform
func (a *Action) GetPlatformCommand(platform Platform) (PlatformCommand, bool) {
	if cmd, ok := a.PlatformCommands[platform]; ok {
		return cmd, true
	}

	return PlatformCommand{}, false
}

func (a *Action) HasCommandForPlatform(platform Platform) bool {
	_, ok := a.GetPlatformCommand(platform)
	if ok {
		return true
	}

	if platform == PlatformDebian || platform == PlatformArch {
		_, ok = a.GetPlatformCommand(PlatformLinux)
		if ok {
			return true
		}
	}

	_, ok = a.GetPlatformCommand(PlatformAny)
	return ok
}

func (a *Action) HasCheckCommand(platform Platform) bool {
	cmd, ok := a.GetPlatformCommand(platform)
	if !ok {
		return false
	}
	return cmd.CheckCommand != ""
}

func (a *Action) GetCheckCommand(platform Platform) (string, bool) {
	cmd, ok := a.GetPlatformCommand(platform)
	if !ok {
		return "", false
	}
	if cmd.CheckCommand == "" {
		return "", false
	}
	return cmd.CheckCommand, true
}
