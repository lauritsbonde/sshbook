package tui

import (
	"strings"

	"sshbook/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// focusable field indexes for the add-connection form.
const (
	fieldName = iota
	fieldHostName
	fieldUser
	fieldPort
	fieldIdentity // selector, not a text input
	fieldCustom   // custom-path input, only focusable when "Custom" is chosen
	fieldCount
)

var textFieldLabels = map[int]string{
	fieldName:     "Name (Host alias)",
	fieldHostName: "HostName",
	fieldUser:     "User",
	fieldPort:     "Port",
}

const (
	identNone   = "None"
	identCustom = "Custom path…"
)

// form is the add-connection input form.
type form struct {
	inputs       map[int]textinput.Model // text fields, keyed by field index
	identOptions []string                // None, <keys...>, Custom
	identIdx     int
	focus        int
	err          string
}

// newForm builds a fresh form; keys are the loaded key basenames offered as
// IdentityFile choices.
func newForm(keys []string) form {
	f := form{
		inputs:       map[int]textinput.Model{},
		identOptions: append(append([]string{identNone}, keys...), identCustom),
	}
	for _, i := range []int{fieldName, fieldHostName, fieldUser, fieldPort, fieldCustom} {
		ti := textinput.New()
		ti.Placeholder = fieldPlaceholder(i)
		f.inputs[i] = ti
	}
	f.setFocus(fieldName)
	return f
}

func fieldPlaceholder(i int) string {
	switch i {
	case fieldName:
		return "myserver"
	case fieldHostName:
		return "1.2.3.4 or example.com"
	case fieldUser:
		return "optional (defaults to local user)"
	case fieldPort:
		return "22 (optional)"
	case fieldCustom:
		return "~/.ssh/some_key"
	}
	return ""
}

// customSelected reports whether the identity selector is on "Custom path…".
func (f form) customSelected() bool {
	return f.identOptions[f.identIdx] == identCustom
}

// identityFile resolves the chosen IdentityFile value.
func (f form) identityFile() string {
	switch opt := f.identOptions[f.identIdx]; opt {
	case identNone:
		return ""
	case identCustom:
		return strings.TrimSpace(f.inputs[fieldCustom].Value())
	default:
		return "~/.ssh/" + opt
	}
}

// connection builds a models.Connection from the current field values.
func (f form) connection() models.Connection {
	return models.Connection{
		Name:         strings.TrimSpace(f.inputs[fieldName].Value()),
		HostName:     strings.TrimSpace(f.inputs[fieldHostName].Value()),
		User:         strings.TrimSpace(f.inputs[fieldUser].Value()),
		Port:         strings.TrimSpace(f.inputs[fieldPort].Value()),
		IdentityFile: f.identityFile(),
	}
}

func (f form) validate() string {
	if strings.TrimSpace(f.inputs[fieldName].Value()) == "" {
		return "Name is required"
	}
	if strings.TrimSpace(f.inputs[fieldHostName].Value()) == "" {
		return "HostName is required"
	}
	if f.customSelected() && strings.TrimSpace(f.inputs[fieldCustom].Value()) == "" {
		return "Custom IdentityFile path is required (or pick None)"
	}
	return ""
}

// formResult signals how a form Update turn ended.
type formResult int

const (
	formPending formResult = iota
	formSubmitted
	formCancelled
)

// update handles a key message; the returned result tells the caller whether
// the form was submitted, cancelled, or is still being edited.
func (f form) update(msg tea.Msg) (form, formResult, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return f.forwardToInput(msg)
	}

	switch key.String() {
	case "esc":
		return f, formCancelled, nil
	case "enter":
		if f.focus == f.lastField() {
			if e := f.validate(); e != "" {
				f.err = e
				return f, formPending, nil
			}
			return f, formSubmitted, nil
		}
		f.focusNext(1)
		return f, formPending, nil
	case "tab", "down":
		f.focusNext(1)
		return f, formPending, nil
	case "shift+tab", "up":
		f.focusNext(-1)
		return f, formPending, nil
	case "left":
		if f.focus == fieldIdentity {
			f.cycleIdent(-1)
			return f, formPending, nil
		}
	case "right":
		if f.focus == fieldIdentity {
			f.cycleIdent(1)
			return f, formPending, nil
		}
	}

	return f.forwardToInput(msg)
}

func (f form) forwardToInput(msg tea.Msg) (form, formResult, tea.Cmd) {
	if ti, ok := f.inputs[f.focus]; ok {
		var cmd tea.Cmd
		ti, cmd = ti.Update(msg)
		f.inputs[f.focus] = ti
		return f, formPending, cmd
	}
	return f, formPending, nil
}

// lastField is the index of the final focusable field (custom path only when
// the Custom option is selected).
func (f form) lastField() int {
	if f.customSelected() {
		return fieldCustom
	}
	return fieldIdentity
}

func (f *form) cycleIdent(dir int) {
	n := len(f.identOptions)
	f.identIdx = (f.identIdx + dir + n) % n
}

func (f *form) setFocus(i int) {
	for idx, ti := range f.inputs {
		if idx == i {
			ti.Focus()
		} else {
			ti.Blur()
		}
		f.inputs[idx] = ti
	}
	f.focus = i
}

func (f *form) focusNext(dir int) {
	last := f.lastField()
	next := f.focus + dir
	if next < 0 {
		next = last
	}
	if next > last {
		next = 0
	}
	f.setFocus(next)
}

func (f form) view(width int) string {
	var b strings.Builder
	b.WriteString(titleStyle.Foreground(lipgloss.Color("42")).Render("Add connection"))
	b.WriteString("\n\n")

	for _, i := range []int{fieldName, fieldHostName, fieldUser, fieldPort} {
		f.writeField(&b, i, textFieldLabels[i], f.inputs[i].View())
	}

	// Identity selector row.
	f.writeField(&b, fieldIdentity, "IdentityFile", f.identView())

	// Custom path input, shown only when Custom is selected.
	if f.customSelected() {
		f.writeField(&b, fieldCustom, "Custom path", f.inputs[fieldCustom].View())
	}

	if f.err != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(f.err) + "\n\n")
	}
	b.WriteString(lipgloss.NewStyle().Faint(true).Render("tab/↑↓ move · ←/→ pick key · enter next/submit · esc cancel"))

	box := activeBorder.Width(width - 2)
	return box.Render(b.String())
}

func (f form) writeField(b *strings.Builder, idx int, label, value string) {
	if idx == f.focus {
		label = selectedRow.Render("➤ " + label)
	} else {
		label = "  " + label
	}
	b.WriteString(label + "\n")
	b.WriteString("  " + value + "\n\n")
}

// identView renders the identity selector as "‹ option ›".
func (f form) identView() string {
	return "‹ " + f.identOptions[f.identIdx] + " ›"
}
