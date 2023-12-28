package commands

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Set = "set" // usual Set
)

type SetCommand struct {
	action string
	key    string
	value  string

	server int32
	mode   uint8
	level  uint8
}

const (
	defaultSetMode = iota
	readOnlySetMode
	notExistingSetMode
	existingSetMode
)

const (
	defaultLevel = iota
	restrictedLevel
	secretLevel
)

func (s SetCommand) Action() string {
	return s.action
}

func (s SetCommand) Args() []string {
	return []string{s.key, s.value}
}

func (s SetCommand) Server() int32 {
	return s.server
}

func (s SetCommand) Mode() uint8 {
	return s.mode
}

func (s SetCommand) Level() uint8 {
	return s.level
}

func (s SetCommand) Extract() SetCommand {
	return s
}

// ParseSet parses set command.
/*
------------------- [ MODE ] --- [    LEVEL     ] - [    SERVER    ]


SET key "value" [ NX | RO | XX ] [ D | R | S ] [ [0-9]+ ]

----------------------------------------------------------------------

MODE - Defines the mode of the operation.

- `NX` - If the key already exists, it won't be overwritten.

- `RO` - If the key already exists, an error will be returned.

- `XX` - If the key doesn't exist, it won't be created.

----------------------------------------------------------------------

LEVEL - Defines the level of permission.

- `R` (Restricted) - NO encryption, ACL validation

- `S` (Secret) - encryption, ACL validation

By default - NO encryption, NO ACL validation

----------------------------------------------------------------------

SERVER - Defines server number to use.

- Automatically saving to a less loaded server by default.

----------------------------------------------------------------------

Examples:

@> SET key "value"

@> SET key "value" XX

@> SET key "value" R

@> SET key "value" XX R 1

*/
func ParseSet(split []string) (sc SetCommand, err error) {
	if len(split) < 2 {
		return SetCommand{}, fmt.Errorf("wrong set signature")
	}

	sc.action = Set
	sc.key = split[0]
	unhandledText := strings.Join(split[1:], " ")

	vSplit := strings.Split(unhandledText, `"`)

	if len(vSplit) < 3 {
		return SetCommand{}, fmt.Errorf("wrong set signature. can't parse value")
	}

	sc.value = strings.Join(vSplit[1:len(vSplit)-1], `"`)

	if vSplit[0] != "" {
		return SetCommand{}, fmt.Errorf("unexpected symbol before value")
	}

	after := vSplit[len(vSplit)-1]
	if len(after) > 0 && after[0] != ' ' {
		return SetCommand{}, fmt.Errorf("unexpected symbol after value")
	}

	split = split[2:]
	unhandledText = strings.Join(split, " ")

	for i := 0; i < len(split); i++ {
		switch split[i] {
		case "RO":
			sc.mode = readOnlySetMode
		case "NX":
			sc.mode = notExistingSetMode
		case "XX":
			sc.mode = existingSetMode
		case "R":
			sc.level = 1
		case "S":
			sc.level = 2
		default:
			num, err := strconv.ParseInt(split[i], 10, 32)
			if err != nil {
				return SetCommand{}, fmt.Errorf("wrong set signature. can't recognize [%s]", split[i])
			}

			sc.server = int32(num)
		}
	}

	return sc, nil
}
