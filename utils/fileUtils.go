package utils

import (
	"log"
	"os"
)

func DirIsExist(dir string) {
	_, err := os.Stat(dir)
	//文件夹或者文件不存在
	if err != nil {
		err := os.Mkdir(dir, 0777)
		if err != nil {
			log.Printf("创建%s文件夹失败 %v", dir, err)
		}
		_ = os.Chmod(dir, 0777)
	}
}