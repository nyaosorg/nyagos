Console I/O Library for Windows
===============================

func GetKey()(rune,uint16)
--------------------------

Read One-letter from keyboard.
Returns the pair of unicode-codepoint and scancode.


func ReadLine() string
----------------------

Read 1-Line with tiny Emacs-like keybind.

Now Support
    Left      Ctrl-F
    Right     Ctrl-B 
    Home      Ctrl-A
    End       Ctrl-E
    BackSpace Ctrl-H
    Delete    Ctrl-D
    Enter     Ctrl-M

Sample
------

    package main
    import "../conio"
    import "fmt"

    func main(){
        fmt.Print("conio.ReadLine>")
        result := conio.ReadLine()
        fmt.Println("Result=" + result)
    }
