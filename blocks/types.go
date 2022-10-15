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
	EventReceivedMessage            BlockType = "cyberpi.cyberpi_wifi_broadcast_when_received_message"
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

	SensorColorDefineColor               BlockType = "mbuild_quad_color_sensor.BLOCK_1626250042594"
	SensorColorDefineColorIndex          BlockType = "mbuild_quad_color_sensor.BLOCK_1626250042594_index_menu"
	SensorColorL1R1Status                BlockType = "mbuild_quad_color_sensor.BLOCK_1618397596925"
	SensorColorL1R1StatusIndex           BlockType = "mbuild_quad_color_sensor.BLOCK_1618397596925_index_menu"
	SensorColorStatus                    BlockType = "mbuild_quad_color_sensor.BLOCK_1618397679204"
	SensorColorStatusIndex               BlockType = "mbuild_quad_color_sensor.BLOCK_1618397679204_index_menu"
	SensorColorGetRGBGrayLight           BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_rgb_gray_light"
	SensorColorGetRGBGrayLightIndex      BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_rgb_gray_light_index_menu"
	SensorColorGetRGBGrayLightInput2     BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_rgb_gray_light_inputMenu_2_menu"
	SensorColorGetRGBGrayLightInput3     BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_rgb_gray_light_inputMenu_3_menu"
	SensorColorGetOffTrack               BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_off_track"
	SensorColorGetOffTrackIndex          BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_off_track_index_menu"
	SensorColorIsStatusL1R1              BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_sta_with_inputMenu"
	SensorColorIsStatusL1R1Index         BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_sta_with_inputMenu_index_menu"
	SensorColorIsStatusL1R1Input         BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_get_sta_with_inputMenu_inputMenu_2_menu"
	SensorColorIsStatus                  BlockType = "mbuild_quad_color_sensor.BLOCK_1618364921511"
	SensorColorIsStatusIndex             BlockType = "mbuild_quad_color_sensor.BLOCK_1618364921511_index_menu"
	SensorColorIsStatusInput             BlockType = "mbuild_quad_color_sensor.BLOCK_1618364921511_inputMenu_2_menu"
	SensorColorIsLineAndBackground       BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_is_line_and_background"
	SensorColorIsLineAndBackgroundIndex  BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_is_line_and_background_index_menu"
	SensorColorIsLineAndBackgroundInput2 BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_is_line_and_background_inputMenu_2_menu"
	SensorColorIsLineAndBackgroundInput3 BlockType = "mbuild_quad_color_sensor.mbuild_quad_color_sensor_is_line_and_background_inputMenu_3_menu"

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

	SensorsResetAxisRotationAngle BlockType = "cyberpi.cyberpi_reset_axis_rotation_angle"
	SensorsResetYaw               BlockType = "cyberpi.cyberpi_reset_yaw"
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

// Ambient light (ultrasonic sensor)
const (
	UltrasonicSetBrightness      BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_set_bri"
	UltrasonicSetBrightnessIndex BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_set_bri_index_menu"
	UltrasonicSetBrightnessOrder BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_set_bri_order_menu"
	UltrasonicAddBrightness      BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_add_bri"
	UltrasonicAddBrightnessIndex BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_add_bri_index_menu"
	UltrasonicAddBrightnessOrder BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_add_bri_order_menu"
	UltrasonicGetBrightness      BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_get_bri"
	UltrasonicGetBrightnessIndex BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_get_bri_index_menu"
	UltrasonicGetBrightnessOrder BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_get_bri_order_menu"
	UltrasonicOffLED             BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_off_led"
	UltrasonicOffLEDIndex        BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_off_led_index_menu"
	UltrasonicOffLEDInput        BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_off_led_inputMenu_3_menu"
	UltrasonicShowEmotion        BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_show_emotion"
	UltrasonicShowEmotionIndex   BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_show_emotion_index_menu"
	UltrasonicShowEmotionMenu    BlockType = "cyberpi_mbuild_ultrasonic2.mbuild_ultrasonic2_show_emotion_emotion_menu"
)

// Bottom light (RGB sensor)
const (
	SensorColorSetFillColor          BlockType = "mbuild_quad_color_sensor.BLOCK_1618382823173"
	SensorColorSetFillColorIndex     BlockType = "mbuild_quad_color_sensor.BLOCK_1618382823173_index_menu"
	SensorColorDisableFillColor      BlockType = "mbuild_quad_color_sensor.BLOCK_1620904215289"
	SensorColorDisableFillColorIndex BlockType = "mbuild_quad_color_sensor.BLOCK_1620904215289_index_menu"
)

// Net
const (
	NetSetWifiBroadcast          BlockType = "cyberpi.cyberpi_set_wifi_broadcast"
	NetSetWifiBroadcastWithValue BlockType = "cyberpi.cyberpi_set_wifi_broadcast_with_value"
	NetSetWifiChannel            BlockType = "cyberpi.cyberpi_set_wifi_channels"
	NetConnectWifi               BlockType = "cyberpi.cyberpi_wifi_set"
	NetWifiIsConnected           BlockType = "cyberpi.cyberpi_wifi_is_connect"
	NetWifiReconnect             BlockType = "cyberpi.cyberpi_wifi_reconnect"
	NetWifiDisconnect            BlockType = "cyberpi.cyberpi_wifi_disconnect"
	NetWifiGetValue              BlockType = "cyberpi.cyberpi_wifi_broadcast_get_value"
)

// Motors
const (
	Mbot2MoveDirectionWithRPM                      BlockType = "mbot2.mbot2_move_direction_with_rpm"
	Mbot2MoveDirectionWithTime                     BlockType = "mbot2.mbot2_move_direction_with_time"
	Mbot2MoveMoveWithCmAndInch                     BlockType = "mbot2.mbot2_move_straight_with_cm_and_inch"
	Mbot2CwAndCcwWithAngle                         BlockType = "mbot2.mbot2_cw_and_ccw_with_angle"
	Mbot2EncoderMotorSet                           BlockType = "mbot2.mbot2_encoder_motor_set"
	Mbot2EncoderMotorSetMenu                       BlockType = "mbot2.mbot2_encoder_motor_set_inputMenu_1_menu"
	Mbot2EncoderMotorSetWithTime                   BlockType = "mbot2.mbot2_encoder_motor_set_with_time"
	Mbot2EncoderMotorSetWithTimeMenu               BlockType = "mbot2.mbot2_encoder_motor_set_with_time_fieldMenu_1_menu"
	Mbot2EncoderMotorStop                          BlockType = "mbot2.mbot2_encoder_motor_stop"
	Mbot2EncoderMotorResetAngle                    BlockType = "mbot2.mbot2_encoder_motor_reset_angle"
	Mbot2EncoderMotorResetAngleMenu                BlockType = "mbot2.mbot2_encoder_motor_reset_angle_inputMenu_1_menu"
	Mbot2EncoderMotorLockUnlock                    BlockType = "mbot2.mbot2_encoder_motor_lock_and_unlock"
	Mbot2EncoderMotorLockUnlockMenu                BlockType = "mbot2.mbot2_encoder_motor_lock_and_unlock_inputMenu_1_menu"
	Mbot2EncoderMotorSetWithTimeAngleAndCircle     BlockType = "mbot2.mbot2_encoder_motor_set_with_time_angle_and_circle"
	Mbot2EncoderMotorSetWithTimeAngleAndCircleMenu BlockType = "mbot2.mbot2_encoder_motor_set_with_time_angle_and_circle_fieldMenu_1_menu"
	Mbot2EncoderMotorGetSpeed                      BlockType = "mbot2.mbot2_encoder_motor_get_speed"
	Mbot2EncoderMotorGetSpeedMenu                  BlockType = "mbot2.mbot2_encoder_motor_get_speed_inputMenu_2_menu"
	Mbot2EncoderMotorGetAngle                      BlockType = "mbot2.mbot2_encoder_motor_get_angle"
	Mbot2EncoderMotorGetAngleMenu                  BlockType = "mbot2.mbot2_encoder_motor_get_angle_inputMenu_1_menu"
	Mbot2EncoderMotorDriveSpeed                    BlockType = "mbot2.mbot2_encoder_motor_drive_speed2"
	Mbot2EncoderMotorDrivePower                    BlockType = "mbot2.mbot2_encoder_motor_drive_power"
)

// MBot2
const (
	Mbot2SetParameters BlockType = "mbot2.mbot2_set_para"
	Mbot2TimerReset    BlockType = "cyberpi.cyberpi_timer_reset"
	Mbot2TimerGet      BlockType = "cyberpi.cyberpi_timer_get"
	Mbot2Hostname      BlockType = "cyberpi.cyberpi_name"
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

// Display+
const (
	SpriteSetBackgroundFillColor     BlockType = "cyberpi_sprite.cyberpi_sprite_background_fill_with_color"
	SpriteSetBackgroundFillColorRGB  BlockType = "cyberpi_sprite.cyberpi_sprite_background_fill_with_rgb"
	SpriteDrawPixelWithMatrix16      BlockType = "cyberpi_sprite.cyberpi_sprite_draw_pixel_with_matrix16"
	SpriteDrawPixelWithIcon          BlockType = "cyberpi_sprite.cyberpi_sprite_draw_pixel_with_icon"
	SpriteDrawPixelWithIconInputMenu BlockType = "cyberpi_sprite.cyberpi_sprite_draw_pixel_with_icon_inputMenu_2_menu"
	SpriteDrawText                   BlockType = "cyberpi_sprite.cyberpi_sprite_draw_text"
	SpriteDrawQR                     BlockType = "cyberpi_sprite.cyberpi_sprite_draw_QR"
	SpriteMirrorWithAxis             BlockType = "cyberpi_sprite.cyberpi_mirror_with_axis"
	SpriteDelete                     BlockType = "cyberpi_sprite.cyberpi_sprite_delete"
	SpriteSetAlign                   BlockType = "cyberpi_sprite.cyberpi_sprite_set_align"
	SpriteSetAlignInputMenu          BlockType = "cyberpi_sprite.cyberpi_sprite_set_align_inputMenu_2_menu"
	SpriteMoveXY                     BlockType = "cyberpi_sprite.cyberpi_sprite_move_x_and_y"
	SpriteMoveTo                     BlockType = "cyberpi_sprite.cyberpi_sprite_move_to"
	SpriteMoveRandom                 BlockType = "cyberpi_sprite.cyberpi_sprite_move_random"
	SpriteRotate                     BlockType = "cyberpi_sprite.cyberpi_sprite_rotate"
	SpriteRotateTo                   BlockType = "cyberpi_sprite.cyberpi_sprite_rotate_to"
	SpriteSetSize                    BlockType = "cyberpi_sprite.cyberpi_sprite_set_size"
	SpriteSetColorWithColor          BlockType = "cyberpi_sprite.cyberpi_sprite_set_color_with_color"
	SpriteSetColorWithRGB            BlockType = "cyberpi_sprite.cyberpi_sprite_set_color_with_rgb"
	SpriteCloseColor                 BlockType = "cyberpi_sprite.cyberpi_sprite_close_color"
	SpriteShowAndHide                BlockType = "cyberpi_sprite.cyberpi_sprite_show_and_hide"
	SpriteZMinMax                    BlockType = "cyberpi_sprite.cyberpi_sprite_z_max_and_min"
	SpriteZUpDown                    BlockType = "cyberpi_sprite.cyberpi_sprite_z_up_and_down"
	SpriteScreenRender               BlockType = "cyberpi_sprite.cyberpi_screen_render"
	SpriteIsTouchOtherSprite         BlockType = "cyberpi_sprite.cyberpi_sprite_is_touch_other_sprite"
	SpriteIsTouchEdge                BlockType = "cyberpi_sprite.cyberpi_sprite_is_touch_edge"
	SpriteGetColorEqualWithRGB       BlockType = "cyberpi_sprite.cyberpi_screen_get_color_equal_with_rgb"
	SpriteGetXYRotationSizeAlign     BlockType = "cyberpi_sprite.cyberpi_sprite_get_x_y_rotation_size_align"
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

// Lists
const (
	ListAdd       BlockType = "data_addtolist"
	ListDelete    BlockType = "data_deleteoflist"
	ListClear     BlockType = "data_deletealloflist"
	ListInsert    BlockType = "data_insertatlist"
	ListReplace   BlockType = "data_replaceitemoflist"
	ListItem      BlockType = "data_itemoflist"
	ListItemIndex BlockType = "data_itemnumoflist"
	ListLength    BlockType = "data_lengthoflist"
	ListContains  BlockType = "data_listcontainsitem"
)

// Variables
const (
	VariableSetTo    BlockType = "data_setvariableto"
	VariableChangeBy BlockType = "data_changevariableby"
)

// Custom blocks
const (
	ProceduresDefinition         BlockType = "procedures_definition"
	ProceduresPrototype          BlockType = "procedures_prototype"
	ArgumentReporterStringNumber BlockType = "argument_reporter_string_number"
	ArgumentReporterBoolean      BlockType = "argument_reporter_boolean"
	ProceduresCall               BlockType = "procedures_call"
)
