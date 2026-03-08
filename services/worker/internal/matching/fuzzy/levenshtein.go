package fuzzy

// LevenshteinDistance computes the Levenshtein distance between two strings.
// The Levenshtein distance is the minimum number of single-character edits
// (insertions, deletions, or substitutions) required to change one string into another.
func LevenshteinDistance(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Create a matrix to store distances
	d := make([][]int, len1+1)
	for i := range d {
		d[i] = make([]int, len2+1)
	}

	// Initialize first column and row
	for i := 0; i <= len1; i++ {
		d[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		d[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			d[i][j] = min(
				d[i-1][j]+1,     // deletion
				d[i][j-1]+1,     // insertion
				d[i-1][j-1]+cost, // substitution
			)
		}
	}

	return d[len1][len2]
}

// SimilarityScore computes a similarity score between 0 and 1 based on Levenshtein distance.
// A score of 1.0 means identical strings, 0.0 means completely different.
func SimilarityScore(s1, s2 string) float64 {
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	if maxLen == 0 {
		return 1.0 // Both strings are empty
	}

	distance := LevenshteinDistance(s1, s2)
	return 1.0 - (float64(distance) / float64(maxLen))
}

// min returns the minimum of three integers.
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
