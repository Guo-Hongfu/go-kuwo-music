package middleware

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"kuwo/request"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)
type FileDownloader struct {
	fileSize       int
	url            string
	saveName string
	totalPart      int //下载线程
	outputDir      string
	doneFilePart   []filePart


}

//filePart 文件分片
type filePart struct {
	Index int    //文件分片的序号
	From  int    //开始byte
	To    int    //解决byte
	Data  []byte //http下载得到的文件内容
}

func NewFileDownloader(url, saveName, outputDir string, totalPart int) *FileDownloader {
	if outputDir == "" {
		wd, err := os.Getwd() //获取当前工作目录
		if err != nil {
			log.Println(err)
		}
		outputDir = wd
	}
	return &FileDownloader{
		fileSize:       0,
		url:            url,
		saveName: saveName,
		outputDir:      outputDir,
		totalPart:      totalPart,
		doneFilePart:   make([]filePart, totalPart),
	}
}

// 开始下载
func (d *FileDownloader) Run() error {
	// 获取 下载文件长度
	d.fileSize = d.getFileSize()
	if 0 == d.fileSize{
		zap.S().Error("head获取失败，长度为0")
		return errors.New("head获取失败，长度为0")
	}
	jobs := make([]filePart, d.totalPart)
	eachSize := d.fileSize / d.totalPart
	// 构建每个分块 需要下载的文件大小
	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To + 1
		}
		if i < d.totalPart-1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			//the last filePart
			jobs[i].To = d.fileSize - 1
		}
	}
	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job filePart) {
			defer wg.Done()
			err := d.downloadPart(job)
			if err != nil {
				zap.S().Error("！！！！ 下载文件失败，",err, job)
			}
		}(j)

	}
	wg.Wait()
	return d.mergeFileParts()



}

func (d *FileDownloader) getFileSize() int {
	req := request.Requests{Method: "HEAD",Url: d.url}
	resp,_ := req.News()
	if length,err := strconv.Atoi(resp.Header.Get("Content-Length")); err == nil {
		return length
	}
	return 0
}

// 下载分片文件
func (d FileDownloader) downloadPart(c filePart) error {
	//zap.S().Infof("开始[%d]下载from:%d to:%d\n", c.Index, c.From, c.To)
	req := request.Requests{Method: "GET",Url: d.url,Headers: request.ReqH{
		"Range":fmt.Sprintf("bytes=%v-%v", c.From, c.To),
	}}
	resp,bs := req.News()

	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("服务器错误状态码: %v", resp.StatusCode))
	}
	bsSize := len(bs)
	if bsSize != (c.To - c.From + 1) {
		return errors.New("下载文件分片长度错误")
	}
	c.Data = bs
	d.doneFilePart[c.Index] = c
	return nil
}

func (d FileDownloader) mergeFileParts() error {
	zap.S().Info("%s 开始合并",d.saveName)
	path := filepath.Join(d.outputDir, d.saveName)
	mergedFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer mergedFile.Close()
	hash := sha256.New()
	totalSize := 0
	for _, s := range d.doneFilePart {

		mergedFile.Write(s.Data)
		hash.Write(s.Data)
		totalSize += len(s.Data)
	}
	if totalSize != d.fileSize {
		zap.S().Errorf("！！！！ %s 文件不完整",d.saveName)
		return errors.New("文件不完整")
	}
	zap.S().Info("~~~~~~ %s 下载完成~~~~~~~",d.saveName)
	return nil

}
