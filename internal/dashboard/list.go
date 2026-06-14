package dashboard

import (
	"encoding/json"
	"fmt"
)

type ObjectMeta struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// Implements [json.Unmarshaler], to be passed to e.g. [json.Unmarshal],
// [Resource] should be set in initialization as the resource name in lowercase, plural form.
type ListUnmarshaler struct {
	ItemsKey    string
	ObjectMetas []ObjectMeta
}

var _ json.Unmarshaler = &ListUnmarshaler{}

func (l *ListUnmarshaler) UnmarshalJSON(data []byte) error {
	firstLayer := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &firstLayer); err != nil {
		return fmt.Errorf("cannot unmarshal first layer: %w", err)
	}

	list, ok := firstLayer[l.ItemsKey]
	if !ok {
		return fmt.Errorf("expected key %s at first layer not found", l.ItemsKey)
	}

	unmarshaledList := []struct {
		ObjectMeta `json:"objectMeta"`
	}{}
	if err := json.Unmarshal(list, &unmarshaledList); err != nil {
		return fmt.Errorf("cannot unmarshal list content: %w", err)
	}

	l.ObjectMetas = make([]ObjectMeta, 0, len(unmarshaledList))
	for _, o := range unmarshaledList {
		l.ObjectMetas = append(l.ObjectMetas, o.ObjectMeta)
	}

	return nil
}
