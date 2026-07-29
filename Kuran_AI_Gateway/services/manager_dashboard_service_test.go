package services

import "testing"

func TestNormalizeDeviceType(t *testing.T) {
	tests := map[string]string{
		"iOS":            "ios",
		"iPhone 15":      "ios",
		"android":        "android",
		"Android Tablet": "android",
		"web":            "web",
		"browser":        "web",
		"":               "unknown",
		"desktop-client": "unknown",
	}
	for input, expected := range tests {
		if actual := normalizeDeviceType(input); actual != expected {
			t.Fatalf("normalizeDeviceType(%q) = %q, want %q", input, actual, expected)
		}
	}
}

func TestNormalizePage(t *testing.T) {
	page, pageSize := normalizePage(0, 1000)
	if page != 1 || pageSize != 100 {
		t.Fatalf("normalizePage(0, 1000) = (%d, %d), want (1, 100)", page, pageSize)
	}
}

func TestCompleteBreakdownAddsMissingKeys(t *testing.T) {
	rows := completeBreakdown([]MetricBreakdown{{Key: "iOS", Count: 4}}, []string{"ios", "android", "web", "unknown"})
	if len(rows) != 4 || rows[0].Key != "ios" || rows[0].Count != 4 || rows[1].Count != 0 {
		t.Fatalf("unexpected completed breakdown: %#v", rows)
	}
}
