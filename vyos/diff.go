package vyos

import (
	"fmt"
	"github.com/r3labs/diff/v2"
	"log"
	"strings"
)

func CommandsFromDiff(changelog diff.Changelog) []string {
	var commands []string
	for _, change := range changelog {
		switch change.Type {
		case "delete":
			fmt.Printf("delete %s %v\n", strings.Join(change.Path, " "), fieldName(change))
		case "create", "update":
			switch change.To.(type) {
			case map[string]interface{}:
				for _, v := range mapPaths(nil, change.To.(map[string]interface{})) {
					commands = append(commands, fmt.Sprintf("set %s %s", strings.Join(change.Path, " "), v))
				}
			case string, int, bool:
				commands = append(commands, fmt.Sprintf("set %s %v", strings.Join(change.Path, " "), change.To))
			default:
				log.Fatalf("unknown type: %T", change.To)
			}
		default:
			log.Fatalf("unknown change type %q", change.Type)
		}
	}

	return commands
}

func getKeys(m map[string]interface{}) []string {
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}

func mapPaths(path []string, m map[string]interface{}) []string {
	var paths []string
	for k, v := range m {
		if s, ok := v.(map[string]interface{}); ok {
			paths = append(paths, mapPaths(append(path, k), s)...)
		} else {
			paths = append(paths, fmt.Sprintf("%s %v", strings.Join(append(path, k), " "), v))
		}
	}

	return paths
}

func fieldName(change diff.Change) interface{} {
	var field interface{}

	switch change.From.(type) {
	case map[string]interface{}:
		m := change.From.(map[string]interface{})
		var path []string
		for {
			keys := getKeys(m)
			if len(keys) == 1 {
				field = keys[0]
				path = append(path, keys[0])
			} else {
				break
			}
			if v, ok := m[field.(string)].(map[string]interface{}); ok {
				m = v
				continue
			}
			break
		}

		field = strings.Join(path, " ")
	default:
		field = change.From
	}

	return field
}
