package sqlrender

import (
	"fmt"
	"sort"
	"strings"
)

func qident(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func sqlString(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `''`) + `'`
}

func fqSchema(s TableSchema) string {
	return qident(s.Catalog) + "." + qident(s.Name)
}

func fqTable(t TableRef) string {
	return fqSchema(t.TableSchema) + "." + qident(t.TableName)
}

func fqTableFromTable(t Table) string {
	return fqSchema(t.TableSchema) + "." + qident(t.TableName)
}

func s3Join(base string, parts ...string) string {
	scheme := ""
	rest := base

	if idx := strings.Index(base, "://"); idx >= 0 {
		scheme = base[:idx+3]
		rest = base[idx+3:]
	}

	cleaned := make([]string, 0)
	rest = strings.Trim(rest, "/")

	if rest != "" {
		cleaned = append(cleaned, rest)
	}

	for _, p := range parts {
		p = strings.Trim(p, "/")
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}

	return scheme + strings.Join(cleaned, "/")
}

func sqlValue(v any) string {
	switch x := v.(type) {
	case nil:
		return "NULL"
	case bool:
		if x {
			return "TRUE"
		}
		return "FALSE"
	case string:
		return sqlString(x)
	case int:
		return fmt.Sprintf("%d", x)
	case int64:
		return fmt.Sprintf("%d", x)
	case float64:
		return fmt.Sprintf("%v", x)
	case []string:
		items := make([]string, 0, len(x))
		for _, item := range x {
			items = append(items, sqlValue(item))
		}
		return "ARRAY[" + strings.Join(items, ", ") + "]"
	default:
		return fmt.Sprintf("%v", x)
	}
}

func propsSQL(props map[string]any) string {
	if len(props) == 0 {
		return ""
	}

	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))

	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}

		lines = append(lines, fmt.Sprintf("  %s = %s%s", k, sqlValue(props[k]), comma))
	}

	return strings.Join(lines, "\n")
}

func tableWithProps(t Table) map[string]any {
	result := make(map[string]any)

	for k, v := range t.TableProp.Extra {
		result[k] = v
	}

	if t.TableProp.Format != "" {
		result["format"] = t.TableProp.Format
	}

	if len(t.TableProp.PartitionedBy) > 0 {
		result["partitioned_by"] = t.TableProp.PartitionedBy
	}

	defaultLocation := s3Join(t.TableSchema.Location, t.TableSchema.Name, t.TableName)
	if defaultLocation != "" {
		result["external_location"] = defaultLocation
	}

	return result
}

func notLastColumn(i int, cols []Column) bool {
	return i < len(cols)-1
}
