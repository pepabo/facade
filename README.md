# façade

This library adds a git-like sub command feature into your command with just one line.

## Usage

Just call `facade.Run()` in your main file like below:

```go
package main

import (
	"github.com/kentaro/facade"
)

func main() {
	facade.Run()
}
```

## Sub Command

1. If you name your command `your-command` and run it with `your-command foo bar baz`, façade regards `foo` as sub command and the rest of the arguments as ones for the sub command.
2. Then execute `your-command-foo` with arguments `bar baz`.

## Logging

façade takes over STDOUT and STDERR of sub command and emits outputs from them in some pretty manner.

## License

MIT
