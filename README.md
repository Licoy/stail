## :rocket: 介绍
<div align="center">
  <h1>Stail</h1>
  <p>通过系统级的命令( 如 `tail -f` )来实时监控文件变动</p>
</div>

## :memo: 安装
```
go get github.com/Licoy/stail
```
## :hammer: 使用
```golang
func useSTail(filepath string, tailLine int) {
	st, err := stail.New(stail.Options{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = st.Tail(filepath, tailLine, func(content string) {
		fmt.Print(fmt.Sprintf("获取到内容: %s", content))
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
```
相当于
```bash
tail -{tailLine}f {filepath}
```
## :bulb: 参数
- `filepath` 需要监听的文件路径
- `tailLine` 只查看末尾的指定行
- `call` 内容回调func，内容类型为string 

## 协议
[MIT](./LICENSE)
