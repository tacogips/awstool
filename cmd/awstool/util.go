package awstool

import "time"

func ToJstFormatFunc(fmt string) func(time.Time) string {
	return func(t time.Time) string {
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		inJst := t.In(jst)

		return inJst.Format(fmt)
	}
}
func AsKiB(size int64) int64 {
	return size / 1024
}

func AsMiB(size int64) int64 {
	return size / (1024 * 1024)
}

func AsGiB(size int64) int64 {
	return size / (1024 * 1024 * 1024)
}
