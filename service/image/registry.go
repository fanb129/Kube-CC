package image

//ListImages 列出public和属于自己的, uid为0，则列出所有
//func ListImages(uid, gid uint) (*responses.ImageListResponse, error) {
//	repos, err := dao.Hub.Repositories()
//	if err != nil {
//		return nil, err
//	}
//	ImageList := make([]responses.ImageInfo, 0)
//	len := 0
//	for _, repo := range repos {
//		tags, err := dao.Hub.Tags(repo)
//		if err != nil {
//			return nil, err
//		}
//
//		for _, tag := range tags {
//			manifest, err := dao.Hub.ManifestV2(repo, tag)
//			if err != nil {
//				return nil, err
//			}
//
//			annotations := manifest.Config.Annotations
//
//			if uid > 0 { // 普通用户查看自己
//				if annotations != nil && (annotations["kind"] == "public" || annotations["uid"] == strconv.Itoa(int(uid))) {
//					image, err := addToImageList(annotations, repo, tag, manifest)
//					if err != nil {
//						zap.S().Errorln(err)
//					}
//					ImageList = append(ImageList, *image)
//					len++
//				}
//			} else {
//				if gid > 0 { // 组管理员查看本组所有
//					users, err := dao.GetGroupUserById(gid)
//					if err != nil {
//						//zap.S().Errorln(err)
//						return nil, err
//					}
//					for _, user := range users {
//						if annotations != nil && (annotations["kind"] == "public" || annotations["uid"] == strconv.Itoa(int(user.ID))) {
//							image, err := addToImageList(annotations, repo, tag, manifest)
//							if err != nil {
//								zap.S().Errorln(err)
//							}
//							ImageList = append(ImageList, *image)
//							len++
//							break
//						}
//					}
//				} else { // 超级管理员查看所有
//					if annotations != nil && (annotations["kind"] == "public" || annotations["kind"] == "private") {
//						image, err := addToImageList(annotations, repo, tag, manifest)
//						if err != nil {
//							zap.S().Errorln(err)
//						}
//						ImageList = append(ImageList, *image)
//						len++
//					}
//				}
//			}
//
//		}
//	}
//
//	return &responses.ImageListResponse{
//		ImageList: ImageList,
//		Length:    len,
//		Response:  responses.OK,
//	}, nil
//}

// DeleteImage 从镜像仓库删除 无效
//func DeleteImage(name, tag string) (*responses.Response, error) {
//	digest, err := dao.Hub.ManifestDigest(name, tag)
//	if err != nil {
//		return nil, err
//	}
//	err = dao.Hub.DeleteManifest(name, digest)
//	if err != nil {
//		return nil, err
//	}
//	return &responses.OK, nil
//}

// EditImage 更改 无效
//func EditImage(name, tag string, uid uint, kind, description string) (*responses.Response, error) {
//	manifest, err := dao.Hub.ManifestV2(name, tag)
//	if err != nil {
//		return nil, err
//	}
//	config := manifest.Target()
//	if config.Annotations == nil {
//		config.Annotations = make(map[string]string)
//	}
//
//	config.Annotations["uid"] = strconv.Itoa(int(uid))
//	config.Annotations["kind"] = kind
//	config.Annotations["description"] = description
//
//	err = dao.Hub.PutManifest(name, tag, manifest)
//	if err != nil {
//		return nil, err
//	}
//	return &responses.OK, nil
//}

// 大小单位为MB
//func getSize(size int64) string {
//	floatSize := float64(size) / 1000000.00
//	decimalSize := strconv.FormatFloat(floatSize, 'f', 1, 64)
//	return decimalSize + "MB"
//}

//func addToImageList(annotations map[string]string, repo, tag string, manifest *schema2.DeserializedManifest) (*responses.ImageInfo, error) {
//	getuid, err := strconv.Atoi(annotations["uid"])
//	if err != nil {
//		//zap.S().Errorln(err)
//		return nil, err
//	}
//	username := ""
//	nickname := ""
//	user, err := dao.GetUserById(uint(getuid))
//	if err != nil {
//		zap.S().Errorln(err)
//	} else {
//		username = user.Username
//		nickname = user.Nickname
//	}
//
//	image := responses.ImageInfo{
//		Username: username,
//		Nickname: nickname,
//		Uid:      uint(getuid),
//
//		Name: repo,
//		Tag:  tag,
//		Kind:    annotations["kind"],
//		Size:    getSize(manifest.Config.Size),
//		ImageId: manifest.Config.Digest.String(),
//
//		Description: annotations["description"],
//	}
//	return &image, nil
//}
