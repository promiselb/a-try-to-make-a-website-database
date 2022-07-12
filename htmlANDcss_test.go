package website

import (
	"testing"
)

func TestHTml(t *testing.T) {
	got := Html(Head("", "", "") + Body("", "", ""))
	want := `<html><head id="" class=""></head><body id="" class=""></body></html>`

	if got != want {
		t.Errorf(
			`want = %s
		got = %s`, want, got)
	}
}
