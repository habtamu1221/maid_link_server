package model

const (
	IBABBYSitting = iota + 1
	ICooking      = iota
	IJanitor      = iota
	ILoundary     = iota
	IHomeCare     = iota
	IAllInclusive = iota
)

const (
	BABBYSitting = "Baby Sitting"
	Cooking      = "Cooking"
	Janitor      = "Janitor"
	Loundary     = "Laundary"
	HomeCare     = "Home Care"
	AllInclusive = "All Inclusive"
)

// working hours
const (
	IFullDayAndNight  = iota + 1
	IFullDay          = iota
	IHalfDayMorning   = iota
	IHalfDayAfternoon = iota
	IFullight         = iota
)

const (
	FullDayAndNight  = "Full day and night"
	FullDay          = "Full day"
	HalfDayMorning   = "Half day morning"
	HalfDayAfternoon = "Half day after noon"
	Fullight         = "Full night"
)

// roles
const (
	ADMIN  = iota //
	MAID   = iota //
	CLIENT = iota //
)

const (
	SADMIN  = "admin" //
	SMAID   = "maid"  //
	SCLIENT = "client"
	SUSER   = "user"
)
