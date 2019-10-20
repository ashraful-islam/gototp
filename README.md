# Go TOTP: TOTP(HOTP) Token Generator

This is a simply TOTP(HOTP) token generator writtin in Go.
Inspired by [MinTOTP](https://github.com/susam/mintotp), the source code is translated
from Python to Go for easier compiled distribution.

## Disclaimer
This is a small hobby project and must not be used in production or secure environment.

## How To

### Requirements

Following are required to work on the source code locally:

- Go 1.12
- Git
- An IDE of choice

### Development

To work on the source code locally, simply clone to your local disk

```
$ git clone https://github.com/ashraful-islam/gototp.git
```

### Compilation

To compile to usable binary, run following command from inside the project directory

```
$ go build -o gototp main.go
```

There should be a new executable binary with name `gototp`(on Windows, it'll be `gototp.exe`)

### Usage

This is a simply cli based program. To check for available configurable parameters
```
$ ./gototp -h
```

Following configuration parameters are available

- `-key` : the secret key, i.e. seed used to generate the token (default: will try to read from stdin)
- `-digits` : number(positive integer) of digits the token will be of (default: `6`)
- `-digest` : the hashing algorithm to use, available `sha1`, `sha256`, `sha512` (default: auto detected if not given)
- `-steps` : time steps(in seconds) to use for counter (default: `30`)

Example using all defaults
```
$ ./gototp -key=ZYTYYE5FOAGW5ML7LRWUL4WTZLNJAMZS
// following is output
TOKEN: 675179
```

Example using more digits and different step size
```
$ ./gototp -key=ZYTYYE5FOAGW5ML7LRWUL4WTZLNJAMZS -digits=8 -digest=sha256
// following is output
TOKEN: 21690441
```

Example using key passed in stdin from a file
```
$ echo "ZYTYYE5FOAGW5ML7LRWUL4WTZLNJAMZS" > mykey.txt
$ cat mykey.txt | ./gototp
// following is output
TOKEN: 675179
```
OR
```
$ ./gototp <<< ZYTYYE5FOAGW5ML7LRWUL4WTZLNJAMZS
// followin is output
TOKEN: 675179
```

Example of passing configuration options as well as passing key in stdin
```
$ echo "ZYTYYE5FOAGW5ML7LRWUL4WTZLNJAMZS" > mykey.txt
$ cat mykey.txt | ./gototp -digits=7
// following is output
TOKEN: 0480559
```
OR
```
$ ./gototp -digits=7 <<< ZYTYYE5FOAGW5ML7LRWUL4WTZLNJAMZS
// followin is output
TOKEN: 0480559
```

*WARN*: For space separated keys(ones used by Dropbox), pass the key in quotes(`"`) else only one block will be detected!
```
$ ./gototp <<< "zyty ye5f oagw 5ml7 lrwu l4wt zlnj amzs"
//following is output
TOKEN: 675179
```

### CREDITS:

The project was inspired by [MinTOTP](https://github.com/susam/mintotp) by [Susam Pal](https://github.com/susam),
posted on [HackerNew](https://news.ycombinator.com/item?id=21297664). 

Some of the minor changes made were introduced following

- [RFC-6238](https://tools.ietf.org/html/rfc6238)
