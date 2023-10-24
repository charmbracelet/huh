package huh

// Run runs a single field.
func Run(f Field) error {
	return NewForm(NewGroup(f)).ShowHelp(false).Run()
}
