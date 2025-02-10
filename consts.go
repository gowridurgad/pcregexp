package pcregexp

const (
	// PCRE2_ZERO_TERMINATED is used to indicate that the pattern is null
	// terminated.
	PCRE2_ZERO_TERMINATED uint64 = 0
	// PCRE2_INFO_CAPTURECOUNT tells pcre2_pattern_info to return the number of
	// capturing subpatterns.
	PCRE2_INFO_CAPTURECOUNT uint32 = 0
)
