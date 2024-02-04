package tools

type Tool interface {
	GetName() string
	Check() bool
	Run()
}
