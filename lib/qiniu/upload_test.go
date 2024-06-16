package qiniu

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"testing"
)

func TestUploadDir(t *testing.T) {
	u := NewQiniu("", "").Uploader()
	err := u.UploadDir(&storage.ZoneHuadong, "nameimtest", "test/", `../qiniu/`, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("ok")
}

func TestUploadFile(t *testing.T) {
	u := NewQiniu("", "").Uploader()
	_, err := u.UploadFile(&storage.ZoneHuadong, "nameimtest", "drone/", `Z:\go_path\src\github.com\bysir-zl\drone-qiniu\Dockerfile`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("ok")
}
