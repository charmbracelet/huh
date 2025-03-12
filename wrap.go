package huh

import "github.com/charmbracelet/x/ansi"

// TODO: replace with cellbuf.Wrap?
func wrap(s string, limit int) string {
	return ansi.Wrap(s, limit, ",.-; ")
}
