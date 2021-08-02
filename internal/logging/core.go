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

func (f fileWritingCoreShim) Enabled(l zapcore.Level) bool {
	return f.bCore.Enabled(l)
}

func (f fileWritingCoreShim) With(fi []zap.Field) zapcore.Core {
	bc := f.bCore.With(fi)
	wc := f.wCore.With(fi)
	return newFileWritingCoreShim(bc, wc)
}

func (f fileWritingCoreShim) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return f.wCore.Check(e, ce)
}

func (f fileWritingCoreShim) Write(e zapcore.Entry, fi []zap.Field) error {
	return f.wCore.Write(e, fi)
}

func (f fileWritingCoreShim) Sync() error {
	return f.wCore.Sync()
}

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
				zap.DebugLevel,
			)

			return newFileWritingCoreShim(c, wc)
		},
	)
}
