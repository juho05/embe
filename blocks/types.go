package blocks

type BlockType string

// Events
const (
	WhenLaunch                     BlockType = "cyberpi.cyberpi_when_launch"
	WhenButtonPress                BlockType = "cyberpi.cyberpi_when_button_press"
	WhenDirectionKeyPress          BlockType = "cyberpi.cyberpi_when_direction_key_press"
	WhenDetectAttitude             BlockType = "cyberpi.cyberpi_when_detect_attitude"
	WhenDetectAction               BlockType = "cyberpi.cyberpi_when_detect_action"
	WhenSensorValueBiggerOrSmaller BlockType = "cyberpi.cyberpi_when_sensor_value_bigger_or_smaller_than"
)

// Audio
const (
	GetVolume BlockType = "cyberpi.cyberpi_get_volume"
	SetVolume BlockType = "cyberpi.cyberpi_set_volume"
	AddVolume BlockType = "cyberpi.cyberpi_add_volume"

	StopAudio                     BlockType = "cyberpi.cyberpi_stop_audio"
	PlayBuzzerTone                BlockType = "cyberpi.cyberpi_play_buzzer_tone"
	PlayBuzzerToneWithTime        BlockType = "cyberpi.cyberpi_play_buzzer_tone_with_time"
	PlayClip                      BlockType = "cyberpi.cyberpi_play_audio_3"
	PlayClipFileNameMenu          BlockType = "cyberpi.cyberpi_play_audio_3_file_name_menu"
	PlayClipUntilDone             BlockType = "cyberpi.cyberpi_play_audio_until_3"
	PlayClipUntilDoneFileNameMenu BlockType = "cyberpi.cyberpi_play_audio_until_3_file_name_menu"
	PlayNote                      BlockType = "cyberpi.cyberpi_play_music_with_tone_and_note_2"
	Note                          BlockType = "note"
	PlayMusicInstrument           BlockType = "cyberpi.cyberpi_play_music_with_note"
	PlayMusicInstrumentMenu       BlockType = "cyberpi.cyberpi_play_music_with_note_fieldMenu_1_menu"

	RecordStart         BlockType = "cyberpi.cyberpi_start_record"
	RecordStop          BlockType = "cyberpi.cyberpi_stop_record"
	PlayRecord          BlockType = "cyberpi.cyberpi_play_record"
	PlayRecordUntilDone BlockType = "cyberpi.cyberpi_play_record_until"
)

// Control
const (
	If            BlockType = "control_if"
	IfElse        BlockType = "control_if_else"
	Wait          BlockType = "control_wait"
	WaitUntil     BlockType = "control_wait_until"
	Repeat        BlockType = "control_repeat"
	RepeatUntil   BlockType = "control_repeat_until"
	RepeatForever BlockType = "control_forever"
	Stop          BlockType = "control_stop"
	Restart       BlockType = "cyberpi.cyberpi_restart"
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
)
