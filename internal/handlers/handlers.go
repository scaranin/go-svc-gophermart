package handlers

import "runtime"

type URLHandler struct {
	configList string
}

// Функция для получения имени текущего метода
func (h URLHandler) GetCurrentMethodName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()

	return fullName
}
