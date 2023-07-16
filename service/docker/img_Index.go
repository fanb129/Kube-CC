package docker

import (
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"errors"
)

// 分页浏览容器信息
func IndexDocker(page int, uid uint) (*responses.ImageListResponse, error) {
	u, total, err := dao.GetImageList(page, conf.PageSize, uid)
	if err != nil {
		return nil, errors.New("获取镜像列表失败")
	}
	// 如果无数据，则返回到第一页
	if len(u) == 0 && page > 1 {
		page = 1
		u, total, err = dao.GetImageList(page, conf.PageSize, uid)
		if err != nil {
			return nil, errors.New("获取镜像列表失败")
		}
	}
	imageListPublic := make([]responses.ImageInfo, len(u))
	imageListPrivate := make([]responses.ImageInfo, len(u))
	imageListAll := make([]responses.ImageInfo, len(u))
	var public = 0
	var private = 0
	var alls = 0
	for _, v := range u {
		tmp := responses.ImageInfo{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			ImageName: v.ImageName,
			ImageId:   v.ImageId,
			UserId:    v.UserId,
			Kind:      v.Kind,
			Tag:       v.Tag,
			Size:      v.Size,
		}
		imageListAll[alls] = tmp
		alls++
		if v.Kind == 2 {
			imageListPrivate[private] = tmp
			private++
		} else {
			imageListPublic[public] = tmp
			public++
		}
	}
	return &responses.ImageListResponse{
		Response:         responses.OK,
		Page:             page,
		Total:            total,
		ImageListPBULIC:  imageListPublic,
		ImageListPRIVATE: imageListPrivate,
		ImageListAll:     imageListAll,
	}, nil
}
