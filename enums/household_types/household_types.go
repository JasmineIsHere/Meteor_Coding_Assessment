package household_types

type HouseholdType int

//go:generate enumer -type=HouseholdType -json
const (
	LANDED HouseholdType = iota + 1
	CONDOMINIUM
	HDB
)

var strs = [...]string{
	"Landed",
	"Condominium",
	"HDB",
}
