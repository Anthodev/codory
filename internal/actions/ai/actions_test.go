package ai

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/test/testutil"
)

func TestAICategoryActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	aiCategory, exists := globalRegistry.GetCategory("ai")
	if !exists {
		t.Fatal("AI category not found")
	}

	testutil.AssertStringEquals(t, aiCategory.ID, "ai", "Category ID")
	testutil.AssertStringEquals(t, aiCategory.Name, "AI", "Category Name")
	testutil.AssertStringEquals(t, aiCategory.Description, "AI tools you can install on your system", "Category Description")

	if len(aiCategory.Actions) != 1 {
		t.Fatalf("Expected 1 action in AI category, got %d", len(aiCategory.Actions))
	}

	installHerdrAction := aiCategory.Actions[0]
	if installHerdrAction.ID != "install_herdr" {
		t.Fatalf("Expected install_herdr action in AI category, got %s", installHerdrAction.ID)
	}

	testutil.AssertStringEquals(t, installHerdrAction.Name, "Install Herdr", "InstallHerdr Action Name")
	testutil.AssertStringEquals(t, installHerdrAction.Description, "Install Herdr agent multiplexer", "InstallHerdr Action Description")

	if installHerdrAction.Type != actions.ActionTypeCommand {
		t.Errorf("Expected InstallHerdr action type '%s', got '%s'", actions.ActionTypeCommand, installHerdrAction.Type)
	}

	if installHerdrAction.PlatformCommands == nil {
		t.Error("Expected PlatformCommands to be set on InstallHerdr action")
	}
}
