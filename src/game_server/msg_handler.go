package main

import (
	"bytes"
	"compress/zlib"
	_ "crypto/aes"
	"crypto/tls"
	"io/ioutil"
	"mm_server_new/libs/log"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	"mm_server_new/src/share_data"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
)

var msg_handler_http_mux map[string]func(http.ResponseWriter, *http.Request)

type MsgHttpHandle struct{}

func (this *MsgHttpHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var act_str, url_str string
	url_str = r.URL.String()
	idx := strings.Index(url_str, "?")
	if -1 == idx {
		act_str = url_str
	} else {
		act_str = string([]byte(url_str)[:idx])
	}
	//log.Debug("hall msg handler ServeHTTP actstr(%s)", act_str)
	if h, ok := msg_handler_http_mux[act_str]; ok {
		h(w, r)
	}

	return
}

//=======================================================

type CLIENT_MSG_HANDLER func([]byte) (int32, *Player)

type CLIENT_PLAYER_MSG_HANDLER func(*Player, []byte) int32

type MsgHandlerInfo struct {
	msg_handler        CLIENT_MSG_HANDLER
	player_msg_handler CLIENT_PLAYER_MSG_HANDLER
	if_player_msg      bool
}

type MsgHandlerMgr struct {
	msg_http_listener net.Listener
	login_http_server http.Server
	msgid2handler     map[int32]*MsgHandlerInfo
}

var msg_handler_mgr MsgHandlerMgr

func (this *MsgHandlerMgr) Init() bool {
	this.msgid2handler = make(map[int32]*MsgHandlerInfo)
	return true
}

func (this *MsgHandlerMgr) SetMsgHandler(msg_code uint16, msg_handler CLIENT_MSG_HANDLER) {
	this.msgid2handler[int32(msg_code)] = &MsgHandlerInfo{msg_handler: msg_handler, if_player_msg: false}
}

func (this *MsgHandlerMgr) SetPlayerMsgHandler(msg_code uint16, msg_handler CLIENT_PLAYER_MSG_HANDLER) {
	this.msgid2handler[int32(msg_code)] = &MsgHandlerInfo{player_msg_handler: msg_handler, if_player_msg: true}
}

func (this *MsgHandlerMgr) StartHttp() bool {
	var err error
	this.reg_http_mux()

	this.msg_http_listener, err = net.Listen("tcp", config.ListenClientInIP)
	if nil != err {
		log.Error("Center StartHttp Failed %s", err.Error())
		return false
	}

	signal_mgr.RegCloseFunc("msg_handler_mgr", this.CloseFunc)

	msg_http_server := http.Server{
		Handler:     &MsgHttpHandle{},
		ReadTimeout: 6 * time.Second,
	}

	log.Info("???????????????????????? IP:%s", config.ListenClientInIP)
	err = msg_http_server.Serve(this.msg_http_listener)
	if err != nil {
		log.Error("?????????????????????????????? %s", err.Error())
		return false
	}

	return true
}

func (this *MsgHandlerMgr) StartHttps(crt_file, key_file string) bool {
	this.reg_http_mux()

	signal_mgr.RegCloseFunc("msg_handler_mgr", this.CloseFunc)

	this.login_http_server = http.Server{
		Addr:        config.ListenClientInIP,
		Handler:     &MsgHttpHandle{},
		ReadTimeout: 6 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	err := this.login_http_server.ListenAndServeTLS(crt_file, key_file)
	if err != nil {
		log.Error("??????????????????????????????%s", err.Error())
		return false
	}

	return true
}

func (this *MsgHandlerMgr) CloseFunc(info *SignalRegRecod) {
	if nil != this.msg_http_listener {
		this.msg_http_listener.Close()
	}

	this.login_http_server.Close()

	info.close_flag = true
	return
}

//=========================================================

func (this *MsgHandlerMgr) reg_http_mux() {
	msg_handler_http_mux = make(map[string]func(http.ResponseWriter, *http.Request))
	msg_handler_http_mux["/client_msg"] = client_msg_handler
}

func _send_error(w http.ResponseWriter, ret_code int32) {
	m := &msg_client_message.S2C_ONE_MSG{ErrorCode: ret_code}
	res2cli := &msg_client_message.S2C_MSG_DATA{MsgList: []*msg_client_message.S2C_ONE_MSG{m}}
	final_data, err := proto.Marshal(res2cli)
	if nil != err {
		log.Error("client_msg_handler marshal 1 client msg failed err(%s)", err.Error())
		return
	}

	data := final_data
	//data = append(data, byte(0))

	iret, err := w.Write(data)
	if nil != err {
		log.Error("client_msg_handler write data 1 failed err[%s] ret %d", err.Error(), iret)
		return
	}
}

func _push_client_msg_res(err_code int32, msg_id int32, data []byte, msg_res *msg_client_message.S2C_MSG_DATA) {
	msg_res.MsgList = append(msg_res.MsgList, &msg_client_message.S2C_ONE_MSG{
		MsgCode:   msg_id,
		ErrorCode: err_code,
		Data:      data,
	})
}

func _process_one_client_msg(p *Player, msg_id int32, msg_data []byte, handlerinfo *MsgHandlerInfo, msg_res *msg_client_message.S2C_MSG_DATA) {
	if msg_id <= 0 {
		_push_client_msg_res(int32(msg_client_message.E_ERR_PLAYER_MSG_ID_INVALID), 0, nil, msg_res)
		log.Error("!!!!!! Invalid Msg Id %v from Player Id %v", msg_id, p.Id)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			atomic.StoreInt32(&p.is_lock, 0)
			log.Stack(err)
		}
	}()

	var ret_code int32
	var data []byte
	if !atomic.CompareAndSwapInt32(&p.is_lock, 0, 1) {
		log.Debug("Player[%v] send msg[%v] cant process, because prev msg is processing", p.Id, msg_id)
		ret_code = int32(msg_client_message.E_ERR_PLAYER_SEND_TOO_FREQUENTLY)
	} else {
		p.b_base_prop_chg = false
		p.OnInit()
		if atomic.LoadInt32(&p.is_login) > 0 || msg_id == int32(msg_client_message.C2SReconnectRequest_ProtoID) {
			ret_code = handlerinfo.player_msg_handler(p, msg_data)
		} else {
			ret_code = int32(msg_client_message.E_ERR_PLAYER_MUST_RECONN_WITH_DISCONN_STATE)
		}
		data = p.PopCurMsgData()
		atomic.CompareAndSwapInt32(&p.is_lock, 1, 0)
	}

	_push_client_msg_res(ret_code, msg_id, data, msg_res)
}

func _process_client_msgs(r *http.Request, p *Player, msg_list []*msg_client_message.C2S_ONE_MSG, msg_res *msg_client_message.S2C_MSG_DATA) (err int32) {
	for _, m := range msg_list {
		msg_id := m.GetMsgCode()
		handlerinfo := msg_handler_mgr.msgid2handler[msg_id]
		if nil == handlerinfo {
			_push_client_msg_res(int32(msg_client_message.E_ERR_PLAYER_MSG_ID_NOT_FOUND), 0, nil, msg_res)
			log.Error("client_msg_handler msg_handler_mgr[%d] nil ", msg_id)
			continue
		}

		msg_data := m.GetData()
		if !handlerinfo.if_player_msg {
			var ret_code int32
			var data []byte
			ret_code, p = handlerinfo.msg_handler(msg_data)
			if p != nil {
				data = p.PopCurMsgData()
				ip_port := strings.Split(r.RemoteAddr, ":")
				if len(ip_port) >= 2 {
					p.pos = position_table.GetPosByIP(ip_port[0])
				}
			} else {
				data = nil
			}
			_push_client_msg_res(ret_code, msg_id, data, msg_res)
		} else {
			if p == nil {
				err = int32(msg_client_message.E_ERR_PLAYER_NOT_FOUND_BY_TOKEN)
				log.Error("!!!!!! Process msg %v not found player", msg_id)
				break
			}
			_process_one_client_msg(p, msg_id, msg_data, handlerinfo, msg_res)
		}
	}

	if p != nil {
		if USE_CONN_TIMER_WHEEL == 0 {
			conn_timer_mgr.Insert(p.Id)
		} else {
			conn_timer_wheel.Insert(p.Id)
		}
	}

	return
}

type ResponseData struct {
	CompressType int32 // ????????????  0 ?????????  1 zlib  2 snappy
	Data         []byte
}

const (
	COMPRESS_TYPE_NONE   = iota
	COMPRESS_TYPE_ZLIB   = 1
	COMPRESS_TYPE_SNAPPY = 2
)

var g_compress_type int32 = COMPRESS_TYPE_SNAPPY

func client_msg_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			debug.PrintStack()
		}
	}()

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if nil != err {
		_send_error(w, -1)
		log.Error("client_msg_handler ReadAll err[%s]", err.Error())
		return
	}

	tmp_msg := &msg_client_message.C2S_MSG_DATA{}
	err = proto.Unmarshal(data, tmp_msg)
	if nil != err {
		_send_error(w, -1)
		log.Error("client_msg_handler proto Unmarshal err[%s]", err.Error())
		return
	}

	token := tmp_msg.GetToken()
	var access_token share_data.AccessTokenInfo
	if !access_token.ParseString(token) {
		_send_error(w, int32(msg_client_message.E_ERR_PLAYER_TOKEN_ERROR))
		log.Error("client_msg_handler parse token %v failed", token)
		return
	}

	token_info := login_token_mgr.GetTokenByUid(access_token.UniqueId)
	if token_info == nil || token_info.token != token {
		if token_info == nil {
			_send_error(w, int32(msg_client_message.E_ERR_PLAYER_TOKEN_NOT_FOUND))
			log.Warn("UniqueId[%v] no token info", access_token.UniqueId)
		} else {
			_send_error(w, int32(msg_client_message.E_ERR_PLAYER_OTHER_PLACE_LOGIN))
			log.Warn("UniqueId[%v] token[%v] invalid, need[%v]", access_token.UniqueId, tmp_msg.GetToken(), token_info.token)
		}
		return
	}

	var res2cli msg_client_message.S2C_MSG_DATA

	msg_list := tmp_msg.GetMsgList()
	if msg_list != nil {
		p := player_mgr.GetPlayerByUid(access_token.UniqueId)
		e := _process_client_msgs(r, p, msg_list, &res2cli)
		if e < 0 {
			_send_error(w, e)
			return
		}
	}

	final_data, err := proto.Marshal(&res2cli)
	if nil != err {
		_send_error(w, -1)
		log.Error("client_msg_handler marshal 2 client msg failed err(%s)", err.Error())
		return
	}

	if false {
		var ct int32
		if len(final_data) > 500 {
			if g_compress_type == COMPRESS_TYPE_ZLIB {
				var in bytes.Buffer
				wr, err := zlib.NewWriterLevel(&in, zlib.BestSpeed)
				if err != nil {
					log.Error("New zlib writer with level %v err %v", zlib.DefaultCompression, err.Error())
					return
				}
				wr.Write(final_data)
				wr.Close()
				data = in.Bytes()
			} else if g_compress_type == COMPRESS_TYPE_SNAPPY {
				data = snappy.Encode(nil, final_data)
				if data == nil {
					log.Error("Snappy encode %v nil", final_data)
					return
				}
			}
			ct = g_compress_type
			log.Debug("Compressed Data len %v from len %v", len(data), len(final_data))
		} else {
			data = final_data
		}

		data = append(data, byte(ct))
	} else {
		data = final_data
	}

	iret, err := w.Write(data)
	if nil != err {
		_send_error(w, -1)
		log.Error("client_msg_handler write data 2 failed err[%s] ret %d", err.Error(), iret)
		return
	}
}
