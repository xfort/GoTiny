package tiny

import (
	"net/http"
	"encoding/base64"
	"io/ioutil"
	"bytes"
	"log"
	"path/filepath"
	"strings"
)

const TinyApiHost = "https://api.tinify.com"

type TinyHandler struct {
	ApiKey      string //api key 必须值
	httpclient  *http.Client
	authortoken string
	outImgDir   string
}

func (tinyhandler *TinyHandler) SetData(apikey string, outImgDir string) {
	tinyhandler.ApiKey = apikey
	tinyhandler.httpclient = http.DefaultClient
	tinyhandler.outImgDir = outImgDir
	tinyhandler.getAuthorCode(tinyhandler.ApiKey)
}

func (tinyhandler *TinyHandler) getAuthorCode(apikey string) string {
	strCode := "api:" + apikey
	tinyhandler.authortoken = "Basic " + base64.StdEncoding.EncodeToString([]byte(strCode))
	return tinyhandler.authortoken
}

/**
上传图片
 */
func (tinyhandler *TinyHandler) UploadFile(imgfilepath string) (error, string) {
	apiurl := TinyApiHost + "/shrink"
	imgBytes, err := ioutil.ReadFile(imgfilepath)
	if err != nil {
		return err, ""
	}
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(imgBytes))
	if err != nil {
		return err, ""
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Authorization", tinyhandler.authortoken)
	res, err := tinyhandler.httpclient.Do(req)
	if err != nil {
		return err, ""
	}
	imgUrl := res.Header.Get("Location")
	log.Println("img url", imgUrl)
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, imgUrl
	}
	log.Println("res", string(resBytes))
	return nil, imgUrl
}

/**
下载图片
 */
func (tinyhandler *TinyHandler) DownloadImg(imgUrl string, outFilePath string) error {

	req, err := http.NewRequest("GET", imgUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", tinyhandler.authortoken)
	response, err := tinyhandler.httpclient.Do(req)
	if err != nil {
		return err
	}
	resByres, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outFilePath, resByres, 0666)
	return err
}

/**
压缩单张图片
 */
func (tinyHandler *TinyHandler) CompressImageFile(imgFile string, outFile string) {
	if outFile == "" {
		fileName := filepath.Base(imgFile)
		outFile = filepath.Join(tinyHandler.outImgDir, fileName)
	}
	err, imgUrl := tinyHandler.UploadFile(imgFile)
	if err != nil {
		log.Println("上传图片错误", err, imgFile)
	}
	if imgUrl == "" {
		log.Println("上传图片错误_图片url为空", imgFile)
		return
	}
	err = tinyHandler.DownloadImg(imgUrl, outFile)
	if err != nil {
		log.Println("下载图片失败", imgFile, err)
		return
	}
	log.Println("下载图片成功", imgFile, outFile)
}

/**
压缩文件夹内的所有图片
 */
func (tinyHandler *TinyHandler) CompressAllImages(imgsDir string, outDir string) error {
	fileinfoList, err := ioutil.ReadDir(imgsDir)
	if err != nil {
		return err
	}
	for index, itemFile := range fileinfoList {
		if itemFile.IsDir() || itemFile.Size() <= 0 {
			continue
		}
		if !strings.Contains(itemFile.Name(), "png") && !strings.Contains(itemFile.Name(), "jpg") {
			continue
		}
		log.Println("压缩图片", index, itemFile.Name())
		tinyHandler.CompressImageFile(filepath.Join(imgsDir, itemFile.Name()), outDir)
	}
	return nil
}
