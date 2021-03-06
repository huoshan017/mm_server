package main

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mm_server_new/libs/log"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
)

const (
	GAME_CONN_STATE_DISCONNECT  = 0
	GAME_CONN_STATE_CONNECTED   = 1
	GAME_CONN_STATE_FORCE_CLOSE = 2
)

// ========================================================================================

type GameConnection struct {
	use_https      bool
	state          int32
	last_conn_time int32
	acc            string
	token          string
	game_ip        string
	playerid       int32

	blogin bool

	last_send_time int64
	server_id      int32
}

var cur_game_conn *GameConnection

func new_game_connect(game_ip, acc, token string, use_https bool) *GameConnection {
	ret_conn := &GameConnection{}
	ret_conn.acc = acc
	ret_conn.game_ip = game_ip
	ret_conn.token = token
	ret_conn.use_https = use_https

	log.Info("new game connection to ip %v", game_ip)

	return ret_conn
}

type ResponseData struct {
	CompressType int32
	Data         []byte
}

func (this *GameConnection) Send(msg_id uint16, msg proto.Message) {
	data, err := proto.Marshal(msg)
	if nil != err {
		log.Error("login unmarshal failed err[%s]", err.Error())
		return
	}

	C2S_MSG := &msg_client_message.C2S_MSG_DATA{}
	C2S_MSG.Token = this.token
	one_msg := &msg_client_message.C2S_ONE_MSG{
		MsgCode: int32(msg_id),
		Data:    data,
	}
	C2S_MSG.MsgList = []*msg_client_message.C2S_ONE_MSG{one_msg}

	data, err = proto.Marshal(C2S_MSG)
	if nil != err {
		log.Error("login C2S_MSG Marshal err(%s) !", err.Error())
		return
	}

	var resp *http.Response
	if this.use_https {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err = client.Post(this.game_ip+"/client_msg", "application/x-www-form-urlencoded", bytes.NewReader(data))
	} else {
		resp, err = http.Post(this.game_ip+"/client_msg", "application/x-www-form-urlencoded", bytes.NewReader(data))
	}
	if nil != err {
		log.Error("login C2S_MSG http post[%s] error[%s]", this.game_ip+"/client_msg", err.Error())
		return
	}

	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Error("HallConnection Send read resp body err [%s]", err.Error())
		return
	}

	//log.Info("???????????????????????? ??????[%v] ??????[%v]", len(data), data)
	if len(data) < 1 {
		log.Error("data length %v invalid", len(data))
		return
	}

	var d []byte
	if false {
		compress_type := int32(data[len(data)-1])
		if len(data) == 1 {
			d = []byte{}
		} else {
			d = data[:len(data)-1]
			if compress_type != 0 {
				if compress_type == 1 {
					in := bytes.NewBuffer(d)
					r, e := zlib.NewReader(in)
					if e != nil {
						log.Error("Zlib New Reader err %v", e.Error())
						return
					}
					var out bytes.Buffer
					io.Copy(&out, r)
					d = out.Bytes()
				} else if compress_type == 2 {
					d, err = snappy.Decode(nil, d)
					if err != nil {
						log.Error("Snappy decode %v err %v", d, err.Error())
						return
					}
				} else {
					log.Error("Compress type %v not supported", compress_type)
					return
				}
			}
		}
	} else {
		d = data
	}

	S2C_MSG := &msg_client_message.S2C_MSG_DATA{}
	err = proto.Unmarshal(d, S2C_MSG)
	if nil != err {
		log.Error("HallConnection unmarshal resp data err(%s) !", err.Error())
		return
	}

	if S2C_MSG.MsgList == nil {
		log.Warn("Server return empty msg list !!!")
		return
	}

	for _, m := range S2C_MSG.MsgList {
		if m.GetErrorCode() < 0 {
			log.Error("????????????????????????[%d]", m.GetErrorCode())
			return
		}

		var msg_code uint16
		var cur_len, sub_len int32
		total_data_len := int32(len(m.Data))
		for cur_len < total_data_len {
			msg_code = uint16(m.Data[cur_len])<<8 + uint16(m.Data[cur_len+1])
			sub_len = int32(m.Data[cur_len+2])<<8 + int32(m.Data[cur_len+3])
			sub_data := m.Data[cur_len+4 : cur_len+4+sub_len]
			cur_len = cur_len + 4 + sub_len

			handler_info := msg_handler_mgr.msgid2handler[int32(msg_code)]
			if nil == handler_info {
				continue
			}

			new_msg := game_conn_msgid2msg(msg_code)
			if new_msg == nil {
				return
			}
			log.Trace("??????[%d:%s]?????????????????????%v:[%s]", this.playerid, this.acc, msg_code, new_msg.String())
			err = proto.Unmarshal(sub_data, new_msg)
			if nil != err {
				log.Error("HallConnection failed unmarshal msg data !", msg_code)
				return
			}

			handler_info(this, new_msg)
		}
	}

	return
}

//========================================================================

type CLIENT_MSG_HANDLER func(*GameConnection, proto.Message)

type NEW_MSG_FUNC func() proto.Message

type MsgHandlerMgr struct {
	msgid2handler map[int32]CLIENT_MSG_HANDLER
}

var msg_handler_mgr MsgHandlerMgr

func (this *MsgHandlerMgr) Init() bool {
	this.msgid2handler = make(map[int32]CLIENT_MSG_HANDLER)
	this.RegisterMsgHandler()
	return true
}

func (this *MsgHandlerMgr) SetMsgHandler(msg_code uint16, msg_handler CLIENT_MSG_HANDLER) {
	log.Info("set msg [%d] handler !", msg_code)
	this.msgid2handler[int32(msg_code)] = msg_handler
}

func (this *MsgHandlerMgr) RegisterMsgHandler() {
	this.SetMsgHandler(uint16(msg_client_message.S2CEnterGameResponse_ProtoID), S2CEnterGameHandler)
}

func game_conn_msgid2msg(msg_id uint16) proto.Message {
	if msg_id == uint16(msg_client_message.S2CEnterGameResponse_ProtoID) {
		return &msg_client_message.S2CEnterGameResponse{}
	} else if msg_id == uint16(msg_client_message.S2CEnterGameCompleteNotify_ProtoID) {
		return &msg_client_message.S2CEnterGameCompleteNotify{}
	} else {
		log.Error("Cant found proto message by msg_id[%v]", msg_id)
	}
	return nil
}

func S2CEnterGameHandler(game_conn *GameConnection, m proto.Message) {
	res := m.(*msg_client_message.S2CEnterGameResponse)
	cur_game_conn = game_conn_mgr.GetGameConnByAcc(res.GetAcc())
	if nil == cur_game_conn {
		log.Error("S2CLoginResponseHandler failed to get cur hall[%s]", res.GetAcc())
		return
	}

	game_conn.playerid = res.GetPlayerId()
	game_conn.blogin = true
	log.Trace("player[%v]????????????????????????????????? %v", res.GetAcc(), res)

	return
}
