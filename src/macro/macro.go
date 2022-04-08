package macro

import "linflow/src/actions"

type Config struct {
	Name    string   `yaml:"name"`
	Code    string   `yaml:"code"`
	Actions []string `yaml:"actions"`
}

func (m *Config) ToMacro() (*Macro, error) {
	var err error
	acts := make([]actions.Action, len(m.Actions))
	for i, action := range m.Actions {
		acts[i], err = actions.Parse(action)
		if err != nil {
			return nil, err
		}
	}

	return NewMacro(m.Name, m.Code, acts), nil
}

type Macro struct {
	Name    string           `yaml:"name"`
	Code    string           `yaml:"code"`
	Actions []actions.Action `yaml:"actions"`
}

func NewMacro(name, code string, actions []actions.Action) *Macro {
	return &Macro{
		Name:    name,
		Code:    code,
		Actions: actions,
	}
}

func (m *Macro) Execute() error {
	for _, action := range m.Actions {
		err := action.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}
