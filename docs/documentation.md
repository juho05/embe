# Documentation

## Events

Event | Parameter | Triggered when...
--- | --- | ---
`@launch` | *none* | the specified button is pressed.
`@joystick` | `left`, `right`, `up`, `down` | the joystick is pulled in the specified direction.
`@tilt` | `left`, `right`, `forward`, `backward` | the robot is tilted in the specified direction.
`@face` | `up`, `down` | the screen of the robot is facing in the specified direction.
`@wave` | `left`, `right` | a waving motion in the specified direction is detected.
`@rotate` | `clockwise`, `anticlockwise` | the robot is rotating in the specified direction.
`@shake` | *none* | the robot is shaken.
`@light` | `>50`, `<3.14`, … | the brightness of the environment fulfills the specified condition.
`@sound` | `>50`, `<3.14`, … | the loudness of the environment fulfills the specified condition.
`@shakeval` | `>50`, `<3.14`, … | the strength with which the robot is shaken fulfills the specified condition.
`@timer` | `>50`, `<3.14`, … | the value of the timer fulfills the specified condition.
`@receive` | message: string | the specified message is received over LAN.

## Namespaces

Name | Purpose
--- | ---
`audio` | play sounds
`display` | show text on the display of the CyberPi
`lights` | control the LED lights of the robot
`lists` | work with lists
`math` | math functions like `random`, `round`, `sin`, `abs`, `floor`, …
`mbot` | get data like the current battery level or whether a specific button is currently pressed
`motors` | control the motors of the robot
`net` | communicate with other robots
`script` | stop the current, all or all other scripts
`sensors` | get data from all the different sensors of the robot
`sprite` | manipulate sprites
`strings` | work with strings
`time` | wait and control timers
