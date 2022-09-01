package blocks

type BlockType string

// Events
const (
	EventLaunch                     BlockType = "cyberpi.cyberpi_when_launch"
	EventButtonPress                BlockType = "cyberpi.cyberpi_when_button_press"
	EventDirectionKeyPress          BlockType = "cyberpi.cyberpi_when_direction_key_press"
	EventDetectAttitude             BlockType = "cyberpi.cyberpi_when_detect_attitude"
	EventDetectAction               BlockType = "cyberpi.cyberpi_when_detect_action"
	EventSensorValueBiggerOrSmaller BlockType = "cyberpi.cyberpi_when_sensor_value_bigger_or_smaller_than"
)

// Sensors
const (
	SensorDetectAttitude                BlockType = "cyberpi.cyberpi_detect_attitude"
	SensorDetectAction                  BlockType = "cyberpi.cyberpi_detect_action"
	SensorBatteryLevelMacAddressAndSoOn BlockType = "cyberpi.cyberpi_battery_macaddress_blename_and_so_on"
	SensorLoudness                      BlockType = "cyberpi.cyberpi_loudness"
	SensorBrightness                    BlockType = "cyberpi.cyberpi_brightness"
	SensorUltrasonicDistance            BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_get_distance"
	SensorUltrasonicDistanceMenu        BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_get_distance_index_menu"
	SensorUltrasonicOutOfRange          BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_out_of_range"
	SensorUltrasonicOutOfRangeMenu      BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_out_of_range_index_menu"

	SensorShakingStrength BlockType = "cyberpi.cyberpi_shaked_value"

	SensorWaveAngle BlockType = "cyberpi.cyberpi_wave_angle"
	SensorWaveSpeed BlockType = "cyberpi.cyberpi_wave_speed"

	SensorTiltDegree BlockType = "cyberpi.cyberpi_tilt_degree"

	SensorAcceleration  BlockType = "cyberpi.cyberpi_axis_acceleration"
	SensorAngleSpeed    BlockType = "cyberpi.cyberpi_axis_angle_speed"
	SensorRotationAngle BlockType = "cyberpi.cyberpi_axis_rotation_angle"

	SensorButtonPress            BlockType = "cyberpi.cyberpi_button_press"
	SensorButtonPressCount       BlockType = "cyberpi.cyberpi_button_count"
	SensorDirectionKeyPress      BlockType = "cyberpi.cyberpi_direction_key_press"
	SensorDirectionKeyPressCount BlockType = "cyberpi.cyberpi_direction_key_count"
)

// Audio
const (
	AudioGetVolume BlockType = "cyberpi.cyberpi_get_volume"
	AudioSetVolume BlockType = "cyberpi.cyberpi_set_volume"
	AudioAddVolume BlockType = "cyberpi.cyberpi_add_volume"

	AudioGetSpeed BlockType = "cyberpi.cyberpi_get_audio_speed"
	AudioSetSpeed BlockType = "cyberpi.cyberpi_set_audio_speed"
	AudioAddSpeed BlockType = "cyberpi.cyberpi_add_audio_speed"

	AudioStop                          BlockType = "cyberpi.cyberpi_stop_audio"
	AudioPlayBuzzerTone                BlockType = "cyberpi.cyberpi_play_buzzer_tone"
	AudioPlayBuzzerToneWithTime        BlockType = "cyberpi.cyberpi_play_buzzer_tone_with_time"
	AudioPlayClip                      BlockType = "cyberpi.cyberpi_play_audio_3"
	AudioPlayClipFileNameMenu          BlockType = "cyberpi.cyberpi_play_audio_3_file_name_menu"
	AudioPlayClipUntilDone             BlockType = "cyberpi.cyberpi_play_audio_until_3"
	AudioPlayClipUntilDoneFileNameMenu BlockType = "cyberpi.cyberpi_play_audio_until_3_file_name_menu"
	AudioPlayNote                      BlockType = "cyberpi.cyberpi_play_music_with_tone_and_note_2"
	AudioNote                          BlockType = "note"
	AudioPlayMusicInstrument           BlockType = "cyberpi.cyberpi_play_music_with_note"
	AudioPlayMusicInstrumentMenu       BlockType = "cyberpi.cyberpi_play_music_with_note_fieldMenu_1_menu"

	AudioRecordStart         BlockType = "cyberpi.cyberpi_start_record"
	AudioRecordStop          BlockType = "cyberpi.cyberpi_stop_record"
	AudioRecordPlay          BlockType = "cyberpi.cyberpi_play_record"
	AudioRecordPlayUntilDone BlockType = "cyberpi.cyberpi_play_record_until"
)

// LED
const (
	LEDPlayAnimation                             BlockType = "cyberpi.cyberpi_play_led_animation_until"
	LEDDisplay                                   BlockType = "cyberpi.cyberpi_show_led"
	LEDDisplaySingleColor                        BlockType = "cyberpi.cyberpi_led_show_single_with_color_2"
	LEDDisplaySingleColorFieldMenu               BlockType = "cyberpi.cyberpi_led_show_single_with_color_2_fieldMenu_1_menu"
	LEDDisplaySingleColorWithTime                BlockType = "cyberpi.cyberpi_led_show_single_with_color_and_time_2"
	LEDDisplaySingleColorWithTimeFieldMenu       BlockType = "cyberpi.cyberpi_led_show_single_with_color_and_time_2_fieldMenu_1_menu"
	LEDDisplaySingleColorWithRGB                 BlockType = "cyberpi.cyberpi_led_show_single_with_rgb_2"
	LEDDisplaySingleColorWithRGBFieldMenu        BlockType = "cyberpi.cyberpi_led_show_single_with_rgb_2_fieldMenu_1_menu"
	LEDDisplaySingleColorWithRGBAndTime          BlockType = "cyberpi.cyberpi_led_show_single_with_rgb_and_time"
	LEDDisplaySingleColorWithRGBAndTimeFieldMenu BlockType = "cyberpi.cyberpi_led_show_single_with_rgb_and_time_fieldMenu_1_menu"
	LEDOff                                       BlockType = "cyberpi.cyberpi_led_off_2"
	LEDOffFieldMenu                              BlockType = "cyberpi.cyberpi_led_off_2_fieldMenu_1_menu"
	LEDMove                                      BlockType = "cyberpi.cyberpi_move_led"
	LEDGetBrightness                             BlockType = "cyberpi.cyberpi_get_led_brightness"
	LEDSetBrightness                             BlockType = "cyberpi.cyberpi_set_led_brightness"
	LEDAddBrightness                             BlockType = "cyberpi.cyberpi_add_led_brightness"
)

// Display
const (
	DisplayPrintln                        BlockType = "cyberpi.cyberpi_display_println"
	DisplayPrint                          BlockType = "cyberpi.cyberpi_display_print"
	DisplaySetFont                        BlockType = "cyberpi.cyberpi_console_set_font"
	DisplaySetFontMenu                    BlockType = "cyberpi.cyberpi_console_set_font_inputMenu_1_menu"
	DisplayLabelShowSomewhereWithSize     BlockType = "cyberpi.cyberpi_display_label_show_at_somewhere_with_size"
	DisplayLabelShowSomewhereWithSizeMenu BlockType = "cyberpi.cyberpi_display_label_show_at_somewhere_with_size_inputMenu_4_menu"
	DisplayLabelShowXYWithSize            BlockType = "cyberpi.cyberpi_display_label_show_label_xy_with_size"
	DisplayLabelShowXYWithSizeMenu        BlockType = "cyberpi.cyberpi_display_label_show_label_xy_with_size_inputMenu_4_menu"
	DisplayLineChartAddData               BlockType = "cyberpi.cyberpi_display_line_chart_add_data"
	DisplayLineChartSetInterval           BlockType = "cyberpi.cyberpi_display_bar_chart_set_interval" // not a typo
	DisplayBarChartAddData                BlockType = "cyberpi.cyberpi_display_bar_chart_add_data"
	DisplayTableAddDataAtRowColumn        BlockType = "cyberpi.cyberpi_display_table_add_data_at_row_column_2"
	DisplayTableAddDataAtRowColumnMenu    BlockType = "cyberpi.cyberpi_display_table_add_data_at_row_column_2_fieldMenu_1_menu"
	DisplaySetBrushColor                  BlockType = "cyberpi.cyberpi_display_set_brush_with_color"
	DisplaySetBrushColorRGB               BlockType = "cyberpi.cyberpi_display_set_brush_with_r_g_b"
	DisplayClear                          BlockType = "cyberpi.cyberpi_display_clear"
	DisplaySetOrientation                 BlockType = "cyberpi.cyberpi_display_rotate_to_2"
	DisplaySetOrientationMenu             BlockType = "cyberpi.cyberpi_display_rotate_to_2_fieldMenu_1_menu"
)

// Control
const (
	ControlIf            BlockType = "control_if"
	ControlIfElse        BlockType = "control_if_else"
	ControlWait          BlockType = "control_wait"
	ControlWaitUntil     BlockType = "control_wait_until"
	ControlRepeat        BlockType = "control_repeat"
	ControlRepeatUntil   BlockType = "control_repeat_until"
	ControlRepeatForever BlockType = "control_forever"
	ControlStop          BlockType = "control_stop"
	ControlRestart       BlockType = "cyberpi.cyberpi_restart"
)

// Operators
const (
	OpEquals      BlockType = "operator_equals"
	OpOr          BlockType = "operator_or"
	OpAnd         BlockType = "operator_and"
	OpLessThan    BlockType = "operator_lt"
	OpGreaterThan BlockType = "operator_gt"
	OpNot         BlockType = "operator_not"

	OpAdd      BlockType = "operator_add"
	OpSubtract BlockType = "operator_subtract"
	OpMultiply BlockType = "operator_multiply"
	OpDivide   BlockType = "operator_divide"
	OpMod      BlockType = "operator_mod"

	OpRound  BlockType = "operator_round"
	OpRandom BlockType = "operator_random"
	OpMath   BlockType = "operator_mathop"
)

// Strings
const (
	OpJoin     BlockType = "operator_join"
	OpLetterOf BlockType = "operator_letter_of"
	OpContains BlockType = "operator_contains"
	OpLength   BlockType = "operator_length"
)

// Variables
const (
	VariableSetTo    BlockType = "data_setvariableto"
	VariableChangeBy BlockType = "data_changevariableby"
)
