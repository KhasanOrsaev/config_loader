package tests

import (
	"encoding/json"
	"testing"
)

// дописать
func TestEventSet(t *testing.T) {
	data := map [string]map[string]string {
		"dwh":{"user_id": "ba10008c-9109-4321-9a38-d228e7fe84e4","type": "event_type_1","source": "event_source_1"},
	}
	j, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
}