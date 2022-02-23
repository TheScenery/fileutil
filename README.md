# file_util
Go lang file operation helper methods

# How to use?
```
go get -u https://github.com/TheScenery/fileutil
```

### Zip
Zip a file or dir use zip method
```go
fileutil.Zip("xxx.zip", "srcDir")
```

### UnZip
Un zip a file
```go
fileutil.UnZip("dstDir", "xxx.zip")
```