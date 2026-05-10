# ANSI Preservation

`ntagcolor` preserves active ANSI SGR styles by default. When it renders a tag,
it still emits its own tag colors and then resets them, but it follows that
reset with the SGR style that was active in the input stream.

This keeps already-colored pipeline input readable:

```sh
fd --color always | ntagcolor
```

To disable preservation and use the legacy full-reset behavior:

```sh
ntagcolor --no-preserve-ansi
ntagcolor --preserve-ansi=false
```

## Supported SGR State

The tracker handles standard semicolon-form SGR sequences:

- full reset: `0`
- attributes: bold, dim, italic, underline, blink, reverse, hidden, strike
- partial resets such as `22`, `23`, `24`, `25`, `27`, `28`, `29`, `39`, and `49`
- 8/16-color foreground and background codes
- 256-color forms such as `38;5;244` and `48;5;236`
- truecolor forms such as `38;2;120;130;140` and `48;2;20;30;40`

Unsupported or malformed escape sequences are passed through unchanged and do
not update the tracked style.
