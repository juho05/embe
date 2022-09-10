package generator

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type ExprFuncCall struct {
	Name       string
	Signatures []Signature
	Fn         func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error)
}

var ExprFuncCalls = make(map[string]ExprFuncCall)

func newExprFuncCall(name string, fn func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error), signatures ...Signature) {
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

	call.Fn = func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		for _, s := range signatures {
			if len(s.Params) == len(expr.Parameters) {
				return fn(g, expr)
			}
		}
		return nil, signatures[0].ReturnType, g.newError("Wrong argument count.", expr.Name)
	}

	ExprFuncCalls[name] = call
}

func init() {
	newExprFuncCall("mbot.isButtonPressed", exprFuncIsButtonPressed, Signature{Params: []Param{{Name: "button", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("mbot.buttonPressCount", exprFuncButtonPressCount, Signature{Params: []Param{{Name: "button", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("mbot.isJoystickPulled", exprFuncIsJoystickPulled, Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("mbot.joystickPullCount", exprFuncJoystickPullCount, Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("lights.front.brightness", exprFuncLEDAmbientBrightness, Signature{Params: []Param{}, ReturnType: parser.DTNumber}, Signature{Params: []Param{{Name: "light", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})

	newExprFuncCall("sensors.isTilted", exprFuncIsTilted, Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isFaceUp", exprFuncIsFaceUp, Signature{Params: []Param{}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isWaving", exprFuncDetectAction([]string{"up", "down", "left", "right"}, "wave"), Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isRotating", exprFuncDetectAction([]string{"clockwise", "anticlockwise"}, ""), Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isFalling", exprFuncDetectSingleAction("freefall"), Signature{Params: []Param{}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.isShaking", exprFuncDetectSingleAction("shake"), Signature{Params: []Param{}, ReturnType: parser.DTBool})

	newExprFuncCall("sensors.tiltAngle", exprFuncTiltAngle([]string{"forward", "backward", "left", "right"}), Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.rotationAngle", exprFuncTiltAngle([]string{"clockwise", "anticlockwise"}), Signature{Params: []Param{{Name: "direction", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("sensors.acceleration", exprFuncAcceleration, Signature{Params: []Param{{Name: "axis", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.rotation", exprFuncRotation, Signature{Params: []Param{{Name: "axis", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.angleSpeed", exprFuncAngleSpeed, Signature{Params: []Param{{Name: "axis", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("sensors.colorStatus", exprFuncColorStatus, Signature{Params: []Param{{Name: "target", Type: parser.DTString}}, ReturnType: parser.DTNumber}, Signature{Params: []Param{{Name: "target", Type: parser.DTString}, {Name: "inner", Type: parser.DTBool}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.getColor", exprFuncGetColor, Signature{Params: []Param{{Name: "sensor", Type: parser.DTString}, {Name: "valueType", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("sensors.isColorStatus", exprFuncIsColorStatus, Signature{Params: []Param{{Name: "target", Type: parser.DTString}, {Name: "status", Type: parser.DTBool}}, ReturnType: parser.DTBool}, Signature{Params: []Param{{Name: "target", Type: parser.DTString}, {Name: "status", Type: parser.DTNumber}, {Name: "inner", Type: parser.DTBool}}, ReturnType: parser.DTBool})
	newExprFuncCall("sensors.detectColor", exprFuncDetectColor, Signature{Params: []Param{{Name: "sensor", Type: parser.DTString}, {Name: "target", Type: parser.DTString}}, ReturnType: parser.DTBool})
	newExprFuncCall("motors.rpm", exprFuncMotorsSpeed("speed"), Signature{Params: []Param{{Name: "motor", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("motors.power", exprFuncMotorsSpeed("power"), Signature{Params: []Param{{Name: "motor", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("motors.angle", exprFuncMotorsAngle, Signature{Params: []Param{{Name: "motor", Type: parser.DTString}}, ReturnType: parser.DTNumber})

	newExprFuncCall("net.receive", exprFuncNetReceive, Signature{Params: []Param{{Name: "message", Type: parser.DTString}}, ReturnType: parser.DTString})

	newExprFuncCall("math.round", exprFuncMathRound, Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.random", exprFuncMathRandom, Signature{Params: []Param{{Name: "from", Type: parser.DTNumber}, {Name: "to", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.abs", exprFuncMathOp("abs"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.floor", exprFuncMathOp("floor"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.ceil", exprFuncMathOp("ceil"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.sqrt", exprFuncMathOp("sqrt"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.sin", exprFuncMathOp("sin"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.cos", exprFuncMathOp("cos"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.tan", exprFuncMathOp("tan"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.asin", exprFuncMathOp("asin"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.acos", exprFuncMathOp("acos"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.atan", exprFuncMathOp("atan"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.ln", exprFuncMathOp("ln"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.log", exprFuncMathOp("log"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.ePowerOf", exprFuncMathOp("e ^"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})
	newExprFuncCall("math.tenPowerOf", exprFuncMathOp("10 ^"), Signature{Params: []Param{{Name: "n", Type: parser.DTNumber}}, ReturnType: parser.DTNumber})

	newExprFuncCall("strings.length", exprFuncStringsLength, Signature{Params: []Param{{Name: "str", Type: parser.DTString}}, ReturnType: parser.DTNumber})
	newExprFuncCall("strings.letter", exprFuncStringsLetter, Signature{Params: []Param{{Name: "str", Type: parser.DTString}, {Name: "index", Type: parser.DTNumber}}, ReturnType: parser.DTString})
	newExprFuncCall("strings.contains", exprFuncStringsLetter, Signature{Params: []Param{{Name: "str", Type: parser.DTString}, {Name: "substr", Type: parser.DTString}}, ReturnType: parser.DTBool})
}

func exprFuncIsButtonPressed(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorButtonPress, false)

	btn, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTBool, err
	}

	buttons := []string{"a", "b"}
	if !slices.Contains(buttons, btn.(string)) {
		return nil, parser.DTBool, g.newError(fmt.Sprintf("Unknown button. Available options: %s", strings.Join(buttons, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["fieldMenu_1"] = []any{btn.(string), nil}

	return block, parser.DTBool, nil
}

func exprFuncButtonPressCount(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorButtonPressCount, false)

	btn, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	buttons := []string{"a", "b"}
	if !slices.Contains(buttons, btn.(string)) {
		return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown button. Available options: %s", strings.Join(buttons, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["fieldMenu_1"] = []any{btn.(string), nil}

	return block, parser.DTNumber, nil
}

func exprFuncIsJoystickPulled(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorDirectionKeyPress, false)

	direction, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTBool, err
	}

	directions := []string{"up", "down", "left", "right", "middle", "any"}
	if !slices.Contains(directions, direction.(string)) {
		return nil, parser.DTBool, g.newError(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(directions, ", ")), parameterToken(expr.Parameters[0]))
	}
	if direction == "any" {
		direction = "any_direction"
	}

	block.Fields["fieldMenu_1"] = []any{direction.(string), nil}

	return block, parser.DTBool, nil
}

func exprFuncJoystickPullCount(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorDirectionKeyPressCount, false)

	direction, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	directions := []string{"up", "down", "left", "right", "middle"}
	if !slices.Contains(directions, direction.(string)) {
		return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(directions, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["fieldMenu_1"] = []any{direction.(string), nil}

	return block, parser.DTNumber, nil
}

func exprFuncLEDAmbientBrightness(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.UltrasonicGetBrightness, false)

	err := selectAmbientLight(g, block, blocks.UltrasonicGetBrightnessOrder, expr.Name, expr.Parameters, 0, "order", "MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX", false)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.UltrasonicGetBrightnessIndex, true)
	indexMenu.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	return block, parser.DTNumber, nil
}

func exprFuncIsTilted(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorDetectAttitude, false)

	param, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTBool, err
	}

	options := []string{"forward", "backward", "left", "right"}
	if !slices.Contains(options, param.(string)) {
		return nil, parser.DTBool, g.newError(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
	}

	if param == "backward" {
		param = "back"
	}

	block.Fields["tilt"] = []any{param.(string), nil}

	return block, parser.DTBool, nil
}

func exprFuncIsFaceUp(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorDetectAttitude, false)

	block.Fields["tilt"] = []any{"faceup", nil}

	return block, parser.DTBool, nil
}

func exprFuncDetectAction(options []string, prefix string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		block := g.NewBlock(blocks.SensorDetectAction, false)

		param, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
		if err != nil {
			return nil, parser.DTBool, err
		}

		if !slices.Contains(options, param.(string)) {
			return nil, parser.DTBool, g.newError(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
		}

		block.Fields["tilt"] = []any{prefix + param.(string), nil}

		return block, parser.DTBool, nil
	}
}

func exprFuncDetectSingleAction(actionName string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		block := g.NewBlock(blocks.SensorDetectAction, false)

		block.Fields["tilt"] = []any{actionName, nil}

		return block, parser.DTBool, nil
	}
}

func exprFuncTiltAngle(options []string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		block := g.NewBlock(blocks.SensorTiltDegree, false)

		param, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
		if err != nil {
			return nil, parser.DTNumber, err
		}

		if !slices.Contains(options, param.(string)) {
			return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
		}

		switch param.(string) {
		case "forward":
			param = "up"
		case "backward":
			param = "down"
		case "anticlockwise":
			param = "counterclockwise"
		}
		block.Fields["rotation"] = []any{param.(string), nil}

		return block, parser.DTNumber, nil
	}
}

func exprFuncAcceleration(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorAcceleration, false)

	axis, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	options := []string{"x", "y", "z"}
	if !slices.Contains(options, axis.(string)) {
		return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["axis"] = []any{axis.(string), nil}

	return block, parser.DTNumber, nil
}

func exprFuncRotation(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorRotationAngle, false)

	axis, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	options := []string{"x", "y", "z"}
	if !slices.Contains(options, axis.(string)) {
		return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["axis"] = []any{axis.(string), nil}

	return block, parser.DTNumber, nil
}

func exprFuncAngleSpeed(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorAngleSpeed, false)

	axis, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	options := []string{"x", "y", "z"}
	if !slices.Contains(options, axis.(string)) {
		return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["axis"] = []any{axis.(string), nil}

	return block, parser.DTNumber, nil
}

func exprFuncColorStatus(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorColorStatus, false)

	target, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	options := []string{"line", "ground", "white", "red", "yellow", "green", "cyan", "blue", "purple", "black", "custom"}
	if !slices.Contains(options, target.(string)) {
		return nil, parser.DTNumber, g.newError(fmt.Sprintf("Unknown target. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["inputMenu_1"] = []any{target, nil}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorStatusIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	if len(expr.Parameters) == 2 {
		inner, err := g.literal(expr.Name, expr.Parameters[1], parser.DTBool)
		if err != nil {
			return nil, parser.DTNumber, err
		}
		if inner.(bool) {
			block.Type = blocks.SensorColorL1R1Status
			indexMenu.Type = blocks.SensorColorL1R1StatusIndex
		}
	}

	return block, parser.DTNumber, nil
}

func exprFuncGetColor(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorColorGetRGBGrayLight, false)

	var err error
	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SensorColorGetRGBGrayLightInput2, "", "MBUILD_QUAD_COLOR_SENSOR_GET_RGB_GRAY_LIGHT_INPUTMENU_2", block.ID, expr.Name, expr.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		sensors := []string{"L1", "L2", "R1", "R2"}
		if !slices.Contains(sensors, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown sensor. Available options: %s", strings.Join(sensors, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, parser.DTNumber, err
	}

	block.Inputs["inputMenu_3"], err = g.fieldMenu(blocks.SensorColorGetRGBGrayLightInput3, "", "MBUILD_QUAD_COLOR_SENSOR_GET_RGB_GRAY_LIGHT_INPUTMENU_3", block.ID, expr.Name, expr.Parameters[1], parser.DTString, func(v any, token parser.Token) error {
		types := []string{"red", "green", "blue", "gray", "light", "color_sta"}
		if !slices.Contains(types, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown value type. Available options: %s", strings.Join(types, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, parser.DTNumber, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorGetRGBGrayLightIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, parser.DTNumber, nil
}

func exprFuncIsColorStatus(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	blockType := blocks.SensorColorIsStatus
	indexType := blocks.SensorColorIsStatusIndex
	inputType := blocks.SensorColorIsStatusInput
	inputMenuKey := "BLOCK_1618364921511_INPUTMENU_2"
	if len(expr.Parameters) == 3 {
		inner, err := g.literal(expr.Name, expr.Parameters[2], parser.DTBool)
		if err != nil {
			return nil, parser.DTBool, err
		}
		if inner.(bool) {
			blockType = blocks.SensorColorIsStatusL1R1
			indexType = blocks.SensorColorIsStatusL1R1Index
			inputType = blocks.SensorColorIsStatusL1R1Input
			inputMenuKey = "MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INPUTMENU_2"
		}
	}

	block := g.NewBlock(blockType, false)

	target, err := g.literal(expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTBool, err
	}

	options := []string{"line", "ground", "white", "red", "yellow", "green", "cyan", "blue", "purple", "black", "custom"}
	if !slices.Contains(options, target.(string)) {
		return nil, parser.DTBool, g.newError(fmt.Sprintf("Unknown target. Available options: %s", strings.Join(options, ", ")), parameterToken(expr.Parameters[0]))
	}

	block.Fields["inputMenu_1"] = []any{target, nil}

	block.Inputs["inputMenu_2"], err = g.fieldMenu(inputType, "", inputMenuKey, block.ID, expr.Name, expr.Parameters[1], parser.DTNumber, func(v any, token parser.Token) error {
		value := int(v.(float64))
		if blockType == blocks.SensorColorIsStatusL1R1 {
			if math.Mod(v.(float64), 1.0) != 0 || value < 0 || value > 3 {
				return g.newError("Invalid status. Available options: 0-3", token)
			}
		} else {
			if math.Mod(v.(float64), 1.0) != 0 || value < 0 || value > 15 {
				return g.newError("Invalid status. Available options: 0-15", token)
			}
		}
		return nil
	})
	if err != nil {
		return nil, parser.DTBool, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(indexType, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	return block, parser.DTBool, nil
}

func exprFuncDetectColor(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.SensorColorIsLineAndBackground, false)

	var err error
	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SensorColorIsLineAndBackgroundInput2, "", "MBUILD_QUAD_COLOR_SENSOR_IS_LINE_AND_BACKGROUND_INPUTMENU_2", block.ID, expr.Name, expr.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		sensors := []string{"any", "L1", "L2", "R1", "R2"}
		if !slices.Contains(sensors, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown sensor. Available options: %s", strings.Join(sensors, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, parser.DTBool, err
	}

	block.Inputs["inputMenu_3"], err = g.fieldMenu(blocks.SensorColorIsLineAndBackgroundInput3, "", "MBUILD_QUAD_COLOR_SENSOR_IS_LINE_AND_BACKGROUND_INPUTMENU_3", block.ID, expr.Name, expr.Parameters[1], parser.DTString, func(v any, token parser.Token) error {
		types := []string{"line", "ground", "white", "red", "green", "blue", "yellow", "cyan", "purple", "black"}
		if !slices.Contains(types, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown target. Available options: %s", strings.Join(types, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, parser.DTBool, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorIsLineAndBackgroundIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, parser.DTBool, nil
}

func exprFuncMotorsSpeed(unit string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		block := g.NewBlock(blocks.Mbot2EncoderMotorGetSpeed, false)

		var err error
		block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.Mbot2EncoderMotorGetSpeedMenu, "", "MBOT2_ENCODER_MOTOR_GET_SPEED_INPUTMENU_2", block.ID, expr.Name, expr.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
			encoderMotor := v.(string)
			if encoderMotor != "EM1" && encoderMotor != "EM2" {
				return g.newError("Unknown encoder motor. Available options: EM1, EM2", token)
			}
			return nil
		})
		if err != nil {
			return nil, parser.DTNumber, err
		}

		block.Fields["fieldMenu_3"] = []any{unit, nil}

		return block, parser.DTNumber, nil
	}
}

func exprFuncMotorsAngle(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorGetAngle, false)

	var err error
	block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorGetAngleMenu, "", "MBOT2_ENCODER_MOTOR_GET_SPEED_INPUTMENU_2", block.ID, expr.Name, expr.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		encoderMotor := v.(string)
		if encoderMotor != "EM1" && encoderMotor != "EM2" {
			return g.newError("Unknown encoder motor. Available options: EM1, EM2", token)
		}
		return nil
	})
	if err != nil {
		return nil, parser.DTNumber, err
	}

	return block, parser.DTNumber, nil
}

func exprFuncNetReceive(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.NetWifiGetValue, false)

	var err error
	block.Inputs["message"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTString, err
	}

	return block, parser.DTString, nil
}

func exprFuncMathRound(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.OpRound, false)

	var err error
	block.Inputs["NUM"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	return block, parser.DTNumber, nil
}

func exprFuncMathRandom(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.OpRandom, false)

	var err error
	block.Inputs["FROM"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	block.Inputs["TO"], err = g.value(block.ID, expr.Name, expr.Parameters[1], parser.DTNumber)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	return block, parser.DTNumber, nil
}

func exprFuncMathOp(operator string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		block := g.NewBlock(blocks.OpMath, false)

		block.Fields["OPERATOR"] = []any{operator, nil}

		var err error
		block.Inputs["NUM"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, parser.DTNumber, err
		}

		return block, parser.DTNumber, nil
	}
}

func exprFuncStringsLength(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.OpLength, false)
	var err error
	block.Inputs["STRING"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	return block, parser.DTNumber, err
}

func exprFuncStringsLetter(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.OpLetterOf, false)
	var err error
	block.Inputs["STRING"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTString, err
	}
	block.Inputs["LETTER"], err = g.value(block.ID, expr.Name, expr.Parameters[1], parser.DTNumber)
	return block, parser.DTString, err
}

func exprFuncStringsContains(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	block := g.NewBlock(blocks.OpContains, false)
	var err error
	block.Inputs["STRING1"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTBool, err
	}
	block.Inputs["STRING2"], err = g.value(block.ID, expr.Name, expr.Parameters[1], parser.DTString)
	return block, parser.DTBool, err
}

func parameterToken(expr parser.Expr) parser.Token {
	if l, ok := expr.(*parser.ExprLiteral); ok {
		return l.Token
	}
	if c, ok := expr.(*parser.ExprIdentifier); ok {
		return c.Name
	}
	panic("expr must be of type *parser.ExprLiteral or *parser.ExprIdentifier.")
}
