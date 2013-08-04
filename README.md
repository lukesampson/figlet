Figlet
======

This is a port of Figlet from C to the Go programming language.

Figlet is a program that makes large letters out of ordinary text.

```
 ,          _   _         __        __         _     _ _
/|   |     | | | |        \ \      / /__  _ __| | __| | |
 |___|  _  | | | |  __     \ \ /\ / / _ \| '__| |/ _` | |
 |   |\|/  |/  |/  /  \_    \ V  V / (_) | |  | | (_| |_|
 |   |/|__/|__/|__/\__/      \_/\_/ \___/|_|  |_|\__,_(_)

 ```

For information about the original FIGlet, see [figlet.org](http://www.figlet.org/).

### Usage
```
figlet [ -lcrhR ] [ -f fontfile ]
       [ -w outputwidth ] [ -m smushmode ]
       [ message ]
```

###### Options
`-h`
Shows help info: really just the usage info above plus the address of this page.

`-l, -c, -r`
These control the alignment of the output: left, center and right accordingly.

`-R`
Reverses the direction of text. So if the font specifies left-to-right, this will make it right-to-left, and vice versa.

`-f fontfile`
Specify a font to use. The fonts come from the "fonts" directory, in the same directory as the `figlet` program. You can see the available fonts with `figlet -list`.

`-w outputwidth`
FIGlet assumes an 80 character wide terminal. Use this to specify a different output width.

`-m smushmode`
Use a different "smush mode". Smush modes control how Figlet "smushes" together the big letters for output. This option is only really useful if you're making a font and need to experiment with the various settings—usually the font author has already specified the smush mode that works best with that font. You can find more information on smush modes in [figfont.txt](https://raw.github.com/lukesampson/figlet/master/figfont.txt), although this version of figfont.txt is written for the C version.

`message`
The message you want to print out. If you don't specify one, Figlet will go into interactive mode where it waits for you to enter a line of text and then prints it out in large letters. You can do this as many times as you like, and use Ctrl-C to quit.

### Why did you port it?

I couldn't get [the C version](https://github.com/cmatsuoka/figlet) to build and run properly on Windows using MSYS. Rather than mess around with lots of things I don't understand, I decided this would be a good opportunity to learn Go instead.

Also, the original version of this program is over 20 years old, and the code shows it. The main loop has a comment that says:

    The following code is complex and thoroughly tested.
    Be careful when modifying!

I like to think that the Go version is a lot clearer, especially with a lot of the legacy options stripped out. Although I admit the Go code is not the best—this is my first time programming in Go. I'd appreciate pull requests that make it better.


### Differences from the original version

###### Control files

Control files aren't supported in this version. They seem like a legacy workaround for something that's not so much a problem any more. I've tested passing unicode characters directly to this version and it seems to work ok, when the font supports the character. Even if I haven't gotten it right, Go has excellent UTF8 support so it shouldn't be too hard to fix this in a way that doesn't involve the complexity of control files.

###### Newline handling

The original version has options for handling newlines, and I think it renders newlines as it receives them from input. This version just treats newlines as whitespace and won't print a new line by default. I might be wrong, but I think this is pretty much what you want in most cases anyway.

###### Unsupported options

These command-line options aren't supported in this version:

`-knopstvxDELNSWX`
Too complicated!

`-f fontdirectory`
This version just uses the "fonts" directory in the same directory as the `figlet` executable.

`-C controlfile`
Control files aren't supported, for reasons given above.

`-I infocode`
Not supported

`-R`
This is supported, but it behaves differently in this version. In the original it meant "Right-to-left" print direction. In this version it means "Reverse" the print direction, as specified in the font file. Most times the font file is what you want, so this is mainly for testing and as a gimmick to confuse people.
