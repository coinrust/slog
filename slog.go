package slog

var defaultLogger *Logger

var (
	logFuncs = map[Level]func(args ...interface{}){
		DebugLevel: _debug,
		InfoLevel:  _info,
		WarnLevel:  _warn,
		ErrorLevel: _error,
		CritLevel:  _crit,
	}
	logfFuncs = map[Level]func(msg string, args ...interface{}){
		DebugLevel: _debugf,
		InfoLevel:  _infof,
		WarnLevel:  _warnf,
		ErrorLevel: _errorf,
		CritLevel:  _critf,
	}
	logPtrs = map[Level]*func(args ...interface{}){
		DebugLevel: &Debug,
		InfoLevel:  &Info,
		WarnLevel:  &Warn,
		ErrorLevel: &Error,
		CritLevel:  &Crit,
	}
	logfPtrs = map[Level]*func(msg string, args ...interface{}){
		DebugLevel: &Debugf,
		InfoLevel:  &Infof,
		WarnLevel:  &Warnf,
		ErrorLevel: &Errorf,
		CritLevel:  &Critf,
	}
)

// SetDefaultLogger set the logger as the defaultLogger.
// The logging functions in this package use it as their logger.
// This function should be called before using the others.
func SetDefaultLogger(l *Logger) {
	defaultLogger = l

	minLevel := l.GetMinLevel()
	for level, f := range logFuncs {
		if minLevel <= level {
			*logPtrs[level] = f
		} else {
			*logPtrs[level] = nop
		}
	}
	for level, f := range logfFuncs {
		if minLevel <= level {
			*logfPtrs[level] = f
		} else {
			*logfPtrs[level] = nopf
		}
	}
}

func nop(args ...interface{})              {}
func nopf(msg string, args ...interface{}) {}

// Debug logs a _debug level message. It uses fmt.Fprint() to format args.
var Debug func(args ...interface{})

// Debugf logs a _debug level message. It uses fmt.Fprintf() to format msg and args.
var Debugf func(msg string, args ...interface{})

// Info logs a _info level message. It uses fmt.Fprint() to format args.
var Info func(args ...interface{})

// Infof logs a _info level message. It uses fmt.Fprintf() to format msg and args.
var Infof func(msg string, args ...interface{})

// Warn logs a _warning level message. It uses fmt.Fprint() to format args.
var Warn func(args ...interface{})

// Warnf logs a _warning level message. It uses fmt.Fprintf() to format msg and args.
var Warnf func(msg string, args ...interface{})

// Error logs an _error level message. It uses fmt.Fprint() to format args.
var Error func(args ...interface{})

// Errorf logs a _error level message. It uses fmt.Fprintf() to format msg and args.
var Errorf func(msg string, args ...interface{})

// Crit logs a _critical level message. It uses fmt.Fprint() to format args.
var Crit func(args ...interface{})

// Critf logs a _critical level message. It uses fmt.Fprintf() to format msg and args.
var Critf func(msg string, args ...interface{})

func _debug(args ...interface{}) {
	defaultLogger.Log(DebugLevel, "", args...)
}

func _debugf(msg string, args ...interface{}) {
	defaultLogger.Log(DebugLevel, msg, args...)
}

func _info(args ...interface{}) {
	defaultLogger.Log(InfoLevel, "", args...)
}

func _infof(msg string, args ...interface{}) {
	defaultLogger.Log(InfoLevel, msg, args...)
}

func _warn(args ...interface{}) {
	defaultLogger.Log(WarnLevel, "", args...)
}

func _warnf(msg string, args ...interface{}) {
	defaultLogger.Log(WarnLevel, msg, args...)
}

func _error(args ...interface{}) {
	defaultLogger.Log(ErrorLevel, "", args...)
}

func _errorf(msg string, args ...interface{}) {
	defaultLogger.Log(ErrorLevel, msg, args...)
}

func _crit(args ...interface{}) {
	defaultLogger.Log(CritLevel, "", args...)
}

func _critf(msg string, args ...interface{}) {
	defaultLogger.Log(CritLevel, msg, args...)
}
