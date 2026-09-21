package email

// IANA / RFC 8621 §4.1.1 keyword constants.
const (
	KeywordSeen      = "$seen"
	KeywordDraft     = "$draft"
	KeywordFlagged   = "$flagged"
	KeywordAnswered  = "$answered"
	KeywordForwarded = "$forwarded"
	KeywordPhishing  = "$phishing"
	KeywordJunk      = "$junk"
	KeywordNotJunk   = "$notjunk"
)
