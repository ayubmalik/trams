package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ayubmalik/trams"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorStart = 130
	lineWidth  = 120
	pad        = 1
)

type FormattedMetrolink struct {
	text string
}

func (fm FormattedMetrolink) String() string {
	return fm.text
}

func FormatMetrolink(m trams.Metrolink, colorIndex int) FormattedMetrolink {
	width := 34
	color := strconv.Itoa(colorStart + colorIndex*6) // ansi color index rainbow effect
	mainStyle := lipgloss.NewStyle().
		Bold(false).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color)).
		Foreground(lipgloss.Color(color)).
		Inline(false).
		PaddingLeft(pad).
		PaddingRight(pad).
		Reverse(false).
		Width(width)

	inline := mainStyle.Copy().
		Inline(true).
		Bold(true).
		Reverse(true)

	text := inline.Render(fmt.Sprintf("%2s %s", m.TLAREF, strings.ToUpper(m.StationLocation)))
	// text += inline.Render(fmt.Sprintf("Platform %s (%s)", m.Platform(), m.Direction))

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

	// inline2 := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("230"))
	// text += "\n\n" + inline2.Render(m.MessageBoard)
	return FormattedMetrolink{text: mainStyle.Render(text)}
}

type FormattedStationID struct {
	text string
}

func (fs FormattedStationID) String() string {
	return fs.text
}

func FormatStationID(stationID string, colorIndex int) string {
	n := strconv.Itoa(colorIndex)
	style1 := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(n))
	style2 := style1.Copy().
		Bold(true)

	if stationID == "" {
		return stationID
	}

	return style1.Render(stationID[:3]) + style2.Render(stationID[3:])
}

func StationIDRows(stationIDs []string) []string {
	cols := 3
	count := len(stationIDs) - 1
	rows := make([]string, 0)
	row := ""
	colorIndex := 161
	for i, s := range stationIDs {
		var left, middle, right string
		left = padRight(s)
		if (i + 1) <= count {
			middle = padRight(stationIDs[i+1])
		}
		if (i + 2) <= count {
			right = padRight(stationIDs[i+2])
		}
		row = lipgloss.JoinHorizontal(lipgloss.Right, FormatStationID(left, colorIndex),
			FormatStationID(middle, colorIndex), FormatStationID(right, colorIndex))

		if i%cols == 0 {
			rows = append(rows, row)
			row = ""
		}

		if i%(cols*6) == 0 {
			colorIndex += 6
			if colorIndex > 231 {
				colorIndex = 161
			}
		}
	}
	return rows
}

func padRight(s string) string {
	return fmt.Sprintf("%-*s", lineWidth/3, s)
}
