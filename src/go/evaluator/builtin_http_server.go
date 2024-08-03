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
var httpServerConfigs map[string]map[string]string

func doHttpServe(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	formatedTime := t.Format("2006-01-02 15:04:05")
	fmt.Println(formatedTime + " request url is: " + r.URL.Path)
	path := r.URL.Path
	function, ok := httpServerRoutes[path]
	if ok {
		result := Eval(function.Body, &initedEnv)
		routeConfig, ok := httpServerConfigs[path]
		if ok {
			contentType, ok := routeConfig["Content-Type"]
			if ok {
				w.Header().Add("Content-Type", contentType)
			}
		}
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
		httpServerConfigs = make(map[string]map[string]string)
		i := 0
		for _, route := range routes {
			path, ok1 := route.Key.(*object.String)
			function, ok2 := route.Value.(*object.Function)
			if ok1 && ok2 {
				httpServerRoutes[path.Value] = function
				i++
			} else {
				config, ok := route.Value.(*object.Hash)
				fmt.Println(route.Value)
				if ok {
					fmt.Println("is hash")
					key := object.String{}
					key.Value = "fn"
					hashKey, ok := object.Object(&key).(object.Hashable)
					if !ok {
						return newError("unusable as hash key: %s,", key.Type())
					}
					hashed := hashKey.HashKey()
					function, ok := config.Pairs[hashed].Value.(*object.Function)
					if ok {
						fmt.Println("add...")
						httpServerRoutes[path.Value] = function
					}
					cfgKey := object.String{}
					cfgKey.Value = "cfg"
					cfgHashKey, ok := object.Object(&cfgKey).(object.Hashable)
					if !ok {
						return newError("unusable as hash key:%s,", key.Type())
					}
					cfgHashed := cfgHashKey.HashKey()
					hashConfig, ok := config.Pairs[cfgHashed].Value.(*object.Hash)
					if ok {

						config := make(map[string]string, len(hashConfig.Pairs))
						for _, pair := range hashConfig.Pairs {
							config[pair.Key.Inspect()] = pair.Value.Inspect()
						}
						httpServerConfigs[path.Value] = config
					}
				}
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
