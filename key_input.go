package huh

import tea "charm.land/bubbletea/v2"

// sanitizeKeyPressForInput clears stray Text from navigation keys before they
// reach text inputs. On Windows, Tab and Enter key presses can include \t or \r
// in Text while also using navigation key codes, which bubbles would otherwise
// insert as literal characters.
func sanitizeKeyPressForInput(msg tea.Msg) tea.Msg {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return msg
	}

	key := keyMsg.Key()
	switch key.Code {
	case tea.KeyTab, tea.KeyEnter:
		if key.Text == "" {
			return msg
		}
		key.Text = ""
		return tea.KeyPressMsg(key)
	default:
		return msg
	}
}
