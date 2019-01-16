package ce

type Status string

const (
	Clean     Status = "clean"
	CatchAll  Status = "catch-all"
	Invalid   Status = "invalid"
	Bounced   Status = "bounced"
	Special   Status = "special"
	BadMX     Status = "bad-mx"
	SpamTrap  Status = "spam-trap"
	Temporary Status = "temporary"
	Unknown   Status = "unknown"
)
