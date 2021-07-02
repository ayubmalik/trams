package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ayubmalik/trams"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorStart = 161
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
	width := 40
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

	text := inline.Render(fmt.Sprintf("%3s %s (Plat.%s)", m.TLAREF, strings.ToUpper(m.StationLocation), m.Platform()))
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

func FormatStationID(s trams.StationID) string {
	return fmt.Sprintf("| %s | %-40s", s.TLAREF, s.StationLocation)
}

func MetrolinkRows(metrolinks []trams.Metrolink) []string {
	groupedMetrolinks := groupMetrolinksByRef(metrolinks)
	rows := make([]string, 0)
	colorIndex := colorStart
	for _, v := range groupedMetrolinks {
		count := len(v) - 1
		var left, middle, right FormattedMetrolink
		left = FormatMetrolink(v[0], colorIndex)
		if count > 1 {
			middle = FormatMetrolink(v[1], colorIndex)
		}
		if count > 2 {
			middle = FormatMetrolink(v[2], colorIndex)
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Right,
			left.String(),
			middle.String(),
			right.String(),
		)
		rows = append(rows, row)
	}
	return rows
}

func StationIDRows(stationIDs []string) []string {
	cols := 3
	count := len(stationIDs) - 1
	rows := make([]string, 0)
	row := ""
	colorIndex := colorStart
	for i, s := range stationIDs {
		var left, middle, right string
		left = padRight(s)
		if (i + 1) <= count {
			middle = padRight(stationIDs[i+1])
		}
		if (i + 2) <= count {
			right = padRight(stationIDs[i+2])
		}
		row = lipgloss.JoinHorizontal(
			lipgloss.Right,
			colorizeStationID(left, colorIndex),
			colorizeStationID(middle, colorIndex),
			colorizeStationID(right, colorIndex),
		)

		if i%cols == 0 {
			rows = append(rows, row)
			row = ""
		}

		if i%(cols*6) == 0 {
			colorIndex += 6
			if colorIndex > 231 {
				colorIndex = colorStart
			}
		}
	}
	return rows
}

func colorizeStationID(stationID string, colorIndex int) string {
	if stationID == "" {
		return stationID
	}

	n := strconv.Itoa(colorIndex)
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(n)).Render(stationID)
}

func padRight(s string) string {
	return fmt.Sprintf("%-*s", lineWidth/3, s)
}

func groupMetrolinksByRef(metrolinks []trams.Metrolink) map[string][]trams.Metrolink {
	gm := make(map[string][]trams.Metrolink)
	for _, m := range metrolinks {
		gm[m.TLAREF] = append(gm[m.TLAREF], m)
	}
	return gm
}
