package context

type Context interface {
	GetOrder() string
	GetGroup() string
	GetLimit() int
	GetSkip() int
}
