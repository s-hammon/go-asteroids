# Asteroids!

A classic Asteroids-like game built with Go and the Ebiten game engine.

Play it [here](https://s-hammon.github.io/asteroids.html)!

## Installation

Check the latest release for Windows/Linux binary.

#### Go

```bash
go install github.com/s-hammon/go-asteroids@latest
```

## Controls

- **W**: Thrust forward
- **A**: Rotate left
- **D**: Rotate right
- **Space**: Shoot

## Todo

- Create sprites for player & asteroids
- Add weapon upgrades
  * randomly generate upgrade sprite in place of an asteroid
  * when player collides with one of these, temporarily replaces default weapon
  * add a gauge for duration of the weapon
- UI enhancements
- Add in-game debug logging

## References

- [Go](https://golang.org/)
- [Ebiten](https://ebiten.org/) - A dead simple 2D game library for Go
