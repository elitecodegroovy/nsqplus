package nsqplus

import (
	"internal/http_api"
	"httprouter"
	"internal/clusterinfo"
	"net/http"
	"encoding/json"
	"strings"
	"nsq"
	"os"
	"log"
	"io/ioutil"
	"strconv"
	"bytes"
)

var nullLogger = log.New(ioutil.Discard, "[producer]", log.Ldate|log.Ltime|log.Lmicroseconds)
const (
	md5_key = "!@#$biostime&mama100!@#$"
)

type httpServer struct {
	ctx    		*Context
	router 		http.Handler
	client 		*http_api.Client
	ci     		*clusterinfo.ClusterInfo
	w      		*nsq.Producer
	responseChan    chan *nsq.ProducerTransaction
}

const TopticNamePrefix string = "biostime"

type ServerResponse struct {
	Code int32 `json:"code"`
	Msg  string `json:"msg"`
}

type AppServerResponse struct {
	Response ServerResponse `json:"response"`
}

type Event struct {
	Platform 		int32 	    `json:"platform"`
	PointCode		string      `json:"point_code"`
	Compaign                string      `json:"campaign"`
	Url 			string      `json:"url"`
	CreatedTime		int64       `json:"created_time"`
	CustomerId		string      `json:"customer_id"`
	UserMark		string      `json:"user_mark"`
	UserSourceId		string      `json:"user_sourceid"`
	Mobile			string      `json:"mobile"`
	MobileInfo		string      `json:"mobile_info"`
	AppVersion		string      `json:"app_version"`
	SourceFrom		string      `json:"source_from"`
	Sku			string 	    `json:"sku"`
	Spu			string 	    `json:"spu"`
	CourseId		string 	    `json:"course_id"`
	AccountId		string 	    `json:"account_id"`
	TerminalCode		string 	    `json:"terminal_code"`
	CouponDefid		string 	    `json:"coupon_defid"`
	ExpertId		string      `json:"expert_id"`
	QuestionId		string 	    `json:"question_id"`
	Ip 			string 	    `json:"ip"`
	VarX			string 	    `json:"var_x"`
	VarY			string 	    `json:"var_y"`
	VarZ			string 	    `json:"var_z"`
	Sign                    string 	    `json:"sign"`
}

type Events []*Event

func getMessages(msg []string) string{
	if len(msg) > 0 {
		return "WARNNING msg: "+ strings.Join(msg, "\n")
	}
	return ""
}

func (s *httpServer)initProducer() {
	config := nsq.NewConfig()
	w, err := nsq.NewProducer(strings.Join(s.ctx.nsqplus.opts.NSQDHTTPAddresses, ","), config)
	if err != nil {
		s.ctx.nsqplus.logf("ERRR can't created producer: ", err.Error())
		os.Exit(1)
	}
	w.SetLogger(nullLogger, nsq.LogLevelInfo)

	s.w = w
	s.responseChan = make(chan *nsq.ProducerTransaction)
}

func (s *httpServer)sendMsg(event Event )error {
	js, err := json.Marshal(event)
	if err != nil {
		s.ctx.nsqplus.logf("json.Marshal(event) ERROR :%s ", err.Error())
		return err
	}
	//s.ctx.nsqplus.logf("event json format :%s ", js)
	return s.w.Publish(TopticNamePrefix + strconv.FormatInt(int64(event.Platform), 10), []byte(js))
}

func (s *httpServer)sendMsgs(events Events)error {
	var platform int32
	var data [][]byte
	for _, event := range events {
		platform = event.Platform
		js, err := json.Marshal(event)
		if err != nil {
			s.ctx.nsqplus.logf("json.Marshal(event) ERROR :%s ", err.Error())
			return err
		}
		data = append(data, js)
		s.ctx.nsqplus.logf("req event: %#v, isSuccess:true ", event)
	}
	//s.ctx.nsqplus.logf("event json format :%s ", js)
	return s.w.MultiPublish(TopticNamePrefix + strconv.FormatInt(int64(platform), 10), data)
}

func NewHTTPServer(ctx *Context) *httpServer{
	log := http_api.Log(ctx.nsqplus.opts.Logger)
	client := http_api.NewClient(ctx.nsqplus.httpClientTLSConfig)

	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.PanicHandler = http_api.LogPanicHandler(ctx.nsqplus.opts.Logger)
	router.NotFound = http_api.LogNotFoundHandler(ctx.nsqplus.opts.Logger)
	router.MethodNotAllowed = http_api.LogMethodNotAllowedHandler(ctx.nsqplus.opts.Logger)
	s := &httpServer{
		ctx:    ctx,
		router: router,
		client: client,
		ci:     clusterinfo.New(ctx.nsqplus.opts.Logger, client),
	}
	s.initProducer()

	router.Handle("GET", "/", http_api.Decorate(s.indexHandler, log,  http_api.V1))
	router.Handle("GET", "/ping", http_api.Decorate(s.pingHandler, log,  http_api.V1))

	router.Handle("GET", "/pointstats", http_api.Decorate(s.pingHandler, log,  http_api.V1))
	//business logic API
	router.Handle("POST", "/pointstats/eventpv", http_api.Decorate(s.handlePVHandler, log,  http_api.V1))

	router.Handle("POST", "/pointstats/eventspv", http_api.Decorate(s.handlePVSHandler, log,  http_api.V1))
	router.Handle("POST", "/pointstats/appeventspv", http_api.Decorate(s.handleAppPVSHandler, log,  http_api.V1))

	return s
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *httpServer) indexHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return ServerResponse{1, "success"}, nil
}

func (s *httpServer) pingHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return ServerResponse{0, "ok"}, nil
}

//Only for Android and iOS native request
func (s *httpServer)handleAppPVSHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	var appEvents Events
	err := req.ParseForm()
	if err != nil{
		response := ServerResponse{0,err.Error()}
		return AppServerResponse{Response:response}, nil
	}

	jsonEvents := req.PostFormValue("pvdata")
	//s.ctx.nsqplus.logf("req body: %#v ", jsonEvents)
	err = json.Unmarshal([]byte(jsonEvents), &appEvents)
	if err != nil {
		response := ServerResponse{0,err.Error()}
		return AppServerResponse{Response:response}, nil
	}
	ip := RealIP(req)
	s.ctx.nsqplus.logf("realIp: %s ", ip)

	var messages []string
	for _, event := range appEvents {
		messages = s.validateEvent(event)
		event.Ip = ip
	}
	if len(messages) > 0 {
		msg := getMessages(messages)
		s.ctx.nsqplus.logf("error req body , messages: ", msg)
		return ServerResponse{0, msg}, nil
	}
	//post command to nsq-plus
	err = s.sendMsgs(appEvents)
	if err != nil {
		msg := "send message: " + err.Error()
		return ServerResponse{0, msg}, nil
	}
	//s.ctx.nsqplus.logf("req event: %#v, isSuccess:true ", event)
	return AppServerResponse{Response:ServerResponse{100, "ok"}}, nil
}



func (s *httpServer)handlePVSHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	events := Events{}
	err := json.NewDecoder(req.Body).Decode(&events)
	s.ctx.nsqplus.logf("req body: %#v ", events)
	if err != nil {
		return ServerResponse{0,err.Error()}, nil
	}
	ip := RealIP(req)
	s.ctx.nsqplus.logf("realIp: %s ", ip)

	var messages []string
	for _, event := range events {
		messages = s.validateEvent(event)
		event.Ip = ip
	}
	if len(messages) > 0 {
		msg := getMessages(messages)
		s.ctx.nsqplus.logf("error req body , messages: ", msg)
		return ServerResponse{0, msg}, nil
	}
	//post command to nsq-plus
	err = s.sendMsgs(events)
	if err != nil {
		msg := "send message: " + err.Error()
		return ServerResponse{0, msg}, nil
	}
	//s.ctx.nsqplus.logf("req event: %#v, isSuccess:true ", event)
	return ServerResponse{100, "ok"}, nil
}

func (s *httpServer) handlePVHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	event := Event{}

	err := json.NewDecoder(req.Body).Decode(&event)
	//s.ctx.nsqplus.logf("req body: %#v ", event)
	if err != nil {
		return ServerResponse{0,err.Error()}, nil
	}

	messages := s.validateEvent(&event)
	if len(messages) > 0 {
		msg := getMessages(messages)
		s.ctx.nsqplus.logf("error req body , messages: ", msg)
		return ServerResponse{0, msg}, nil
	}
	event.Ip = RealIP(req)
	s.ctx.nsqplus.logf("realIp: %s ", event.Ip)

	//post command to nsq-plus
	err = s.sendMsg(event)
	if err != nil {
		msg := "send message: " + err.Error()
		return ServerResponse{0, msg}, nil
	}
	s.ctx.nsqplus.logf("req event: %#v, isSuccess:true ", event)
	return ServerResponse{100, "ok"}, nil
}



//validate the request body parameters
func  (s *httpServer)validateEvent(event *Event) []string{
	var messages []string
	if event.Platform < 0 || event.Platform> 9 {
		messages = append(messages, "field 'platform' must be the range 0~99 ")
	}
	if len(event.PointCode) > 12  || len(event.PointCode) == 0 {
		messages = append(messages, "field 'point_code' must be not empty and length must be less than 13 ")
	}
	if event.Compaign != "" && len(event.Compaign) > 100 {
		messages = append(messages, "field 'compaign' must be not empty and length must be less than 100 ")
	}
	if event.Url != "" && len(event.Url) > 2000 {
		messages = append(messages, "field 'url' must be  less than 2000 ")
	}
	if event.CreatedTime == 0 {
		messages = append(messages, "field 'created_time' must be not empty and length must not be 0 ")
	}
	if event.CustomerId != "" && len(event.CustomerId) > 12 {
		messages = append(messages, "field 'customer_id' must be  less than 13 ")
	}
	if event.UserMark != "" && len(event.UserMark) > 50 {
		messages = append(messages, "field 'user_mark' must be  less than 50 ")
	}
	if event.UserSourceId != "" && len(event.UserSourceId) > 50 {
		messages = append(messages, "field 'user_source_id' must be  less than 50 ")
	}
	if event.Mobile != "" && len(event.Mobile) > 11 {
		messages = append(messages, "field 'created_time' must be 11 bit ")
	}
	if event.MobileInfo != "" && len(event.MobileInfo) > 100 {
		messages = append(messages, "field 'mobile_info' must be  less than 100 ")
	}
	if event.AppVersion != "" && len(event.AppVersion) > 10 {
		messages = append(messages, "field 'app_version' must be  less than 100 ")
	}
	if event.SourceFrom != "" && len(event.SourceFrom) > 2000 {
		messages = append(messages, "field 'source_from' must be  less than 3000 ")
	}
	if event.Sku != "" && len(event.Sku) > 10 {
		messages = append(messages, "field 'sku' must be less than 11 ")
	}
	if event.Spu != "" && len(event.Spu) > 10 {
		messages = append(messages, "field 'spu' must be less than 11 ")
	}
	if event.CourseId != "" && len(event.CourseId) > 10 {
		messages = append(messages, "field 'cource_id' must be less than 11 ")
	}
	if event.AccountId != "" && len(event.AccountId) > 10 {
		messages = append(messages, "field 'account_id' must be less than 11 ")
	}
	if event.TerminalCode != "" && len(event.TerminalCode) > 10 {
		messages = append(messages, "field 'terminal_code' must be less than 11 ")
	}
	if event.CouponDefid != "" && len(event.CouponDefid) > 10 {
		messages = append(messages, "field 'coupon_defid' must be less than 11 ")
	}
	if event.ExpertId != "" && len(event.ExpertId) > 10 {
		messages = append(messages, "field 'expert_id' must be less than 11 ")
	}
	if event.QuestionId != "" && len(event.QuestionId) > 10 {
		messages = append(messages, "field 'question_id' must be less than 11 ")
	}
	if event.VarX != "" && len(event.VarX) > 300 {
		messages = append(messages, "field 'VarX' must be less than 300 ")
	}
	if event.VarY != "" && len(event.VarY) > 300 {
		messages = append(messages, "field 'VarY' must be less than 300 ")
	}
	if event.VarZ != "" && len(event.VarZ) > 300 {
		messages = append(messages, "field 'VarZ' must be less than 300 ")
	}
	sign := buildMD5Digest(event)
	s.ctx.nsqplus.logf("req sign: %s ", sign)
	if event.Sign != sign {
		messages = append(messages, "field 'sign' is invalid. real sign:"+ sign + ",your sign:"+ event.Sign)
	}
	return messages
}

/**
   {
    "platform":1,
    "point_code":"1110100",
    "campaign":null,
    "url":null,
    "created_time":1492309920193,
    "customer_id":null,
    "user_mark":null,
    "user_sourceid":null,
    "mobile":null,
    "mobile_info":"iphone 6 plus",
    "app_version":"5.6.0",
    "source_from":"妈妈100App",
    "sku":null,
    "course_id":null,
    "account_id":null,
    "terminal_code":null,
    "coupon_defid":null,
    "expert_id":null,
    "question_id":null,
    "var_x":null,
    "var_y":null,
    "var_z":null
    }
 */
func buildMD5Digest(event *Event)string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(int64(event.Platform), 10))     //1:platform
	buffer.WriteString(event.PointCode)
	buffer.WriteString(event.Compaign)
	buffer.WriteString(event.Url)
	buffer.WriteString(strconv.FormatInt(int64(event.CreatedTime), 10))
	buffer.WriteString(event.CustomerId)
	buffer.WriteString(event.UserMark)
	buffer.WriteString(event.UserSourceId)
	buffer.WriteString(event.Mobile)
	buffer.WriteString(event.MobileInfo)
	buffer.WriteString(event.AppVersion)
	buffer.WriteString(event.SourceFrom)
	buffer.WriteString(event.Sku)
	buffer.WriteString(event.Spu)
	buffer.WriteString(event.CourseId)
	buffer.WriteString(event.AccountId)
	buffer.WriteString(event.TerminalCode)
	buffer.WriteString(event.CouponDefid)
	buffer.WriteString(event.ExpertId)
	buffer.WriteString(event.QuestionId)
	buffer.WriteString(event.VarX)
	buffer.WriteString(event.VarY)
	buffer.WriteString(event.VarZ)
	return GetMD5Digest(buffer.String())
}

