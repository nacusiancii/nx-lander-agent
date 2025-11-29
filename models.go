package main

type Model map[string]string

func (m Model) Name() string {
	return m["name"]
}

var (
	MINIMAX_M2 = Model{
		"name":    "minimax/minimax-m2",
		"Minimax": "minimax/fp8",
		"Google":  "google-vertex",
	}
	KIMI_K2_THINKING = Model{
		"name":   "moonshotai/kimi-k2-thinking",
		"Google": "google-vertex",
	}
)
