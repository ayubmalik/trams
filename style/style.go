package style

import (
	"fmt"
	"strconv"
	"strings"

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

func FormatMetrolink(m trams.Metrolink, colorIndex int) string {
	pad := 1
	width := 34
	color := strconv.Itoa(161 + colorIndex*6) // ansi color index rainbow effect
	style := lipgloss.NewStyle().
		Bold(false).
		PaddingLeft(pad).
		PaddingRight(pad).
		Foreground(lipgloss.Color(color)).
		// Background(lipgloss.Color("234")).
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color))

	inline := style.Copy().
		Inline(true).
		Bold(true).
		Reverse(true)

	text := inline.Render(fmt.Sprintf("%2s %s", m.TLAREF, strings.ToUpper(m.StationLocation))) + "\n"
	text += inline.Render(fmt.Sprintf("Platform %s (%s)", m.Platform(), m.Direction))

	if m.Status0 == "" {
		text += "\nNo information available"
	}
	if m.Status0 != "" {
		text += fmt.Sprintf("\n%2sm %s", m.Wait0, m.Dest0)
	}
	if m.Status1 != "" {
		text += fmt.Sprintf("\n%2sm %s", m.Wait1, m.Dest1)
	}

	if m.Status2 != "" {
		text += fmt.Sprintf("\n%2sm %s", m.Wait2, m.Dest2)
	}

	inline2 := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("230"))
	text += "\n\n" + inline2.Render(m.MessageBoard)
	return style.Render(text)
}
