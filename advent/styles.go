package advent

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2025/advent/color"
)

var (
	// the style of the solution text
	SolutionStyle        = lipgloss.NewStyle().Foreground(color.Aquamarine86)
	StepStyle            = lipgloss.NewStyle().Foreground(color.LightYellow11)
	DataStyle            = lipgloss.NewStyle().Foreground(color.AmericanOrange208)
	correctResultStyle   = lipgloss.NewStyle().Foreground(color.StrongLimeGreen40)
	incorrectResultStyle = lipgloss.NewStyle().Foreground(color.LightRed9)

	// some common styles
	WallStyle    = lipgloss.NewStyle().Foreground(color.BlazeOrange202)
	VisitedStyle = lipgloss.NewStyle().Foreground(color.Aqua14).Background(color.FreshEggplant90)

	BoxStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color(color.BrightGreen82))
	BoxHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color.BrightGreen82)).Background(color.CarnationPink19)
)
