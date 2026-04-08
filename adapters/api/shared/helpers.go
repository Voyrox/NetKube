package shared

import (
	"sort"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Fallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	return value
}

func CloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}

	clone := make(map[string]string, len(values))
	for key, value := range values {
		clone[key] = value
	}

	return clone
}

func ReadyRatio(current, total int32) string {
	return strings.TrimSpace(strings.Join([]string{Int32String(current), Int32String(total)}, "/"))
}

func Int32String(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func Int32PointerString(value *int32) string {
	if value == nil {
		return "-"
	}
	return Int32String(*value)
}

func IntOrStringValue(value *intstr.IntOrString) string {
	if value == nil {
		return "-"
	}
	return value.String()
}

func DesiredReplicas(value *int32) int32 {
	if value == nil {
		return 1
	}

	return *value
}

func ReplicasOrZero(value *int32) int32 {
	if value == nil {
		return 0
	}

	return *value
}

func StringPointer(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func BoolPointer(value *bool) bool {
	return value != nil && *value
}

func YesNo(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MapToSelector(values map[string]string) string {
	if len(values) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(values))
	for key, value := range values {
		parts = append(parts, key+"="+value)
	}
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}

func FormatTime(value *metav1.Time) string {
	if value == nil || value.IsZero() {
		return "-"
	}
	return value.Time.Format("2006-01-02 15:04:05 MST")
}
