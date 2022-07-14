package occupation_types

type OccupationType int

//go:generate enumer -type=OccupationType -json
const (
	UNEMPLOYED OccupationType = iota + 1
	STUDENT
	EMPLOYED
)

var strs = [...]string{
	"Unemployed",
	"Student",
	"Employed",
}
