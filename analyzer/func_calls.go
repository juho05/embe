package analyzer

import (
	"fmt"

	"github.com/Bananenpro/embe/parser"
)

type Param struct {
	Name string
	Type parser.DataType
}

type Signature struct {
	FuncName   string
	Params     []Param
	ReturnType parser.DataType
}

func (s Signature) String() string {
	signature := s.FuncName + "("

	for i, p := range s.Params {
		if i > 0 {
			signature += ", "
		}
		signature = fmt.Sprintf("%s%s: %s", signature, p.Name, p.Type)
	}

	signature += ")"

	if s.ReturnType != "" {
		signature += " : " + string(s.ReturnType)
	}

	return signature
}

type FuncCall struct {
	Name       string
	Signatures []Signature
}

var FuncCalls = make(map[string]FuncCall)

func newFuncCall(name string, signatures ...[]Param) {
	if len(signatures) == 0 {
		signatures = append(signatures, []Param{})
	}

	call := FuncCall{
		Name:       name,
		Signatures: make([]Signature, len(signatures)),
	}

	for i, s := range signatures {
		call.Signatures[i].FuncName = name
		call.Signatures[i].Params = s
	}

	FuncCalls[name] = call
}

func init() {
	newFuncCall("audio.stop")
	newFuncCall("audio.playBuzzer", []Param{{Name: "frequency", Type: parser.DTNumber}}, []Param{{Name: "frequency", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("audio.playClip", []Param{{Name: "name", Type: parser.DTString}}, []Param{{Name: "name", Type: parser.DTString}, {Name: "block", Type: parser.DTBool}})
	newFuncCall("audio.playInstrument", []Param{{Name: "name", Type: parser.DTString}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("audio.playNote", []Param{{Name: "name", Type: parser.DTString}, {Name: "octave", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "note", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("audio.record.start")
	newFuncCall("audio.record.stop")
	newFuncCall("audio.record.play", []Param{}, []Param{{Name: "block", Type: parser.DTBool}})

	newFuncCall("lights.deactivate")
	newFuncCall("lights.back.playAnimation", []Param{{Name: "name", Type: parser.DTString}})
	newFuncCall("lights.front.setBrightness", []Param{{Name: "value", Type: parser.DTNumber}}, []Param{{Name: "light", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lights.front.addBrightness", []Param{{Name: "value", Type: parser.DTNumber}}, []Param{{Name: "light", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lights.front.displayEmotion", []Param{{Name: "emotion", Type: parser.DTString}})
	newFuncCall("lights.front.deactivate", []Param{}, []Param{{Name: "light", Type: parser.DTNumber}})
	newFuncCall("lights.bottom.deactivate")
	newFuncCall("lights.bottom.setColor", []Param{{Name: "color", Type: parser.DTString}})
	newFuncCall("lights.back.display", []Param{{Name: "color1", Type: parser.DTString}, {Name: "color2", Type: parser.DTString}, {Name: "color3", Type: parser.DTString}, {Name: "color4", Type: parser.DTString}, {Name: "color5", Type: parser.DTString}})
	newFuncCall("lights.back.displayColor", []Param{{Name: "color", Type: parser.DTString}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "color", Type: parser.DTString}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("lights.back.displayColorFor", []Param{{Name: "color", Type: parser.DTString}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "color", Type: parser.DTString}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("lights.back.deactivate", []Param{}, []Param{{Name: "led", Type: parser.DTNumber}})
	newFuncCall("lights.back.move", []Param{{Name: "n", Type: parser.DTNumber}})

	newFuncCall("display.print", []Param{{Name: "text", Type: parser.DTString}})
	newFuncCall("display.println", []Param{{Name: "text", Type: parser.DTString}})
	newFuncCall("display.setFontSize", []Param{{Name: "size", Type: parser.DTNumber}})
	newFuncCall("display.setColor", []Param{{Name: "color", Type: parser.DTString}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("display.showLabel", []Param{{Name: "label", Type: parser.DTNumber}, {Name: "text", Type: parser.DTString}, {Name: "location", Type: parser.DTString}, {Name: "size", Type: parser.DTNumber}}, []Param{{Name: "label", Type: parser.DTString}, {Name: "text", Type: parser.DTString}, {Name: "x", Type: parser.DTNumber}, {Name: "y", Type: parser.DTNumber}, {Name: "size", Type: parser.DTNumber}})
	newFuncCall("display.lineChart.addData", []Param{{Name: "value", Type: parser.DTNumber}})
	newFuncCall("display.lineChart.setInterval", []Param{{Name: "interval", Type: parser.DTNumber}})
	newFuncCall("display.barChart.addData", []Param{{Name: "value", Type: parser.DTNumber}})
	newFuncCall("display.table.addData", []Param{{Name: "text", Type: parser.DTString}, {Name: "row", Type: parser.DTNumber}, {Name: "column", Type: parser.DTNumber}})
	newFuncCall("display.setOrientation", []Param{{Name: "orientation", Type: parser.DTNumber}})
	newFuncCall("display.clear")

	newFuncCall("display.setBackgroundColor", []Param{{Name: "color", Type: parser.DTString}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("display.render")

	newFuncCall("sprite.fromIcon", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "name", Type: parser.DTString}})
	newFuncCall("sprite.fromText", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "text", Type: parser.DTString}})
	newFuncCall("sprite.fromQR", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "url", Type: parser.DTString}})
	newFuncCall("sprite.flipH", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.flipV", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.delete", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.setAnchor", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "location", Type: parser.DTString}})
	newFuncCall("sprite.moveLeft", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("sprite.moveRight", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("sprite.moveUp", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("sprite.moveDown", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("sprite.moveTo", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "x", Type: parser.DTNumber}, {Name: "y", Type: parser.DTNumber}})
	newFuncCall("sprite.moveRandom", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.rotate", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "angle", Type: parser.DTNumber}})
	newFuncCall("sprite.rotateTo", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "angle", Type: parser.DTNumber}})
	newFuncCall("sprite.setScale", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "scale", Type: parser.DTNumber}})
	newFuncCall("sprite.setColor", []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "color", Type: parser.DTString}}, []Param{{Name: "sprite", Type: parser.DTImage}, {Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("sprite.resetColor", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.show", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.hide", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.toFront", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.toBack", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.layerUp", []Param{{Name: "sprite", Type: parser.DTImage}})
	newFuncCall("sprite.layerDown", []Param{{Name: "sprite", Type: parser.DTImage}})

	newFuncCall("net.broadcast", []Param{{Name: "message", Type: parser.DTString}}, []Param{{Name: "message", Type: parser.DTString}, {Name: "value", Type: parser.DTString}})
	newFuncCall("net.setChannel", []Param{{Name: "channel", Type: parser.DTNumber}})
	newFuncCall("net.connect", []Param{{Name: "ssid", Type: parser.DTString}, {Name: "password", Type: parser.DTString}})
	newFuncCall("net.reconnect")
	newFuncCall("net.disconnect")

	newFuncCall("sensors.resetAngle", []Param{{Name: "axis", Type: parser.DTString}})
	newFuncCall("sensors.resetYawAngle")
	newFuncCall("sensors.defineColor", []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}, {Name: "tolerance", Type: parser.DTNumber}})
	newFuncCall("sensors.calibrateColors")
	newFuncCall("sensors.enhancedColorDetection", []Param{{Name: "enable", Type: parser.DTBool}})

	newFuncCall("motors.run", []Param{{Name: "rpm", Type: parser.DTNumber}}, []Param{{Name: "rpm", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.runBackward", []Param{{Name: "rpm", Type: parser.DTNumber}}, []Param{{Name: "rpm", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.moveDistance", []Param{{Name: "distance", Type: parser.DTNumber}})
	newFuncCall("motors.moveDistanceBackward", []Param{{Name: "distance", Type: parser.DTNumber}})
	newFuncCall("motors.turnLeft", []Param{{Name: "angle", Type: parser.DTNumber}})
	newFuncCall("motors.turnRight", []Param{{Name: "angle", Type: parser.DTNumber}})
	newFuncCall("motors.rotateRPM", []Param{{Name: "motor", Type: parser.DTString}, {Name: "rpm", Type: parser.DTNumber}}, []Param{{Name: "motor", Type: parser.DTString}, {Name: "rpm", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.rotatePower", []Param{{Name: "motor", Type: parser.DTString}, {Name: "power", Type: parser.DTNumber}}, []Param{{Name: "motor", Type: parser.DTString}, {Name: "power", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.rotateAngle", []Param{{Name: "motor", Type: parser.DTString}, {Name: "angle", Type: parser.DTNumber}})
	newFuncCall("motors.driveRPM", []Param{{Name: "em1RPM", Type: parser.DTNumber}, {Name: "em2RPM", Type: parser.DTNumber}})
	newFuncCall("motors.drivePower", []Param{{Name: "em1Power", Type: parser.DTNumber}, {Name: "em2Power", Type: parser.DTNumber}})
	newFuncCall("motors.stop", []Param{}, []Param{{Name: "motor", Type: parser.DTString}})
	newFuncCall("motors.resetAngle", []Param{}, []Param{{Name: "motor", Type: parser.DTString}})
	newFuncCall("motors.lock", []Param{}, []Param{{Name: "motor", Type: parser.DTString}})
	newFuncCall("motors.unlock", []Param{}, []Param{{Name: "motor", Type: parser.DTString}})

	newFuncCall("time.wait", []Param{{Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "continueCondition", Type: parser.DTBool}})
	newFuncCall("time.resetTimer")

	newFuncCall("mbot.restart")
	newFuncCall("mbot.resetParameters")
	newFuncCall("mbot.calibrateParameters")

	newFuncCall("script.stop")
	newFuncCall("script.stopAll")
	newFuncCall("script.stopOther")

	newFuncCall("lists.append", []Param{{Name: "list", Type: parser.DTStringList}, {Name: "value", Type: parser.DTString}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lists.remove", []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "index", Type: parser.DTNumber}})
	newFuncCall("lists.clear", []Param{{Name: "list", Type: parser.DTStringList}}, []Param{{Name: "list", Type: parser.DTNumberList}})
	newFuncCall("lists.insert", []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTString}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lists.replace", []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTString}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})

	newFuncCall("draw.begin")
	newFuncCall("draw.finish")
	newFuncCall("draw.clear")
	newFuncCall("draw.setColor", []Param{{Name: "color", Type: parser.DTString}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("draw.setThickness", []Param{{Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("draw.setSpeed", []Param{{Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("draw.rotate", []Param{{Name: "angle", Type: parser.DTNumber}})
	newFuncCall("draw.rotateTo", []Param{{Name: "angle", Type: parser.DTNumber}})
	newFuncCall("draw.line", []Param{{Name: "length", Type: parser.DTNumber}})
	newFuncCall("draw.circle", []Param{{Name: "angle", Type: parser.DTNumber}, {Name: "radius", Type: parser.DTNumber}})
	newFuncCall("draw.moveUp", []Param{{Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("draw.moveDown", []Param{{Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("draw.moveLeft", []Param{{Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("draw.moveRight", []Param{{Name: "pixels", Type: parser.DTNumber}})
	newFuncCall("draw.moveTo", []Param{{Name: "x", Type: parser.DTNumber}, {Name: "y", Type: parser.DTNumber}})
	newFuncCall("draw.moveToCenter")
	newFuncCall("draw.save", []Param{{Name: "img", Type: parser.DTImage}})
}
