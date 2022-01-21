package main

import (
  "flag"
  "fmt"
  "os"

  "github.com/midbel/textwrap"
)

func main() {
  limit := flag.Int("n", 72, "limit")
  flag.Parse()

  b, err := os.ReadFile(flag.Arg(0))
  if err != nil {
    os.Exit(1)
  }
  str := textwrap.WrapN(string(b), *limit)
  fmt.Println(str)
}
