package difftest

import "testing"

func TestIt(t *testing.T) {
	got := map[string]interface{}{
		"message": "hello world",
		"id":      1,
	}
	want := `
{ "id": 1, "message": "hello world"}
`
	AssertGotAndWantString(t, got, want)
}
