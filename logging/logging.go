package logging

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/weaveworks/common/user"
)

// Setup configures logging output to stderr, sets the log level and sets the formatter.
func Setup(logLevel string) error {
	log.SetOutput(os.Stderr)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("Error parsing log level: %v", err)
	}
	log.SetLevel(level)
	log.SetFormatter(&textFormatter{})
	return nil
}

type textFormatter struct{}

// Based off logrus.TextFormatter, which behaves completely
// differently when you don't want colored output
func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	levelText := strings.ToUpper(entry.Level.String())[0:4]
	timeStamp := entry.Time.Format("2006/01/02 15:04:05.000000")
	if len(entry.Data) > 0 {
		fmt.Fprintf(b, "%s: %s %-44s ", levelText, timeStamp, entry.Message)
		for k, v := range entry.Data {
			fmt.Fprintf(b, " %s=%v", k, v)
		}
	} else {
		// No padding when there's no fields
		fmt.Fprintf(b, "%s: %s %s", levelText, timeStamp, entry.Message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// WithContext returns a logger that has information from the context as fields.
//
// e.g.
//     logger := logging.WithContext(ctx)
//     logger.Errorf("Some error")
func WithContext(ctx context.Context) *log.Entry {
	fields := log.Fields{}
	userID, err := user.ExtractUserID(ctx)
	if err != nil {
		fields["userID"] = userID
	}
	orgID, err := user.ExtractOrgID(ctx)
	if err != nil {
		fields["orgID"] = orgID
	}
	return log.WithFields(fields)
}
