# go-eds-logger
## 设置环境变量
**Windows 10**
1. 搜索框中输入“环境变量”
2. 选择“编辑系统环境变量”
3. 添加环境变量 `EDS_USR_ID` 和 `EDS_USR_PWD`
4. **重启电脑**（重启命令窗口或IDE不起作用）

**MacOS**
1. 编辑 `vi ~/.bash_profile`
2. 新增变量
   ```
   export EDS_USR_ID="your id"
   export EDS_USR_PWD="your password"
   ```
3. 生效 `source ~/.bash_profile`
4. 重启Goland

## Dev Notes
### v0.4.0
使用config配置项
### v0.3.0
将执行结果通过邮件告知
### v0.2.0
周末如果有调休也要填日志
### v0.1.0
略