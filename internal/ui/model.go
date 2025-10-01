package ui

import (
	"context"
	"fmt"
	"strings"

	"anthodev/codory/internal/actions"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateMenu state = iota
	stateExecuting
	stateResult
	statePackageManagerPrompt
)

type Model struct {
	state    state
	registry *actions.Registry
	executor *actions.Executor

	// Navigation
	currentCategory *actions.Category
	categoryStack   []*actions.Category
	cursor          int

	// Execution
	executingAction *actions.Action
	result          string
	err             error

	// Package manager installation
	needsPackageManager bool
	packageManagerType  actions.PackageSource
	pendingAction       *actions.Action

	width  int
	height int
}

func NewModel() Model {
	registry := actions.GlobalRegistry()

	return Model{
		state:           stateMenu,
		registry:        registry,
		executor:        actions.NewExecutor(),
		currentCategory: registry.GetRoot(),
		categoryStack:   make([]*actions.Category, 0),
		cursor:          0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type actionCompleteMsg struct {
	result string
	err    error
}

func executeAction(action *actions.Action, executor *actions.Executor) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		result, err := executor.Execute(ctx, action)
		return actionCompleteMsg{result: result, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateMenu:
			return m.updateMenu(msg)
		case stateResult:
			return m.updateResult(msg)
		}

	case actionCompleteMsg:
		m.state = stateResult
		m.result = msg.result
		m.err = msg.err

		// Si c'était l'installation d'un package manager, exécuter l'action en attente
		if m.pendingAction != nil && m.needsPackageManager && msg.err == nil {
			action := m.pendingAction
			m.pendingAction = nil
			m.needsPackageManager = false
			m.packageManagerType = ""
			m.executingAction = action
			m.state = stateExecuting
			return m, executeAction(action, m.executor)
		}

		m.pendingAction = nil
		m.needsPackageManager = false
		m.packageManagerType = ""
		return m, nil
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		visibleSubCats := m.currentCategory.GetVisibleSubCategories(actions.Platform(m.executor.GetPlatformInfo().OS))
		visibleActions := m.currentCategory.GetVisibleActions(actions.Platform(m.executor.GetPlatformInfo().OS))
		maxItems := len(visibleSubCats) + len(visibleActions)
		if m.cursor < maxItems-1 {
			m.cursor++
		}

	case "enter":
		return m.selectItem()

	case "esc", "backspace":
		return m.goBack()
	}

	return m, nil
}

func (m Model) selectItem() (tea.Model, tea.Cmd) {
	platformOS := actions.Platform(m.executor.GetPlatformInfo().OS)
	visibleSubCats := m.currentCategory.GetVisibleSubCategories(platformOS)
	visibleActions := m.currentCategory.GetVisibleActions(platformOS)

	subCatCount := len(visibleSubCats)

	// Selection of a sub-category
	if m.cursor < subCatCount {
		selectedCat := visibleSubCats[m.cursor]
		m.categoryStack = append(m.categoryStack, m.currentCategory)
		m.currentCategory = selectedCat
		m.cursor = 0
		return m, nil
	}

	// Selection an action
	actionIndex := m.cursor - subCatCount
	if actionIndex < len(visibleActions) {
		action := visibleActions[actionIndex]

		m.executingAction = action
		m.state = stateExecuting
		return m, executeAction(action, m.executor)
	}

	return m, nil
}

func (m Model) goBack() (tea.Model, tea.Cmd) {
	if len(m.categoryStack) == 0 {
		return m, nil
	}

	// Retour à la catégorie parente
	m.currentCategory = m.categoryStack[len(m.categoryStack)-1]
	m.categoryStack = m.categoryStack[:len(m.categoryStack)-1]
	m.cursor = 0

	return m, nil
}

func (m Model) updateResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter", "esc":
		m.state = stateMenu
		m.result = ""
		m.err = nil
		m.executingAction = nil
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateExecuting:
		return m.viewExecuting()
	case stateResult:
		return m.viewResult()
	}
	return ""
}

func (m Model) viewMenu() string {
	var s strings.Builder

	platformInfo := m.executor.GetPlatformInfo()
	platformOS := actions.Platform(platformInfo.OS)

	// Titre
	s.WriteString(titleStyle.Render("🛠️  Codory"))
	s.WriteString(" ")
	s.WriteString(platformStyle.Render(fmt.Sprintf("(%s)", platformInfo.OS.DisplayName())))
	s.WriteString("\n\n")

	// Breadcrumb
	if len(m.categoryStack) > 0 {
		breadcrumb := make([]string, 0, len(m.categoryStack)+1)
		for _, cat := range m.categoryStack {
			breadcrumb = append(breadcrumb, cat.Name)
		}
		breadcrumb = append(breadcrumb, m.currentCategory.Name)
		s.WriteString(breadcrumbStyle.Render(strings.Join(breadcrumb, " > ")))
		s.WriteString("\n\n")
	}

	visibleSubCats := m.currentCategory.GetVisibleSubCategories(platformOS)
	visibleActions := m.currentCategory.GetVisibleActions(platformOS)

	// Subcategories
	for i, cat := range visibleSubCats {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s.WriteString(selectedStyle.Render(fmt.Sprintf("%s 📁 %s", cursor, cat.Name)))
		} else {
			s.WriteString(categoryStyle.Render(fmt.Sprintf("%s 📁 %s", cursor, cat.Name)))
		}
		s.WriteString("\n")
	}

	// Actions
	offset := len(visibleSubCats)
	for i, action := range visibleActions {
		cursor := " "
		idx := offset + i

		// Vérifier si l'action nécessite un package manager
		actionText := fmt.Sprintf("%s ▶️  %s", cursor, action.Name)

		if m.cursor == idx {
			cursor = ">"
			s.WriteString(selectedStyle.Render(actionText))
		} else {
			s.WriteString(actionStyle.Render(actionText))
		}
		s.WriteString("\n")
	}

	// Message if no actions available
	if len(visibleSubCats) == 0 && len(visibleActions) == 0 {
		s.WriteString(warningStyle.Render("No actions available for this platform"))
		s.WriteString("\n")
	}

	// Help
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"))

	return s.String()
}

func (m Model) viewExecuting() string {
	return fmt.Sprintf("\n  Executing: %s...\n\n", m.executingAction.Name)
}

func (m Model) viewResult() string {
	var s strings.Builder

	s.WriteString("\n")
	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
		s.WriteString("\n\n")
	} else {
		s.WriteString(successStyle.Render("✅ Success!"))
		s.WriteString("\n\n")
	}

	if m.result != "" {
		s.WriteString(resultStyle.Render("Result:"))
		s.WriteString("\n")
		s.WriteString(m.result)
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Press enter to continue..."))

	return s.String()
}
