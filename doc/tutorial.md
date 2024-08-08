### introduce
 z is simple language for build web application.

### install

install from source code
```shell
git clone git@github.com:pantingwen/zlang.git
cd zlang/src/go
make
cd ../../ # found the dist/z
```
### hello world

```shell
echo 'puts("hello world", "\n")' > hello.z
dist/z hello.z ## the output is "hello world"
```