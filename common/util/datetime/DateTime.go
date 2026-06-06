package datetime

import "time"

var DATE_FORMATS = []string{"2006-01-02"}
var DATE_TIME_FORMATS = []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "Today 15:04:05", time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano}

func GetDateTimeFormatted() string {
	return time.Now().Format(DATE_TIME_FORMATS[0])
}

func ParseDate(input *string) (result time.Time, err error) {
	for _, dateFormat := range DATE_FORMATS {
		result, err = time.Parse(dateFormat, *input)
		if err == nil {
			return result, nil
		}
	}
	return result, err
}

func ParseDateTime(input *string) (result time.Time, err error) {
	for _, dateFormat := range DATE_TIME_FORMATS {
		result, err = time.Parse(dateFormat, *input)
		if err == nil {
			return result, nil
		}
	}
	return result, err
}

func FormatDate(input time.Time) string {
	return input.Format(DATE_FORMATS[0])
}

func FormatDateTime(input time.Time) string {
	/*if IsToday(input) {
		return input.Format(DATE_TIME_FORMATS[1])
	}*/
	return input.Format(DATE_TIME_FORMATS[0])
}

func IsToday(input time.Time) bool {
	y1, m1, d1 := input.Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
