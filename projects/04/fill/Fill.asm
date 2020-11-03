// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// screen limit: 512 * 256 / 16 - 1
@8191
D=A
@R3
M=D

(LOOP)
@KBD
D=M

@BLACK
D;JGT

@WHITE
0;JMP

@LOOP
0;JMP
(END)

(BLACK)
@R4 // counter
M=0
    (BLACKLOOP)
    @R3
    D=M
    @R4
    D=D-M
    @BLACKEND
    D;JLE

    // fill
    @R4
    D=M
    @SCREEN
    A=D+A
    D=0
    M=!D

    // increment
    @R4
    M=M+1

    @BLACKLOOP
    0;JMP
    (BLACKEND)
@LOOP
0;JMP

(WHITE)
@R4 // counter
M=0
    (WHITELOOP)
    @R3
    D=M
    @R4
    D=D-M
    @WHITEEND
    D;JLE

    // fill
    @R4
    D=M
    @SCREEN
    A=D+A
    M=0

    // increment
    @R4
    M=M+1

    @WHITELOOP
    0;JMP
    (WHITEEND)
@LOOP
0;JMP
