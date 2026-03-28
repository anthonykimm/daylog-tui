package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderColor   = lipgloss.Color("12")
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	doneStyle     = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("8"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type mode int

const (
	modeNormal mode = iota
	modeInput
)

type model struct {
	db       *sql.DB
	tasks    []Task
	cursor   int
	mode     mode
	input    textinput.Model
	width    int
	height   int
	err      error
	quitting bool
}

func newModel(db *sql.DB, tasks []Task) model {
	ti := textinput.New()
	ti.Placeholder = "New task..."
	ti.CharLimit = 256

	return model{
		db:    db,
		tasks: tasks,
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeInput {
			return m.handleInputMode(msg)
		}
		return m.handleNormalMode(msg)
	}

	return m, nil
}

func (m model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.tasks)-1 {
			m.cursor++
		}

	case "a":
		m.mode = modeInput
		m.input.Reset()
		m.input.Focus()
		return m, m.input.Cursor.BlinkCmd()

	case " ", "enter":
		if len(m.tasks) > 0 {
			task := m.tasks[m.cursor]
			if err := toggleTask(m.db, task.ID); err != nil {
				m.err = err
				return m, nil
			}
			m.tasks[m.cursor].Completed = !m.tasks[m.cursor].Completed
		}

	case "d":
		if len(m.tasks) > 0 {
			task := m.tasks[m.cursor]
			if err := deleteTask(m.db, task.ID); err != nil {
				m.err = err
				return m, nil
			}
			m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
			if m.cursor >= len(m.tasks) && m.cursor > 0 {
				m.cursor--
			}
		}
	}

	return m, nil
}

func (m model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		title := m.input.Value()
		if title != "" {
			task, err := addTask(m.db, title)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.tasks = append(m.tasks, task)
			m.cursor = len(m.tasks) - 1
		}
		m.mode = modeNormal
		m.input.Blur()
		return m, nil

	case "esc":
		m.mode = modeNormal
		m.input.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	today := time.Now().Format("2006-01-02")
	s := titleStyle.Render(fmt.Sprintf("daylog — %s", today)) + "\n\n"

	if len(m.tasks) == 0 && m.mode != modeInput {
		s += "No tasks for today. Press 'a' to add one.\n"
	}

	for i, task := range m.tasks {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		check := "[ ]"
		if task.Completed {
			check = "[x]"
		}

		line := fmt.Sprintf("%s%s %s", cursor, check, task.Title)

		if task.Completed {
			line = doneStyle.Render(line)
		} else if i == m.cursor {
			line = selectedStyle.Render(line)
		}

		s += line + "\n"
	}

	if m.mode == modeInput {
		s += "\n" + m.input.View() + "\n"
	}

	if m.err != nil {
		s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v", m.err))
	}

	help := helpStyle.Render("a: add • space/enter: toggle • d: delete • j/k: navigate • q: quit")

	// Border takes 2 cols and 2 rows, padding takes 2 cols and 2 rows
	innerWidth := m.width - 4
	innerHeight := m.height - 4
	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	// Count content lines and pad to push help to the bottom
	contentLines := strings.Count(s, "\n")
	padding := innerHeight - contentLines - 1 // -1 for the help line
	if padding > 0 {
		s += strings.Repeat("\n", padding)
	}
	s += help

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 1).
		Width(innerWidth).
		Height(innerHeight)

	return box.Render(s)
}
