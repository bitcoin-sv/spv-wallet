package filter

import "time"

// TimeRange defines a range between two points in time.
type TimeRange struct {
	// From specifies the start time of the range. It's optional and can be nil.
	From *time.Time `json:"from,omitempty" example:"2024-02-26T11:01:28Z"`
	// To specifies the end time of the range. It's optional and can be nil.
	To *time.Time `json:"to,omitempty" example:"2024-02-26T11:01:28Z"`
}

// ToDbConditions converts the TimeRange to a set of database query conditions.
// Returns nil if both From and To are nil or zero, indicating no conditions.
func (tr *TimeRange) ToDbConditions() map[string]interface{} {
	if tr == nil || tr.isEmpty() {
		return nil
	}

	conditions := make(map[string]interface{})

	if tr.hasFrom() {
		conditions["$gte"] = *tr.From
	}
	if tr.hasTo() {
		conditions["$lte"] = *tr.To
	}

	return conditions
}

// isEmpty checks if both From and To are not set (nil or zero).
func (tr *TimeRange) isEmpty() bool {
	return !tr.hasFrom() && !tr.hasTo()
}

func (tr *TimeRange) hasFrom() bool {
	return !(tr.From == nil || tr.From.IsZero())
}

func (tr *TimeRange) hasTo() bool {
	return !(tr.To == nil || tr.To.IsZero())
}
