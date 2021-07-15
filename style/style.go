package style

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ayubmalik/trams"
	"github.com/charmbracelet/lipgloss"
)

const (
	boxWidth    = 40
	lineWidth   = 120
	nl          = "\n"
	pad         = 1
	space       = " "
	startColour = 161
)

type FormattedMetrolink struct {
	text string
}

func (fm FormattedMetrolink) String() string {
	return fm.text
}

func FormatMetrolink(m trams.Metrolink, colorIndex, height int) FormattedMetrolink {
	color := strconv.Itoa(colorIndex) // ansi color index rainbow effect
	mainStyle := lipgloss.NewStyle().
		Bold(false).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color)).
		Foreground(lipgloss.Color(color)).
		Inline(false).
		PaddingLeft(pad).
		PaddingRight(pad).
		Reverse(false).
		Height(height).
		Width(boxWidth)

	inline := mainStyle.Copy().
		Inline(true).
		Bold(true).
		Height(1).
		Reverse(true)

	text := inline.Render(fmt.Sprintf("%3s %s (Plat.%s)", m.TLAREF, strings.ToUpper(m.StationLocation), m.Platform()))

	if m.Status0 == "" {
		text += "\nNo information available"
	} else {
		text += fmt.Sprintf("\n%02sm %s", m.Wait0, m.Dest0)
	}

	if m.Status1 == "" {
		text += nl
	} else {
		text += fmt.Sprintf("\n%02sm %s", m.Wait1, m.Dest1)
	}

	if m.Status2 == "" {
		text += nl
	} else {
		text += fmt.Sprintf("\n%02sm %s", m.Wait2, m.Dest2)
	}

	if m.Status3 == "" {
		text += nl
	} else {
		text += fmt.Sprintf("\n%02sm %s", m.Wait3, m.Dest3)
	}

	text += FormatMessageBoard(m.MessageBoard, boxWidth-2)
	return FormattedMetrolink{text: mainStyle.Render(text)}
}

func FormatMessageBoard(msg string, width int) string {
	if msg == "<no message>" {
		return ""
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Faint(true).Italic(true).Inline(true)
	chunks := chunk(msg, width)
	text := nl
	for _, c := range chunks {
		re := regexp.MustCompile(`^\s+`)
		c := re.ReplaceAllString(c, "")
		text += style.Render(c) + nl
	}
	return text
}

func FormatStationID(s trams.StationID) string {
	return fmt.Sprintf("| %s | %-40s", s.TLAREF, s.StationLocation)
}

func MetrolinkRows(metrolinks []trams.Metrolink) []string {
	cols := 3
	groupedMetrolinks := groupMetrolinksByRef(metrolinks)
	rows := make([]string, 0)
	colorIndex := startColour

	for _, v := range groupedMetrolinks {
		count := len(v)
		formattetMetrolinks := make([]string, 0, cols)

		maxHeight := 1
		for _, metrolink := range v {
			renderHeight := lipgloss.Height(FormatMetrolink(metrolink, colorIndex, 1).String())
			if renderHeight > maxHeight {
				maxHeight = renderHeight - 2
			}
		}

		// render again with same heights
		for _, metrolink := range v {
			formattetMetrolinks = append(formattetMetrolinks, FormatMetrolink(metrolink, colorIndex, maxHeight).String()+space)

		}

		rowCount := 1 + (count-1)/cols
		for i := 0; i < rowCount; i++ {
			lower := i * cols
			upper := (1 + i) * cols
			row := lipgloss.JoinHorizontal(lipgloss.Top, formattetMetrolinks[lower:upper]...)
			rows = append(rows, row)
		}

		rows = append(rows, "")
		colorIndex = nextColour(colorIndex)
	}
	return rows
}

func StationIDRows(stationIDs []string) []string {
	cols := 3
	count := len(stationIDs) - 1
	rows := make([]string, 0)
	row := ""
	colourIndex := startColour
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
			colorizeStationID(left, colourIndex),
			colorizeStationID(middle, colourIndex),
			colorizeStationID(right, colourIndex),
		)

		if i%cols == 0 {
			rows = append(rows, row)
			row = ""
		}

		if i%(cols*6) == 0 {
			colourIndex = nextColour(colourIndex)
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

func nextColour(n int) int {
	c := n + 6
	if c > 231 {
		c = startColour
	}
	return c
}

func chunk(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
