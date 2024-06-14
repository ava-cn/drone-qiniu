# drone-qiniu
上传文件到七牛的drone插件

## Build
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o qiniu
```

## Params

```
SET ACCESS_KEY=xxx
SET SECRET_KEY=xxx
SET ZONE=huadong
SET BUCKET=nameimtest
SET PREFIX=drone/
SET PATH=e:\drone-qiniu
```

## Usage
```
steps:
  - name: upload-static
    image: bysir/drone-qiniu:master
    pull: if-not-exists
    privileged: true
    settings:
      access_key:
        from_secret: qiniu_access_key
      SECRET_KEY:
        from_secret: qiniu_secret_key
      zone: huadong
      bucket: creght-sys
      prefix: render/
      dir: ./internal/render/static

```