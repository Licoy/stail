[中文](./README.md) | English

## :rocket: Introduction
<div align="center">
  <h1>Stail</h1>
  <p>real-time monitoring of file changes through system-level commands (such as `tail -f`)</p>
</div>

## :wrench: Install
```
go get github.com/Licoy/stail
```
## :hammer: Use
```golang
func useSTail(filepath string, tailLine int) {
    st, err := stail.New(stail.Options{})
    if err != nil {
        fmt.Println(err)
        return
    }
    si, err := st.Tail(filepath, tailLine, func(content string) {
        fmt.Print(fmt.Sprintf("get content: %s", content))
    })
    if err != nil {
        fmt.Println(err)
        return
    }
    time.AfterFunc(time.Second*10, func() {
        si.Close() // close the acquisition channel after 10 s
    }
    si.Watch()
}
```
equivalent to
```bash
tail -{tailLine}f {filepath}
```
## :bulb: Parameter
- `filepath` file path that needs to be monitored
- `tailLine` view only the specified line at the end
- `call` the content callback func, the content type is string

## :memo: License
[MIT](./LICENSE)
