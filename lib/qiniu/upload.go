package qiniu

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
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

type Assets struct {
	Path string
	Key  string
}

// UploadDir 上传一个文件夹
func (u Uploader) UploadDir(zone *storage.Zone, bucket string, keyPrefix string, dirPath string, parallel int) (err error) {
	if parallel == 0 {
		parallel = 5
	}
	fmt.Printf("upload dir: '%s' to bucket '%s', prefix: '%s', parallel: '%d'\n", dirPath, bucket, keyPrefix, parallel)

	as := make(chan Assets, 10)

	// 上传
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range as {
				if err != nil {
					return
				}
				start := time.Now()
				rsp, e := u.UploadFile(zone, bucket, item.Key, item.Path)
				if e != nil {
					err = fmt.Errorf("upload file '%s' error: %w", item.Key, err)
					return
				}
				fmt.Printf("uploaded file '%s' success, size: %s, spend time: %v\n", item.Key, humanize.IBytes(uint64(rsp.Fsize)), time.Since(start))
			}
		}()
	}

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

		as <- Assets{
			Path: path,
			Key:  keyPrefix + reaPath,
		}

		return nil
	})
	if err != nil {
		return
	}
	close(as)

	wg.Wait()
	return err
}

// UploadFile 服务端表单直传 + 自定义回  JSON
// key: 自定义上传文件名称 可以说是时间+string.后缀的形式
// localFile: 填入你本地图片的绝对地址，你也可以把图片放入项目文件中
func (u Uploader) UploadFile(zone *storage.Zone, bucket string, key string, localFile string) (ret PutRsp, err error) {
	// 上传文件自定义返回值结构体
	putPolicy := storage.PutPolicy{
		// 覆盖上传，参见：https://developer.qiniu.com/kodo/1206/put-policy
		Scope:      fmt.Sprintf("%s:%s", bucket, key),
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
