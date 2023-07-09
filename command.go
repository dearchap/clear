package clear

import (
	"context"
	"fmt"
)

// Cmd is the basic unit of execution
type Cmd struct {
	name        string
	args        []Arg
	posMap      map[string]Arg
	action      func(*Cmd) error
	pre         func(*Cmd) error
	post        func(*Cmd)
	shortoption bool
}

// CmdOptions are functions to configure the command
type CmdOptions func(*Cmd)

// NoShortOption disables short option processing which
// is enabled by default
func NoShortOption() func(*Cmd) {
	return func(c *Cmd) {
		c.shortoption = false
	}
}

// CmdArgs adds an arg to the cmd
func CmdArgs(args ...Arg) func(*Cmd) {
	return func(c *Cmd) {
		c.AddArgs(args...)
	}
}

// NewCmd returns a new instance of command with given options
func NewCmd(name string, options ...CmdOptions) *Cmd {
	c := &Cmd{
		name:        name,
		shortoption: true,
		posMap:      map[string]Arg{},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Cmd) Run(ctx context.Context, tzer Tokenizer) (err error) {

	if c.pre != nil {
		if err := c.pre(c); err != nil {
			return err
		}
	}

	if c.post != nil {
		defer c.post(c)
	}

	for tzer.HasNext() {

		token := tzer.Peek()

		// -- ends all processing
		if token == "--" {
			return nil
		}

		// check positional args
		if token[0] == '-' {
			if len(token) < 2 {
				return nil
			}

			// okay valid token consume
			_ = tzer.Next()
			argName := token[1:]

			// first loop through all positional args
			if arg, ok := c.posMap[argName]; ok {
				if err := arg.Consume(tzer); err != nil {
					return err
				}
			}

			// not a positional arg. check if we can break up
			// as a short option
			if c.shortoption {
				var split []string
				for _, c := range argName {
					split = append(split, fmt.Sprintf("-%c", c))
				}
				tzer = &chainTzer{
					tzers: []Tokenizer{
						&ssTzer{
							tokens: split,
						},
						tzer,
					},
				}
				continue
			}

			return fmt.Errorf("no valid arg %s", argName)
		}

		// process normal args
		consumed := false
		for _, arg := range c.args {
			if arg.Saturated() {
				continue
			}
			consumed = true
			if err := arg.Consume(tzer); err != nil {
				return err
			}
		}
		if !consumed {
			break
		}
	}

	if c.action != nil {
		return c.action(c)
	}

	return nil
}

func (c *Cmd) AddArgs(args ...Arg) {
	for _, arg := range args {
		c.args = append(c.args, arg)
		if arg.Positional() {
			for _, name := range arg.Names() {
				c.posMap[name] = arg
			}
		}
	}
}
