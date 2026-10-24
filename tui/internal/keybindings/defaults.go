package keybindings

func Common() Bindings {
	return Bindings{
		"quit":          {"ctrl+c"},
		"back":          {"esc"},
		"switch-view":   {"ctrl+b"},
		"doctor":        {"ctrl+/", "ctrl+_", "f1"},
		"help":          {"?"},
		"open":          {"enter"},
		"previous":      {"up", "ctrl+k"},
		"next":          {"down", "ctrl+j"},
		"mark-next":     {"tab"},
		"mark-previous": {"shift+tab"},
		"search-all":    {"ctrl+v"},
		"preview-up":    {"pgup"},
		"preview-down":  {"pgdown"},
	}
}

func Tmux() Bindings {
	return Bindings{
		"focus-windows":   {"ctrl+l"},
		"focus-sessions":  {"ctrl+h"},
		"previous-window": {"k"},
		"next-window":     {"j"},
		"create":          {"ctrl+n"},
		"rename":          {"ctrl+r"},
		"move-window":     {"ctrl+o"},
		"delete":          {"ctrl+x"},
	}
}

func Codex() Bindings {
	return Bindings{
		"archive":         {"ctrl+a"},
		"unarchive":       {"ctrl+u"},
		"delete":          {"ctrl+x"},
		"toggle-archived": {"ctrl+g"},
		"toggle-ignored":  {"ctrl+t"},
	}
}

func Doctor() Bindings {
	return Bindings{"refresh": {"ctrl+r"}}
}
