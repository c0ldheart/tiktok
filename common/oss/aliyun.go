package oss

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"tikapp/common/config"
	"tikapp/common/log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func AliyunInit() {
	client, err := oss.New(config.AliyunCfg.Endpoint, config.AliyunCfg.AccessKeyID, config.AliyunCfg.AccessKeySecret)
	if err != nil {
		fmt.Println(err)
		return
	}
	AliyunClient = client
	log.Logger.Info("aliyun oss 初始化成功")
}

func CreateBucket(name string) {
	err := AliyunClient.CreateBucket(name, oss.ACL(oss.ACLPublicReadWrite))
	if err != nil {
		exist, err := AliyunClient.IsBucketExist(name)
		if err == nil && exist {
			log.Logger.Info(fmt.Sprintf("We already own %s\n", name))
		} else {
			logrus.Error("create bucket error", err)
			return
		}
	}
	log.Logger.Info(fmt.Sprintf("Successfully created %s\n", name))
}

func UploadVideoToOss(bucketName string, objectName string, reader multipart.File) (bool, error) {
	bucket, err := AliyunClient.Bucket(bucketName)
	if err != nil {
		return false, err
	}
	err = bucket.PutObject(objectName, reader)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func GetOssVideoUrlAndImgUrl(bucketName string, objectName string) (string, string, error) {
	url := "https://" + bucketName + "." + config.AliyunCfg.Endpoint + "/" + objectName
	return url, url + "?x-oss-process=video/snapshot,t_0,f_jpg,w_0,h_0,m_fast,ar_auto", nil
}
