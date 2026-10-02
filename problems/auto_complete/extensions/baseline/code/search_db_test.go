package code

import "testing"

func TestAddEntry(t *testing.T) {
	sdb := NewSearchDB()

	if err := sdb.AddEntry("apple"); err != nil {
		t.Errorf("%v | unable to add entry in DB: %v", t.Name(), err)
	}

	if err := sdb.AddEntry(""); err == nil {
		t.Errorf("%v | expected error to occur when adding empty entry in DB", t.Name())
	}
}

func TestSearchNSimilarByFrequency(t *testing.T) {
	sdb := NewSearchDB()

	// 1. Seed test data using a clean, readable loop
	entries := []string{"apple", "appleOne", "appleTwo", "apple", "appleOne", "apple"}
	for _, entry := range entries {
		if err := sdb.AddEntry(entry); err != nil {
			t.Fatalf("failed to add entry %q to DB: %v", entry, err)
		}
	}

	// 2. Execute search
	searchEntry := "app"
	threshold := 3
	retEntries, err := sdb.SearchNSimilarByFrequency(searchEntry, threshold)
	if err != nil {
		t.Fatalf("failed to search entry %q: %v", searchEntry, err)
	}

	// 3. Assert results match expected length and content
	// Note: Adjust the expected order slice based on how you simulated frequencies in your setup!
	expectedOrder := []string{"apple", "appleone", "appletwo"}

	if len(retEntries) != len(expectedOrder) {
		t.Errorf("expected %d results, got %d", len(expectedOrder), len(retEntries))
	}

	for i, entry := range retEntries {
		if i < len(expectedOrder) && entry != expectedOrder[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expectedOrder[i], entry)
		}
	}
}
