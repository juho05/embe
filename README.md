# embe

Embe is a programming language that compiles to [mBlock](https://makeblock.com) Scratch code for the [mBot2](https://education.makeblock.com/mbot2/) robot.

**NOTE:** Due to some incompatibilities between the desktop and the web version of the mBlock IDE it is recommended to use the [web version](https://ide.mblock.cc/) to display and run output files of embe.

## [Documentation](docs/documentation.md)

## Installation

### Windows

[Download](https://github.com/juho05/embe/releases/latest/download/install.bat) and execute `install.bat`.

### macOS/Linux

Paste one of the following commands into a terminal window:

#### curl

```bash
curl -L https://raw.githubusercontent.com/juho05/embe/main/install.sh | bash
```

#### wget (in case curl is not installed)

```bash
wget -q --show-progress https://raw.githubusercontent.com/juho05/embe/main/install.sh -O- | bash
```

## Uninstallation

To remove embe from your system run:

```bash
embe uninstall
```

## Editor Support

- LSP: [embe-ls](https://github.com/juho05/embe/blob/main/cmd/embe-ls/README.md)
- VS Code: [vscode-embe](https://github.com/juho05/vscode-embe)
- Vim: [vim-embe](https://github.com/juho05/vim-embe)

## Building

### Prerequisites

- [Go](https://go.dev) 1.19+

```sh
git clone https://github.com/juho05/embe
cd embe
go build ./cmd/embe
go build ./cmd/embe-ls
```

## License

Copyright (c) 2022-2023 Julian Hofmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
