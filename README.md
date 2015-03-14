# facade

This library provides a git-like sub command feature.

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

facade searches for sub command in the order below and execute it if found.

1. If you run your command as `your-command foo bar baz`, facade regards `foo` as sub command and the rest of the arguments as ones for the sub command
2. Then execute `your-command-foo` with arguments `bar baz`.

## Logging

facade takes over STDOUT and STDERR of sub command and emits it in some pretty manner.

## License

MIT
