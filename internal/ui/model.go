package ui

import (
	"context"
	"fmt"
	"strings"

	"anthodev/codory/internal/actions"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateMenu state = iota
	stateInput
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

	// Input handling
	inputFields   []InputField
	currentInput  int
	textInput     textinput.Model
	collectedArgs []string

	width  int
	height int
}

type InputField struct {
	Name        string
	Description string
	Required    bool
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

// executeInteractiveCommand uses tea.ExecProcess to run interactive commands
func executeInteractiveCommand(cmdStr string, action *actions.Action, executor *actions.Executor) tea.Cmd {
	// Check if the command exists already using the CheckCommand field
	platform := actions.Platform(executor.GetPlatformInfo().OS)
	if checkCmd, found := action.GetCheckCommand(platform); found {
		// Use the executor's CommandExists method to check if the command exists
		if executor.CommandExists(checkCmd) {
			return func() tea.Msg {
				return actionCompleteMsg{
					result: "Command or files already exist, skipping installation",
					err:    nil,
				}
			}
		}
	}

	c := actions.NewShellCommand(context.Background(), cmdStr)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return actionCompleteMsg{
				result: "",
				err:    err,
			}
		}
		successMsg := "Interactive command completed successfully"
		if action.SuccessMessage != "" {
			successMsg = action.SuccessMessage
		}
		return actionCompleteMsg{
			result: successMsg,
			err:    nil,
		}
	})
}

func executePackageManagerInstall(source actions.PackageSource, executor *actions.Executor) tea.Cmd {
	return func() tea.Msg {
		result, err := executor.InstallPackageManager(context.Background(), source)
		return actionCompleteMsg{result: result, err: err}
	}
}

func executeActionWithArgs(action *actions.Action, executor *actions.Executor, args []string) tea.Cmd {
	return func() tea.Msg {
		// Store args in context using your custom key
		ctx := context.WithValue(context.Background(), actions.ArgsContextKey, args)

		// Execute with the context containing the args
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
		case stateInput:
			return m.updateInput(msg)
		case statePackageManagerPrompt:
			return m.updatePackageManagerPrompt(msg)
		case stateResult:
			return m.updateResult(msg)
		}

	case actionCompleteMsg:
		m.state = stateResult
		m.result = msg.result
		m.err = msg.err

		// If package manager was installed, execute pending action
		if m.pendingAction != nil && m.needsPackageManager && msg.err == nil {
			action := m.pendingAction
			m.pendingAction = nil
			m.needsPackageManager = false
			m.packageManagerType = ""
			m.executingAction = action
			m.state = stateExecuting

			// Check if action is interactive
			if m.executor.IsInteractiveCommand(action) {
				if cmdStr, found := m.executor.GetCommandString(action); found {
					return m, executeInteractiveCommand(cmdStr, action, m.executor)
				}
			}

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

	case "enter", "right":
		return m.selectItem()

	case "esc", "backspace", "left":
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

		if needs, source := m.executor.NeedsPackageManagerInstallation(action); needs {
			m.pendingAction = action
			m.needsPackageManager = true
			m.packageManagerType = source
			m.state = statePackageManagerPrompt
			return m, nil
		}

		if len(action.Arguments) > 0 {
			// Initialize input state
			m.executingAction = action
			m.state = stateInput
			m.currentInput = 0
			m.collectedArgs = make([]string, 0, len(action.Arguments))

			// Initialize text input
			ti := textinput.New()
			ti.Placeholder = action.Arguments[0].Description
			ti.Focus()
			ti.CharLimit = 36
			m.textInput = ti

			return m, nil
		}

		// Check if this is an interactive command
		if m.executor.IsInteractiveCommand(action) {
			if cmdStr, found := m.executor.GetCommandString(action); found {
				m.executingAction = action
				m.state = stateExecuting
				return m, executeInteractiveCommand(cmdStr, action, m.executor)
			}
		}

		m.executingAction = action
		m.state = stateExecuting
		return m, executeAction(action, m.executor)
	}

	return m, nil
}

func (m Model) updatePackageManagerPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "n", "N":
		m.pendingAction = nil
		m.needsPackageManager = false
		m.packageManagerType = ""
		m.state = stateMenu
		return m, nil
	case "enter", "y", "Y":
		m.executingAction = &actions.Action{Name: fmt.Sprintf("Install %s", packageManagerName(m.packageManagerType))}
		m.state = stateExecuting
		return m, executePackageManagerInstall(m.packageManagerType, m.executor)
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

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		// Cancel and go back to menu
		m.state = stateMenu
		m.executingAction = nil
		m.collectedArgs = nil
		return m, nil

	case "enter":
		// Save current input
		m.collectedArgs = append(m.collectedArgs, m.textInput.Value())
		m.currentInput++

		// Check if we have all arguments
		if m.currentInput >= len(m.executingAction.Arguments) {
			if needs, source := m.executor.NeedsPackageManagerInstallation(m.executingAction); needs {
				m.pendingAction = m.executingAction
				m.needsPackageManager = true
				m.packageManagerType = source
				m.state = statePackageManagerPrompt
				return m, nil
			}

			// Execute the action with collected args
			m.state = stateExecuting

			// Check if this is an interactive command
			if m.executor.IsInteractiveCommand(m.executingAction) {
				if cmdStr, found := m.executor.GetCommandString(m.executingAction); found {
					return m, executeInteractiveCommand(cmdStr, m.executingAction, m.executor)
				}
			}

			return m, executeActionWithArgs(m.executingAction, m.executor, m.collectedArgs)
		}

		// Move to next argument
		m.textInput.SetValue("")
		m.textInput.Placeholder = m.executingAction.Arguments[m.currentInput].Description
		return m, nil
	}

	// Update text input
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
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
	var content string

	switch m.state {
	case stateMenu:
		content = m.viewMenu()
	case stateInput:
		content = m.viewInput()
	case stateExecuting:
		content = m.viewExecuting()
	case statePackageManagerPrompt:
		content = m.viewPackageManagerPrompt()
	case stateResult:
		content = m.viewResult()
	default:
		content = ""
	}

	rendered := appStyle.
		Width(safeSize(m.width-2, 90)).
		Render(content)

	return capRenderedHeight(rendered, m.height)
}

func (m Model) viewMenu() string {
	platformInfo := m.executor.GetPlatformInfo()
	platformOS := actions.Platform(platformInfo.OS)

	visibleSubCats := m.currentCategory.GetVisibleSubCategories(platformOS)
	visibleActions := m.currentCategory.GetVisibleActions(platformOS)
	contentWidth := safeSize(m.width-22, 90)
	paneHeight := terminalPaneHeight(m.height, 18)
	navWidth := clampInt(contentWidth/4, 22, 30)
	detailWidth := clampInt(contentWidth/3, 28, 42)
	listWidth := contentWidth - navWidth - detailWidth - 4
	if listWidth < 28 {
		listWidth = 28
		detailWidth = maxInt(24, contentWidth-navWidth-listWidth-4)
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		titleStyle.Render("🛠️  Codory"),
		" ",
		platformStyle.Render(fmt.Sprintf("%s", platformInfo.OS.DisplayName())),
	)

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.viewNavigationPane(navWidth, paneHeight),
		"  ",
		m.viewActionListPane(listWidth, paneHeight, visibleSubCats, visibleActions),
		"  ",
		m.viewDetailPane(detailWidth, paneHeight, visibleSubCats, visibleActions),
	)

	footer := helpStyle.Render("↑/↓ or j/k navigate • enter/→ select • esc/← back • q quit")
	return lipgloss.JoinVertical(lipgloss.Left, header, "", panes, "", footer)
}

func (m Model) viewNavigationPane(width, height int) string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Path"))
	s.WriteString("\n")

	if len(m.categoryStack) == 0 {
		s.WriteString(selectedStyle.Render("❯ " + m.currentCategory.Name))
	} else {
		for _, cat := range m.categoryStack {
			s.WriteString(breadcrumbStyle.Render("  " + cat.Name))
			s.WriteString("\n")
		}
		s.WriteString(selectedStyle.Render("❯ " + m.currentCategory.Name))
	}

	return panelStyle.Width(width).Render(capLines(s.String(), paneContentLines(height)))
}

func (m Model) viewActionListPane(width, height int, visibleSubCats []*actions.Category, visibleActions []*actions.Action) string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Actions"))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render(m.currentCategory.Description))
	s.WriteString("\n\n")

	for i, cat := range visibleSubCats {
		if m.cursor == i {
			s.WriteString(selectedStyle.Render(fmt.Sprintf("❯ 📁 %s", cat.Name)))
		} else {
			s.WriteString(categoryStyle.Render(fmt.Sprintf("  📁 %s", cat.Name)))
		}
		s.WriteString("\n")
	}

	offset := len(visibleSubCats)
	for i, action := range visibleActions {
		idx := offset + i
		if m.cursor == idx {
			s.WriteString(selectedStyle.Render(fmt.Sprintf("❯ ▶ %s", action.Name)))
		} else {
			s.WriteString(actionStyle.Render(fmt.Sprintf("  ▶ %s", action.Name)))
		}
		s.WriteString("\n")
	}

	if len(visibleSubCats) == 0 && len(visibleActions) == 0 {
		s.WriteString(warningStyle.Render("No actions available for this platform"))
	}

	return activePanelStyle.Width(width).Render(capLines(s.String(), paneContentLines(height)))
}

func (m Model) viewDetailPane(width, height int, visibleSubCats []*actions.Category, visibleActions []*actions.Action) string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Details"))
	s.WriteString("\n\n")

	if cat := selectedCategory(m.cursor, visibleSubCats); cat != nil {
		s.WriteString(categoryStyle.Render("Category"))
		s.WriteString("\n")
		s.WriteString(cat.Name)
		s.WriteString("\n\n")
		s.WriteString(subtitleStyle.Render(emptyFallback(cat.Description, "Open category")))
	} else if action := selectedAction(m.cursor, visibleSubCats, visibleActions); action != nil {
		s.WriteString(subtitleStyle.Render(emptyFallback(action.Description, "No description")))
		if cmd, ok := action.ResolvePlatformCommand(actions.Platform(m.executor.GetPlatformInfo().OS)); ok && cmd.PackageSource != "" {
			s.WriteString("\n\n")
			s.WriteString(promptStyle.Render(fmt.Sprintf("Source: %s", cmd.PackageSource)))
		}
		if len(action.Arguments) > 0 {
			s.WriteString("\n")
			s.WriteString(promptStyle.Render(fmt.Sprintf("Inputs: %d", len(action.Arguments))))
		}
	} else {
		s.WriteString(warningStyle.Render("Nothing selected"))
	}

	return panelStyle.Width(width).Render(capLines(s.String(), paneContentLines(height)))
}

func (m Model) viewInput() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(m.executingAction.Name))
	s.WriteString("\n\n")

	arg := m.executingAction.Arguments[m.currentInput]
	s.WriteString(promptStyle.Render(fmt.Sprintf("Enter %s", arg.Name)))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render(arg.Description))
	s.WriteString("\n\n")

	s.WriteString(inputStyle.Width(safeSize(m.width-12, 50)).Render(m.textInput.View()))
	s.WriteString("\n\n")

	s.WriteString(helpStyle.Render("enter: confirm • esc: cancel"))

	return centeredPanel(s.String(), m.width, m.height)
}

func (m Model) viewExecuting() string {
	message := "Executing..."
	if m.executingAction == nil {
		return centeredPanel(promptStyle.Render(message), m.width, m.height)
	}
	message = fmt.Sprintf("Executing: %s...", m.executingAction.Name)
	return centeredPanel(promptStyle.Render(message), m.width, m.height)
}

func (m Model) viewPackageManagerPrompt() string {
	var s strings.Builder

	actionName := "selected action"
	if m.pendingAction != nil {
		actionName = m.pendingAction.Name
	}

	s.WriteString(titleStyle.Render("Package manager required"))
	s.WriteString("\n\n")
	s.WriteString(subtitleStyle.Render(fmt.Sprintf("%s requires %s.", actionName, packageManagerName(m.packageManagerType))))
	s.WriteString("\n")
	s.WriteString(promptStyle.Render("Install it now?"))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("y/enter: install • n/esc: cancel • q: quit"))

	return centeredPanel(s.String(), m.width, m.height)
}

func packageManagerName(source actions.PackageSource) string {
	switch source {
	case actions.PackageSourceAUR:
		return "Yay"
	case actions.PackageSourceBrew:
		return "Homebrew"
	case actions.PackageSourceWinget:
		return "Winget"
	default:
		return string(source)
	}
}

func (m Model) viewResult() string {
	var s strings.Builder

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
	s.WriteString(helpStyle.Render("Press enter to continue or q/ctrl+c to quit..."))

	return centeredPanel(s.String(), m.width, m.height)
}

func centeredPanel(content string, width, height int) string {
	panelWidth := safeSize(width-12, 60)
	paneHeight := terminalPaneHeight(height, 18)
	return activePanelStyle.Width(panelWidth).Render(capLines(content, paneContentLines(paneHeight)))
}

func terminalPaneHeight(terminalHeight, fallback int) int {
	if terminalHeight <= 0 {
		return fallback
	}
	return maxInt(1, terminalHeight-6)
}

func capRenderedHeight(content string, maxHeight int) string {
	if maxHeight <= 0 || lipgloss.Height(content) <= maxHeight {
		return content
	}
	return capLines(content, maxHeight)
}

func paneContentLines(height int) int {
	return maxInt(1, height-4)
}

func capLines(content string, maxLines int) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) <= maxLines {
		return content
	}
	return strings.Join(lines[:maxLines], "\n")
}

func selectedCategory(cursor int, cats []*actions.Category) *actions.Category {
	if cursor < 0 || cursor >= len(cats) {
		return nil
	}
	return cats[cursor]
}

func selectedAction(cursor int, cats []*actions.Category, actionList []*actions.Action) *actions.Action {
	idx := cursor - len(cats)
	if idx < 0 || idx >= len(actionList) {
		return nil
	}
	return actionList[idx]
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func safeSize(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func clampInt(value, minValue, maxValue int) int {
	return minInt(maxInt(value, minValue), maxValue)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
