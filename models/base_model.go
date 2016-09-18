package models

type Model interface {
	Learn(issues []Issue) void
}
