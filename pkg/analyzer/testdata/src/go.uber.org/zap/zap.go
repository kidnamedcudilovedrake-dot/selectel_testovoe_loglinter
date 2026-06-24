package zap

type Logger struct{}

func (l *Logger) Info(msg string, fields ...any) {}
func (l *Logger) Error(msg string, fields ...any) {}
func (l *Logger) Warn(msg string, fields ...any) {}
func (l *Logger) Debug(msg string, fields ...any) {}

type SugaredLogger struct{}

func (s *SugaredLogger) Info(args ...any) {}
func (s *SugaredLogger) Infof(template string, args ...any) {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...any) {}

func (s *SugaredLogger) Error(args ...any) {}
func (s *SugaredLogger) Errorf(template string, args ...any) {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...any) {}

func (s *SugaredLogger) Warn(args ...any) {}
func (s *SugaredLogger) Warnf(template string, args ...any) {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...any) {}

func (s *SugaredLogger) Debug(args ...any) {}
func (s *SugaredLogger) Debugf(template string, args ...any) {}
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...any) {}
