package ui

import (
	"atlas.sql/internal/db"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("205"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type focus int

const (
	focusInput focus = iota
	focusTable
)

type Model struct {
	db        db.Database
	input     textinput.Model
	table     table.Model
	focus     focus
	err       error
	width     int
	height     int
	result     *db.Result
	colOffset  int
	showDetail bool
	showHelp   bool
	viewport   viewport.Model
}

func NewModel(database db.Database) Model {
	ti := textinput.New()
	ti.Placeholder = "SELECT * FROM users;"
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 60

	t := table.New(
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Model{
		db:    database,
		input: ti,
		table: t,
		focus: focusInput,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.focus == focusTable {
				m.focus = focusInput
				m.input.Focus()
				return m, nil
			}
			return m, tea.Quit

		case "tab":
			if m.focus == focusInput {
				m.focus = focusTable
				m.input.Blur()
			} else {
				m.focus = focusInput
				m.input.Focus()
			}

		case "left":
			if m.showDetail || m.showHelp {
				return m, nil
			}
			if m.focus == focusTable && m.colOffset > 0 {
				m.colOffset--
				m.updateTable(m.result)
			}

		case "right", "l":
			if m.showDetail || m.showHelp {
				return m, nil
			}
			if m.focus == focusTable && m.result != nil && m.colOffset < len(m.result.Columns)-1 {
				m.colOffset++
				m.updateTable(m.result)
			}

		case "h":
			if m.focus == focusInput {
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
			m.showHelp = !m.showHelp
			if m.showHelp {
				m.showDetail = false
			}

		case "c":
			if m.focus == focusTable && m.result != nil {
				m.copyAsCSV()
			}

		case "v":
			if m.focus == focusTable && m.result != nil && len(m.result.Rows) > 0 {
				m.showDetail = !m.showDetail
				if m.showDetail {
					m.updateViewport()
				}
			}

		case "ctrl+t":
			m.showDetail = false
			res, err := m.db.ListTables()
			if err != nil {
				m.err = err
			} else {
				m.err = nil
				m.colOffset = 0
				m.result = res
				m.updateTable(res)
				m.focus = focusTable
				m.input.Blur()
				m.table.Focus()
			}

		case "ctrl+s":
			res, err := m.db.ListSchemas()
			if err != nil {
				m.err = err
			} else {
				m.err = nil
				m.colOffset = 0
				m.result = res
				m.updateTable(res)
				m.focus = focusTable
				m.input.Blur()
				m.table.Focus()
			}

		case "enter":
			if m.focus == focusInput {
				query := m.input.Value()
				if query != "" {
					res, err := m.db.Query(query)
					if err != nil {
						m.err = err
					} else {
						m.err = nil
						m.colOffset = 0
						m.result = res
						m.updateTable(res)
						m.focus = focusTable
						m.input.Blur()
						m.table.Focus()
					}
				}
			}

		case "esc":
			if m.showDetail {
				m.showDetail = false
				return m, nil
			}
			if m.focus == focusTable {
				m.focus = focusInput
				m.input.Focus()
			} else {
				m.input.SetValue("")
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = m.width - 4
		m.table.SetWidth(m.width - 4)
		// Leave space for Title (3), Input (5), Help (3), Error/Status (2) + Padding
		tableHeight := m.height - 15
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)
		m.viewport.Width = m.width - 4
		m.viewport.Height = tableHeight
	}

	if m.showDetail {
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	if m.focus == focusInput {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m *Model) updateViewport() {
	if m.result == nil {
		return
	}
	
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.result.Rows) {
		return
	}
	
	row := m.result.Rows[cursor]
	var content string
	
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	
	for i, col := range m.result.Columns {
		val := "NULL"
		if i < len(row) {
			val = row[i]
		}
		content += headerStyle.Render(col + ":") + "\n"
		content += valueStyle.Render(val) + "\n\n"
	}
	
	m.viewport.SetContent(content)
}

func (m *Model) updateTable(res *db.Result) {
	if res == nil {
		return
	}

	const minColWidth = 15
	const maxColWidth = 30

	// Determine how many columns can fit
	availableWidth := m.width - 10
	if availableWidth < 20 {
		availableWidth = 20
	}

	// We'll show columns starting from m.colOffset
	displayCols := []table.Column{}
	currentWidth := 0
	var visibleColIndices []int

	for i := m.colOffset; i < len(res.Columns); i++ {
		w := len(res.Columns[i]) + 4
		if w < minColWidth {
			w = minColWidth
		}
		if w > maxColWidth {
			w = maxColWidth
		}

		if currentWidth+w > availableWidth && len(displayCols) > 0 {
			break
		}

		displayCols = append(displayCols, table.Column{Title: res.Columns[i], Width: w})
		visibleColIndices = append(visibleColIndices, i)
		currentWidth += w
	}

	// If we couldn't even fit one column, force it
	if len(displayCols) == 0 && len(res.Columns) > 0 {
		idx := m.colOffset
		if idx >= len(res.Columns) {
			idx = len(res.Columns) - 1
		}
		displayCols = append(displayCols, table.Column{Title: res.Columns[idx], Width: availableWidth})
		visibleColIndices = append(visibleColIndices, idx)
	}

	rows := []table.Row{}
	for _, resRow := range res.Rows {
		row := table.Row{}
		for _, idx := range visibleColIndices {
			if idx < len(resRow) {
				row = append(row, resRow[idx])
			} else {
				row = append(row, "")
			}
		}
		rows = append(rows, row)
	}

	m.table.SetRows([]table.Row{})
	m.table.SetColumns(displayCols)
	m.table.SetRows(rows)
}

func (m *Model) copyAsCSV() {
	if m.result == nil {
		return
	}

	var b strings.Builder
	w := csv.NewWriter(&b)

	// Write header
	if err := w.Write(m.result.Columns); err != nil {
		m.err = err
		return
	}

	// Write rows
	if err := w.WriteAll(m.result.Rows); err != nil {
		m.err = err
		return
	}

	w.Flush()
	if err := clipboard.WriteAll(b.String()); err != nil {
		m.err = err
	}
}

func (m Model) View() string {
	var s string

	// Title
	s += lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render(" Atlas SQL ")
	s += "\n\n"

	if m.showHelp {
		s += focusedStyle.Render(m.helpView())
		s += "\n\n"
		s += helpStyle.Render(" Press 'h' or 'Esc' to return")
		return lipgloss.NewStyle().Padding(1, 2).Render(s)
	}

	// Input Section
	inputStyle := baseStyle
	if m.focus == focusInput {
		inputStyle = focusedStyle
	}
	
	s += inputStyle.Render(m.input.View())
	s += "\n\n"

	// Result Section
	if m.err != nil {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v", m.err))
		s += "\n"
	} else if m.result != nil {
		if m.showDetail {
			s += focusedStyle.Render(m.viewport.View())
			s += "\n"
			s += helpStyle.Render(fmt.Sprintf(" Detail View: %d columns", len(m.result.Columns)))
			s += "\n"
		} else {
			tableStyle := baseStyle
			if m.focus == focusTable {
				tableStyle = focusedStyle
			}
			s += tableStyle.Render(m.table.View())
			s += "\n"
			
			colRange := fmt.Sprintf("Cols: %d-%d of %d", m.colOffset+1, m.colOffset+len(m.table.Columns()), len(m.result.Columns))
			s += helpStyle.Render(fmt.Sprintf(" %d rows returned • %s", len(m.result.Rows), colRange))
			s += "\n"
		}
	}

	// Help Section
	s += "\n"
	s += helpStyle.Render(" Tab: Switch focus • Enter: Run Query • h: Help • c: Copy CSV")
	s += "\n"
	s += helpStyle.Render(" ←/→/l: Scroll Columns • j/k/↑/↓: Scroll Rows • v: Detail • q: Quit")

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}

func (m Model) helpView() string {
	header := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	key := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Width(15)
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	tutorial := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	lines := []string{
		header.Render("Navigation"),
		key.Render("Tab") + desc.Render("Switch focus between Input and Table"),
		key.Render("j/k or ↑/↓") + desc.Render("Scroll through result rows"),
		key.Render("l or →") + desc.Render("Scroll columns right"),
		key.Render("←") + desc.Render("Scroll columns left"),
		"",
		header.Render("Querying"),
		key.Render("Enter") + desc.Render("Execute SQL query in input"),
		key.Render("Ctrl+T") + desc.Render("List all tables"),
		key.Render("Ctrl+S") + desc.Render("List all schemas (databases)"),
		"",
		header.Render("Data Operations"),
		key.Render("v") + desc.Render("Toggle Detail View for selected row"),
		key.Render("c") + desc.Render("Copy entire result set as CSV to clipboard"),
		"",
		header.Render("SQL Quick Tutorial"),
		tutorial.Render("  -- Select all columns from a table"),
		tutorial.Render("  SELECT * FROM table_name;"),
		tutorial.Render("  -- Filter results"),
		tutorial.Render("  SELECT * FROM users WHERE age > 18;"),
		tutorial.Render("  -- Join tables"),
		tutorial.Render("  SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id;"),
		tutorial.Render("  -- SQLite Meta: List Tables / Schema"),
		tutorial.Render("  SELECT name FROM sqlite_master WHERE type='table';"),
		tutorial.Render("  PRAGMA table_info(table_name); -- Describe table"),
		tutorial.Render("  -- Postgres Meta: List Tables"),
		tutorial.Render("  SELECT table_name FROM information_schema.tables;"),
		"",
		header.Render("General"),
		key.Render("h") + desc.Render("Toggle this help screen"),
		key.Render("Esc") + desc.Render("Clear input or return from detail/help"),
		key.Render("q or Ctrl+C") + desc.Render("Quit application"),
	}

	return strings.Join(lines, "\n")
}
