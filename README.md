# nand2tetris

- [Home \| nand2tetris](https://www.nand2tetris.org/)
- [O'Reilly Japan \- コンピュータシステムの理論と実装](https://www.oreilly.co.jp/books/9784873117126/)

Results of 12 projects of nand2tetris.
Implemented assembler, vm translator and jack compiler in Go.

## How To Run

See the guides in each chapter of the book.
For example, the final product `Pong` game can be launched by the following procedure.

```
sh ./tools/VMEmulator.sh

# in VM Emulator
- Load ./projects/11/Pong directory
- Set `Animate: no animation`
- Press `Run` Button
```

## What I haven't done yet

- Optimization
    - In particular, the project 12 are implemented with non-optimized (easy to write) algorithms.
        - `Memory.alloc()` does not support `freeList`
        - `Math.XXX` remains a naive implementation
    - As a result, the `Pong` game is slow
- Writing test code
    - I haven't written any test code except where I thought it would be more efficient to have tests at the time of implementation
- Error handling
- Project 13's Additional challenges
