package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"errors"
)

// 分页浏览容器信息
func IndexDocker(page int) (*responses.ImageListResponse, error) {
	u, total, err := dao.GetImageList(page, conf.PageSize)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}
	// 如果无数据，则返回到第一页
	if len(u) == 0 && page > 1 {
		page = 1
		u, total, err = dao.GetImageList(page, conf.PageSize)
		if err != nil {
			return nil, errors.New("获取用户列表失败")
		}
	}
	imageList := make([]responses.ImageInfo, len(u))
	for i, v := range u {
		tmp := responses.ImageInfo{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			ImageId:   v.ImageId,
			UserId:    v.UserId,
			Kind:      v.Kind,
		}
		imageList[i] = tmp
	}
	return &responses.ImageListResponse{
		Response:  responses.OK,
		Page:      page,
		Total:     total,
		ImageList: imageList,
	}, nil
}
