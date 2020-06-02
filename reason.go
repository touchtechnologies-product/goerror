package goerror

type Reason struct {
	FieldName string      `json:"fieldName"`
	Reason    string      `json:"reason"`
	Value     interface{} `json:"value,omitempty"`
}

func (e *GoError) GetReasons() []*Reason {
	if len(e.reasons) == 0 {
		return []*Reason{}
	}

	return e.reasons
}

func (e *GoError) AddReason(fieldName, reason string, value interface{}) {
	if len(e.reasons) == 0 {
		e.reasons = []*Reason{}
	}

	e.reasons = append(e.reasons, &Reason{
		FieldName: fieldName,
		Reason:    reason,
		Value:     value,
	})
}
