package tool

import (
	"context"
	"errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
)

var GlobalOssConf map[string]string

type Oss struct {
	client *cos.Client
}

func (o *Oss) getConf()(map[string]string,error) {
	if GlobalOssConf == nil {
		cfg := new(Config)
		OssConf ,err := cfg.GetSection("oss")
		if err != nil {
			return nil,err
		}
		GlobalOssConf = OssConf
	}
	return GlobalOssConf,nil
}


func (o *Oss) NewClient() error  {
	//1、获取配置项
	cfg ,err :=o.getConf()
	if err!=nil {
		return err
	}

	rawUrl,ok := cfg["url"]
	if !ok{
		return errors.New("oss.url配置读取失败")
	}

	secretId,ok := cfg["secretId"]
	if !ok{
		return errors.New("oss.secretId配置读取失败")
	}

	secretKey,ok := cfg["secretKey"]
	if !ok{
		return errors.New("oss.secretKey配置读取失败")
	}

	u, _ := url.Parse(rawUrl)
	b := &cos.BaseURL{
		BucketURL: u,
	}
	ossClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
			Transport: &debug.DebugRequestTransport{
				//RequestHeader:  true,
				//RequestBody:    true,
				//ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	o.client = ossClient
	return nil
}

//上传文件至Oss
//cosPath远程OSS路径
//outFileName 本地目录路径
func (o *Oss) UploadFile(cosPath string,outFileName string)(string,error){
	//1、获取授权客户端
	if o.client == nil {
		err := o.NewClient()
		if err != nil {
			return "",err
		}
	}

	_, err := o.client.Object.PutFromFile(context.Background(), cosPath, outFileName, nil)
	if err != nil {
		return "",err
	}

	rawUrl := GlobalOssConf["url"]

	return rawUrl+cosPath,nil
}

