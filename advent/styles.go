package advent

import "github.com/charmbracelet/lipgloss"

var (
	// the style of the solution text
	solutionStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	correctResultStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("40"))
	incorrectResultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	// day6 map
	guardStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	obstacleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	pathStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))

	// day7 numbers and operators
	numberStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	operatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))

	// day8 map
	antennaStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	antennaWithAntinodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("87"))
	antinodeStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
)
