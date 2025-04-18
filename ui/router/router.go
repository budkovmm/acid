package router

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Router struct {
	window       tea.WindowSizeMsg
	models       map[string]tea.Model
	currentModel tea.Model
	currentRoute Message
	startRoute   tea.Cmd
}

func NewRouter(m map[string]tea.Model, startRoute tea.Cmd) *Router {
	return &Router{
		models:       m,
		startRoute:   startRoute,
		currentModel: nil,
	}
}

func (r *Router) Init() tea.Cmd {
	return r.startRoute
}

func (r *Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return r, tea.Quit
		}
	case tea.WindowSizeMsg:
		r.window = msg
	case Message:
		if msg != r.currentRoute {
			r.currentRoute = msg
			r.currentModel = r.models[msg.Model]
			cmds := []tea.Cmd{r.currentModel.Init(), returnMsg(msg)}
			if r.window.Height != 0 {
				cmds = append(cmds, returnMsg(r.window))
			}

			return r, tea.Sequence(cmds...)
		}
	}

	if r.currentModel == nil {
		return r, nil
	}

	m, cmd := r.currentModel.Update(msg)
	r.currentModel = m

	return r, cmd
}

func (r *Router) View() string {
	if r.currentModel == nil {
		return "Initializing application..."
	}

	return r.currentModel.View()
}

type Message struct {
	Model   string
	DataInt int
}

func Route(model string, dataInt ...int) tea.Cmd {
	r := Message{
		Model: model,
	}

	if len(dataInt) > 0 {
		r.DataInt = dataInt[0]
	}

	return returnMsg(r)
}

func returnMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
