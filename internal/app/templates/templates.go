package templates

import "html/template"

// Templates ..
//type Templates map[string]*template.Template
type Templates struct {
	Root *template.Template
}

// New ...
func New(config *Config) (*Templates, error) {

	var err error

	temp := new(Templates)

	temp.Root, err = template.ParseFiles(config.Root)
	if err != nil {
		return nil, err
	}

	return temp, nil
}
