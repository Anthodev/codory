package actions

import (
	"fmt"
	"sync"
)

// Registry handles the registration of categories and actions
type Registry struct {
	mu         sync.RWMutex
	categories map[string]*Category
	root       *Category
}

var globalRegistry = NewRegistry()

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		categories: make(map[string]*Category),
		root: &Category{
			ID:            "root",
			Name:          "Main Menu",
			Description:   "Select a category",
			SubCategories: make([]*Category, 0),
			Actions:       make([]*Action, 0),
		},
	}
}

// GlobalRegistry returns the global registry
func GlobalRegistry() *Registry {
	return globalRegistry
}

// RegisterCategory registers a category
func (r *Registry) RegisterCategory(parent string, category *Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[category.ID]; exists {
		return fmt.Errorf("category %s already exists", category.ID)
	}

	r.categories[category.ID] = category

	if parent == "" || parent == "root" {
		r.root.SubCategories = append(r.root.SubCategories, category)
	} else {
		parentCat, exists := r.categories[parent]
		if !exists {
			return fmt.Errorf("parent category %s not found", parent)
		}
		parentCat.SubCategories = append(parentCat.SubCategories, category)
	}

	return nil
}

// RegisterAction registers an action in a category
func (r *Registry) RegisterAction(categoryID string, action *Action) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	category, exists := r.categories[categoryID]
	if !exists {
		return fmt.Errorf("category %s not found", categoryID)
	}

	category.Actions = append(category.Actions, action)
	return nil
}

// GetRoot returns the root category
func (r *Registry) GetRoot() *Category {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.root
}

// GetCategory returns a category by its ID
func (r *Registry) GetCategory(id string) (*Category, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cat, exists := r.categories[id]
	return cat, exists
}

// GetVisibleSubCategories returns the visible subcategories for a platform
func (c *Category) GetVisibleSubCategories(platform Platform) []*Category {
	visible := make([]*Category, 0)
	for _, subCat := range c.SubCategories {
		if subCat.IsVisibleOnPlatform(platform) {
			visible = append(visible, subCat)
		}
	}
	return visible
}

// GetVisibleActions returns the visible actions for a platform
func (c *Category) GetVisibleActions(platform Platform) []*Action {
	visible := make([]*Action, 0)
	for _, action := range c.Actions {
		if action.IsVisibleOnPlatform(platform) {
			visible = append(visible, action)
		}
	}
	return visible
}
