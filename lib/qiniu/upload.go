package qiniu

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Uploader struct {
	mac *qbox.Mac
}

// 自定义返回值结构体
type PutRsp struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
}

// UploadDir 上传一个文件夹
func (u Uploader) UploadDir(zone *storage.Zone, bucket string, keyPrefix string, dirPath string) (err error) {
	fmt.Printf("upload dir: '%s' to bucket '%s', prefix: '%s'\n", dirPath, bucket, keyPrefix)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		reaPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		start := time.Now()
		_, err = u.UploadFile(zone, bucket, keyPrefix+reaPath, path)
		if err != nil {
			return fmt.Errorf("upload file '%s' error: %w", keyPrefix+reaPath, err)
		}
		fmt.Printf("uploaded file '%s' success, spend time: %v\n", keyPrefix+reaPath, time.Since(start))

		return nil
	})
	if err != nil {
		return
	}
	return
}

// 服务端表单直传 + 自定义回 JSON
// key: 自定义上传文件名称 可以说是时间+string.后缀的形式
// localFile: 填入你本地图片的绝对地址，你也可以把图片放入项目文件中
func (u Uploader) UploadFile(zone *storage.Zone, bucket string, key string, localFile string) (ret PutRsp, err error) {
	// 上传文件自定义返回值结构体
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`,
	}
	upToken := putPolicy.UploadToken(u.mac)

	cfg := storage.Config{
		Zone:          zone,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)

	putExtra := storage.PutExtra{} // 可选配置 自定义返回字段
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return
	}

	return
}

// 服务端上传一个Reader
func (u Uploader) UploadReader(zone *storage.Zone, bucket string, key string, reader io.Reader, size int64) (ret PutRsp, err error) {
	putExtra := storage.PutExtra{} // 可选配置 自定义返回字段
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`,
	}
	upToken := putPolicy.UploadToken(u.mac)

	cfg := storage.Config{
		Zone:          zone,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)

	if key == "" {
		err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, reader, size, &putExtra)
	} else {
		err = formUploader.Put(context.Background(), &ret, upToken, key, reader, size, &putExtra)
	}
	if err != nil {
		return
	}

	return
}
