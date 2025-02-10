package pcregexp

var (
	// pcre2_compile_8 signature:
	//   pcre2_code *pcre2_compile_8(PCRE2_SPTR pattern, PCRE2_SIZE length,
	//       uint32_t options, int *errorcode, PCRE2_SIZE *erroroffset,
	//       pcre2_compile_context *ccontext);
	pcre2_compile func(pattern *uint8, length uint64, options uint32, errorcode *int32, erroroffset *uint64, compileContext uintptr) uintptr

	// pcre2_code_free_8: void pcre2_code_free_8(pcre2_code *code);
	pcre2_code_free func(code uintptr)

	// pcre2_pattern_info_8: int pcre2_pattern_info_8(const pcre2_code *code,
	//    uint32_t what, void *where);
	pcre2_pattern_info func(code uintptr, what uint32, where uintptr) int32

	// pcre2_match_8: int pcre2_match_8(const pcre2_code *code,
	//    PCRE2_SPTR subject, PCRE2_SIZE length, PCRE2_SIZE startoffset,
	//	  uint32_t options, pcre2_match_data *match_data,
	// 	  pcre2_match_context *mcontext);
	pcre2_match func(code uintptr, subject *uint8, length uint64, startoffset uint64, options uint32, matchData uintptr, matchContext uintptr) int32

	// pcre2_match_data_create_from_pattern_8:
	// 	  pcre2_match_data *pcre2_match_data_create_from_pattern_8(
	// 	  	  const pcre2_code *code, pcre2_general_context *gcontext);
	pcre2_match_data_create_from_pattern func(code uintptr, generalContext uintptr) uintptr

	// pcre2_match_data_free_8:
	// 	  void pcre2_match_data_free_8(pcre2_match_data *match_data);
	pcre2_match_data_free func(matchData uintptr)

	// pcre2_get_ovector_pointer_8:
	// 	  PCRE2_SIZE *pcre2_get_ovector_pointer_8(pcre2_match_data *match_data);
	pcre2_get_ovector_pointer func(matchData uintptr) *uint64
)
