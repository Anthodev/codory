package actions

import (
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if registry.root == nil {
		t.Error("Expected root category to be initialized")
	}
	if registry.root.ID != "root" {
		t.Errorf("Expected root category ID to be 'root', got '%s'", registry.root.ID)
	}
	if registry.categories == nil {
		t.Error("Expected categories map to be initialized")
	}
}

func TestGlobalRegistry(t *testing.T) {
	registry1 := GlobalRegistry()
	registry2 := GlobalRegistry()

	if registry1 != registry2 {
		t.Error("GlobalRegistry() should return the same instance")
	}
}

func TestRegistry_RegisterCategory(t *testing.T) {
	tests := []struct {
		name        string
		parent      string
		category    *Category
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
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:   "register with non-existent parent",
			parent: "non-existent",
			category: &Category{
				ID:   "test-cat2",
				Name: "Test Category 2",
			},
			wantErr:     true,
			errContains: "parent category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry()

			// First registration (for duplicate test)
			if tt.name == "register duplicate category" {
				registry.RegisterCategory(tt.parent, tt.category)
			}

			err := registry.RegisterCategory(tt.parent, tt.category)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify category was registered
				if _, exists := registry.categories[tt.category.ID]; !exists {
					t.Error("Category was not registered in categories map")
				}

				// Verify category was added to parent's subcategories
				parentCat := registry.root
				if tt.parent != "" && tt.parent != "root" {
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
					t.Error("Category was not added to parent's subcategories")
				}
			}
		})
	}
}

func TestRegistry_RegisterAction(t *testing.T) {
	tests := []struct {
		name        string
		categoryID  string
		action      *Action
		setup       func(*Registry)
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
			setup: func(r *Registry) {
				r.RegisterCategory("", &Category{
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
			registry := NewRegistry()

			if tt.setup != nil {
				tt.setup(registry)
			}

			err := registry.RegisterAction(tt.categoryID, tt.action)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify action was added to category
				category, _ := registry.categories[tt.categoryID]
				found := false
				for _, action := range category.Actions {
					if action.ID == tt.action.ID {
						found = true
						break
					}
				}
				if !found {
					t.Error("Action was not added to category")
				}
			}
		})
	}
}

func TestRegistry_GetRoot(t *testing.T) {
	registry := NewRegistry()
	root := registry.GetRoot()

	if root == nil {
		t.Fatal("GetRoot() returned nil")
	}
	if root.ID != "root" {
		t.Errorf("Expected root category ID to be 'root', got '%s'", root.ID)
	}
}

func TestRegistry_GetCategory(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setup      func(*Registry)
		wantExists bool
	}{
		{
			name: "get existing category",
			id:   "test-cat",
			setup: func(r *Registry) {
				r.RegisterCategory("", &Category{
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
			registry := NewRegistry()

			if tt.setup != nil {
				tt.setup(registry)
			}

			category, exists := registry.GetCategory(tt.id)

			if tt.wantExists {
				if !exists {
					t.Error("Expected category to exist")
				}
				if category == nil {
					t.Error("Expected category to be non-nil")
				}
			} else {
				if exists {
					t.Error("Expected category to not exist")
				}
				if category != nil {
					t.Error("Expected category to be nil")
				}
			}
		})
	}
}

func TestCategory_GetVisibleSubCategories(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.category.GetVisibleSubCategories(tt.platform)
			if len(got) != tt.want {
				t.Errorf("GetVisibleSubCategories() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestCategory_GetVisibleActions(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.category.GetVisibleActions(tt.platform)
			if len(got) != tt.want {
				t.Errorf("GetVisibleActions() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}
