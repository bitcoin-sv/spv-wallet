package lox

// Catch is a utility function that catches an error and collects it in an ErrorCollector.
func Catch[T any](collector *ErrorCollector, callback func() (T, error)) T {
	val, err := callback()
	if err != nil {
		collector.Collect(err)
	}
	return val
}

// CatchFn is a utility function that catches an error and collects it in an ErrorCollector.
// It returns a function that can be called to execute the callback.
func CatchFn[T any](collector *ErrorCollector, callback func() (T, error)) func() T {
	return func() T {
		return Catch(collector, callback)
	}
}
