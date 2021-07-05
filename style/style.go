package style

import (
	"fmt"
	"strconv"

	"github.com/ayubmalik/trams"
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

type FormattedMetrolink struct {
	trams.Metrolink
	Header  string
	Details string
}

func FormatMetrolink(m trams.Metrolink) string {
	pad := 1
	width := 52

	style := lipgloss.NewStyle().
		Bold(false).
		PaddingLeft(pad).
		PaddingRight(pad).
		Foreground(lipgloss.Color("215")).
		Background(lipgloss.Color("234")).
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("199"))

	inline := lipgloss.NewStyle().
		Inline(true).
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("215")).
		Width(width)

	text := fmt.Sprintf("%s %s (Platform %s %s)", m.TLAREF, m.StationLocation, m.Platform(), m.Direction)
	text = inline.Render(text)

	if m.Status0 == "" {
		text += "\nNo information available"
	}
	if m.Status0 != "" {
		text += fmt.Sprintf("\n%sm %s", m.Wait0, m.Dest0)
	}
	if m.Status1 != "" {
		text += fmt.Sprintf("\n%sm %s", m.Wait1, m.Dest1)
	}

	if m.Status2 != "" {
		text += fmt.Sprintf("\n%sm %s", m.Wait2, m.Dest2)
	}

	inline2 := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("230"))
	text += "\n\n" + inline2.Render(m.MessageBoard)
	return style.Render(text)
}
