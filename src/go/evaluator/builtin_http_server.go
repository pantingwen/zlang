package evaluator

import (
	"encoding/json"
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
	fmt.Println(r.URL.Query())
	body, _ := io.ReadAll(r.Body)
	fmt.Println(string(body))

	postJson := make(map[string]interface{})
	err := json.Unmarshal(body, &postJson)
	pairs := make(map[object.HashKey]object.HashPair)
	if err == nil {
		handlePostData(postJson, pairs)
	}
	handleGetData(r, pairs)
	request := object.Hash{Pairs: pairs}
	initedEnv.Set("request", object.Object(&request)) // pass request parameter

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
		io.WriteString(w, result.Json())
	} else {
		io.WriteString(w, "path not found")
	}
}

func handleGetData(r *http.Request, pairs map[object.HashKey]object.HashPair) {
	getPirs := make(map[object.HashKey]object.HashPair)
	for query, value := range r.URL.Query() {
		getItemName := object.String{Value: query}
		getItemNameHashKey, _ := object.Object(&getItemName).(object.Hashable)
		getItemNameHashed := getItemNameHashKey.HashKey()
		getItemValue := object.String{Value: ""}
		if len(value) > 0 {
			getItemValue.Value = value[0]
		}
		getPirs[getItemNameHashed] = object.HashPair{Key: object.Object(&getItemName), Value: object.Object(&getItemValue)}
	}
	getHash := object.Hash{Pairs: getPirs}
	getName := object.String{}
	getName.Value = "get"
	getNameHashKey, _ := object.Object(&getName).(object.Hashable)
	getNameHashed := getNameHashKey.HashKey()
	pairs[getNameHashed] = object.HashPair{Key: object.Object(&getName), Value: object.Object(&getHash)}
}

func handlePostData(postJson map[string]interface{}, pairs map[object.HashKey]object.HashPair) {
	postPirs := make(map[object.HashKey]object.HashPair)
	for post, value := range postJson {
		postItemName := object.String{Value: post}
		postItemNameHashKey, _ := object.Object(&postItemName).(object.Hashable)
		postItemNameHashed := postItemNameHashKey.HashKey()
		postItemValue := object.String{Value: ""}
		valueStr, ok := value.(string)
		if ok {
			postItemValue.Value = valueStr
		}
		postPirs[postItemNameHashed] = object.HashPair{Key: &postItemName, Value: &postItemValue}
	}
	postHash := object.Hash{Pairs: postPirs}
	postName := object.String{}
	postName.Value = "post"
	postNameHashKey, _ := object.Object(&postName).(object.Hashable)
	postNameHashed := postNameHashKey.HashKey()
	pairs[postNameHashed] = object.HashPair{Key: &postName, Value: object.Object(&postHash)}
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
				if ok {
					key := object.String{}
					key.Value = "fn"
					hashKey, ok := object.Object(&key).(object.Hashable)
					if !ok {
						return newError("unusable as hash key: %s,", key.Type())
					}
					hashed := hashKey.HashKey()
					function, ok := config.Pairs[hashed].Value.(*object.Function)
					if ok {
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
