package advent

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2025/advent/color"
)

var (
	// the style of the solution text
	solutionStyle        = lipgloss.NewStyle().Foreground(color.Aquamarine86)
	StepStyle            = lipgloss.NewStyle().Foreground(color.LightYellow11)
	data1Style           = lipgloss.NewStyle().Foreground(color.LightCobaltBlue110)
	data2Style           = lipgloss.NewStyle().Foreground(color.Coral209)
	correctResultStyle   = lipgloss.NewStyle().Foreground(color.StrongLimeGreen40)
	incorrectResultStyle = lipgloss.NewStyle().Foreground(color.LightRed9)

	// some common styles
	wallStyle    = lipgloss.NewStyle().Foreground(color.BlazeOrange202)
	visitedStyle = lipgloss.NewStyle().Foreground(color.Aqua14).Background(color.FreshEggplant90)

	boxStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color(color.BrightGreen82))
	boxHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color.BrightGreen82)).Background(color.BlazeOrange202)
)
