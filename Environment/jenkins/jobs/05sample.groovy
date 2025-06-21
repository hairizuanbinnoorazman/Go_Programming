listView('example') {
    jobFilters {
        regex {
            matchType(MatchType.EXCLUDE_UNMATCHED)
            matchValue(RegexMatchValue.NAME)
            regex('Hola')
        }
    }
}