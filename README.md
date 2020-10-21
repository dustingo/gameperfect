### GamePerfect
- 游戏服务器初始化工具
- Usage  
```text
usage: gameperfect [<flags>]

Flags:
  --help       Show context-sensitive help (also try --help-long and --help-man).
  --yaml       only do yaml action，default is only do basic check[sys,cpu,disk ...]
  --all        do basic check and yaml action
  --usage      The explanation of the config of this tool
  --version    show the version of the tools
  --mode=MODE  chose witch mode you want to run,must use with --yaml

```
```text
>>目录：config 存放yaml配置文件和第三方脚本配置
>>目录:tools 存放自定义脚本
>>config.yaml: 
  >1.modeYum: action仅为'install'时，执行安装
  >2.modeDir: action为'chown'时，根据para [用户.属组],更改目录属性，不存在则创建; action为'chmod'时，根据perm[]执行更改目录|文件权限
  >3.modeScripts: action为'run'且匹配'hostname' 时，根据env [shell|python]和name [path + scripts_name]执行自定义脚本
ATTENTION！注意：
the scripts should add 'x' permission
如果scripts是在win编辑的，在linux下记得使用dos转换
perm注意格式，如：0755 | 0600
host匹配服务器hostname时，只在匹配的服务器执行，当host为'all' 或 host为 '' 空时，任何服务器都执行
>>dns.yaml: 
  >1.https的完整url，执行https GET
  >2.域名，执行域名解析
```