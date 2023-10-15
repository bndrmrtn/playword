package helpers

// MergeMaps merges 2 maps into one
func MergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
