# Prayer-Times CLI

I fancy myself as something of a command-line junkie. And as a muslim, prayer is important to me. If you know me, you know I don't want to mess with my browser or any GUI app, unless I must. You are also probably like me. So, I built prayertimes to quickly retrieve prayertimes for my city on the go. The tool does what it says, and it does it well and fast.

## Features

- Retrieves daily prayer times for the city you're in.

## Installation

### Using Homebrew (For MacOS & Linux)

```
brew install prayertimes
```

### Or just install it with `go` (If you have `go` installed)

```
go install github.com/lanrey-waju/prayer-times
```

### Build (requires Go 1.24+)

```
git clone https://github.com/lanrey-waju/prayer-times.git
cd prayer-times
go build
```

## Usage

```
prayertimes
```

Example Output:
![prayer times output](./assets/prayertimes.png)
