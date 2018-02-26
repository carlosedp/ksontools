package jsonnetutil

var (
	log = logger{}
)

// Logger is a logging interface.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

// SetLogger sets the logger for the package.
func SetLogger(l Logger) {
	log.logger = l
}

// Logger is a logger.
type logger struct {
	logger Logger
}

var _ Logger = (*logger)(nil)

// Print prints to the configured logger. Arguments are handled in the same manner of fmt.Print.
func (l *logger) Print(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Print(args...)
}

// Printf prints to the configured logger. Arguments are handled in the same manner of fmt.Printf.
func (l *logger) Printf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Printf(format, args...)
}

// Println prints to the configured logger. Arguments are handled in the same manner of fmt.Println.
func (l *logger) Println(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Println(args...)
}

// Debug prints a debug message to the configured logger. Arguments are handled in the same manner
// of fmt.Print.
func (l *logger) Debug(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Debug(args...)
}

// Debugf prints a debug message to the configured logger. Arguments are handled in the same manner
// of fmt.Printf.
func (l *logger) Debugf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Debugf(format, args...)
}

// Debugln prints a debug message to the configured logger. Arguments are handled in the same manner
// of fmt.Println.
func (l *logger) Debugln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Debugln(args...)
}

// Info prints to the configured logger. Arguments are handled in the same manner of fmt.Print.
func (l *logger) Info(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Info(args...)
}

// Infof prints to the configured logger. Arguments are handled in the same manner of fmt.Printf.
func (l *logger) Infof(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Infof(format, args...)
}

// Infoln prints to the configured logger. Arguments are handled in the same manner of fmt.Println.
func (l *logger) Infoln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Infoln(args...)
}

// Warn prints a warning to the configured logger. Arguments are handled in the same manner of
// fmt.Print.
func (l *logger) Warn(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Warn(args...)
}

// Warnf prints a warning to the configured logger. Arguments are handled in the same manner of
// fmt.Printf.
func (l *logger) Warnf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Warnf(format, args...)
}

// Warnln prints a warning to the configured logger. Arguments are handled in the same manner of
// fmt.Println.
func (l *logger) Warnln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Warnln(args...)
}

// Warning prints a warning to the configured logger. Arguments are handled in the same manner of
// fmt.Print.
func (l *logger) Warning(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Warning(args...)
}

// Warningf prints a warning to the configured logger. Arguments are handled in the same manner of
// fmt.Printf.
func (l *logger) Warningf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Warningf(format, args...)
}

// Warningln prints a warning to the configured logger. Arguments are handled in the same manner of
// fmt.Println.
func (l *logger) Warningln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Warningln(args...)
}

// Error prints an error to the configured logger. Arguments are handled in the same manner of
// fmt.Print.
func (l *logger) Error(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Error(args...)
}

// Errorf prints an error to the configured logger. Arguments are handled in the same manner of
// fmt.Printf.
func (l *logger) Errorf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Errorf(format, args...)
}

// Errorln prints an error to the configured logger. Arguments are handled in the same manner of
// fmt.Println.
func (l *logger) Errorln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Errorln(args...)
}

// Fatal is a equivalent to Print() followed by a call to os.Exit(1).
func (l *logger) Fatal(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Fatal(args...)
}

// Fatalf is a equivalent to Printf() followed by a call to os.Exit(1).
func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Fatalf(format, args...)
}

// Fatalln is a equivalent to Println() followed by a call to os.Exit(1).
func (l *logger) Fatalln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Fatalln(args...)
}

// Panic is a equivalent to Print() followed by a call to panic().
func (l *logger) Panic(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Panic(args...)
}

// Panicf is a equivalent to Printf() followed by a call to panic().
func (l *logger) Panicf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Panicf(format, args...)
}

// Panicln is a equivalent to Println() followed by a call to panic().
func (l *logger) Panicln(args ...interface{}) {
	if l.logger == nil {
		return
	}

	l.logger.Panicln(args...)
}
