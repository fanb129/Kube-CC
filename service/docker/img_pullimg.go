package docker

// PullImage
// 拉取指定镜像
// TODO 后续将方法规范重构
/*func PullImage(imageName string) (*responses.Response, error) {
	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}

	return nil
}

// PullPrivateImage
// 拉取私有仓库的镜像
func PullPrivateImage(imageName string, username string, passwd string) (*responses.Response, error) {
	ctx := context.Background()
	authConf := registry.AuthConfig{
		Username: username,
		Password: passwd,
	}
	encodeJson, _ := json.Marshal(authConf)
	authStr := base64.StdEncoding.EncodeToString(encodeJson)
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}
*/
