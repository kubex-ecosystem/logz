package logz

// import "os"

// // Logger global, estilo log padrão, mas Kubex-grade.
// var (
// 	globalDynamicWriter = NewDynamicWriter(NewIOWriter(os.Stdout))
// 	globalLogger        = NewLogger(
// 		WithWriter(globalDynamicWriter),
// 		WithFormatter(NewPrettyFormatter()),
// 		WithMinLevel(LevelDebug),
// 	)
// )

// // SetGlobalWriter permite trocar o destino em runtime.
// func SetGlobalWriter(w Writer) {
// 	globalDynamicWriter.Set(w)
// }

// // SetGlobalFormatter troca o formatter global.
// func SetGlobalFormatter(f Formatter) {
// 	globalLogger.mu.Lock()
// 	defer globalLogger.mu.Unlock()
// 	globalLogger.formatter = f
// }

// // SetGlobalMinLevel define o min level global.
// func SetGlobalMinLevel(lv Level) {
// 	globalLogger.mu.Lock()
// 	defer globalLogger.mu.Unlock()
// 	globalLogger.minLevel = lv
// }

// // Facade global conveniência.

// func Debug(msg string, fields ...Field) error { return globalLogger.Debug(msg, fields...) }
// func Info(msg string, fields ...Field) error  { return globalLogger.Info(msg, fields...) }
// func Warn(msg string, fields ...Field) error  { return globalLogger.Warn(msg, fields...) }
// func Error(msg string, fields ...Field) error { return globalLogger.Error(msg, fields...) }
// func Fatal(msg string, fields ...Field) error { return globalLogger.Fatal(msg, fields...) }

// // Entry chainable usando o logger global.
// func Entry() *EntryBuilder { return globalLogger.Entry() }
