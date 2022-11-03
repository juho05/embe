package generator

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type ExprFuncCall func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error)

var ExprFuncCalls = map[string]ExprFuncCall{
	"mbot.isButtonPressed":   exprFuncIsButtonPressed,
	"mbot.buttonPressCount":  exprFuncButtonPressCount,
	"mbot.isJoystickPulled":  exprFuncIsJoystickPulled,
	"mbot.joystickPullCount": exprFuncJoystickPullCount,

	"lights.front.brightness": exprFuncLEDAmbientBrightness,

	"sensors.isTilted":   exprFuncIsTilted,
	"sensors.isFaceUp":   exprFuncIsFaceUp,
	"sensors.isWaving":   exprFuncDetectAction([]string{"up", "down", "left", "right"}, "wave"),
	"sensors.isRotating": exprFuncDetectAction([]string{"clockwise", "anticlockwise"}, ""),
	"sensors.isFalling":  exprFuncDetectSingleAction("freefall"),
	"sensors.isShaking":  exprFuncDetectSingleAction("shake"),

	"sensors.tiltAngle":     exprFuncTiltAngle([]string{"forward", "backward", "left", "right"}),
	"sensors.rotationAngle": exprFuncTiltAngle([]string{"clockwise", "anticlockwise"}),

	"sensors.acceleration": exprFuncAcceleration,
	"sensors.rotation":     exprFuncRotation,
	"sensors.angleSpeed":   exprFuncAngleSpeed,

	"sensors.colorStatus":   exprFuncColorStatus,
	"sensors.getColorValue": exprFuncGetColorValue,
	"sensors.getColorName":  exprFuncGetColorName,
	"sensors.isColorStatus": exprFuncIsColorStatus,
	"sensors.detectColor":   exprFuncDetectColor,
	"motors.rpm":            exprFuncMotorsSpeed("speed"),
	"motors.power":          exprFuncMotorsSpeed("power"),
	"motors.angle":          exprFuncMotorsAngle,

	"net.receive": exprFuncNetReceive,

	"math.round":      exprFuncMathRound,
	"math.random":     exprFuncMathRandom,
	"math.abs":        exprFuncMathOp("abs"),
	"math.floor":      exprFuncMathOp("floor"),
	"math.ceil":       exprFuncMathOp("ceil"),
	"math.sqrt":       exprFuncMathOp("sqrt"),
	"math.sin":        exprFuncMathOp("sin"),
	"math.cos":        exprFuncMathOp("cos"),
	"math.tan":        exprFuncMathOp("tan"),
	"math.asin":       exprFuncMathOp("asin"),
	"math.acos":       exprFuncMathOp("acos"),
	"math.atan":       exprFuncMathOp("atan"),
	"math.ln":         exprFuncMathOp("ln"),
	"math.log":        exprFuncMathOp("log"),
	"math.ePowerOf":   exprFuncMathOp("e ^"),
	"math.tenPowerOf": exprFuncMathOp("10 ^"),

	"strings.length":   exprFuncStringsLength,
	"strings.letter":   exprFuncStringsLetter,
	"strings.contains": exprFuncStringsContains,

	"lists.get":      exprFuncListsGet,
	"lists.indexOf":  exprFuncListsIndexOf,
	"lists.length":   exprFuncListsLength,
	"lists.contains": exprFuncListsContains,

	"display.pixelIsColor": exprFuncDisplayPixelIsColor,
	"sprite.touchesSprite": exprFuncSpriteTouchesSprite,
	"sprite.touchesEdge":   exprFuncSpriteTouchesEdge,
	"sprite.positionX":     exprFuncSpritePosition("get_x"),
	"sprite.positionY":     exprFuncSpritePosition("get_y"),
	"sprite.rotation":      exprFuncSpritePosition("get_rotation"),
	"sprite.scale":         exprFuncSpritePosition("get_size"),
	"sprite.anchor":        exprFuncSpritePosition("get_align"),
}

func exprFuncIsButtonPressed(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorButtonPress, false)

	btn, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	buttons := []string{"a", "b"}
	if !slices.Contains(buttons, btn.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown button. Available options: %s", strings.Join(buttons, ", ")), expr.Parameters[0])
	}

	block.Fields["fieldMenu_1"] = []any{btn.(string), nil}

	return block, nil
}

func exprFuncButtonPressCount(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorButtonPressCount, false)

	btn, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	buttons := []string{"a", "b"}
	if !slices.Contains(buttons, btn.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown button. Available options: %s", strings.Join(buttons, ", ")), expr.Parameters[0])
	}

	block.Fields["fieldMenu_1"] = []any{btn.(string), nil}

	return block, nil
}

func exprFuncIsJoystickPulled(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorDirectionKeyPress, false)

	direction, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	directions := []string{"up", "down", "left", "right", "middle", "any"}
	if !slices.Contains(directions, direction.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(directions, ", ")), expr.Parameters[0])
	}
	if direction == "any" {
		direction = "any_direction"
	}

	block.Fields["fieldMenu_1"] = []any{direction.(string), nil}

	return block, nil
}

func exprFuncJoystickPullCount(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorDirectionKeyPressCount, false)

	direction, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	directions := []string{"up", "down", "left", "right", "middle"}
	if !slices.Contains(directions, direction.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(directions, ", ")), expr.Parameters[0])
	}

	block.Fields["fieldMenu_1"] = []any{direction.(string), nil}

	return block, nil
}

func exprFuncLEDAmbientBrightness(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.UltrasonicGetBrightness, false)

	err := selectAmbientLight(g, block, blocks.UltrasonicGetBrightnessOrder, expr.Name, expr.Parameters, 0, "order", "MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX", false)
	if err != nil {
		return nil, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.UltrasonicGetBrightnessIndex, true)
	indexMenu.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	return block, nil
}

func exprFuncIsTilted(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorDetectAttitude, false)

	param, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	options := []string{"forward", "backward", "left", "right"}
	if !slices.Contains(options, param.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
	}

	if param == "backward" {
		param = "back"
	}

	block.Fields["tilt"] = []any{param.(string), nil}

	return block, nil
}

func exprFuncIsFaceUp(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorDetectAttitude, false)

	block.Fields["tilt"] = []any{"faceup", nil}

	return block, nil
}

func exprFuncDetectAction(options []string, prefix string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SensorDetectAction, false)

		param, err := g.literal(expr.Parameters[0])
		if err != nil {
			return nil, err
		}

		if !slices.Contains(options, param.(string)) {
			return nil, g.newErrorExpr(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
		}

		block.Fields["tilt"] = []any{prefix + param.(string), nil}

		return block, nil
	}
}

func exprFuncDetectSingleAction(actionName string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SensorDetectAction, false)

		block.Fields["tilt"] = []any{actionName, nil}

		return block, nil
	}
}

func exprFuncTiltAngle(options []string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SensorTiltDegree, false)

		param, err := g.literal(expr.Parameters[0])
		if err != nil {
			return nil, err
		}

		if !slices.Contains(options, param.(string)) {
			return nil, g.newErrorExpr(fmt.Sprintf("Unknown direction. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
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

		return block, nil
	}
}

func exprFuncAcceleration(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorAcceleration, false)

	axis, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	options := []string{"x", "y", "z"}
	if !slices.Contains(options, axis.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
	}

	block.Fields["axis"] = []any{axis.(string), nil}

	return block, nil
}

func exprFuncRotation(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorRotationAngle, false)

	axis, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	options := []string{"x", "y", "z"}
	if !slices.Contains(options, axis.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
	}

	block.Fields["axis"] = []any{axis.(string), nil}

	return block, nil
}

func exprFuncAngleSpeed(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorAngleSpeed, false)

	axis, err := g.literal(expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	options := []string{"x", "y", "z"}
	if !slices.Contains(options, axis.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
	}

	block.Fields["axis"] = []any{axis.(string), nil}

	return block, nil
}

func exprFuncColorStatus(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorStatus, false)

	target, err := g.literal(expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	} else {
		options := []string{"line", "ground", "white", "red", "yellow", "green", "cyan", "blue", "purple", "black", "custom"}
		if !slices.Contains(options, target.(string)) {
			return nil, g.newErrorExpr(fmt.Sprintf("Unknown target. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
		}
		block.Fields["inputMenu_1"] = []any{target, nil}
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorStatusIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	if len(expr.Parameters) == 2 {
		inner, err := g.literal(expr.Parameters[1])
		if err != nil {
			return nil, err
		}
		if inner.(bool) {
			block.Type = blocks.SensorColorL1R1Status
			indexMenu.Type = blocks.SensorColorL1R1StatusIndex
		}
	}

	return block, nil
}

func exprFuncGetColorValue(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorGetRGBGrayLight, false)

	var err error
	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SensorColorGetRGBGrayLightInput2, "", "MBUILD_QUAD_COLOR_SENSOR_GET_RGB_GRAY_LIGHT_INPUTMENU_2", block.ID, expr.Parameters[0], func(v any, token parser.Token) error {
		sensors := []string{"L1", "L2", "R1", "R2"}
		if !slices.Contains(sensors, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown sensor. Available options: %s", strings.Join(sensors, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["inputMenu_3"], err = g.fieldMenu(blocks.SensorColorGetRGBGrayLightInput3, "", "MBUILD_QUAD_COLOR_SENSOR_GET_RGB_GRAY_LIGHT_INPUTMENU_3", block.ID, expr.Parameters[1], func(v any, token parser.Token) error {
		types := []string{"red", "green", "blue", "gray", "light"}
		if !slices.Contains(types, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown value type. Available options: %s", strings.Join(types, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorGetRGBGrayLightIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func exprFuncGetColorName(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorGetRGBGrayLight, false)

	var err error
	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SensorColorGetRGBGrayLightInput2, "", "MBUILD_QUAD_COLOR_SENSOR_GET_RGB_GRAY_LIGHT_INPUTMENU_2", block.ID, expr.Parameters[0], func(v any, token parser.Token) error {
		sensors := []string{"L1", "L2", "R1", "R2"}
		if !slices.Contains(sensors, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown sensor. Available options: %s", strings.Join(sensors, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["inputMenu_3"], err = g.fieldMenu(blocks.SensorColorGetRGBGrayLightInput3, "", "MBUILD_QUAD_COLOR_SENSOR_GET_RGB_GRAY_LIGHT_INPUTMENU_3", block.ID, &parser.ExprLiteral{
		Token: parser.Token{
			Literal:  "color_sta",
			Type:     parser.TkLiteral,
			DataType: parser.DTString,
		},
		ReturnType: parser.DTString,
	}, func(v any, token parser.Token) error {
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorGetRGBGrayLightIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func exprFuncIsColorStatus(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	blockType := blocks.SensorColorIsStatus
	indexType := blocks.SensorColorIsStatusIndex
	inputType := blocks.SensorColorIsStatusInput
	inputMenuKey := "BLOCK_1618364921511_INPUTMENU_2"
	if len(expr.Parameters) == 3 {
		inner, err := g.literal(expr.Parameters[2])
		if err != nil {
			return nil, err
		}
		if inner.(bool) {
			blockType = blocks.SensorColorIsStatusL1R1
			indexType = blocks.SensorColorIsStatusL1R1Index
			inputType = blocks.SensorColorIsStatusL1R1Input
			inputMenuKey = "MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INPUTMENU_2"
		}
	}

	block := g.NewBlock(blockType, false)

	target, err := g.literal(expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	} else {
		options := []string{"line", "ground", "white", "red", "yellow", "green", "cyan", "blue", "purple", "black", "custom"}
		if !slices.Contains(options, target.(string)) {
			return nil, g.newErrorExpr(fmt.Sprintf("Unknown target. Available options: %s", strings.Join(options, ", ")), expr.Parameters[0])
		}

		block.Fields["inputMenu_1"] = []any{target, nil}
	}

	block.Inputs["inputMenu_2"], err = g.fieldMenu(inputType, "", inputMenuKey, block.ID, expr.Parameters[1], func(v any, token parser.Token) error {
		value := int(v.(float64))
		if blockType == blocks.SensorColorIsStatusL1R1 {
			if math.Mod(v.(float64), 1.0) != 0 || value < 0 || value > 3 {
				return g.newErrorTk("Invalid status. Available options: 0-3", token)
			}
		} else {
			if math.Mod(v.(float64), 1.0) != 0 || value < 0 || value > 15 {
				return g.newErrorTk("Invalid status. Available options: 0-15", token)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(indexType, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	return block, nil
}

func exprFuncDetectColor(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorIsLineAndBackground, false)

	var err error
	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SensorColorIsLineAndBackgroundInput2, "", "MBUILD_QUAD_COLOR_SENSOR_IS_LINE_AND_BACKGROUND_INPUTMENU_2", block.ID, expr.Parameters[0], func(v any, token parser.Token) error {
		sensors := []string{"any", "L1", "L2", "R1", "R2"}
		if !slices.Contains(sensors, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown sensor. Available options: %s", strings.Join(sensors, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["inputMenu_3"], err = g.fieldMenu(blocks.SensorColorIsLineAndBackgroundInput3, "", "MBUILD_QUAD_COLOR_SENSOR_IS_LINE_AND_BACKGROUND_INPUTMENU_3", block.ID, expr.Parameters[1], func(v any, token parser.Token) error {
		types := []string{"line", "ground", "white", "red", "green", "blue", "yellow", "cyan", "purple", "black"}
		if !slices.Contains(types, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown target. Available options: %s", strings.Join(types, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorIsLineAndBackgroundIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func exprFuncMotorsSpeed(unit string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2EncoderMotorGetSpeed, false)

		var err error
		block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.Mbot2EncoderMotorGetSpeedMenu, "", "MBOT2_ENCODER_MOTOR_GET_SPEED_INPUTMENU_2", block.ID, expr.Parameters[0], func(v any, token parser.Token) error {
			encoderMotor := v.(string)
			if encoderMotor != "EM1" && encoderMotor != "EM2" {
				return g.newErrorTk("Unknown encoder motor. Available options: EM1, EM2", token)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_3"] = []any{unit, nil}

		return block, nil
	}
}

func exprFuncMotorsAngle(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorGetAngle, false)

	var err error
	block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorGetAngleMenu, "", "MBOT2_ENCODER_MOTOR_GET_SPEED_INPUTMENU_2", block.ID, expr.Parameters[0], func(v any, token parser.Token) error {
		encoderMotor := v.(string)
		if encoderMotor != "EM1" && encoderMotor != "EM2" {
			return g.newErrorTk("Unknown encoder motor. Available options: EM1, EM2", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func exprFuncNetReceive(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetWifiGetValue, false)

	var err error
	block.Inputs["message"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func exprFuncMathRound(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.OpRound, false)

	var err error
	block.Inputs["NUM"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func exprFuncMathRandom(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.OpRandom, false)

	var err error
	block.Inputs["FROM"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["TO"], err = g.value(block.ID, expr.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func exprFuncMathOp(operator string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.OpMath, false)

		block.Fields["OPERATOR"] = []any{operator, nil}

		var err error
		block.Inputs["NUM"], err = g.value(block.ID, expr.Parameters[0])
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func exprFuncStringsLength(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.OpLength, false)
	var err error
	block.Inputs["STRING"], err = g.value(block.ID, expr.Parameters[0])
	return block, err
}

func exprFuncStringsLetter(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.OpLetterOf, false)
	var err error
	block.Inputs["STRING"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["LETTER"], err = g.value(block.ID, expr.Parameters[1])
	return block, err
}

func exprFuncStringsContains(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.OpContains, false)
	var err error
	block.Inputs["STRING1"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["STRING2"], err = g.value(block.ID, expr.Parameters[1])
	return block, err
}

func exprFuncListsGet(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListItem, false)
	err := selectList(g, block, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["INDEX"], err = g.valueWithValidator(block.ID, expr.Parameters[1], nil, 7, "")
	if err != nil {
		return nil, err
	}
	return block, nil
}

func exprFuncListsIndexOf(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListItemIndex, false)
	err := selectList(g, block, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["ITEM"], err = g.value(block.ID, expr.Parameters[1])
	if err != nil {
		return nil, err
	}
	return block, nil
}

func exprFuncListsLength(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListLength, false)
	err := selectList(g, block, expr.Parameters[0])
	if err != nil {
		return nil, err
	}
	return block, nil
}

func exprFuncListsContains(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListContains, false)
	err := selectList(g, block, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["ITEM"], err = g.value(block.ID, expr.Parameters[1])
	if err != nil {
		return nil, err
	}
	return block, nil
}

func exprFuncDisplayPixelIsColor(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteGetColorEqualWithRGB, false)

	var err error
	block.Inputs["number_2"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_3"], err = g.value(block.ID, expr.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_4"], err = g.value(block.ID, expr.Parameters[2])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_5"], err = g.value(block.ID, expr.Parameters[3])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_6"], err = g.value(block.ID, expr.Parameters[4])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func exprFuncSpriteTouchesSprite(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteIsTouchOtherSprite, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["string_2"], err = g.value(block.ID, expr.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func exprFuncSpriteTouchesEdge(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteIsTouchEdge, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, expr.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func exprFuncSpritePosition(getter string) func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
	return func(g *generator, expr *parser.ExprFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SpriteGetXYRotationSizeAlign, false)

		var err error
		block.Inputs["string_1"], err = g.value(block.ID, expr.Parameters[0])
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_2"] = []any{getter, nil}

		return block, nil
	}
}
