package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// fileWritingCoreShim implements the zapcore.Core interface; It uses bCore to get whether logging is enabled for a particular level,
//   and wCore is used for all writing functions.
type fileWritingCoreShim struct {
	bCore zapcore.Core
	wCore zapcore.Core
}

// Enabled tells if the core is enabled (able to log) an entry with the given level.
// Part of the zapcore.Core interface
func (f fileWritingCoreShim) Enabled(l zapcore.Level) bool {
	return f.bCore.Enabled(l)
}

// With adds structured fields to the core. The returned core is also a newFileWritingCoreShim,
// still using wCore as the writing core and bCore as the level enabled core.
// Part of the zapcore.Core interface
func (f fileWritingCoreShim) With(fi []zap.Field) zapcore.Core {
	bc := f.bCore.With(fi)
	wc := f.wCore.With(fi)
	return newFileWritingCoreShim(bc, wc)
}

// Check determines if the supplied Entry should be logged.
// Part of the zapcore.Core interface
func (f fileWritingCoreShim) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if f.bCore.Enabled(e.Level) {
		return ce.AddCore(e, f.wCore)
	}
	return ce
}

// Write writes the entry to wCore with the given fields
// Part of the zapcore.Core interface
func (f fileWritingCoreShim) Write(e zapcore.Entry, fi []zap.Field) error {
	return f.wCore.Write(e, fi)
}

// Sync flushes buffered logs on wCore
// Part of the zapcore.Core interface
func (f fileWritingCoreShim) Sync() error {
	return f.wCore.Sync()
}

// newFileWritingCoreShim returns a new fileWritingCoreShim with the given bCore and wCore
func newFileWritingCoreShim(bCore zapcore.Core, wCore zapcore.Core) *fileWritingCoreShim {
	return &fileWritingCoreShim{
		bCore: bCore,
		wCore: wCore,
	}
}

// FileLoggingCoreOption returns a zap option that will log to rotating log files,
// 	using the specified filepath and rotating by adding a timestamp to the filename.
func FileLoggingCoreOption(filePath string) zap.Option {
	return zap.WrapCore(
		func(c zapcore.Core) zapcore.Core {
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   filePath,
				MaxBackups: 3,
				MaxSize:    1,
				MaxAge:     7,
			})

			zapCfg := zap.NewProductionEncoderConfig()
			zapCfg.EncodeTime = zapcore.ISO8601TimeEncoder

			wc := zapcore.NewCore(
				zapcore.NewConsoleEncoder(zapCfg),
				w,
				zap.FatalLevel,
			)

			return newFileWritingCoreShim(c, wc)
		},
	)
}
