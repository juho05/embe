package analyzer

import (
	"fmt"

	"github.com/Bananenpro/embe/parser"
)

type Var struct {
	Name     string
	DataType parser.DataType
}

func (v Var) String() string {
	return fmt.Sprintf("var %s: %s", v.Name, v.DataType)
}

var Variables = make(map[string]Var)

func newVar(name string, dataType parser.DataType) {
	Variables[name] = Var{
		Name:     name,
		DataType: dataType,
	}
}

func init() {
	newVar("audio.volume", parser.DTNumber)
	newVar("audio.speed", parser.DTNumber)

	newVar("lights.back.brightness", parser.DTNumber)

	newVar("time.timer", parser.DTNumber)

	newVar("mbot.battery", parser.DTNumber)
	newVar("mbot.mac", parser.DTString)
	newVar("mbot.hostname", parser.DTString)

	newVar("sensors.wavingAngle", parser.DTNumber)
	newVar("sensors.wavingSpeed", parser.DTNumber)
	newVar("sensors.shakingStrength", parser.DTNumber)
	newVar("sensors.brightness", parser.DTNumber)
	newVar("sensors.loudness", parser.DTNumber)
	newVar("sensors.distance", parser.DTNumber)
	newVar("sensors.outOfRange", parser.DTBool)
	newVar("sensors.lineDeviation", parser.DTNumber)

	newVar("net.connected", parser.DTBool)

	newVar("draw.positionX", parser.DTNumber)
	newVar("draw.positionY", parser.DTNumber)
	newVar("draw.rotation", parser.DTNumber)
	newVar("draw.thickness", parser.DTNumber)

	newVar("math.e", parser.DTNumber)
	newVar("math.pi", parser.DTNumber)
	newVar("math.phi", parser.DTNumber)
}
