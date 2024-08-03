package evaluator

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"z/object"
)

var httpServerRoutes map[string]*object.Function

func doHttpServe(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	formatedTime := t.Format("2006-01-02 15:04:05")
	fmt.Println(formatedTime + " request url is: " + r.URL.Path)
	path := r.URL.Path
	function, ok := httpServerRoutes[path]
	if ok {
		result := Eval(function.Body, &initedEnv)
		io.WriteString(w, result.Inspect())
	} else {
		io.WriteString(w, "path not found")
	}
}

func init_builtin_http_server() *object.Builtin {
	http_server := &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if args[0].Type() != object.STRING_OBJ {
			return newError("argument 1 to `mysql_init` must be String, got=%s", args[0].Type())
		}
		server := args[0].(*object.String).Value
		if args[1].Type() != object.HASH_OBJ {
			return newError("argument 1 to `mysql_init` must be Hash, got=%s", args[1].Type())
		}
		routes := args[1].(*object.Hash).Pairs
		httpServerRoutes = make(map[string]*object.Function)
		i := 0
		for _, route := range routes {
			path, ok1 := route.Key.(*object.String)
			function, ok2 := route.Value.(*object.Function)
			if ok1 && ok2 {
				httpServerRoutes[path.Value] = function
				i++
			}
		}
		fmt.Println("begin start serve, server address is:", server)
		fmt.Println("url list as follow:")
		for route := range httpServerRoutes {
			fmt.Println(route)
		}
		http.HandleFunc("/", doHttpServe)
		fmt.Println("control + c to end the server")
		err := http.ListenAndServe(server, nil)
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
	}
	return http_server
}
