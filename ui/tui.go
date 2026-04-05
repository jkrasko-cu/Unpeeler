package ui

import (
    "fmt"
    "math"
    "os"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/filepicker"
    "github.com/charmbracelet/bubbles/spinner"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/detector"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/entropy"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/hexdump"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/inspector"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/stringSearch"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/compressor"
)

type tuiState int

const (
    statePicking tuiState = iota
    stateMenu
    stateRunning
    stateResults
)

var commands = []string{
    "inspect",
    "peel",
    "wrap",
    "entropy",
    "histogram",
    "strings",
    "hexdump",
}

var tuiTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("34"))
var tuiDimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))
var tuiSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true)
var tuiGreenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
var tuiBoxStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("34"))
var tuiOffsetStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
var tuiHexStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))
var tuiAsciiStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("64"))

type TUIModel struct {
    state      tuiState
    filepicker filepicker.Model
    spinner    spinner.Model
    selected   string
    cursor     int
    result     string
    err        error
}

type resultMsg string
type errMsg error

func NewTUIModel() TUIModel {
    fp := filepicker.New()
    home, _ := os.UserHomeDir()
    fp.CurrentDirectory = home + "/Documents/unpeeler-demo"
    fp.ShowHidden = false
    fp.Styles.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true)
    fp.Styles.Directory = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))
    fp.Styles.File = lipgloss.NewStyle().Foreground(lipgloss.Color("64"))
    fp.Styles.Symlink = lipgloss.NewStyle().Foreground(lipgloss.Color("64"))
    fp.Styles.DisabledFile = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))
    fp.Styles.EmptyDirectory = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))

    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))

    return TUIModel{
        state:      statePicking,
        filepicker: fp,
        spinner:    s,
    }
}

func (m TUIModel) Init() tea.Cmd {
    return m.filepicker.Init()
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch m.state {
        case statePicking:
            if msg.String() == "ctrl+c" || msg.String() == "q" {
                return m, tea.Quit
            }
        case stateMenu:
            switch msg.String() {
            case "ctrl+c", "q":
                return m, tea.Quit
            case "up", "k":
                if m.cursor > 0 {
                    m.cursor--
                }
            case "down", "j":
                if m.cursor < len(commands)-1 {
                    m.cursor++
                }
            case "enter":
                m.state = stateRunning
                return m, tea.Batch(m.spinner.Tick, runCommand(commands[m.cursor], m.selected))
            case "b":
                m.state = statePicking
                m.selected = ""
            }
        case stateResults:
            switch msg.String() {
            case "ctrl+c", "q":
                return m, tea.Quit
            case "b":
                m.state = stateMenu
                m.result = ""
            }
        }

    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd

    case resultMsg:
        m.result = string(msg)
        m.state = stateResults
        return m, nil

    case errMsg:
        m.err = msg
        m.state = stateResults
        return m, nil
    }

    if m.state == statePicking {
        var cmd tea.Cmd
        m.filepicker, cmd = m.filepicker.Update(msg)
        if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
            m.selected = path
            m.state = stateMenu
            m.cursor = 0
        }
        return m, cmd
    }

    return m, nil
}

func (m TUIModel) View() string {
    header := tuiTitleStyle.Render("🧅 unpeeler") + tuiDimStyle.Render("  forensics toolkit\n")

    switch m.state {
    case statePicking:
        return tuiBoxStyle.Render(header + "\n  Select a file:\n\n" + m.filepicker.View())

    case stateMenu:
        body := fmt.Sprintf("  File: %s\n\n", tuiGreenStyle.Render(m.selected))
        for i, cmd := range commands {
            if i == m.cursor {
                body += fmt.Sprintf("  %s %s\n", tuiSelectedStyle.Render("▶"), tuiSelectedStyle.Render(cmd))
            } else {
                body += fmt.Sprintf("    %s\n", tuiDimStyle.Render(cmd))
            }
        }
        body += "\n" + tuiDimStyle.Render("  ↑/↓ navigate  enter select  b back  q quit")
        return tuiBoxStyle.Render(header + "\n" + body)

    case stateRunning:
        return tuiBoxStyle.Render(header + "\n  " + m.spinner.View() + tuiDimStyle.Render("  running..."))

    case stateResults:
        if m.err != nil {
            return tuiBoxStyle.Render(header + "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: "+m.err.Error()) + "\n\n  " + tuiDimStyle.Render("b back  q quit"))
        }
        return tuiBoxStyle.Render(header + "\n" + m.result + "\n\n  " + tuiDimStyle.Render("b back  q quit"))
    }

    return ""
}

func runCommand(cmd, path string) tea.Cmd {
    return func() tea.Msg {
        switch cmd {
        case "inspect":
            r, err := inspector.Inspect(path)
            if err != nil {
                return errMsg(err)
            }
            scored := tuiScoreStyle(r.Entropy.Score)
            bar := tuiRenderBar(r.Entropy.Score)
            out := fmt.Sprintf("  %-10s %s\n", "File:", path)
            out += fmt.Sprintf("  %-10s %s\n", "Format:", tuiHexStyle.Render(string(r.Format)))
            out += fmt.Sprintf("  %-10s %s\n", "Size:", formatBytes(int(r.Size)))
            out += fmt.Sprintf("  %-10s %s\n", "Magic:", tuiDimStyle.Render(r.Magic))
            out += fmt.Sprintf("  %-10s %s  %s / 8.0\n", "Entropy:", bar, scored.Render(fmt.Sprintf("%.2f", r.Entropy.Score)))
            out += fmt.Sprintf("  %-10s %s\n", "Verdict:", scored.Render(r.Entropy.Label))
            return resultMsg(out)

        case "entropy":
            r, err := entropy.Analyze(path)
            if err != nil {
                return errMsg(err)
            }
            scored := tuiScoreStyle(r.Score)
            bar := tuiRenderBar(r.Score)
            out := fmt.Sprintf("  %-10s %s\n", "File:", path)
            out += fmt.Sprintf("  %-10s %s\n", "Size:", formatBytes(r.ByteCount))
            out += fmt.Sprintf("  %-10s %s  %s / 8.0\n", "Entropy:", bar, scored.Render(fmt.Sprintf("%.2f", r.Score)))
            out += fmt.Sprintf("  %-10s %s\n", "Verdict:", scored.Render(r.Label))
            return resultMsg(out)

        case "histogram":
            r, err := entropy.Histogram(path, 256)
            if err != nil {
                return errMsg(err)
            }
            out := fmt.Sprintf("  %-10s %-22s %s\n\n", "Offset", "Entropy", "Verdict")
            for _, chunk := range r.Chunks {
                bar := tuiRenderBar(chunk.Score)
                out += fmt.Sprintf("  %08x  %s %s\n", chunk.Offset, bar, tuiDimStyle.Render(fmt.Sprintf("%.2f  %s", chunk.Score, chunk.Label)))
            }
            return resultMsg(out)

        case "strings":
            r, err := stringSearch.Analyze(path)
            if err != nil {
                return errMsg(err)
            }
            out := fmt.Sprintf("  Found: %s\n\n", tuiGreenStyle.Render(fmt.Sprintf("%d strings", r.Count)))
            for _, s := range r.Strings {
                out += fmt.Sprintf("  %s\n", tuiDimStyle.Render(s))
            }
            return resultMsg(out)

        case "hexdump":
            r, err := hexdump.Analyze(path)
            if err != nil {
                return errMsg(err)
            }
            out := ""
            for _, line := range r.Lines {
                out += fmt.Sprintf("  %s  %s  %s\n",
                    tuiOffsetStyle.Render(line.Offset),
                    tuiHexStyle.Render(line.Hex),
                    tuiAsciiStyle.Render("|"+line.ASCII+"|"),
                )
            }
            return resultMsg(out)

        case "peel":
            m := NewModel(path, "./output/"+fmt.Sprintf("%d", time.Now().Unix()))
            p := tea.NewProgram(m)
            if _, err := p.Run(); err != nil {
                return errMsg(err)
            }
            return resultMsg(tuiGreenStyle.Render("Peel complete — check ./output"))

        case "wrap":
            out, err := compressor.Compress(path, "./wrapped", detector.Gzip, 1)
            if err != nil {
                return errMsg(err)
            }
            return resultMsg(tuiGreenStyle.Render("Wrapped → " + out))
        }

        return resultMsg("Unknown command") 
    }
}

func tuiRenderBar(score float64) string {
    const width = 20
    filled := int(math.Round(score / 8.0 * width))
    if filled > width {
        filled = width
    }
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    return tuiScoreStyle(score).Render(bar)
}

func tuiScoreStyle(score float64) lipgloss.Style {
    switch {
    case score < 4.0:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
    case score < 6.5:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("64"))
    default:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
    }
}

func tuiFormatStr(f detector.Format) string {
    if f == detector.Unknown {
        return tuiDimStyle.Render("unknown")
    }
    return tuiHexStyle.Render(string(f))
}

var _ = tuiFormatStr // suppress unused warning

func formatBytes(n int) string {
    switch {
    case n < 1024:
        return fmt.Sprintf("%d B", n)
    case n < 1024*1024:
        return fmt.Sprintf("%.1f KB", float64(n)/1024)
    default:
        return fmt.Sprintf("%.1f MB", float64(n)/(1024*1024))
    }
}
