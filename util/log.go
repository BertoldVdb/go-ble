package bleutil

import "github.com/sirupsen/logrus"

// LogWithPrefix extends a logrus prefix with another one.
func LogWithPrefix(input *logrus.Entry, extra string) *logrus.Entry {
	if input == nil {
		return nil
	}

	base := ""
	if value, ok := input.Data["prefix"]; ok {
		base = value.(string) + "/"
	}

	return input.WithField("prefix", base+extra)
}
