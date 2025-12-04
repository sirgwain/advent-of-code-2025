package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2025/advent/color"
)

var (
	// base styles
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
	viewportStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	solutionStyle = lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(color.CornflowerBlue63)
)

type Model struct {
	ready        bool
	viewport     viewport.Model
	solution     string
	title        string
	minWidth     int
	windowWidth  int
	windowHeight int
}

// custom messages
type (
	updateViewportMsg struct {
		content string
		width   int
		height  int
	}
	updateSolutionMsg struct {
		solution string
	}
)

func NewModel(title string) Model {
	return Model{title: title}
}

func (m Model) WithMinWidth(minWidth int) Model {
	m.minWidth = minWidth
	return m
}

func NewViewportProgram(initialModel Model) *tea.Program {
	return tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
}

func UpdateViewport(content string, width int) tea.Msg {
	return updateViewportMsg{content: content, width: width}
}

func UpdateSolution(solution string) tea.Msg {
	return updateSolutionMsg{solution: solution}
}

func (m Model) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	line := strings.Repeat("─", max(0, m.viewport.Width))
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func (m Model) solutionView() string {
	return solutionStyle.Render(m.solution)
}

// default init, does nothing
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key presses, window size, and viewport/tick updates.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		solutionHeight := lipgloss.Height(m.solutionView())
		verticalMarginHeight := headerHeight + footerHeight + solutionHeight

		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		if m.minWidth == 0 {
			m.minWidth = m.windowWidth
		}

		if !m.ready {
			m.viewport = viewport.New(
				msg.Width,
				msg.Height-verticalMarginHeight, // no extra -1
			)
			// Render viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case updateViewportMsg:
		if m.ready {
			if msg.width != 0 {
				m.viewport.Width = max(m.minWidth, min(m.windowWidth, msg.width)) +
					viewportStyle.GetHorizontalMargins()
			}
			if msg.height != 0 {
				m.viewport.Height = max(m.minWidth, min(m.windowHeight, msg.height))
			}
			m.viewport.SetContent(msg.content)
		}

	case updateSolutionMsg:
		m.solution = msg.solution
	}

	var vcmd tea.Cmd
	m.viewport, vcmd = m.viewport.Update(msg)
	if vcmd != nil {
		cmds = append(cmds, vcmd)
	}

	return m, tea.Batch(cmds...)
}

// The main view renders the header, viewport and footer
func (m Model) View() string {
	return mainStyle.Render(fmt.Sprintf("%s\n%s\n%s\n%s",
		m.headerView(),
		viewportStyle.Render(m.viewport.View()),
		m.solutionView(),
		m.footerView(),
	))
}
