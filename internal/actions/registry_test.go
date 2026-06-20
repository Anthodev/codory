package actions

import (
	"strings"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *Registry
	}{
		{
			name: "creates valid registry",
			want: &Registry{
				categories: make(map[string]*Category),
				root: &Category{
					ID:            "root",
					Name:          "Main Menu",
					Description:   "Select a category",
					SubCategories: make([]*Category, 0),
					Actions:       make([]*Action, 0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewRegistry()

			if got == nil {
				t.Fatal("NewRegistry() returned nil")
			}

			if got.root == nil {
				t.Fatal("root category is nil")
			}

			if got.root.ID != tt.want.root.ID {
				t.Errorf("root.ID = %q, want %q", got.root.ID, tt.want.root.ID)
			}

			if got.root.Name != tt.want.root.Name {
				t.Errorf("root.Name = %q, want %q", got.root.Name, tt.want.root.Name)
			}

			if got.categories == nil {
				t.Fatal("categories map is nil")
			}
		})
	}
}

func TestGlobalRegistry(t *testing.T) {
	t.Parallel()

	registry1 := GlobalRegistry()
	registry2 := GlobalRegistry()

	if registry1 != registry2 {
		t.Error("GlobalRegistry() should return the same instance")
	}
}

func TestRegistry_RegisterCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		parent      string
		category    *Category
		setup       func(*Registry) error
		wantErr     bool
		errContains string
	}{
		{
			name:   "register root category",
			parent: "",
			category: &Category{
				ID:   "test-cat",
				Name: "Test Category",
			},
			wantErr: false,
		},
		{
			name:   "register category with parent",
			parent: "root",
			category: &Category{
				ID:   "test-cat",
				Name: "Test Category",
			},
			wantErr: false,
		},
		{
			name:   "register duplicate category",
			parent: "root",
			category: &Category{
				ID:   "test-cat",
				Name: "Test Category",
			},
			setup: func(r *Registry) error {
				return r.RegisterCategory("root", &Category{
					ID:   "test-cat",
					Name: "Test Category",
				})
			},
			wantErr:     true,
			errContains: "category test-cat already exists",
		},
		{
			name:   "register with non-existent parent",
			parent: "non-existent",
			category: &Category{
				ID:   "test-cat2",
				Name: "Test Category 2",
			},
			wantErr:     true,
			errContains: "parent category non-existent not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := NewRegistry()

			if tt.setup != nil {
				if err := tt.setup(registry); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := registry.RegisterCategory(tt.parent, tt.category)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got none")
				}

				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify category was registered
			got, exists := registry.categories[tt.category.ID]
			if !exists {
				t.Fatal("category was not registered in categories map")
			}

			if got != tt.category {
				t.Error("registered category does not match input category")
			}

			// Verify category was added to parent's subcategories
			var parentCat *Category
			if tt.parent == "" || tt.parent == "root" {
				parentCat = registry.root
			} else {
				parentCat, _ = registry.categories[tt.parent]
			}

			found := false
			for _, subCat := range parentCat.SubCategories {
				if subCat.ID == tt.category.ID {
					found = true
					break
				}
			}
			if !found {
				t.Error("category was not added to parent's subcategories")
			}
		})
	}
}

func TestRegistry_RegisterAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		categoryID  string
		action      *Action
		setup       func(*Registry) error
		wantErr     bool
		errContains string
	}{
		{
			name:       "register action in existing category",
			categoryID: "test-cat",
			action: &Action{
				ID:   "test-action",
				Name: "Test Action",
			},
			setup: func(r *Registry) error {
				return r.RegisterCategory("", &Category{
					ID:   "test-cat",
					Name: "Test Category",
				})
			},
			wantErr: false,
		},
		{
			name:       "register action in non-existent category",
			categoryID: "non-existent",
			action: &Action{
				ID:   "test-action",
				Name: "Test Action",
			},
			wantErr:     true,
			errContains: "category non-existent not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := NewRegistry()

			if tt.setup != nil {
				if err := tt.setup(registry); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := registry.RegisterAction(tt.categoryID, tt.action)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got none")
				}

				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify action was added to category
			category, exists := registry.categories[tt.categoryID]
			if !exists {
				t.Fatal("category not found")
			}

			found := false
			for _, action := range category.Actions {
				if action.ID == tt.action.ID {
					found = true
					if action != tt.action {
						t.Error("registered action does not match input action")
					}
					break
				}
			}
			if !found {
				t.Error("action was not added to category")
			}
		})
	}
}

func TestRegistry_GetRoot(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	root := registry.GetRoot()

	if root == nil {
		t.Fatal("GetRoot() returned nil")
	}

	if root.ID != "root" {
		t.Errorf("root.ID = %q, want %q", root.ID, "root")
	}
}

func TestRegistry_GetCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		id         string
		setup      func(*Registry) error
		wantExists bool
	}{
		{
			name: "get existing category",
			id:   "test-cat",
			setup: func(r *Registry) error {
				return r.RegisterCategory("", &Category{
					ID:   "test-cat",
					Name: "Test Category",
				})
			},
			wantExists: true,
		},
		{
			name:       "get non-existent category",
			id:         "non-existent",
			setup:      nil,
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := NewRegistry()

			if tt.setup != nil {
				if err := tt.setup(registry); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			category, exists := registry.GetCategory(tt.id)

			if tt.wantExists {
				if !exists {
					t.Fatal("expected category to exist")
				}

				if category == nil {
					t.Fatal("expected category to be non-nil")
				}

				if category.ID != tt.id {
					t.Errorf("category.ID = %q, want %q", category.ID, tt.id)
				}
			} else {
				if exists {
					t.Fatal("expected category to not exist")
				}

				if category != nil {
					t.Fatal("expected category to be nil")
				}
			}
		})
	}
}

func TestCategory_GetVisibleSubCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category *Category
		platform Platform
		want     int
	}{
		{
			name: "all subcategories visible",
			category: &Category{
				SubCategories: []*Category{
					{ID: "cat1", HiddenOnPlatforms: []Platform{PlatformWindows}},
					{ID: "cat2", HiddenOnPlatforms: []Platform{PlatformMacOS}},
				},
			},
			platform: PlatformLinux,
			want:     2,
		},
		{
			name: "some subcategories hidden",
			category: &Category{
				SubCategories: []*Category{
					{ID: "cat1", HiddenOnPlatforms: []Platform{PlatformLinux}},
					{ID: "cat2", HiddenOnPlatforms: []Platform{PlatformMacOS}},
				},
			},
			platform: PlatformLinux,
			want:     1,
		},
		{
			name: "all subcategories hidden",
			category: &Category{
				SubCategories: []*Category{
					{ID: "cat1", HiddenOnPlatforms: []Platform{PlatformLinux}},
					{ID: "cat2", HiddenOnPlatforms: []Platform{PlatformLinux}},
				},
			},
			platform: PlatformLinux,
			want:     0,
		},
		{
			name:     "no subcategories",
			category: &Category{},
			platform: PlatformLinux,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.category.GetVisibleSubCategories(tt.platform)

			if len(got) != tt.want {
				t.Errorf("GetVisibleSubCategories() returned %d items, want %d", len(got), tt.want)
			}

			// Verify all returned categories are visible on the platform
			for _, cat := range got {
				if cat.IsHiddenOnPlatform(tt.platform) {
					t.Errorf("returned category %q is hidden on platform %q", cat.ID, tt.platform)
				}
			}
		})
	}
}

func TestCategory_GetVisibleActions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category *Category
		platform Platform
		want     int
	}{
		{
			name: "all actions visible",
			category: &Category{
				Actions: []*Action{
					{ID: "action1", HiddenOnPlatforms: []Platform{PlatformWindows}},
					{ID: "action2", HiddenOnPlatforms: []Platform{PlatformMacOS}},
				},
			},
			platform: PlatformLinux,
			want:     2,
		},
		{
			name: "some actions hidden",
			category: &Category{
				Actions: []*Action{
					{ID: "action1", HiddenOnPlatforms: []Platform{PlatformLinux}},
					{ID: "action2", HiddenOnPlatforms: []Platform{PlatformMacOS}},
				},
			},
			platform: PlatformLinux,
			want:     1,
		},
		{
			name: "all actions hidden",
			category: &Category{
				Actions: []*Action{
					{ID: "action1", HiddenOnPlatforms: []Platform{PlatformLinux}},
					{ID: "action2", HiddenOnPlatforms: []Platform{PlatformLinux}},
				},
			},
			platform: PlatformLinux,
			want:     0,
		},
		{
			name:     "no actions",
			category: &Category{},
			platform: PlatformLinux,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.category.GetVisibleActions(tt.platform)

			if len(got) != tt.want {
				t.Errorf("GetVisibleActions() returned %d items, want %d", len(got), tt.want)
			}

			// Verify all returned actions are visible on the platform
			for _, action := range got {
				if action.IsHiddenOnPlatform(tt.platform) {
					t.Errorf("returned action %q is hidden on platform %q", action.ID, tt.platform)
				}
			}
		})
	}
}

func TestCategory_GetVisibleActions_HidesCommandsWithoutResolvedPlatform(t *testing.T) {
	t.Parallel()

	category := &Category{
		Actions: []*Action{
			{
				ID:   "linux-command",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo linux"},
				},
			},
			{
				ID:   "mac-command",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformMacOS: {Command: "echo mac"},
				},
			},
			{ID: "function", Type: ActionTypeFunction},
		},
	}

	visible := category.GetVisibleActions(PlatformDebian)
	if len(visible) != 2 {
		t.Fatalf("GetVisibleActions() returned %d items, want 2", len(visible))
	}

	for _, action := range visible {
		if action.ID == "mac-command" {
			t.Fatal("command without resolved platform should be hidden")
		}
	}
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	done := make(chan bool)

	// Concurrent category registration
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 10; i++ {
			cat := &Category{
				ID:   string(rune('a' + i)),
				Name: "Category " + string(rune('a'+i)),
			}
			_ = registry.RegisterCategory("", cat)
		}
	}()

	// Concurrent action registration
	go func() {
		defer func() { done <- true }()
		// Wait a bit to ensure categories are registered
		for i := range 10 {
			action := &Action{
				ID:   "action" + string(rune('a'+i)),
				Name: "Action " + string(rune('a'+i)),
			}
			_ = registry.RegisterAction(string(rune('a'+i)), action)
		}
	}()

	// Concurrent reads
	go func() {
		defer func() { done <- true }()
		for i := range 20 {
			_ = registry.GetRoot()
			_, _ = registry.GetCategory(string(rune('a' + i%10)))
		}
	}()

	// Wait for all goroutines
	for range 3 {
		<-done
	}

	// Verify final state
	root := registry.GetRoot()
	if root == nil {
		t.Fatal("GetRoot() returned nil after concurrent access")
	}

	// Should have at least some categories registered
	if len(registry.categories) == 0 {
		t.Error("no categories registered after concurrent access")
	}
}
