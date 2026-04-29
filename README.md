# Pokedex Fetch Go

*v1.0.0*

Single-binary Go CLI that prints random Pokemon ANSI sprites in the terminal, tracks catches in a persistent Pokedex, and integrates with fastfetch.

## Features

- Random encounter on `pokedex-fetch-go` or `pokedex-fetch-go catch`
- Show a specific Pokemon by name or National Dex number
- Regular/shiny and small/large sprite variants
- Gen 1-8 filtering
- Persistent JSON Pokedex and trainer state under `~/.config/pokedex-fetch-go`
- XP, levels, streaks, daily catch count, and 50 achievements worth 1000 total points
- Interactive Pokedex, trainer, and Hall of Fame screens
- fastfetch-compatible raw sprite output

## Installation

### From GitHub Releases

Download the archive for your OS/CPU from the [latest release](https://github.com/brunoorsolon/pokedex-fetch-go/releases/latest).

Release assets are named by platform, for example:

```text
pokedex-fetch-go_1.0.0_linux_amd64.tar.gz
pokedex-fetch-go_1.0.0_linux_arm64.tar.gz
pokedex-fetch-go_1.0.0_darwin_amd64.tar.gz
pokedex-fetch-go_1.0.0_darwin_arm64.tar.gz
pokedex-fetch-go_1.0.0_windows_amd64.zip
checksums.txt
```

The executable inside the archive is named `pokedex-fetch-go` (`pokedex-fetch-go.exe` on Windows).

There are no distro packages yet, so install the binary somewhere in your `PATH`. A good user-local location is `~/.local/bin`:

```bash
mkdir -p ~/.local/bin
```

```bash
# Linux/macOS example. Replace the archive name with the one you downloaded.
tar -xzf pokedex-fetch-go_1.0.0_linux_amd64.tar.gz
cp pokedex-fetch-go_1.0.0_linux_amd64/pokedex-fetch-go ~/.local/bin/pokedex-fetch-go
```

```bash
chmod +x ~/.local/bin/pokedex-fetch-go
```

Make sure `~/.local/bin` is in your shell `PATH`. For zsh, add something like this to `~/.zshrc` if you do not already have it:

```zsh
typeset -U path
[[ -d "$HOME/.local/bin" ]] && path=("$HOME/.local/bin" $path)
```

Then open a new terminal, or run:

```bash
source ~/.zshrc
```

Verify the install:

```bash
which pokedex-fetch-go
pokedex-fetch-go --help
```

Optional checksum verification:

```bash
sha256sum -c checksums.txt
```

On Windows, extract the `.zip` file and place `pokedex-fetch-go.exe` somewhere in your `PATH`.

### Build From Source

Requires Go 1.22+.

```bash
git clone https://github.com/brunoorsolon/pokedex-fetch-go.git
cd pokedex-fetch-go
go build -o bin/pokedex-fetch-go ./cmd/pokedex-fetch-go
```

Run locally:

```bash
./bin/pokedex-fetch-go --help
```

Install to `~/.local/bin`:

```bash
mkdir -p ~/.local/bin
cp ./bin/pokedex-fetch-go ~/.local/bin/pokedex-fetch-go
```

Or install with Go:

```bash
go install ./cmd/pokedex-fetch-go
```

If you use `go install`, make sure your Go bin directory is in `PATH`, usually:

```zsh
[[ -d "$HOME/go/bin" ]] && path=("$HOME/go/bin" $path)
```

## Usage

```bash
pokedex-fetch-go
pokedex-fetch-go catch --gen 1 --raw
pokedex-fetch-go show --name pikachu --shiny
pokedex-fetch-go show --number 25 --big
pokedex-fetch-go list --gen 1-3
pokedex-fetch-go pokedex 3
pokedex-fetch-go trainer --name Ash --region Kanto --hometown "Pallet Town"
pokedex-fetch-go achievements
pokedex-fetch-go config init
pokedex-fetch-go setup bash
```

## Achievements

The Hall of Fame includes 50 achievements across Pokedex progress, catch milestones, trainer rank, and special collections. Easier goals award fewer points and rare/endgame goals award more, for a total score of 1000 achievement points.

```bash
pokedex-fetch-go achievements
```

## fastfetch

```bash
pokedex-fetch-go catch --raw | fastfetch --logo-type file-raw --logo -
```

## Configuration

Run:

```bash
pokedex-fetch-go config init
```

The app also bootstraps this file automatically on first run if it does not exist.

This creates `~/.config/pokedex-fetch-go/config.toml`. See `config.example.toml` for all options.

## Data

Sprites come from `pokemon-colorscripts` and are embedded as gzip-compressed ANSI files.

## Credits

- Pokemon designs, names, and branding: The Pokemon Company.
- Heavily inspired by the [poketerm](https://github.com/chris-wood-mo/poketerm) project. It's basically the same idea but rewritten in Go and using different architecture under the hood.
- Sprite source: [pokemon-colorscripts](https://gitlab.com/phoneybadger/pokemon-colorscripts) (MIT).
