package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
)

type Model struct {
	ready         bool
	viewport      viewport.Model
	title         string
	viewportLines []string
	minWidth      int
	windowWidth   int
	windowHeight  int
}

// custom messages
type (
	updateViewport struct {
		content string
		width   int
		height  int
	}
	updateViewportLine struct {
		lineNum int
		line    string
	}
)

func NewModel(title string) Model {
	return Model{title: title}
}

func (m Model) WithViewport(lines []string) Model {
	m.viewportLines = lines
	return m
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
	return updateViewport{content: content, width: width}
}

func UpdateViewportLine(lineNum int, line string) tea.Msg {
	return updateViewportLine{lineNum: lineNum, line: line}
}

func minInt(nums ...int) int {
	result := math.MaxInt
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}

func maxInt(nums ...int) int {
	result := math.MinInt
	for _, value := range nums {
		if value > result {
			result = value
		}
	}

	return result
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

// default init, does nothing
func (m Model) Init() tea.Cmd {
	return nil
}

// update handles key presses and readying the viewport
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.windowWidth, m.windowHeight = msg.Width, msg.Height
			if m.minWidth == 0 {
				m.minWidth = m.windowWidth
			}
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
			if m.viewportLines != nil {
				m.viewport.SetContent(strings.Join(m.viewportLines, "\n"))
			}

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	case updateViewport:
		m.viewportLines = strings.Split(msg.content, "\n")
		if m.ready {
			if msg.width != 0 {
				m.viewport.Width = maxInt(m.minWidth, minInt(m.windowWidth, msg.width)) + viewportStyle.GetHorizontalMargins() // add two for the margin
			}
			if msg.height != 0 {
				m.viewport.Height = maxInt(m.minWidth, minInt(m.windowHeight, msg.height))
			}
			m.viewport.SetContent(msg.content)
		}
	case updateViewportLine:
		if m.ready && msg.lineNum >= 0 && msg.lineNum < len(m.viewportLines) {
			m.viewportLines[msg.lineNum] = msg.line
			m.viewport.SetContent(strings.Join(m.viewportLines, "\n"))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// The main view renders the header, viewport and footer
func (m Model) View() string {
	return mainStyle.Render(fmt.Sprintf("%s\n%s\n%s", m.headerView(), viewportStyle.Render(m.viewport.View()), m.footerView()))
}
