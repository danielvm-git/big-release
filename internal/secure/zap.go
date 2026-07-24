// story: e08s01
package secure

import (
	"go.uber.org/zap/zapcore"
)

type redactingCore struct {
	zapcore.Core
}

// WrapCore returns a zapcore.Core that redacts sensitive data from log output.
func WrapCore(core zapcore.Core) zapcore.Core {
	return &redactingCore{Core: core}
}

func (c *redactingCore) With(fields []zapcore.Field) zapcore.Core {
	return &redactingCore{Core: c.Core.With(fields)}
}

func (c *redactingCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *redactingCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	ent.Message = RedactKnownSecrets(ent.Message)
	redacted := make([]zapcore.Field, len(fields))
	for i, f := range fields {
		redacted[i] = redactField(f)
	}
	return c.Core.Write(ent, redacted)
}

func redactField(f zapcore.Field) zapcore.Field {
	switch f.Type {
	case zapcore.StringType:
		f.String = RedactKnownSecrets(f.String)
	case zapcore.ErrorType:
		if f.Interface != nil {
			if err, ok := f.Interface.(error); ok {
				f.String = RedactKnownSecrets(err.Error())
				f.Type = zapcore.StringType
				f.Interface = nil
			}
		}
	}
	return f
}
