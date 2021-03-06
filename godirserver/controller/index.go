package controller

import (
	"../config"
	"../model"
	"../utils"
	"../utils/uuid"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	
	p := r.FormValue("p")
	if p == "" {
		p = "/"
	}
	p = strings.Replace(p, "../", "/", -1)
	s, _ := os.Stat(path.Join(config.BASEPATH, p)) //os.Stat获取文件信息
	if s.IsDir() {
		dirs, files := utils.ScanDir(p)
		resData := model.JsonResult{Code: 200, Msg: "", Data: model.Res{Dirs: dirs, Files: files}}
		
		msg, _ := json.Marshal(resData)
		
		//fmt.Fprintf(w, "hello world")
		w.Header().Set("content-type", "text/json")
		w.Write(msg)
	} else {
		//w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(p)))
		fData, _ := ioutil.ReadFile(path.Join(config.BASEPATH, p))
		w.Write(fData)
	}
	
}
func Upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "*") //允许访问所有域
	
	err := r.ParseForm() //解析表单                 //获取文件内容
	if err != nil {
		log.Println(err)
		return
	}
	imgFile, _, err := r.FormFile("f") //获取文件内容
	if err != nil {
		log.Println(err)
		return
	}
	defer imgFile.Close()
	imgName := ""
	files := r.MultipartForm.File //获取表单中的信息
	for _, v := range files {
		for _, vv := range v {
			imgName = vv.Filename
		}
	}
	p := r.FormValue("p")
	if p == "" {
		p = "/"
	}
	p = strings.Replace(p, "../", "/", -1)
	p = path.Join(config.BASEPATH, p)
	_uuid := uuid.New()
	fname := fmt.Sprintf("%x%x%x%x%x", _uuid[0:4], _uuid[4:6], _uuid[6:8], _uuid[8:10], _uuid[10:16])
	imgPath := path.Join(p, fname+path.Ext(imgName))
	
	saveFile, err := os.Create(imgPath)
	if err != nil {
		log.Println(err)
		return
	}
	if saveFile != nil {
		defer saveFile.Close()
	}
	
	s, err := io.Copy(saveFile, imgFile) //保存
	if err != nil {
		log.Println(err)
		return
	}
	
	if s == 0 {
		log.Println("保存不成功")
		return
	}
	
	resData := model.JsonResult{Code: 200, Msg: ""}
	msg, _ := json.Marshal(resData)
	w.Write(msg)
}
func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	p := r.FormValue("p")
	if p == "" {
		p = "/"
	}
	p = strings.Replace(p, "../", "/", -1)
	p = path.Join(config.BASEPATH, p)
	
	err := os.Remove(p)
	fmt.Println(err)
	
	resData := model.JsonResult{Code: 200, Msg: ""}
	msg, _ := json.Marshal(resData)
	w.Write(msg)
}
