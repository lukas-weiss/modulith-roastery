package eventlog

import "testing"

func TestModuleOf(t *testing.T) {
	tests := map[string]string{
		"lweiss.dev/northbean/inventory/internal.OrderPlacedListener.On": "inventory",
		"lweiss.dev/northbean/notification.Listener.On":                  "notification",
		"lweiss.dev/northbean/orphan":                                    "orphan",
		"other/pkg.Listener.On":                                          "other",
		"":                                                               "",
	}
	for listenerID, want := range tests {
		if got := moduleOf(listenerID); got != want {
			t.Errorf("moduleOf(%q) = %q, want %q", listenerID, got, want)
		}
	}
}
