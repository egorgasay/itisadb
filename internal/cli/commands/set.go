package commands

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Set   = "set"   // usual Set
	uset  = "uset"  // unique Set
	rset  = "rset"  // read only Set
	urset = "ruset" // read only unique Set
)

type SetCommand struct {
	action string
	key    string
	value  string

	server   int32
	unique   bool
	readOnly bool
	level    uint8
}

func (s SetCommand) Action() string {
	return s.action
}

func (s SetCommand) Args() []string {
	return []string{s.key, s.value}
}

func (s SetCommand) Server() int32 {
	return s.server
}

func (s SetCommand) Unique() bool {
	return s.unique
}

func (s SetCommand) ReadOnly() bool {
	return s.readOnly
}

func (s SetCommand) Level() uint8 {
	return s.level
}

func (s SetCommand) Extract() SetCommand {
	return s
}

// ParseSet parses set command.
// [mode]set <key> "<value>" [on <server>] [level <level>]
func ParseSet(action string, split []string) (sc SetCommand, err error) {
	if len(split) < 2 {
		return SetCommand{}, fmt.Errorf("wrong set signature")
	}

	sc.key = split[0]
	unhandledText := strings.Join(split[1:], " ")

	vsplit := strings.Split(unhandledText, `"`)

	if len(vsplit) < 3 {
		return SetCommand{}, fmt.Errorf("wrong set signature. can't parse value")
	}

	sc.value = strings.Join(vsplit[1:len(vsplit)-1], `"`)

	if vsplit[0] != "" {
		return SetCommand{}, fmt.Errorf("unexpected symbol before value")
	}

	after := vsplit[len(vsplit)-1]
	if len(after) > 0 && after[0] != ' ' {
		return SetCommand{}, fmt.Errorf("unexpected symbol after value")
	}

	split = split[2:]
	unhandledText = strings.Join(split, " ")

	for i := 0; i < len(split); i++ {
		switch split[i] {
		case level:
			if i+1 >= len(split) {
				return SetCommand{}, fmt.Errorf("wrong set signature. can't missing level")
			}

			lvl, err := strconv.ParseInt(split[i+1], 10, 8)
			if err != nil {
				return SetCommand{}, fmt.Errorf("wrong set signature. can't parse level")
			}

			sc.level = uint8(lvl)
			split = split[i+2:]
		case on:
			if i+1 >= len(split) {
				return SetCommand{}, fmt.Errorf("wrong set signature. can't missing server")
			}

			server, err := strconv.ParseInt(split[i+1], 10, 32)
			if err != nil {
				return SetCommand{}, fmt.Errorf("wrong set signature. can't parse server")
			}

			sc.server = int32(server)
			split = split[i+2:]
		default:
			return SetCommand{}, fmt.Errorf("unexpected symbol %s", split[i])
		}

		i--
	}

	sc.action = Set
	switch action {
	case Set:
	case uset:
		sc.unique = true
	case rset:
		sc.readOnly = true
	case urset:
		sc.readOnly = true
		sc.unique = true
	}

	return sc, nil
}
