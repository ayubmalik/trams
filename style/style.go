package style

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func StationName(name string) string {
	i := int(name[0]) + 134
	n := strconv.Itoa(i)
	style := lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color(n)).
		Width(30).
		Inline(true)

	return style.Render(name)
}
