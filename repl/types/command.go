package types

type Command interface {
	Help() string
	Do(args []string) error
}
