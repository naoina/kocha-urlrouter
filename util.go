package urlrouter

// NextSeparator returns an index of next separator in path.
func NextSeparator(path string, start int) int {
	for start < len(path) && path[start] != '/' && path[start] != '.' {
		start++
	}
	return start
}

// isMetaChar returns whether the meta character.
func IsMetaChar(c byte) bool {
	return c == ParamCharacter || c == WildcardCharacter
}
