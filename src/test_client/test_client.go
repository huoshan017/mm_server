package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"mm_server_new/libs/log"
	"mm_server_new/libs/timer"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type TestClient struct {
	start_time         time.Time
	quit               bool
	shutdown_lock      *sync.Mutex
	shutdown_completed bool
	ticker             *timer.TickTimer
	initialized        bool
	last_heartbeat     int32
	cmd_chan           chan *msg_client_message.C2S_TEST_COMMAND
}

var test_client TestClient

func (this *TestClient) Init() (ok bool) {
	this.start_time = time.Now()
	this.shutdown_lock = &sync.Mutex{}
	this.cmd_chan = make(chan *msg_client_message.C2S_TEST_COMMAND)
	this.initialized = true

	return true
}

func (this *TestClient) Start() (err error) {

	log.Event("客户端已启动", nil)
	log.Trace("**************************************************")

	this.Run()

	return
}

func (this *TestClient) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}

		this.shutdown_completed = true
	}()

	this.ticker = timer.NewTickTimer(1000)
	this.ticker.Start()
	defer this.ticker.Stop()

	go this.SendCmd()

	for {
		select {
		case d, ok := <-this.ticker.Chan:
			{
				if !ok {
					return
				}

				this.OnTick(d)
			}
		}
	}
}

func (this *TestClient) Shutdown() {
	if !this.initialized {
		return
	}

	this.shutdown_lock.Lock()
	defer this.shutdown_lock.Unlock()

	if this.quit {
		return
	}
	this.quit = true

	log.Trace("关闭游戏主循环")

	begin := time.Now()

	if this.ticker != nil {
		this.ticker.Stop()
	}

	for {
		if this.shutdown_completed {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	log.Trace("关闭游戏主循环耗时 %v 秒", time.Now().Sub(begin).Seconds())
}

var cur_hall_conn *GameConnection

func get_res(url string) []byte {
	return nil
}

func _send_func(msg_id int32, msg_data []byte) *msg_client_message.S2C_ONE_MSG {
	var send_msg = msg_client_message.C2S_ONE_MSG{
		MsgCode: msg_id,
		Data:    msg_data,
	}

	data, err := proto.Marshal(&send_msg)
	if nil != err {
		log.Error("C2S_ONE_MSG Marshal err(%s)", err.Error())
		return nil
	}

	var resp *http.Response
	if config.UseHttps {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err = client.Post("https://"+config.LoginServerIP+"/client", "application/x-www-form-urlencoded", bytes.NewReader(data))
	} else {
		resp, err = http.Post("http://"+config.LoginServerIP+"/client", "application/x-www-form-urlencoded", bytes.NewReader(data))
	}

	if nil != err {
		log.Error("Post[%s] C2S_ONE_MSG error[%s]", config.LoginServerIP+"/client_msg", err.Error())
		return nil
	}

	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Error("Read post resp body err [%s]", err.Error())
		return nil
	}

	var resp_msg msg_client_message.S2C_ONE_MSG
	err = proto.Unmarshal(data, &resp_msg)
	if err != nil {
		log.Error("S2C_ONE_MSG unmarshal err %v", err.Error())
		return nil
	}

	return &resp_msg
}

func register_func(account, password string, is_guest int32) {
	var register_msg = msg_client_message.C2SRegisterRequest{
		Account:  account,
		Password: password,
		IsGuest: func() bool {
			if is_guest > 0 {
				return true
			} else {
				return false
			}
		}(),
	}
	data, err := proto.Marshal(&register_msg)
	if err != nil {
		log.Error("C2SRegisterRequest marshal err %v", err.Error())
		return
	}

	resp_msg := _send_func(int32(msg_client_message.C2SRegisterRequest_ProtoID), data)
	if resp_msg.GetErrorCode() < 0 {
		return
	}

	if resp_msg.GetMsgCode() != int32(msg_client_message.S2CRegisterResponse_ProtoID) {
		log.Warn("returned msg_id[%v] is not correct")
		return
	}

	var msg msg_client_message.S2CRegisterResponse
	err = proto.Unmarshal(resp_msg.GetData(), &msg)
	if err != nil {
		log.Error("unmarshal error[%v]", err.Error())
		return
	}

	log.Debug("Account[%v] registered, password is %v", msg.GetAccount(), msg.GetPassword())
}

func bind_new_account_func(server_id int32, account, password, new_account, new_password, new_channel string) {
	var bind_msg = msg_client_message.C2SGuestBindNewAccountRequest{
		ServerId:    server_id,
		Account:     account,
		Password:    password,
		NewAccount:  new_account,
		NewPassword: new_password,
		NewChannel:  new_channel,
	}
	data, err := proto.Marshal(&bind_msg)
	if err != nil {
		log.Error("C2SGuestBindNewAccountRequest marshal err %v", err.Error())
		return
	}

	resp_msg := _send_func(int32(msg_client_message.C2SGuestBindNewAccountRequest_ProtoID), data)
	if resp_msg.GetErrorCode() < 0 {
		return
	}

	if resp_msg.GetMsgCode() != int32(msg_client_message.S2CGuestBindNewAccountResponse_ProtoID) {
		log.Warn("returned msg_id[%v] is not correct")
		return
	}

	var msg msg_client_message.S2CGuestBindNewAccountResponse
	err = proto.Unmarshal(resp_msg.GetData(), &msg)
	if err != nil {
		log.Error("unmarshal error[%v]", err.Error())
		return
	}

	log.Debug("Account[%v] bind new account %v, password is %v", msg.GetAccount(), msg.GetNewAccount(), msg.GetNewPassword())
}

func login_func(account, password, channel, client_os string) {
	var login_msg = msg_client_message.C2SLoginRequest{
		Acc:      account,
		Password: password,
		Channel:  channel,
	}
	data, err := proto.Marshal(&login_msg)
	if err != nil {
		log.Error("C2SLoginRequest marshal err %v", err.Error())
		return
	}

	resp_msg := _send_func(int32(msg_client_message.C2SLoginRequest_ProtoID), data)
	if resp_msg.GetErrorCode() < 0 {
		return
	}

	if resp_msg.GetMsgCode() != int32(msg_client_message.S2CLoginResponse_ProtoID) {
		log.Warn("returned msg_id[%v] is not correct")
		return
	}

	var msg msg_client_message.S2CLoginResponse
	err = proto.Unmarshal(resp_msg.GetData(), &msg)
	if err != nil {
		log.Error("unmarshal error[%v]", err.Error())
		return
	}

	select_server_func(account, msg.GetToken(), msg.GetGameIP(), msg.PlayerInfo.GetServerId())
}

func select_server_func(account, token, game_ip string, server_id int32) {
	strs := strings.Split( /*msg.GetIP()*/ game_ip, ":")
	if len(strs) == 0 {
		log.Error("cant get game server port from ip %v" /*msg.GetIP()*/, game_ip)
		return
	}

	var ip string
	if config.UseHttps {
		ip = fmt.Sprintf("https://%v:%v", strings.Split(config.LoginServerIP, ":")[0], strs[len(strs)-1])
	} else {
		ip = fmt.Sprintf("http://%v:%v", strings.Split(config.LoginServerIP, ":")[0], strs[len(strs)-1])
	}
	cur_game_conn := new_game_connect(ip, account, token, config.UseHttps)
	cur_game_conn.server_id = server_id
	game_conn_mgr.AddGameConn(cur_game_conn)
	req2s := &msg_client_message.C2SEnterGameRequest{}
	req2s.Acc = account
	cur_game_conn.Send(uint16(msg_client_message.C2SEnterGameRequest_ProtoID), req2s)
}

func set_password_func(account, password, new_password string) {
	var pass_msg = msg_client_message.C2SSetLoginPasswordRequest{
		Account:     account,
		Password:    password,
		NewPassword: new_password,
	}
	data, err := proto.Marshal(&pass_msg)
	if err != nil {
		log.Error("C2SSetLoginPasswordRequest marshal err %v", err.Error())
		return
	}

	resp_msg := _send_func(int32(msg_client_message.C2SSetLoginPasswordRequest_ProtoID), data)
	if resp_msg.GetErrorCode() < 0 {
		return
	}

	if resp_msg.GetMsgCode() != int32(msg_client_message.S2CSetLoginPasswordResponse_ProtoID) {
		log.Warn("returned msg_id[%v] is not correct")
		return
	}

	var msg msg_client_message.S2CSetLoginPasswordResponse
	err = proto.Unmarshal(resp_msg.GetData(), &msg)
	if err != nil {
		log.Error("unmarshal error[%v]", err.Error())
		return
	}

	log.Debug("Account[%v] set password[%v] replace old password[%v]", account, new_password, password)
}

func (this *TestClient) cmd_register(use_https bool) {
	fmt.Printf("请输入账号: ")
	var acc, pwd, is_guest string
	fmt.Scanf("%s\n", &acc)
	fmt.Printf("请输入密码: ")
	fmt.Scanf("%s\n", &pwd)
	fmt.Printf("是否游客: y/n? ")
	fmt.Scanf("%s\n", &is_guest)

	var ig int32
	if is_guest == "y" || is_guest == "Y" || is_guest == "" {
		ig = 1
	}

	if config.AccountNum == 0 {
		config.AccountNum = 1
	}

	for i := int32(0); i < config.AccountNum; i++ {
		account := acc
		if config.AccountNum > 1 {
			account = fmt.Sprintf("%v_%s", acc, i)
		}

		register_func(acc, pwd, ig)

		if config.AccountNum > 1 {
			log.Debug("Account[%v] registered, total count %v", account, i+1)
		}
	}
}

func (this *TestClient) cmd_bind_new_account(use_https bool) {
	var server_id, account, password, new_account, new_password string
	fmt.Printf("输入ServerId: ")
	fmt.Scanf("%s\n", &server_id)
	fmt.Printf("输入旧帐号: ")
	fmt.Scanf("%s\n", &account)
	fmt.Printf("输入旧密码: ")
	fmt.Scanf("%s\n", &password)
	fmt.Printf("输入新账号: ")
	fmt.Scanf("%s\n", &new_account)
	fmt.Printf("输入新密码: ")
	fmt.Scanf("%s\n", &new_password)

	for i := int32(0); i < config.AccountNum; i++ {
		acc := account
		if config.AccountNum > 1 {
			acc = fmt.Sprintf("%v_%s", acc, i)
		}
		sid, err := strconv.Atoi(server_id)
		if err != nil {
			continue
		}
		bind_new_account_func(int32(sid), account, password, new_account, new_password, "")
		if config.AccountNum > 1 {
			log.Debug("Account[%v] bind new account %v, total count %v", account, new_account, i+1)
		}
	}
}

func (this *TestClient) cmd_login(use_https bool) {
	var acc, pwd, chl, cos string
	fmt.Printf("请输入账号: ")
	fmt.Scanf("%s\n", &acc)
	fmt.Printf("请输入密码: ")
	fmt.Scanf("%s\n", &pwd)
	fmt.Printf("请输入渠道: ")
	fmt.Scanf("%s\n", &chl)
	fmt.Printf("请输入客户端系统: ")
	fmt.Scanf("%s\n", &cos)
	cur_game_conn = game_conn_mgr.GetGameConnByAcc(acc)
	if nil != cur_game_conn && cur_game_conn.blogin {
		log.Info("[%s] already login", acc)
		return
	}

	if config.AccountNum == 0 {
		config.AccountNum = 1
	}
	for i := int32(0); i < config.AccountNum; i++ {
		account := acc
		if config.AccountNum > 1 {
			account = fmt.Sprintf("%s_%v", acc, i)
		}

		login_func(account, pwd, chl, cos)

		if config.AccountNum > 1 {
			log.Debug("Account[%v] logined, total count[%v]", account, i+1)
		}
	}
}

func (this *TestClient) cmd_set_password(use_https bool) {
	var acc, pwd, new_pwd string
	fmt.Printf("请输入账号: ")
	fmt.Scanf("%s\n", &acc)
	fmt.Printf("请输入密码: ")
	fmt.Scanf("%s\n", &pwd)
	fmt.Printf("请输入新密码: ")
	fmt.Scanf("%s\n", &new_pwd)
	cur_game_conn = game_conn_mgr.GetGameConnByAcc(acc)
	if nil != cur_game_conn && cur_hall_conn.blogin {
		log.Info("[%s] already login", acc)
		return
	}

	if config.AccountNum == 0 {
		config.AccountNum = 1
	}
	for i := int32(0); i < config.AccountNum; i++ {
		account := acc
		if config.AccountNum > 1 {
			account = fmt.Sprintf("%s_%v", acc, i)
		}

		set_password_func(account, pwd, new_pwd)

		if config.AccountNum > 1 {
			log.Debug("Account[%v] set password, total count[%v]", account, i+1)
		}
	}
}

var is_test bool

func (this *TestClient) OnTick(t timer.TickTime) {
	if !is_test {
		fmt.Printf("请输入命令:\n")
		var cmd_str string
		fmt.Scanf("%s\n", &cmd_str)
		switch cmd_str {
		case "register":
			{
				this.cmd_register(true)
			}
		case "bind_new_account":
			{
				this.cmd_bind_new_account(true)
			}
		case "login":
			{
				this.cmd_login(true)
				is_test = true
			}
		case "set_password":
			{
				this.cmd_set_password(true)
			}
		case "enter_test":
			{
				is_test = true
			}
		}
	} else {
		fmt.Printf("请输入测试命令:\n")
		var cmd_str string
		fmt.Scanln(&cmd_str, "\n")
		switch cmd_str {
		case "leave_test":
			{
				is_test = false
			}
		default:
			{
				if cmd_str != "" {
					strs := strings.Split(cmd_str, ",")
					fmt.Printf("strs[%v] length is %v\n", strs, len(strs))
					if len(strs) == 1 {
						//fmt.Printf("命令[%v]参数不够，至少一个\n", strs[0])
						//return
					} else if len(strs) == 0 {
						fmt.Printf("没有输入命令\n")
						return
					}
					req := &msg_client_message.C2S_TEST_COMMAND{}
					req.Cmd = strs[0]
					if len(strs) > 1 {
						req.Args = strs[1:]
					} else {
						req.Args = make([]string, 0)
					}
					this.cmd_chan <- req
				}
			}
		}
	}
	this._heartbeat()
}

func (this *TestClient) _heartbeat() {
	now_time := int32(time.Now().Unix())
	if this.last_heartbeat == 0 {
		this.last_heartbeat = now_time
	}
	if now_time-this.last_heartbeat >= 50 {
		var heartbeat msg_client_message.C2SHeartbeat
		if config.AccountNum > 1 {
			for i := int32(0); i < config.AccountNum; i++ {
				if game_conn_mgr.acc_arr == nil || len(game_conn_mgr.acc_arr) < int(i)+1 {
					break
				}
				c := game_conn_mgr.acc_arr[i]
				if c != nil {
					c.Send(uint16(msg_client_message.C2SHeartbeat_ProtoID), &heartbeat)
				}
			}
		} else {
			if cur_game_conn != nil {
				cur_game_conn.Send(uint16(msg_client_message.C2SHeartbeat_ProtoID), &heartbeat)
			}
		}
		this.last_heartbeat = now_time
	}
}

func (this *TestClient) _cmd(cmd *msg_client_message.C2S_TEST_COMMAND) {
	if config.AccountNum > 1 {
		log.Debug("############## hall conns length %v, config.AccountNum %v", len(game_conn_mgr.acc_arr), config.AccountNum)
		for i := int32(0); i < config.AccountNum; i++ {
			c := game_conn_mgr.acc_arr[i]
			if c == nil {
				continue
			}
			go func(conn *GameConnection) {
				defer func() {
					if err := recover(); err != nil {
						log.Stack(err)
					}

					this.shutdown_completed = true
				}()
				conn.Send(uint16(msg_client_message.C2S_TEST_COMMAND_ProtoID), cmd)
			}(c)
		}
	} else {
		if cur_game_conn == nil {
			log.Error("hall connection is not estabulished")
			return
		}
		cur_game_conn.Send(uint16(msg_client_message.C2S_TEST_COMMAND_ProtoID), cmd)
	}
}

// 发送消息
func (this *TestClient) SendCmd() {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()
	for {
		is_break := false
		for !is_break {
			select {
			case cmd, ok := <-this.cmd_chan:
				{
					if !ok {
						log.Error("cmd chan receive invalid !!!!!")
						return
					}
					this._cmd(cmd)
				}
			default:
				{
					is_break = true
				}
			}
		}

		this._heartbeat()
		time.Sleep(time.Second * 1)
	}
}

//=================================================================================
