package cmd

import "testing"

func TestMaskApiToken(t *testing.T) {
	mapping := map[string]string{
		"sk-x3FTKEJ4DJ^MiVuoQDlGT3BlnkFJwHJD4ex6esKpdxI5qMeA": "sk-x3*****MeA",
	}
	var actual string
	for k, m := range mapping {
		actual = maskApiToken(k)
		if actual != m {
			t.Errorf("API Key %s, expected %s, got %s", k, m, actual)
		}
	}
}
