package gqlgen

// GetIntParam helper
func GetIntParam(param *int) int {
	var _param int
	if param == nil {
		_param = 0
	} else {
		_param = *param
	}
	return _param
}

// GetUint64Param helper
func GetUint64Param(param *uint64) uint64 {
	var _param uint64
	if param == nil {
		_param = 0
	} else {
		_param = *param
	}
	return _param
}

// GetBoolParam helper
func GetBoolParam(param *bool) bool {
	var _param bool
	if param == nil {
		_param = false
	} else {
		_param = *param
	}
	return _param
}

// GetStringParam helper
func GetStringParam(param *string) string {
	var _param string
	if param == nil {
		_param = ""
	} else {
		_param = *param
	}
	return _param
}
