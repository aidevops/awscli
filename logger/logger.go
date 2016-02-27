// Package logger -
package logger

// Imports -
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/cli"
)

// Logger -
type Logger struct {
	Context   []string
	Format    string
	Level     string
	File      string
	_file     *os.File
	_canWrite bool
	_buffer   []LogLine
	_ui       cli.Ui
}

// Levels -
var Levels = map[string]int{
	"DEBUG": 5,
	"INFO":  4,
	"WARN":  3,
	"ERROR": 2,
	"CRIT":  1,
	"NIL":   0,
}

// Channels -
var Channels = map[string]int{
	"DEBUG": 5,
	"INFO":  4,
	"WARN":  3,
	"ERROR": 2,
	"CRIT":  1,
	"NIL":   0,
}

// LogLine -
type LogLine struct {
	Context string      `json:"context"`
	Level   string      `json:"level"`
	Time    time.Time   `json:"time"`
	Msg     interface{} `json:"message"` // TODO(JT): refactor this to a []byte slice or bytes.Buffer
}

// Colors -
var (
	BoldRed       = color.New(color.FgRed).Add(color.Underline).Add(color.Bold).SprintFunc()
	Red           = color.New(color.FgRed).Add(color.Bold).SprintFunc()
	BoldMagenta   = color.New(color.FgMagenta).Add(color.Bold).SprintFunc()
	Magenta       = color.New(color.FgMagenta).Add(color.Bold).SprintFunc()
	Green         = color.New(color.FgGreen).SprintFunc()
	Yellow        = color.New(color.FgYellow).SprintFunc()
	BlueHighlight = color.New(color.FgBlue).Add(color.BgWhite).SprintFunc()
	Blue          = color.New(color.FgBlue).SprintFunc()
	Cyan          = color.New(color.FgCyan).SprintFunc()
	White         = color.New(color.FgWhite).Add(color.Underline).Add(color.Bold).SprintFunc()
)

// NilLogger - Mock Logger
func NilLogger() *Logger {
	c := []string{""}
	file := "/dev/null"
	f, _ := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	return &Logger{
		Context:   c,
		Format:    "",
		Level:     "NIL",
		File:      file,
		_file:     f,
		_canWrite: false,
		_buffer:   []LogLine{},
	}
}

// NewLogger -
func NewLogger(level, file, context, format string) *Logger {
	canWrite := true
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open %s in append mode, not writing to file", file)
		canWrite = false
	}
	// this needs to happen somewhere else
	// defer f.Close()

	c := []string{context}
	return &Logger{
		Context:   c,
		Format:    format,
		Level:     strings.ToUpper(level),
		File:      file,
		_file:     f,
		_canWrite: canWrite,
		_buffer:   []LogLine{},
	}
}

// NewCLILogger -
func NewCLILogger(level, file, context, format string, ui cli.Ui) *Logger {
	canWrite := true
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		ui.Error(fmt.Sprintf("failed to open %s in append mode, not writing to file: %s", file, err))
		canWrite = false
	}
	// this needs to happen somewhere else
	// defer f.Close()

	c := []string{context}
	l := &Logger{
		Context:   c,
		Format:    format,
		Level:     strings.ToUpper(level),
		File:      file,
		_file:     f,
		_canWrite: canWrite,
		_ui:       ui,
		_buffer:   []LogLine{},
	}
	return l
}

// SetLevel -
func (l *Logger) SetLevel(level string) {
	l.Level = strings.ToUpper(level)
}

// SetFormat -
func (l *Logger) SetFormat(format string) {
	if format == "" {
		return
	}
	l.Format = format
}

// AddContext -
func (l *Logger) AddContext(prefix string) {
	l.Context = append(l.Context, prefix)
}

// RemoveContext -
func (l *Logger) RemoveContext() (prefix string) {
	prefix, l.Context = l.Context[len(l.Context)-1], l.Context[:len(l.Context)-1]
	return prefix
}

// GetContext - return a string formatted context seperated by >'s
func (l *Logger) GetContext() string {
	return strings.Join(l.Context, ">")
}

// LookupLevel -
func LookupLevel(verbosity string) (channel string, level int) {
	channel = strings.ToUpper(verbosity)
	level = Channels[channel]

	return channel, level
}

// Debugf -
func (l *Logger) Debugf(format string, a ...interface{}) {
	channel, level := LookupLevel("DEBUG")
	if level > Levels[l.Level] {
		return
	}
	l.Logf(channel, format, a...)
}

// Errorf -
func (l *Logger) Errorf(format string, a ...interface{}) {
	channel, level := LookupLevel("ERROR")
	if level > Levels[l.Level] {
		return
	}
	l.Logf(channel, format, a...)
}

// Infof -
func (l *Logger) Infof(format string, a ...interface{}) {
	channel, level := LookupLevel("INFO")
	if level > Levels[l.Level] {
		return
	}
	l.Logf(channel, format, a...)
}

// Warnf -
func (l *Logger) Warnf(format string, a ...interface{}) {
	channel, level := LookupLevel("WARN")
	if level > Levels[l.Level] {
		return
	}
	l.Logf(channel, format, a...)
}

// Logf -
func (l *Logger) Logf(channel string, format string, a ...interface{}) {
	line := fmt.Sprintf(format, a...)
	context := l.GetContext()

	switch outputFormat := strings.ToUpper(l.Format); outputFormat {
	case "JSON", "PRETTY":
		l.LogJSON(channel, context, format, a...)
	case "TEXT":
		l.LogText(channel, context, line)
	default:
		l.LogText(channel, context, line)
	}
}

// directOutput -
func directOutput(ui cli.Ui, channel, prefix, line string) {

	channelType := strings.ToUpper(channel)
	if ui == nil {
		fmt.Printf("%s%s", prefix, line)
	} else {
		switch channelType {
		case "NIL":
			// do nothing...
		case "ERROR":
			ui.Error(fmt.Sprintf("%s%s", prefix, line))
		case "INFO":
			ui.Info(fmt.Sprintf("%s%s", prefix, line))
		case "WARN":
			ui.Warn(fmt.Sprintf("%s%s", prefix, line))
		default:
			ui.Output(fmt.Sprintf("%s%s", prefix, line))
		}
	}
}

// ChannelColor -
func ChannelColor(channel string) string {
	channelType := strings.ToUpper(channel)
	switch channelType {
	case "DEBUG":
		return Cyan(channel)
	case "ERROR":
		return Red(channel)
	case "INFO":
		return Green(channel)
	case "WARN":
		return Yellow(channel)
	default:
		return channel
	}
}

// LogText -
func (l *Logger) LogText(channel, context, line string) {

	prefix := fmt.Sprintf("[%s:%s] ", ChannelColor(channel), context)
	directOutput(l._ui, channel, prefix, line)

	if err := l.writeLine(prefix, line); err != nil {
		directOutput(l._ui, channel, prefix, line)
	}
}

// writeLine -
func (l *Logger) writeLine(prefix, line string) error {
	if _, err := l._file.WriteString(fmt.Sprintf("%s %s", prefix, line)); err != nil {
		l._canWrite = false
		return fmt.Errorf("failed to write data to '%s', skipping", l.File)
	}
	return nil
}

// writeString -
func (l *Logger) writeString(line string) error {
	if _, err := l._file.WriteString(line); err != nil {
		l._canWrite = false
		return fmt.Errorf("failed to write data to '%s', skipping", l.File)
	}
	return nil
}

// String - return a TEXT string version of logline
func (ll *LogLine) String() string {
	out := fmt.Sprintf("Context: %s, Level: %s, Time: %s, Msg: %s", ll.Context, ChannelColor(ll.Level), ll.Time, ll.Msg)
	return out
}

// Bytes - return a []byte slice formatted version of logline
func (ll *LogLine) Bytes() []byte {
	//out := fmt.Sprintf("{ \n  \"Context\": \"%s\", \n  \"Level\": \"%s\", \n  \"Time\": \"%s\", \n  \"Msg\": %s \n}", ll.Context, ll.Level, ll.Time, ll.Msg)
	out, err := json.Marshal(ll)
	if err != nil {
		fmt.Printf("returning zero byte value '%s'\n", err)
	}
	return out
}

// NewLogLine -
func NewLogLine(channel, context string, msg interface{}) *LogLine {
	return &LogLine{Context: context, Level: channel, Time: time.Now(), Msg: msg}
}

// LogJSON -
func (l *Logger) LogJSON(channel, context, format string, a ...interface{}) {
	line := fmt.Sprintf(format, a...)
	jsonLine := NewLogLine(channel, context, line)
	l._buffer = append(l._buffer, *jsonLine)
}

// Flush - Dump the contents of the buffer
func (l *Logger) Flush() {

	var out []byte
	var err error

	switch outputFormat := strings.ToUpper(l.Format); outputFormat {
	case "JSON":
		out, err = json.Marshal(l._buffer)
	case "PRETTY":
		out, err = json.MarshalIndent(l._buffer, "", " ")
	default:
		out, err = json.Marshal(l._buffer)
	}
	if err != nil {
		fmt.Printf("unable to flush json log buffer '%s'\n", err)
	}

	if len(l._buffer) > 0 {
		fmt.Print(string(out))
		fmt.Println()
		if err := l.writeString(string(out)); err != nil {
			context := l.GetContext()
			warnLine := NewLogLine("ERROR", context, fmt.Sprintf("failed to write log '%s', %s", l.File, err))
			fmt.Print(string(warnLine.Bytes()))
			fmt.Println()
		}
	}
}

// GetUI -
func (l *Logger) GetUI() cli.Ui {
	return l._ui
}
