package shared

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FormatAge(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "-"
	}

	return FormatDuration(time.Since(timestamp.Time))
}

func FormatOptionalAge(timestamp *metav1.Time) string {
	if timestamp == nil || timestamp.IsZero() {
		return "-"
	}

	return FormatAge(*timestamp)
}

func EventTimestamp(event *metav1.ObjectMeta) time.Time {
	if event == nil {
		return time.Time{}
	}

	return event.CreationTimestamp.Time
}

func FormatDuration(duration time.Duration) string {
	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds < 0 {
			seconds = 0
		}
		return fmt.Sprintf("%ds", seconds)
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	if hours == 0 {
		return fmt.Sprintf("%dd", days)
	}
	return fmt.Sprintf("%dd %dh", days, hours)
}

func FormatCPU(milli int64) string {
	if milli < 1000 {
		return fmt.Sprintf("%dm", milli)
	}
	value := float64(milli) / 1000
	return TrimFloat(value) + " cores"
}

func FormatBytes(bytes int64) string {
	const (
		ki = 1024
		mi = ki * 1024
		gi = mi * 1024
	)

	switch {
	case bytes >= gi:
		return TrimFloat(float64(bytes)/float64(gi)) + "Gi"
	case bytes >= mi:
		return TrimFloat(float64(bytes)/float64(mi)) + "Mi"
	case bytes >= ki:
		return TrimFloat(float64(bytes)/float64(ki)) + "Ki"
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func PercentageFloat(used, total int64) float64 {
	if total <= 0 {
		return 0
	}
	return (float64(used) / float64(total)) * 100
}

func TrimFloat(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)
	for len(formatted) > 0 && formatted[len(formatted)-1] == '0' {
		formatted = formatted[:len(formatted)-1]
	}
	if len(formatted) > 0 && formatted[len(formatted)-1] == '.' {
		formatted = formatted[:len(formatted)-1]
	}
	return formatted
}
