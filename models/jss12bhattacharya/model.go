package models

type Model interface {
	Train(Issue *issues) void
	Learn(Issue *issues) void
}
