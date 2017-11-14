package version

// Versions represents a set of versions
type Versions map[int64]struct{}

// Max returns max version
func (versions Versions) Max() int64 {
	result := int64(0)
	for v := range versions {
		if v > result {
			result = v
		}
	}

	return result
}
