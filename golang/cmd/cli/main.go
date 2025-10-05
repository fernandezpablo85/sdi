package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fernandezpablo85/sdi/internal/client"
	"github.com/fernandezpablo85/sdi/internal/env"
)

const maxHistory = 45

var (
	successStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4a4"))
	rateLimitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a44"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#555"))
)

type model struct {
	client          *client.HttpClient
	requestsSent    int
	requestsSuccess int
	requests429     int
	requestErr      int
	history         []int
	isRunning       bool
}

type tickMsg struct {
	time time.Time
}

func (m model) Init() tea.Cmd {
	return nil
}

func tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{time: t}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case tea.KeySpace.String():
			m.isRunning = !m.isRunning
			if m.isRunning {
				return m, tick()
			}
		}
	case tickMsg:
		if m.isRunning {
			m.requestsSent++
			assetName := "BTCUSDT"
			res, err := m.client.GetAssetPrice(assetName)
			var statusCode int
			if err != nil {
				m.requestErr++
				statusCode = 0
			} else if res.StatusCode == http.StatusOK {
				statusCode = 200
				m.requestsSuccess++
			} else {
				// we assume for now that <> 200 is 429
				statusCode = 429
				m.requests429++
			}
			m.history = append(m.history, statusCode)
			if len(m.history) > maxHistory {
				m.history = m.history[1:]
			}
			return m, tick()
		}
	}

	return m, nil
}

func renderHistory(history []int) string {

	var result string
	for _, code := range history {
		switch code {
		case 200:
			result += successStyle.Render("█")
		case 429:
			result += rateLimitStyle.Render("█")
		case 0:
			result += errorStyle.Render("█")
		}
	}
	return result
}

func (m model) View() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")). // Cyan
		Padding(0, 0)

	subtitleStyle := lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#00bfff90")). // Cyan
		Padding(0, 0).Underline(true)

	statusStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")) // Green when running

	if !m.isRunning {
		statusStyle = statusStyle.Foreground(lipgloss.Color("9")) // Red when stopped
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Purple border
		Padding(1, 2).
		Width(50)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")) // Gray

	// Build the content
	status := "STOPPED"
	if m.isRunning {
		status = "RUNNING"
	}

	content := fmt.Sprintf(
		"%s\n\n"+
			"%s\n\n"+
			"Status: %s\n\n"+
			"Requests Sent:    %d\n"+
			"Success (200):    %s\n"+
			"Rate Limited:     %s\n\n"+
			"Network Errors:   %s\n\n"+
			"History (last %d):\n%s\n\n"+
			"%s",
		subtitleStyle.Render("System Design Interview"),
		titleStyle.Render("Rate Limiter Visualizer"),
		statusStyle.Render(status),
		m.requestsSent,
		successStyle.Render(strconv.Itoa(m.requestsSuccess)),
		rateLimitStyle.Render(strconv.Itoa(m.requests429)),
		errorStyle.Render(strconv.Itoa(m.requestErr)),
		len(m.history),
		renderHistory(m.history),
		helpStyle.Render("[space] start/stop  [q] quit"),
	)

	return boxStyle.Render(content)
}

func main() {
	port := env.GetIntOrElse("PORT", 8080)
	baseUrl := fmt.Sprintf("http://localhost:%d", port)
	cli := client.NewClient(baseUrl)

	pong := cli.Ping()
	if !pong {
		log.Fatalf("could not ping api at %s", baseUrl)
	}
	p := tea.NewProgram(model{isRunning: false, client: cli})
	p.Run()
}
