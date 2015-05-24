# Go/Migemo socket server

Socket server of [Go/Migemo]
[Go/Migemo]:https://github.com/koron/gomigemo

## Install

```sh
go get github.com/mizyoukan/gmigemo-socketserver
```

## Usage

Run server:

```sh
./gmigemo-socketserver
```

And connect with own socket client like following (Python3):

```python
import socket
from contextlib import closing

def main():
  s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
  with closing(sock):
    s.connect(('127.0.0.1', 13730))
    s.sendall(b'meido')
    p = s.recv(1024)
    print(p)

if __name__ == '__main__':
  main()
```

## Option

Flag | Description | Default
:-- |:-- |:--
-h | Server host | 127.0.0.1
-p | Server port | 13730
-v | Use vim style regexp | false
-e | Use emacs style regexp | false
-n | Don't use newline match | false


## License

NYSL Version 0.9982
