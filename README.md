# Filtering rectangle boxes Project
In this project we implement a web server (http-server) in order to filtered inputs. These inputs are some retangles with their coordinates. If their coordinates does not have intersection with the coordinates of domain, then they would be filtered. otherwise they will be saved in a `Text` file.

This project is implemented in `golang`.


## Author
Name: `Saeedeh (Javaneh) Bahrami`

Email: bahramisaeede@gmail.com


## Requirements
List of packages:
1. "fmt"
2. "errors"
3. "os"
4. "time"
5. "net/http"
6. "reflect"
7. "io/ioutil"
8. "io"
9. "encoding/json"
10. "log"
11. "bufio"
12. "github.com/mitchellh/mapstructure"
13. "github.com/gorilla/mux"


the last 2 packages must be installed by these commands:

```bash
go get -u github.com/mitchellh/mapstructure
go get -u github.com/gorilla/mux

```

## Running the Http Server
<br>for running the model:
1. first create a mod init file and select a name for your module
```bash
go mod init http_server
```

2. then compile the file
```bash
go build
```

3. finally run the program using module name
```bash
./http_server
```


## Requesting Data Format
In order to feed request data to the http server, plase use `curl post` as below:
```bash
curl -X POST -s localhost:8080 -d '{"main": {"x": 0, "y": 0, "width": 10, "height": 20}, "input": [{"x": 2, "y": 18, "width": 5, "height":4},{"x": -1, "y": -1, "width": 5, "height": 4}]}'

```

in the above command, the domain coordinate nameed as `"main"`.
the rest of coordinates are coordinates of input rectangles named `"input"`

`Note`: we set `localhost` as our endpoint.

`Note`: we set `8080` port. one can change it to another free port.


## Response Format
In order to get the response of data from the http server, plase use `curl get` as below:

```bash
curl -X GET -s localhost:8080
```

`Note`: the output will the list of rectangles which are saved in the `Text file`:
