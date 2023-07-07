package docker

import (
	"errors"
	"math"
	"strconv"
)

// 某人大小单位为MB
func GetIDandSize(repositoryName string) (string, string, error) {
	err, imageList := GetImageListAll()
	if err != nil {
		return "", "-1", err
	}
	id := ""
	var size int64
	for _, image := range imageList {
		if image.RepoTags[0] == repositoryName {
			id = image.ID[7:19]
			size = image.Size
			// 实现占用大小的四舍五入
			floatSize := float64(size) / 1000000.00
			return id, strconv.Itoa(int(math.Floor(floatSize+0.5))) + "MB", nil
		}
	}
	return "", "-1`", errors.New("not found")
}
