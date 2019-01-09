package ce

type Status string

const (
	Clean     Status = "clean"
	CatchAll         = "catch-all"
	Invalid          = "invalid"
	Bounced          = "bounced"
	Special          = "special"
	BadMX            = "bad-mx"
	SpamTrap         = "spam-trap"
	Temporary        = "temporary"
	Unknown          = "unknown"
)
