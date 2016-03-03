package _test_util

const (
	NotFindString = "404 page not found\n"
)

var (
	ListenPort int = 4000
)

type TestUser struct {
	Name  string
	Age   int
	Money float64
	Notes []TestNote
	Alive bool
}

type TestNote struct {
	Text string
}