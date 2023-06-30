package docker

/*func SaveImage(imglist []string) (*responses.Response, error) {
	// 获取同步
	ctx := context.Background()

	reader, err := cli.ImageSave(ctx, imglist)

	if err != nil {
		return nil, errors.New("镜像备份失败")
	}
	// 读取所需要复制的内容

	// TODO written方法需要后续经过测试进行完善
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return nil, err
	}

	// 获取并更新镜像
	for _, img := range imglist {
		image, _ := dao.GetImgById(img)
		row, err := dao.UpdateImage(image)
		if err != nil || row == 0 {
			return nil, errors.New("镜像在进行更新备份时更新失败")
		}
	}

	return &responses.OK, nil
}
*/
