package client

import (
	"encoding/json"
	"fmt"
)

type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch v := raw.(type) {
	case string:
		*f = FlexString(v)
	case float64:
		*f = FlexString(fmt.Sprintf("%v", v))
	case int:
		*f = FlexString(fmt.Sprintf("%d", v))
	case nil:
		*f = ""
	default:
		*f = FlexString(fmt.Sprintf("%v", v))
	}
	return nil
}

type FlexInt int64

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch v := raw.(type) {
	case float64:
		*f = FlexInt(int64(v))
	case int:
		*f = FlexInt(v)
	case string:
		var fl float64
		fmt.Sscanf(v, "%f", &fl)
		*f = FlexInt(int64(fl))
	default:
		*f = 0
	}
	return nil
}
