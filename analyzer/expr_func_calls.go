package analyzer

import (
	"github.com/Bananenpro/embe/parser"
)

type ExprFuncCall struct {
	Name       string
	Signatures []Signature
}

var ExprFuncCalls = make(map[string]ExprFuncCall)

func newExprFuncCall(name string, signatures ...Signature) {
	if len(signatures) == 0 {
		signatures = append(signatures, Signature{Params: []Param{}})
	}

	call := ExprFuncCall{
		Name:       name,
		Signatures: make([]Signature, len(signatures)),
	}

	for i, s := range signatures {
		call.Signatures[i].FuncName = name
		call.Signatures[i].Params = s.Params
		call.Signatures[i].ReturnType = s.ReturnType
	}

	ExprFuncCalls[name] = call
}

func init() {
	newExprFuncCall("mbot.isButtonPressed", Signature{Params: []Param{{Name: "button", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("mbot.buttonPressCount", Signature{Params: []Param{{Name: "button", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("mbot.isJoystickPulled", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("mbot.joystickPullCount", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("lights.front.brightness", Signature{Params: []Param{}, ReturnType: parser.DTNumber}, Signature{Params: []Param{{Name: "light", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})

	newExprFuncCall("sensors.isTilted", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isFaceUp", Signature{Params: []Param{}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isWaving", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isRotating", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isFalling", Signature{Params: []Param{}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isShaking", Signature{Params: []Param{}, ReturnType: parser.DTBool})

	newExprFuncCall("sensors.tiltAngle", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.rotationAngle", Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("sensors.acceleration", Signature{Params: []Param{{Name: "axis", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.rotation", Signature{Params: []Param{{Name: "axis", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.angleSpeed", Signature{Params: []Param{{Name: "axis", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("sensors.colorStatus", Signature{Params: []Param{{Name: "target", Type: parser.DTString}}, ReturnType: parser.DTNumber}, Signature{Params: []Param{{Name: "target", Type: parser.DTString}, {Name: "inner", Type: parser.DTBool}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.getColorValue", Signature{Params: []Param{{Name: "sensor", Type: parser.DTString}, {Name: "valueType", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.getColorName", Signature{Params: []Param{{Name: "sensor", Type: parser.DTString}}, ReturnType: parser.DTString})
	newExprFuncCall("sensors.isColorStatus", Signature{Params: []Param{{Name: "target", Type: parser.DTString}, {Name: "status", Type: parser.DTNumber}}, ReturnType: parser.DTBool}, Signature{Params: []Param{{Name: "target", Type: parser.DTString}, {Name: "status", Type: parser.DTNumber}, {Name: "inner", Type: parser.DTBool}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.detectColor", Signature{Params: []Param{{Name: "sensor", Type: parser.DTString}, {Name: "target", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("motors.rpm", Signature{Params: []Param{{Name: "motor", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("motors.power", Signature{Params: []Param{{Name: "motor", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("motors.angle", Signature{Params: []Param{{Name: "motor", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("net.receive", Signature{Params: []Param{{Name: "message", Type: parser.DTString}}, ReturnType: parser.DTString})

	newExprFuncCall("math.round", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.random", Signature{Params: []Param{{Name: "from", Type: parser.DTNumber}, {Name: "to", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.abs", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.floor", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.ceil", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.sqrt", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.sin", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.cos", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.tan", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.asin", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.acos", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.atan", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.ln", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.log", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.ePowerOf", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.tenPowerOf", Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})

	newExprFuncCall("strings.length", Signature{Params: []Param{{Name: "str", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("strings.letter", Signature{Params: []Param{{Name: "str", Type: parser.DTString}, {Name: "index", Type: parser.DTNumber}}, ReturnType: parser.DTString})
	newExprFuncCall("strings.contains", Signature{Params: []Param{{Name: "str", Type: parser.DTString}, {Name: "substr", Type: parser.DTString}}, ReturnType: parser.DTBool})

	newExprFuncCall("lists.get", Signature{Params: []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}}, ReturnType: parser.DTString}, Signature{Params: []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "index", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("lists.indexOf", Signature{Params: []Param{{Name: "list", Type: parser.DTStringList}, {Name: "value", Type: parser.DTString}}, ReturnType: parser.DTNumber}, Signature{Params: []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "value", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("lists.length", Signature{Params: []Param{{Name: "list", Type: parser.DTStringList}}, ReturnType: parser.DTNumber})
	newExprFuncCall("lists.contains", Signature{Params: []Param{{Name: "list", Type: parser.DTStringList}, {Name: "value", Type: parser.DTString}}, ReturnType: parser.DTBool}, Signature{Params: []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "value", Type: parser.DTNumber}}, ReturnType: parser.DTBool})

	newExprFuncCall("display.pixelIsColor", Signature{Params: []Param{{Name: "x", Type: parser.DTNumber}, {Name: "y", Type: parser.DTNumber}, {Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}}, ReturnType: parser.DTBool})
	newExprFuncCall("sprite.touchesSprite", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "other", Type: parser.DTImage}}, ReturnType: parser.DTBool})
	newExprFuncCall("sprite.touchesEdge", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}}, ReturnType: parser.DTBool})
	newExprFuncCall("sprite.positionX", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sprite.positionY", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sprite.rotation", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sprite.scale", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sprite.anchor", Signature{Params: []Param{{Name: "sprite", Type: parser.DTImage}}, ReturnType: parser.DTString})
}
