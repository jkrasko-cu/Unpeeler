package ui

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/spinner"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/decompressor"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/detector"
)

var (
    titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Padding(0, 1)
    doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
    layerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
    dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
    successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
    errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
    boxStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))
)

type Layer struct {
    Format detector.Format
    Output string
    Done   bool
}

type Model struct {
    inputFile string
    finalFile string
    outputDir string
    layers    []Layer
    spinner   spinner.Model
    current   string
    done      bool
    err       error
    startTime time.Time
}

type processMsg struct {
    layer Layer
    next  string
    err   error
    done  bool
}

func NewModel(input, output string) Model {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    return Model{
        inputFile: input,
        outputDir: output,
        spinner:   s,
        current:   input,
        startTime: time.Now(),
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(m.spinner.Tick, processNext(m.current, m.outputDir, 0))
}

func processNext(path, outputDir string, depth int) tea.Cmd {
    return func() tea.Msg {
        if detector.IsHumanReadable(path) {
            return processMsg{done: true, next: path}
        }
        format := detector.Detect(path)
        if format == detector.Unknown {
            return processMsg{done: true, next: path}
        }
        layerDir := fmt.Sprintf("%s/layer_%d", outputDir, depth)
        os.MkdirAll(layerDir, 0755)
        out, err := decompressor.Decompress(path, layerDir, format)
        return processMsg{
            layer: Layer{Format: format, Output: out, Done: true},
            next:  out,
            err:   err,
            done:  false,
        }
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" || msg.String() == "ctrl+c" {
            return m, tea.Quit
        }
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    case processMsg:
        if msg.err != nil {
            m.err = msg.err
            m.done = true
            return m, tea.Quit
        }
        if msg.done {
            m.done = true
            m.finalFile = msg.next
            return m, tea.Quit
        }
        m.layers = append(m.layers, msg.layer)
        m.current = msg.next
        return m, processNext(m.current, m.outputDir, len(m.layers))
    }
    return m, nil
}

func (m Model) View() string {
    header := titleStyle.Render("🧅 unpeeler") + dimStyle.Render("  recursive decompressor\n")

    body := ""
    for i, layer := range m.layers {
        mark := doneStyle.Render("✓")
        body += fmt.Sprintf("  %s  Layer %-2d  %s%-8s%s\n",
            mark,
            i+1,
            layerStyle.Render("["),
            string(layer.Format),
            layerStyle.Render("]"),
        )
    }

    if !m.done {
        body += fmt.Sprintf("  %s  Layer %-2d  detecting...\n", m.spinner.View(), len(m.layers)+1)
    }

    footer := ""
    elapsed := time.Since(m.startTime).Round(time.Millisecond)
    if m.done && m.err == nil {
        content, _ := os.ReadFile(m.finalFile)
        footer += "\n" + successStyle.Render(fmt.Sprintf("  ✅ Done! %d layers in %s", len(m.layers), elapsed))
        footer += "\n" + dimStyle.Render(fmt.Sprintf("  💾 Final file: %s", m.finalFile))
        footer += "\n\n" + dimStyle.Render("  "+strings.TrimSpace(string(content)))
    } else if m.err != nil {
        footer = "\n" + errorStyle.Render("  ✗ Error: "+m.err.Error())
    }

    return boxStyle.Render(header+"\n"+body+footer) + "\n"
}
