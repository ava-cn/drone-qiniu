# drone-qiniu
上传文件到七牛的 drone 插件

## Build
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o qiniu
```

## Params

```
SET ACCESS_KEY=xxx
SET SECRET_KEY=xxx
SET ZONE=huadong
SET BUCKET=mybucket
SET PREFIX=drone/
SET dir=./dist
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