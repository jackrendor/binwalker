# binwalker

Binwalker solves the problem of extracting "hidden files" and such from another file and not having them "sorted properly" (see blog post: TODO) by their time of extraction.

That's literally all you need to know.


## Example
[![asciicast](https://asciinema.org/a/UuRj6R8SMbYHFq3p8DhvNKeqb.svg)](https://asciinema.org/a/UuRj6R8SMbYHFq3p8DhvNKeqb)


## Requirements
- Have `go 1.24` installed
- Have `binwalk3` installed and that's located in any of the directories in `$PATH`

## Installation
```bash
go install github.com/jackrendor/binwalker@latest
```

Check where it got installed, it could be in one of those two paths:
- `~/.local/bin`
- `~/go/bin/`

If those paths are not inside the `$PATH` variable, then you won't be able to run the command from any directory and you'd have to specify the full path of the binary.

So check your `$PATH` variable and make sure that everything is set up properly.