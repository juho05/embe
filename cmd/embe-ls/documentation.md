// events
@launch
This event is triggered when the script is loaded.
---
@button
This event is triggered when the specified button is pressed.

Buttons: `a`, `b`
---
@joystick
This event is triggered when the joystick is pulled in the specified direction.

Directions: `left`, `right`, `up`, `down`, `middle`
---
@tilt
This event is triggered when the robot is tilted in the specified direction.

Directions: `left`, `right`, `forward`, `backward`
---
@face
This event is triggered when the screen of the robot is facing in the specified direction.

Directions: `up`, `down`
---
@wave
This event is triggered when a waving motion in the specified direction is detected.

Directions: `left`, `right`
---
@rotate
This event is triggered when the robot is rotating in the specified direction.

Directions: `clockwise`, `anticlockwise`
---
@fall
This event is triggered when the robot is falling.
---
@shake
This event is triggered when the robot is shaken.
---
@light
This event is triggered when the brightness of the environment fulfills the specified condition.

Example conditions: `>50`, `<3`
---
@sound
This event is triggered when the loudness of the environment fulfills the specified condition.

Example conditions: `>50`, `<3`
---
@shakeval
This event is triggered when the strength with which the robot is shaken fulfills the specified condition.

Example conditions: `>50`, `<3`
---
@timer
This event is triggered when the value of the timer fulfills the specified condition.

Example conditions: `>50`, `<3`
---
@receive
This event is triggered when the specified message is received over LAN.
---
// variables
audio.volume
The volume audio should be played at.

Values: 0-100
---
audio.speed
The speed audio should be played at.

Default: 100
---
lights.back.brightness
The brightness of the LED lights at the back of the robot.

Values: 0-100
---
time.timer
The timer starts to count from 0 each time the robot is powered on or when the timer is reset.
---
mbot.battery
The current battery level of the robot.

Values: 0-100
---
mbot.mac
The MAC address of the robot.
---
mbot.hostname
The hostname of the robot.
---
sensors.wavingAngle
The angle of a waving motion.
---
sensors.wavingSpeed
The speed of a waving motion.
---
shakingStrength
The strength with which the robot is shaken.
---
sensors.brightness
The current brightness of the environment.
---
sensors.loudness
The current loudness of the environment.
---
sensors.distance
The distance to the nearest object in front of the robot.
---
sensors.outOfRange
Whether the area in front of the robot is completely empty.
---
sensors.lineDeviation
The line deviation of the robot.

-100: far left,
0: center,
100: far right
---
net.connected
Whether the robot is connected to the internet.
---
math.e
Euler's number = 2.718281
---
math.pi
π = 3.141592
---
math.phi
Golden ratio = 1.618033
---
// functions
audio.stop
Stop all audio.
---
audio.playBuzzer
Play a tone with the specified frequency.

`duration`: seconds
---
audio.playClip
Play an audio clip.

Clips: `hi`, `bye`, `yeah`, `wow`, `laugh`, `hum`, `sad`, `sigh`, `annoyed`, `angry`, `surprised`, `yummy`, `curious`, `embarrassed`, `ready`, `sprint`, `sleepy`, `meow`, `start`, `switch`, `beeps`, `buzzing`, `jump`, `level-up`, `low-energy`, `prompt`, `right`, `wrong`, `ring`, `score`, `wake`, `warning`, `metal-clash`, `glass-clink`, `inflator`, `running-water`, `clockwork`, `click`, `current`, `wood-hit`, `iron`, `drop`, `bubble`, `wave`, `magic`, `spitfire`, `heartbeat`
---
audio.playInstrument
Play an instrument for a specified number of beats.

Instruments: `snare`, `bass-drum`, `side-stick`, `crash-cymbal`, `open-hi-hat`, `closed-hi-hat`, `tambourine`, `hand-clap`, `claves`
---
audio.playNote
Play a single note for a specified number of beats.

Note names: `c`, `c#`, `db`, `d`, `d#`, …
---
audio.record.start
Start recording.
---
audio.record.stop
Stop recording.
---
audio.record.play
Play the recording.
---
lights.deactivate
Deactivate all lights of the robot.
---
lights.back.playAnimation
Play a light animation.

Animations: `rainbow`, `spindrift`, `meteor_blue`, `meteor_green`, `flash_red`, `flash_orange`, `firefly`
---
lights.front.setBrightness
Set the brightness of the front lights.

Values: 0-100
---
lights.front.addBrightness
Increase the brightness of the front lights.
---
lights.front.displayEmotion
Display an emotion.

Emotions: `sleepy`, `wink`, `happy`, `dizzy`, `thinking`
---
lights.front.deactivate
Deactivate the front lights.
---
lights.bottom.deactivate
Deactivate the color sensor fill lights.
---
lights.bottom.setColor
Set the color sensor fill color.
---
lights.back.display
Set the color of the LED strip at the back of the robot.
---
lights.back.displayColor
Set the color of an LED at the back of the robot.
---
lights.back.displayColorFor
Display a color with an LED at the back of the robot for a specified number of seconds.
---
lights.back.deactivate
Deactivate the lights at the back of the robot.
---
lights.back.move
Moves every color `n` LEDs to the right.
---
display.print
Print a message on the display.
---
display.println
Print a message followed by a new line on the display.
---
display.setFontSize
Set the font size for printing.
---
display.setColor
Set the color for printing.
---
display.showLabel
Show a label on the display at a specified location.

Locations: `top_left`, `top_mid`, `top_right`, `mid_left`, `center`, `mid_right`, `bottom_left`, `bottom_mid`, `bottom_right`
---
display.lineChart.addData
Add a data point to the line chart.
---
display.lineChart.setInterval
Set the interval of the line chart.
---
display.barChart.addData
Add a data point to the bar chart.
---
display.table.addData
Add an entry to the table.
---
display.setOrientation
Set the orientation of the display.

Orientations: `-90`, `0`, `90`, `180`
---
display.clear
Clear the display.
---
display.setBackgroundColor
Set the background color of the screen.
---
display.render
Refresh the screen and show sprite changes.
---
sprite.fromIcon
Set the sprite to the icon.

Icons: `Music`, `Image`, `Video`, `Clock`, `Play`, `Pause`, `Next`, `Prev`, `Sound`, `Temperature`, `Light`, `Motion`, `Home`, `Gear`, `List`, `Right`, `Wrong`, `Shut_down`, `Refresh`, `Trash_can`, `Download`, `Cloudy`, `Rain`, `Snow`, `Train`, `Rocket`, `Truck`, `Car`, `Droplet`, `Distance`, `Fire`, `Magnetic`, `Gas`, `Vision`, `Color`, `Overcast`, `Sandstorm`, `Foggy`
---
sprite.fromText
Set the sprite to the text.
---
sprite.fromQR
Set the sprite to a QR code pointing to `url`.
---
sprite.flipH
Flip the sprite horizontally.
---
sprite.flipV
Flip the sprite vertically.
---
sprite.delete
Delete the sprite.
---
sprite.setAnchor
Set the anchor of the sprite.

Anchors: `top_left`, `top_mid`, `top_right`, `mid_left`, `center`, `mid_right`, `bottom_left`, `bottom_mid`, `bottom_right`
---
sprite.moveLeft
Move the sprite to the left.
---
sprite.moveRight
Move the sprite to the right.
---
sprite.moveUp
Move the sprite up.
---
sprite.moveDown
Move the sprite up.
---
sprite.moveTo
Move the sprite to the coordinates.
---
sprite.moveRandom
Move the sprite to a random position on the screen.
---
sprite.rotate
Rotate the sprite by `angle`.

`angle`: degrees
---
sprite.rotateTo
Set the rotation of the sprite to `angle`.

`angle`: degrees
---
sprite.setScale
Set the scale of the sprite.

Default: 100
---
sprite.setColor
Set the tint of the sprite.

Default: #ffffff
---
sprite.resetColor
Reset the tint of the sprite to the default.
---
sprite.show
Make the sprite visible.
---
sprite.hide
Make the sprite invisible.
---
sprite.toFront
Move the sprite to highest layer.
---
sprite.toBack
Move the sprite to the lowest layer.
---
sprite.layerUp
Move the sprite one layer up.
---
sprite.layerDown
Move the sprite one layer down.
---
net.broadcast
Broadcast a message over LAN.
---
net.setChannel
The the channel to use for LAN messages.

Default: 6
---
net.connect
Connect to WIFI.
---
net.reconnect
Reconnect to a previously connected WIFI.
---
net.disconnect
Disconnect from WIFI.
---
sensors.resetAngle
Reset the angle sensor.
---
sensors.resetYawAngle
Reset the yaw angle sensor.
---
sensors.defineColor
Define a custom color to use for the color sensor.
---
sensors.calibrateColors
Calibrate the color sensor.
---
sensors.enhancedColorDetection
Enable or disable the enhanced color detection algorithm.

Default: off
---
motors.run
Run the motors with a specific RPM.

`duration`: seconds
---
motors.runBackward
Run the motors backwards with a specific RPM.

`duration`: seconds
---
motors.moveDistance
Move the robot a specific distance.

`distance`: centimeter
---
motors.moveDistanceBackward
Move the robot a specific distance backwards.

`distance`: centimeter
---
motors.turnLeft
Turn the robot to the left.

`angle`: degrees
---
motors.turnRight
Turn the robot to the right.

`angle`: degrees
---
motors.rotateRPM
Run a single motor with a specific RPM.
---
motors.rotatePower
Run a single motor with a specific power.

Values: 0-100
---
motors.rotateAngle
Turn a single motor by an angle.

`angle`: degrees
---
motors.driveRPM
Set the speed of the two encoder motors.
---
motors.drivePower
Set the power of the two encoder motors.

Values: 0-100
---
motors.stop
Stop a motor.
---
motors.resetAngle
Reset the angle sensor of a motor.
---
motors.lock
Lock a motor.
---
motors.unlock
Unlock a motor.
---
time.wait
Wait a specific number of seconds or until a condition is met.
---
time.resetTimer
Reset the timer.
---
mbot.restart
Restart the robot.
---
mbot.resetParameters
Reset all chassis parameters of the robot.
---
mbot.calibrateParameters
Calibrate all chassis parameters of the robot.
---
script.stop
Stop this script.
---
script.stopAll
Stop all scripts.
---
script.stopOther
Stop all other scripts.
---
lists.append
Append an item to the list.
---
lists.remove
Remove an item from the list.
---
lists.clear
Remove all items from the list.
---
lists.insert
Insert the item at `index`.
---
lists.replace
Replace the item at `index` with `value`.
---
// expression functions
mbot.isButtonPressed
Whether the specified button is pressed.

Buttons: `a`, `b`
---
mbot.buttonPressCount
The amount of times the specified button was pressed.

Buttons: `a`, `b`
---
mbot.isJoystickPulled
Whether the joystick is pulled in the specified direction.

Directions: `left`, `right`, `up`, `down`, `middle`
---
mbot.joystickPullCount
The amount of times the joystick was pulled in the specified direction.

Directions: `left`, `right`, `up`, `down`, `middle`
---
lights.front.brightness
The brightness of a front light.
---
sensors.isTilted
Whether the robot is tilted in the specified direction.

Directions: `forward`, `backward`, `left`, `right`
---
sensors.isFaceUp
Whether the display of the robot is facing in the specified direction.

Directions: `up`, `down`
---
sensors.isWaving
Whether a waving motion in the specified direction is detected.

Directions: `up`, `down`, `left`, `right`
---
sensors.isRotating
Whether the robot is rotated in the specified direction.

Directions: `clockwise`, `anticlockwise`
---
sensors.isFalling
Whether the robot is falling.
---
sensors.isShaking
Whether the robot is shaking.
---
sensors.tiltAngle
The angle in degrees the robot is tilted in the specified direction.

Directions: `forward`, `backward`, `left`, `right`
---
sensors.rotationAngle
The angle in degrees the robot is rotated in the specified direction.

Directions: `clockwise`, `anticlockwise`
---
sensors.acceleration
The amount of acceleration that is measured on the specified axis.

Axes: `x`, `y`, `z`
---
sensors.rotation
The angle in degrees the robot is rotated on the specified axis.

Axes: `x`, `y`, `z`
---
sensors.angleSpeed
The speed at which the robot is rotated on the specified axis.

Axes: `x`, `y`, `z`
---
sensors.colorStatus
The status of the color sensor.

Targets: `line`, `ground`, `white`, `red`, `yellow`, `green`, `cyan`, `blue`, `purple`, `black`, `custom`

Result:
  - when inner == true: 0b00 - 0b11
  - when inner == false: 0b0000 - 0b1111
---
sensors.getColorValue
The color value measured by the color sensor.

Types: `red`, `green`, `blue`, `gray`, `light`
---
sensors.getColorName
The name of the color currently detected by the color sensor.

Possible values: `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `purple`, `black`
---
sensors.isColorStatus
Whether the color sensor detects the specified status.

Targets: `line`, `ground`, `white`, `red`, `yellow`, `green`, `cyan`, `blue`, `purple`, `black`, `custom`

Status:
  - when inner == true: 0b00 - 0b11
  - when inner == false: 0b0000 - 0b1111
---
sensors.detectColor
Whether the color sensor detects the specified target.

Targets: `line`, `ground`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `purple`, `black`
---
motors.rpm
The current RPM of the motor.
---
motors.power
The current power of the motor.
---
motors.angle
The current angle of the motor.
---
net.receive
Receive a value over LAN.
---
math.round
Round the `n` to an integer.
---
math.random
Returns a random value between `from` (inclusive) and `to` (inclusive).
---
math.abs
Returns the absolute value of `n`.
---
math.floor
Returns `n` rounded down to the nearest integer.
---
math.ceil
Returns `n` rounded up to the nearest integer.
---
math.sqrt
Returns the square root of `n`.
---
math.sin
Returns the sine of `n` (degrees).
---
math.cos
Returns the cosine of `n` (degrees).
---
math.tan
Returns the tangent of `n` (degrees).
---
math.asin
Returns the arcus sine of `n` in degrees.
---
math.acos
Returns the arcus cosine of `n` in degrees.
---
math.atan
Returns the arcus tangent of `n` in degrees.
---
math.ln
Returns the natural logarithm of `n`.
---
math.log
Returns the logarithm of `n` with base 10.
---
math.ePowerOf
Returns e ^ `n`.
---
math.tenPowerOf
Returns 10 ^ `n`.
---
strings.length
Returns the length of the string.
---
strings.letter
Returns the letter at the index in the string.
---
strings.contains
Whether the string contains the substring.
---
lists.get
Get the item at `index` in the list.
---
lists.indexOf
Get the index of `value`.
---
lists.length
Get the number of items in the list.
---
lists.contains
Check whether the list contains `value`.
---
display.pixelIsColor
Check whether the color of the pixel matches `r`, `g`, `b`.
---
sprite.touchesSprite
Check whether the two sprites touch each other.
---
sprite.touchesEdge
Check whether the sprite is on the edge.
---
sprite.positionX
The x coordinate of the sprite.
---
sprite.positionY
The y coordinate of the sprite.
---
sprite.rotation
The rotation of the sprite.
---
sprite.scale
The scale of the sprite.
---
sprite.anchor
The anchor of the sprite.

Anchors: `top_left`, `top_mid`, `top_right`, `mid_left`, `center`, `mid_right`, `bottom_left`, `bottom_mid`, `bottom_right`
