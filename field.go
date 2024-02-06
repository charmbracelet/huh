package huh

import "github.com/charmbracelet/lipgloss"

type options struct {
	renderer *lipgloss.Renderer
}

func getOptions(opts []ThemeOption) *options {
	options := &options{
		renderer: lipgloss.DefaultRenderer(),
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type ThemeOption func(*options)

func WithRenderer(r *lipgloss.Renderer) ThemeOption {
	return func(o *options) {
		o.renderer = r
	}
}
