package generator

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

var ExprFuncCalls = map[string]func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error){
	"mbot.isButtonPressed":   exprFuncIsButtonPressed,
	"mbot.buttonPressCount":  exprFuncButtonPressCount,
	"mbot.isJoystickPulled":  exprFuncIsJoystickPulled,
	"mbot.joystickPullCount": exprFuncJoystickPullCount,

	"sensors.isTilted": exprFuncIsTilted,
	"sensors.isFaceUp": exprFuncIsFaceUp,

	"sensors.isWaving":   exprFuncDetectAction("sensors.isWaving", []string{"up", "down", "left", "right"}, "wave"),
	"sensors.isRotating": exprFuncDetectAction("sensors.isRotating", []string{"clockwise", "anticlockwise"}, ""),
	"sensors.isFalling":  exprFuncDetectSingleAction("sensors.isFalling", "freefall"),
	"sensors.isShaking":  exprFuncDetectSingleAction("sensors.isShaking", "shake"),

	"sensors.tiltAngle":     exprFuncTiltAngle("sensors.tiltAngle", []string{"forward", "backward", "left", "right"}),
	"sensors.rotationAngle": exprFuncTiltAngle("sensors.rotationAngle", []string{"clockwise", "anticlockwise"}),

	"sensors.acceleration": exprFuncAcceleration,
	"sensors.rotation":     exprFuncRotation,
	"sensors.angleSpeed":   exprFuncAngleSpeed,

	"motors.rpm":   exprFuncMotorsSpeed("speed"),
	"motors.power": exprFuncMotorsSpeed("power"),
	"motors.angle": exprFuncMotorsAngle,

	"net.receive": exprFuncNetReceive,

	"math.random":     exprFuncMathRandom,
	"math.round":      exprFuncMathRound,
	"math.abs":        exprFuncMathOp("math.abs", "abs"),
	"math.floor":      exprFuncMathOp("math.floor", "floor"),
	"math.ceil":       exprFuncMathOp("math.ceil", "ceiling"),
	"math.sqrt":       exprFuncMathOp("math.sqrt", "sqrt"),
	"math.sin":        exprFuncMathOp("math.sin", "sin"),
	"math.cos":        exprFuncMathOp("math.cos", "cos"),
	"math.tan":        exprFuncMathOp("math.tan", "tan"),
	"math.asin":       exprFuncMathOp("math.asin", "asin"),
	"math.acos":       exprFuncMathOp("math.acos", "acos"),
	"math.atan":       exprFuncMathOp("math.atan", "atan"),
	"math.ln":         exprFuncMathOp("math.ln", "ln"),
	"math.log":        exprFuncMathOp("math.log", "log"),
	"math.ePowerOf":   exprFuncMathOp("math.ePowerOf", "e ^"),
	"math.tenPowerOf": exprFuncMathOp("math.tenPowerOf", "10 ^"),

	"strings.letter":   exprFuncStringsLetter,
	"strings.length":   exprFuncStringsLength,
	"strings.contains": exprFuncStringsContains,
}

func exprFuncIsButtonPressed(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	if len(expr.Parameters) != 1 {
		return nil, parser.DTBool, g.newError("The `mbot.isButtonPressed` function takes 1 argument: mbot.isButtonPressed(button: string)", expr.Name)
	}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `mbot.buttonPressCount` function takes 1 argument: mbot.buttonPressCount(button: string)", expr.Name)
	}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTBool, g.newError("The `mbot.isJoystickPulledb function takes 1 argument: mbot.isJoystickPulled(direction: string)", expr.Name)
	}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `mbot.joystickPullCount` function takes 1 argument: mbot.joystickPullCount(direction: string)", expr.Name)
	}
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

func exprFuncIsTilted(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	if len(expr.Parameters) != 1 {
		return nil, parser.DTBool, g.newError("The `sensors.isTilted` function takes 1 argument: sensors.isTilted(direction: string)", expr.Name)
	}
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
	if len(expr.Parameters) > 0 {
		return nil, parser.DTBool, g.newError("The `sensors.isFaceUp` function takes no arguments.", expr.Name)
	}
	block := g.NewBlock(blocks.SensorDetectAttitude, false)

	block.Fields["tilt"] = []any{"faceup", nil}

	return block, parser.DTBool, nil
}

func exprFuncDetectAction(name string, options []string, prefix string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		if len(expr.Parameters) != 1 {
			return nil, parser.DTBool, g.newError(fmt.Sprintf("The `%s` function takes 1 argument: %s(direction: string)", name, name), expr.Name)
		}
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

func exprFuncDetectSingleAction(name, actionName string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		if len(expr.Parameters) > 0 {
			return nil, parser.DTBool, g.newError(fmt.Sprintf("The `%s` function takes no arguments.", name), expr.Name)
		}
		block := g.NewBlock(blocks.SensorDetectAction, false)

		block.Fields["tilt"] = []any{actionName, nil}

		return block, parser.DTBool, nil
	}
}

func exprFuncTiltAngle(name string, options []string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		if len(expr.Parameters) != 1 {
			return nil, parser.DTNumber, g.newError(fmt.Sprintf("The `%s` function takes 1 argument: %s(direction: string)", name, name), expr.Name)
		}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `sensors.acceleration` function takes 1 argument: sensors.acceleration(axis: string)", expr.Name)
	}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `sensors.rotation` function takes 1 argument: sensors.rotation(axis: string)", expr.Name)
	}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `sensors.angleSpeed` function takes 1 argument: sensors.angleSpeed(axis: string)", expr.Name)
	}
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

func exprFuncMotorsSpeed(unit string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		if len(expr.Parameters) != 1 {
			return nil, parser.DTNumber, g.newError("The `motors.rpm` function takes 1 argument: motors.rpm(motor: string)", expr.Name)
		}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `motors.angle` function takes 1 argument: motors.angle(motor: string)", expr.Name)
	}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTString, g.newError("The `net.receive` function takes 1 argument: net.receive(message: string).", expr.Name)
	}
	block := g.NewBlock(blocks.NetWifiGetValue, false)

	var err error
	block.Inputs["message"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, parser.DTString, err
	}

	return block, parser.DTString, nil
}

func exprFuncMathRound(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `math.round` function takes 1 argument: math.round(n: number)", expr.Name)
	}
	block := g.NewBlock(blocks.OpRound, false)

	var err error
	block.Inputs["NUM"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, parser.DTNumber, err
	}

	return block, parser.DTNumber, nil
}

func exprFuncMathRandom(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	if len(expr.Parameters) != 2 {
		return nil, parser.DTNumber, g.newError("The `math.random` function takes 2 arguments: math.random(from: number, to: number)", expr.Name)
	}
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

func exprFuncMathOp(name, operator string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
		if len(expr.Parameters) != 1 {
			return nil, parser.DTNumber, g.newError(fmt.Sprintf("The `%s` function takes 1 argument: %s(n: number)", name, name), expr.Name)
		}
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
	if len(expr.Parameters) != 1 {
		return nil, parser.DTNumber, g.newError("The `strings.length` function takes 1 argument: strings.length(str: string)", expr.Name)
	}
	block := g.NewBlock(blocks.OpLength, false)
	var err error
	block.Inputs["STRING"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	return block, parser.DTNumber, err
}

func exprFuncStringsLetter(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	if len(expr.Parameters) != 2 {
		return nil, parser.DTString, g.newError("The `strings.letter` function takes 2 arguments: strings.letter(str: string, index: number)", expr.Name)
	}
	block := g.NewBlock(blocks.OpLetterOf, false)
	var err error
	block.Inputs["STRING"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, "", err
	}
	block.Inputs["LETTER"], err = g.value(block.ID, expr.Name, expr.Parameters[1], parser.DTNumber)
	return block, parser.DTString, err
}

func exprFuncStringsContains(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, parser.DataType, error) {
	if len(expr.Parameters) != 2 {
		return nil, parser.DTBool, g.newError("The `strings.contains` function takes 2 arguments: strings.contains(str: string, substr: string)", expr.Name)
	}
	block := g.NewBlock(blocks.OpContains, false)
	var err error
	block.Inputs["STRING1"], err = g.value(block.ID, expr.Name, expr.Parameters[0], parser.DTString)
	if err != nil {
		return nil, "", err
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
