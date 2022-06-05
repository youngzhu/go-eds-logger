# go-eds-logger
## 设置环境变量
**Windows 10**
1. 搜索框中输入“环境变量”
2. 选择“编辑系统环境变量”
3. 添加环境变量 `EDS_USR_ID` 和 `EDS_USR_PWD`
4. **重启电脑**

**MacOS**
1. 编辑 `vi ~/.bash_profile`
2. 新增变量
   ```
   export EDS_USR_ID="your id"
   export EDS_USR_PWD="your password"
   ```
3. 生效 `source ~/.bash_profile`
4. 重启Goland
