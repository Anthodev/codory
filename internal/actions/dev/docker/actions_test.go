package dev

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	// Since init() is automatically called when the package is imported,
	// we need to test that the global registry has been properly initialized

	// Get the global registry
	globalRegistry := actions.GlobalRegistry()

	// Check if the docker category exists (it should be registered under dev parent)
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found in global registry")
	}

	// Verify docker category properties
	if dockerCategory.ID != "docker" {
		t.Errorf("Expected docker category ID to be 'docker', got '%s'", dockerCategory.ID)
	}

	if dockerCategory.Name != "Docker" {
		t.Errorf("Expected docker category name to be 'Docker', got '%s'", dockerCategory.Name)
	}

	if dockerCategory.Description != "Docker commands to setup it" {
		t.Errorf("Expected docker category description to be 'Docker commands to setup it', got '%s'", dockerCategory.Description)
	}

	// Check if the expected actions are registered
	if len(dockerCategory.Actions) != 4 {
		t.Fatalf("Expected 4 actions in docker category, got %d", len(dockerCategory.Actions))
	}

	// Verify all actions are present
	var installDockerAction *actions.Action
	var installDockerComposeAction *actions.Action
	var addUserToDockerGroupAction *actions.Action
	var installLazydockerAction *actions.Action
	for _, action := range dockerCategory.Actions {
		if action.ID == "install_docker" {
			installDockerAction = action
		} else if action.ID == "install_docker_compose" {
			installDockerComposeAction = action
		} else if action.ID == "docker_add_user_to_group" {
			addUserToDockerGroupAction = action
		} else if action.ID == "install_lazydocker" {
			installLazydockerAction = action
		}
	}

	if installDockerAction == nil {
		t.Error("install_docker action not found in docker category")
	}

	if installDockerComposeAction == nil {
		t.Error("install_docker_compose action not found in docker category")
	}

	if addUserToDockerGroupAction == nil {
		t.Error("docker_add_user_to_group action not found in docker category")
	}

	if installLazydockerAction == nil {
		t.Error("install_lazydocker action not found in docker category")
	}

	if installDockerAction.Name != "Install Docker engine" {
		t.Errorf("Expected install_docker action name to be 'Install Docker engine', got '%s'", installDockerAction.Name)
	}

	if installDockerComposeAction.Name != "Install Docker Compose" {
		t.Errorf("Expected install_docker_compose action name to be 'Install Docker Compose', got '%s'", installDockerComposeAction.Name)
	}

	if addUserToDockerGroupAction.Name != "Add Current User to Docker Group" {
		t.Errorf("Expected docker_add_user_to_group action name to be 'Add Current User to Docker Group', got '%s'", addUserToDockerGroupAction.Name)
	}
}

func TestDockerCategoryIsSubCategoryOfDev(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the docker category directly from registry
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found in global registry")
	}

	// Verify docker category exists in registry
	if dockerCategory.ID != "docker" {
		t.Errorf("Expected docker category ID to be 'docker', got '%s'", dockerCategory.ID)
	}

	// Get the dev category to verify hierarchy
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Skip("Dev category not found - dev package may not be initialized in test environment")
		return
	}

	// Check if docker category is a subcategory of dev
	dockerCategoryFound := false
	for _, subCat := range devCategory.SubCategories {
		if subCat.ID == "docker" {
			dockerCategoryFound = true
			// Verify it's the same category instance
			if subCat != dockerCategory {
				t.Error("Docker category in subcategories is not the same instance as the registered one")
			}
			break
		}
	}

	if !dockerCategoryFound {
		t.Error("Docker category is not a subcategory of dev")
	}
}

func TestDockerCategoryActionsAreProperlyRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the docker category
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found")
	}

	// Verify that all actions in the docker category are properly initialized
	for _, action := range dockerCategory.Actions {
		if action.ID == "" {
			t.Error("Action ID is empty")
		}
		if action.Name == "" {
			t.Errorf("Action name is empty for ID: %s", action.ID)
		}
		if action.Type == "" {
			t.Errorf("Action type is empty for ID: %s", action.ID)
		}
		if action.Type == actions.ActionTypeFunction && action.Handler == nil {
			t.Errorf("Function action has no handler for ID: %s", action.ID)
		}
		if action.Type == actions.ActionTypeCommand && action.PlatformCommands == nil {
			t.Errorf("Command action has no platform commands for ID: %s", action.ID)
		}
	}
}

func TestDockerCategoryHasNoSubCategories(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the docker category
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found")
	}

	// Verify that docker category has no subcategories (it's a leaf category)
	if !dockerCategory.IsLeaf() {
		t.Error("Docker category should be a leaf category with no subcategories")
	}

	if len(dockerCategory.SubCategories) != 0 {
		t.Errorf("Docker category should have no subcategories, got %d", len(dockerCategory.SubCategories))
	}
}

func TestDockerCategoryHasActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the docker category
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found")
	}

	// Verify that docker category has actions
	if !dockerCategory.HasActions() {
		t.Error("Docker category should have actions")
	}

	if len(dockerCategory.Actions) == 0 {
		t.Error("Docker category should have at least one action")
	}
}

func TestDockerCategoryPlatformVisibility(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the docker category
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found")
	}

	// Test visibility on different platforms
	platforms := []actions.Platform{
		actions.PlatformAny,
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range platforms {
		if !dockerCategory.IsVisibleOnPlatform(platform) {
			t.Errorf("Docker category should be visible on platform %s", platform)
		}
		if dockerCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Docker category should not be hidden on platform %s", platform)
		}
	}
}

func TestDockerCategoryActionsHaveCorrectTypes(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the docker category
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found")
	}

	// Verify that actions have correct types
	for _, action := range dockerCategory.Actions {
		switch action.ID {
		case "install_docker":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		case "install_docker_compose":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		case "docker_add_user_to_group":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		case "install_lazydocker":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		default:
			t.Errorf("Unknown action %s with type %s", action.ID, action.Type)
		}
	}
}

func TestDockerCategoryRegistrationVerification(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Test that docker category is properly registered in the global registry
	dockerCategory, exists := globalRegistry.GetCategory("docker")
	if !exists {
		t.Fatal("Docker category not found in global registry")
	}

	// Verify docker category properties
	if dockerCategory.ID != "docker" {
		t.Errorf("Expected docker category ID to be 'docker', got '%s'", dockerCategory.ID)
	}

	if dockerCategory.Name != "Docker" {
		t.Errorf("Expected docker category name to be 'Docker', got '%s'", dockerCategory.Name)
	}

	if dockerCategory.Description != "Docker commands to setup it" {
		t.Errorf("Expected docker category description to be 'Docker commands to setup it', got '%s'", dockerCategory.Description)
	}

	// Verify docker category has the expected actions
	if len(dockerCategory.Actions) != 4 {
		t.Fatalf("Expected 4 actions in docker category, got %d", len(dockerCategory.Actions))
	}

	// Verify all actions are present
	var installDockerAction *actions.Action
	var installDockerComposeAction *actions.Action
	var addUserToDockerGroupAction *actions.Action
	var installLazydockerAction *actions.Action
	for _, action := range dockerCategory.Actions {
		if action.ID == "install_docker" {
			installDockerAction = action
		} else if action.ID == "install_docker_compose" {
			installDockerComposeAction = action
		} else if action.ID == "docker_add_user_to_group" {
			addUserToDockerGroupAction = action
		} else if action.ID == "install_lazydocker" {
			installLazydockerAction = action
		}
	}

	if installDockerAction == nil {
		t.Error("install_docker action not found in docker category")
	}

	if installDockerComposeAction == nil {
		t.Error("install_docker_compose action not found in docker category")
	}

	if addUserToDockerGroupAction == nil {
		t.Error("docker_add_user_to_group action not found in docker category")
	}

	if installLazydockerAction == nil {
		t.Error("install_lazydocker action not found in docker category")
	}

	// Verify docker category is visible on all platforms
	platforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range platforms {
		if dockerCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Docker category should not be hidden on %s platform", platform)
		}
	}
}

func TestDockerCategoryRegistrationIdempotency(t *testing.T) {
	// This test is disabled because the registry prevents duplicate category registration
	// and the init() function in actions.go would cause a panic if registration fails
	t.Skip("Registration idempotency test disabled - registry prevents duplicate registration")
}

func TestDevCategoryCreation(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Test that dev category exists and has correct properties
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found in global registry")
	}

	// Verify dev category properties
	if devCategory.ID != "dev" {
		t.Errorf("Expected dev category ID to be 'dev', got '%s'", devCategory.ID)
	}

	if devCategory.Name != "Development" {
		t.Errorf("Expected dev category name to be 'Development', got '%s'", devCategory.Name)
	}

	if devCategory.Description != "Development tools and utilities" {
		t.Errorf("Expected dev category description to be 'Development tools and utilities', got '%s'", devCategory.Description)
	}

	// Verify dev category has subcategories (should include docker)
	if len(devCategory.SubCategories) == 0 {
		t.Error("Dev category should have subcategories")
	}

	// Verify docker is in dev subcategories
	dockerFound := false
	for _, subCat := range devCategory.SubCategories {
		if subCat.ID == "docker" {
			dockerFound = true
			break
		}
	}

	if !dockerFound {
		t.Error("Docker category should be a subcategory of dev")
	}
}

func TestInitHandlesMissingDevCategory(t *testing.T) {
	// This test verifies that the init() function properly handles the case
	// where the dev category doesn't exist yet by creating it
	// Since we can't easily reset the global registry in tests, we verify
	// that the dev category exists and has the expected properties

	globalRegistry := actions.GlobalRegistry()

	// Verify dev category exists (it should have been created by init() if missing)
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category should exist - init() should have created it if missing")
	}

	// Verify dev category has the expected structure
	if devCategory.SubCategories == nil {
		t.Error("Dev category should have initialized SubCategories slice")
	}

	if devCategory.Actions == nil {
		t.Error("Dev category should have initialized Actions slice")
	}

	// Verify docker category exists as a subcategory
	dockerFound := false
	for _, subCat := range devCategory.SubCategories {
		if subCat.ID == "docker" {
			dockerFound = true
			if subCat.Actions == nil {
				t.Error("Docker category should have initialized Actions slice")
			}
			// Note: Docker category doesn't have SubCategories initialized since it's a leaf category
			break
		}
	}

	if !dockerFound {
		t.Error("Docker category should be registered as a subcategory of dev")
	}
}
