package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mm_server_new/libs/log"
	"mm_server_new/libs/timer"
	"mm_server_new/libs/utils"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	msg_server_message "mm_server_new/proto/gen_go/server_message"
	"mm_server_new/src/server_config"
	"mm_server_new/src/share_data"
	"net"
	"net/http"
	"regexp"
	"runtime/debug"

	//"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

type WaitCenterInfo struct {
	res_chan    chan *msg_server_message.C2LPlayerAccInfo
	create_time int32
}

type LoginServer struct {
	start_time         time.Time
	quit               bool
	shutdown_lock      *sync.Mutex
	shutdown_completed bool
	ticker             *timer.TickTimer
	initialized        bool

	login_http_listener net.Listener
	login_http_server   http.Server
	use_https           bool

	redis_conn *utils.RedisConn

	acc2c_wait      map[string]*WaitCenterInfo
	acc2c_wait_lock *sync.RWMutex
}

var server *LoginServer

func (this *LoginServer) Init() (ok bool) {
	this.start_time = time.Now()
	this.shutdown_lock = &sync.Mutex{}
	this.acc2c_wait = make(map[string]*WaitCenterInfo)
	this.acc2c_wait_lock = &sync.RWMutex{}
	this.redis_conn = &utils.RedisConn{}
	account_mgr_init()

	if db_use_new {
		//account_record_mgr.OnInitSelectRecords(select_account_records_func)
		//account_record_mgr.RegisterSelectRecordFunc(select_account_record_func)
		//ban_mgr.OnInitSelectRecords(ban_select_records)
		//ban_mgr.RegisterSelectRecordFunc(ban_select_record)
	}

	this.initialized = true

	return true
}

func (this *LoginServer) Start(use_https bool) bool {
	if !this.redis_conn.Connect(config.RedisServerIP) {
		return false
	}

	go server_list.Run()

	if use_https {
		go this.StartHttps(server_config.GetConfPathFile("server.crt"), server_config.GetConfPathFile("server.key"))
	} else {
		go this.StartHttp()
	}

	this.use_https = use_https
	log.Event("??????????????????", nil, log.Property{"IP", config.ListenClientIP})
	log.Trace("**************************************************")

	this.Run()

	return true
}

func (this *LoginServer) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}

		this.shutdown_completed = true
	}()

	this.ticker = timer.NewTickTimer(1000)
	this.ticker.Start()
	defer this.ticker.Stop()

	go this.redis_conn.Run(100)

	for {
		select {
		case d, ok := <-this.ticker.Chan:
			{
				if !ok {
					return
				}

				begin := time.Now()
				this.OnTick(d)
				time_cost := time.Now().Sub(begin).Seconds()
				if time_cost > 1 {
					log.Trace("?????? %v", time_cost)
					if time_cost > 30 {
						log.Error("?????? %v", time_cost)
					}
				}
			}
		}
	}
}

func (this *LoginServer) Shutdown() {
	if !this.initialized {
		return
	}

	this.shutdown_lock.Lock()
	defer this.shutdown_lock.Unlock()

	if this.quit {
		return
	}
	this.quit = true

	this.redis_conn.Close()

	log.Trace("?????????????????????")

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

	this.login_http_listener.Close()
	center_conn.ShutDown()
	game_agent_manager.net.Shutdown()

	if !db_use_new {
		dbc.Save(false)
		dbc.Shutdown()
	} else {
		//db_new.Save()
	}

	log.Trace("??????????????????????????? %v ???", time.Now().Sub(begin).Seconds())
}

func (this *LoginServer) OnTick(t timer.TickTime) {
}

func (this *LoginServer) add_to_c_wait(acc string, c_wait *WaitCenterInfo) {
	this.acc2c_wait_lock.Lock()
	defer this.acc2c_wait_lock.Unlock()

	this.acc2c_wait[acc] = c_wait
}

func (this *LoginServer) remove_c_wait(acc string) {
	this.acc2c_wait_lock.Lock()
	defer this.acc2c_wait_lock.Unlock()

	delete(this.acc2c_wait, acc)
}

func (this *LoginServer) get_c_wait_by_acc(acc string) *WaitCenterInfo {
	this.acc2c_wait_lock.RLock()
	defer this.acc2c_wait_lock.RUnlock()

	return this.acc2c_wait[acc]
}

func (this *LoginServer) pop_c_wait_by_acc(acc string) *WaitCenterInfo {
	this.acc2c_wait_lock.Lock()
	defer this.acc2c_wait_lock.Unlock()

	cur_wait := this.acc2c_wait[acc]
	if nil != cur_wait {
		delete(this.acc2c_wait, acc)
		return cur_wait
	}

	return nil
}

//=================================================================================

type LoginHttpHandle struct{}

func (this *LoginServer) StartHttp() bool {
	var err error
	this.reg_http_mux()

	this.login_http_listener, err = net.Listen("tcp", config.ListenClientIP)
	if nil != err {
		log.Error("LoginServer StartHttp Failed %s", err.Error())
		return false
	}

	login_http_server := http.Server{
		Handler:     &LoginHttpHandle{},
		ReadTimeout: 6 * time.Second,
	}

	err = login_http_server.Serve(this.login_http_listener)
	if err != nil {
		log.Error("??????Login Http Server %s", err.Error())
		return false
	}

	return true
}

func (this *LoginServer) StartHttps(crt_file, key_file string) bool {
	this.reg_http_mux()

	this.login_http_server = http.Server{
		Addr:        config.ListenClientIP,
		Handler:     &LoginHttpHandle{},
		ReadTimeout: 6 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	err := this.login_http_server.ListenAndServeTLS(crt_file, key_file)
	if err != nil {
		log.Error("??????https server error[%v]", err.Error())
		return false
	}

	return true
}

var login_http_mux map[string]func(http.ResponseWriter, *http.Request)

func (this *LoginServer) reg_http_mux() {
	login_http_mux = make(map[string]func(http.ResponseWriter, *http.Request))
	login_http_mux["/client"] = client_http_handler
}

func (this *LoginHttpHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var act_str, url_str string
	url_str = r.URL.String()
	idx := strings.Index(url_str, "?")
	if -1 == idx {
		act_str = url_str
	} else {
		act_str = string([]byte(url_str)[:idx])
	}
	log.Info("ServeHTTP url(%v) actstr(%s)", url_str, act_str)
	if h, ok := login_http_mux[act_str]; ok {
		h(w, r)
	}
	return
}

type JsonRequestData struct {
	MsgId   int32  // ??????ID
	MsgData []byte // ?????????
}

type JsonResponseData struct {
	Code    int32  // ?????????
	MsgId   int32  // ??????ID
	MsgData []byte // ?????????
}

func _check_register(account, password string) (err_code int32) {
	if b, err := regexp.MatchString(`^[\.a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, account); !b {
		if err != nil {
			log.Error("account[%v] not valid account, err %v", account, err.Error())
		} else {
			log.Error("account[%v] not match", account)
		}
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_IS_INVALID)
		return
	}

	if password == "" {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_PASSWORD_INVALID)
		return
	}

	err_code = 1
	return
}

func _generate_account_uuid(account string) string {
	uid := uuid.NewV1()
	return uid.String()
}

func register_handler(account, password string, is_guest bool) (err_code int32, resp_data []byte) {
	if len(account) > 128 {
		log.Error("Account[%v] length %v too long", account, len(account))
		return -1, nil
	}

	if len(password) > 32 {
		log.Error("Account[%v] password[%v] length %v too long", account, password, len(password))
		return -1, nil
	}

	if dbc.Accounts.GetRow(account) != nil {
		log.Error("Account[%v] already exists", account)
		return int32(msg_client_message.E_ERR_ACCOUNT_ALREADY_REGISTERED), nil
	}

	if !is_guest {
		err_code = _check_register(account, password)
		if err_code < 0 {
			return
		}
	}

	uid := _generate_account_uuid(account)
	if uid == "" {
		err_code = -1
		return
	}

	row := dbc.Accounts.AddRow(account)
	if row == nil {
		err_code = -1
		return
	}
	row.SetUniqueId(uid)
	row.SetPassword(password)
	row.SetRegisterTime(int32(time.Now().Unix()))
	if is_guest {
		row.SetChannel("guest")
	}

	var response msg_client_message.S2CRegisterResponse = msg_client_message.S2CRegisterResponse{
		Account:  account,
		Password: password,
		IsGuest:  is_guest,
	}

	var err error
	resp_data, err = proto.Marshal(&response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("login_handler marshal response error: %v", err.Error())
		return
	}

	log.Debug("Account[%v] password[%v] registered", account, password)

	err_code = 1
	return
}

func bind_new_account_handler(server_id int32, account, password, new_account, new_password, new_channel string) (err_code int32, resp_data []byte) {
	if len(new_account) > 128 {
		log.Error("Account[%v] length %v too long", new_account, len(new_account))
		return -1, nil
	}

	if new_channel != "facebook" && len(new_password) > 32 {
		log.Error("Account[%v] password[%v] length %v too long", new_account, new_password, len(new_password))
		return -1, nil
	}

	if account == new_account {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NAME_MUST_DIFFRENT_TO_OLD)
		log.Error("Account %v can not bind same new account", account)
		return
	}

	row := dbc.Accounts.GetRow(account)
	if row == nil {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NOT_REGISTERED)
		log.Error("Account %v not registered, cant bind new account", account)
		return
	}

	ban_row := dbc.BanPlayers.GetRow(row.GetUniqueId())
	if ban_row != nil && ban_row.GetStartTime() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	channel := row.GetChannel()
	if channel != "guest" && channel != "facebook" {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NOT_GUEST)
		log.Error("Account %v not guest and not facebook user", account)
		return
	}

	if channel == "facebook" {
		err_code = _verify_facebook_login(account, password)
		if err_code < 0 {
			return
		}
	}

	if channel != "facebook" && row.GetBindNewAccount() != "" {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_ALREADY_BIND)
		log.Error("Account %v already bind", account)
		return
	}

	if dbc.Accounts.GetRow(new_account) != nil {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NEW_BIND_ALREADY_EXISTS)
		log.Error("New Account %v to bind already exists", new_account)
		return
	}

	if new_channel != "" {
		if new_channel == "facebook" {
			err_code = _verify_facebook_login(new_account, new_password)
			if err_code < 0 {
				return
			}
		} else {
			err_code = -1
			log.Error("Account %v bind a unsupported channel %v account %v", account, new_channel, new_account)
			return
		}
	} else {
		err_code = _check_register(new_account, new_password)
		if err_code < 0 {
			return
		}
	}

	row.SetBindNewAccount(new_account)
	register_time := row.GetRegisterTime()
	uid := row.GetUniqueId()

	last_server_id := row.GetServerId()

	row = dbc.Accounts.AddRow(new_account)
	if row == nil {
		err_code = -1
		log.Error("Account %v bind new account %v database error", account, new_account)
		return
	}

	if new_channel == "" {
		row.SetPassword(new_password)
	}
	row.SetRegisterTime(register_time)
	row.SetUniqueId(uid)
	row.SetOldAccount(account)
	row.SetServerId(last_server_id)

	//dbc.Accounts.RemoveRow(account) // ???????????????

	game_agent := game_agent_manager.GetAgentByID(server_id)
	if nil == game_agent {
		err_code = int32(msg_client_message.E_ERR_PLAYER_SELECT_SERVER_NOT_FOUND)
		log.Error("login_http_handler get hall_agent failed")
		return
	}

	req := &msg_server_message.L2GBindNewAccountRequest{
		UniqueId:   uid,
		Account:    account,
		NewAccount: new_account,
	}
	game_agent.Send(uint16(msg_server_message.MSGID_L2G_BIND_NEW_ACCOUNT_REQUEST), req)

	response := &msg_client_message.S2CGuestBindNewAccountResponse{
		Account:     account,
		NewAccount:  new_account,
		NewPassword: new_password,
		NewChannel:  new_channel,
	}

	var err error
	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("login_handler marshal response error: %v", err.Error())
		return
	}

	log.Debug("Account[%v] bind new account[%v]", account, new_account)
	err_code = 1
	return
}

func _verify_facebook_login(user_id, input_token string) int32 {
	type _facebook_data struct {
		AppID     string `json:"app_id"`
		IsValid   bool   `json:"is_valid"`
		UserID    string `json:"user_id"`
		IssuedAt  int    `json:"issued_at"`
		ExpiresAt int    `json:"expires_at"`
	}

	type facebook_data struct {
		Data _facebook_data `json:"data"`
	}

	var resp *http.Response
	var err error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url_str := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%v&access_token=%v|%v", input_token, config.FacebookAppID, config.FacebookAppSecret)
	log.Debug("verify facebook url: %v", url_str)

	client := &http.Client{Transport: tr}
	resp, err = client.Get(url_str)
	if nil != err {
		log.Error("Facebook verify error %s", err.Error())
		return -1
	}

	if resp.StatusCode != 200 {
		log.Error("Facebook verify response code %v", resp.StatusCode)
		return -1
	}

	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Error("Read facebook verify result err(%s) !", err.Error())
		return -1
	}

	log.Debug("facebook verify result data: %v", string(data))

	var fdata facebook_data
	err = json.Unmarshal(data, &fdata)
	if nil != err {
		log.Error("Facebook verify ummarshal err(%s)", err.Error())
		return -1
	}

	if !fdata.Data.IsValid {
		log.Error("Facebook verify input_token[%v] failed", input_token)
		return -1
	}

	if fdata.Data.UserID != user_id {
		log.Error("Facebook verify client user_id[%v] different to result user_id[%v]", user_id, fdata.Data.UserID)
		return -1
	}

	log.Debug("Facebook verify user_id[%v] and input_token[%v] success", user_id, input_token)

	return 1
}

func login_handler(account, password, channel string) (err_code int32, resp_data []byte) {
	now_time := time.Now()
	var err error
	acc_row := dbc.Accounts.GetRow(account)
	if config.VerifyAccount {
		if channel == "" {
			if acc_row == nil {
				err_code = int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR)
				log.Error("Account %v not exist", account)
				return
			}
			if acc_row.GetPassword() != password {
				err_code = int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR)
				log.Error("Account %v password %v invalid", account, password)
				return
			}
		} else if channel == "facebook" {
			err_code = _verify_facebook_login(account, password)
			if err_code < 0 {
				return
			}
			if acc_row == nil {
				acc_row = dbc.Accounts.AddRow(account)
				if acc_row == nil {
					log.Error("Account %v add row with channel facebook failed")
					return -1, nil
				}
				acc_row.SetChannel("facebook")
				acc_row.SetRegisterTime(int32(now_time.Unix()))
			}
			acc_row.SetPassword(password)
		} else if channel == "guest" {
			if acc_row == nil {
				acc_row = dbc.Accounts.AddRow(account)
				if acc_row == nil {
					log.Error("Account %v add row with channel guest failed")
					return -1, nil
				}
				acc_row.SetChannel("guest")
				acc_row.SetRegisterTime(int32(now_time.Unix()))
			} else {
				if acc_row.GetPassword() != password {
					err_code = int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR)
					log.Error("Account %v password %v invalid", account, password)
					return
				}
			}
		} else {
			log.Error("Account %v use unsupported channel %v login", account, channel)
			return -1, nil
		}
	} else {
		if acc_row == nil {
			acc_row = dbc.Accounts.AddRow(account)
			if acc_row == nil {
				log.Error("Account %v add row without verify failed")
				return -1, nil
			}
			acc_row.SetRegisterTime(int32(now_time.Unix()))
		}
	}

	if acc_row.GetUniqueId() == "" {
		uid := _generate_account_uuid(account)
		if uid != "" {
			acc_row.SetUniqueId(uid)
		}
	}

	ban_row := dbc.BanPlayers.GetRow(acc_row.GetUniqueId())
	if ban_row != nil && ban_row.GetStartTime() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	// --------------------------------------------------------------------------------------------
	// ???????????????
	server_id := acc_row.GetServerId()
	if server_id <= 0 {
		server := server_list.RandomOneServer()
		if server == nil {
			err_code = int32(msg_client_message.E_ERR_INTERNAL)
			log.Error("Server List random null !!!")
			return
		}
		server_id = server.Id
		acc_row.SetServerId(server_id)
	}

	var hall_ip, token string
	err_code, hall_ip, token = _select_server(acc_row.GetUniqueId(), account, server_id)
	if err_code < 0 {
		return
	}
	// --------------------------------------------------------------------------------------------

	account_login(account, token, "")

	acc_row.SetToken(token)

	response := &msg_client_message.S2CLoginResponse{
		Acc:    account,
		Token:  token,
		GameIP: hall_ip,
	}

	if channel == "guest" {
		response.BoundAccount = acc_row.GetBindNewAccount()
	}

	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("login_handler marshal response error: %v", err.Error())
		return
	}

	log.Debug("Account[%v] channel[%v] logined", account, channel)

	return
}

func _select_server(unique_id, account string, server_id int32) (err_code int32, game_ip, access_token string) {
	sinfo := server_list.GetServerById(server_id)
	if sinfo == nil {
		err_code = int32(msg_client_message.E_ERR_PLAYER_SELECT_SERVER_NOT_FOUND)
		log.Error("select_server_handler player[%v] select server[%v] not found")
		return
	}

	game_agent := game_agent_manager.GetAgentByID(server_id)
	if nil == game_agent {
		err_code = int32(msg_client_message.E_ERR_PLAYER_SELECT_SERVER_NOT_FOUND)
		log.Error("login_http_handler account %v get game_agent failed by server_id %v", account, server_id)
		return
	}

	access_token = share_data.GenerateAccessToken(unique_id)
	game_agent.Send(uint16(msg_server_message.MSGID_L2G_SYNC_ACCOUNT_TOKEN), &msg_server_message.L2GSyncAccountToken{
		UniqueId: unique_id,
		Account:  account,
		Token:    access_token,
	})

	game_ip = sinfo.IP

	err_code = 1

	return
}

func set_password_handler(account, password, new_password string) (err_code int32, resp_data []byte) {
	row := dbc.Accounts.GetRow(account)
	if row == nil {
		err_code = int32(msg_client_message.E_ERR_PLAYER_NOT_EXIST)
		log.Error("set_password_handler account[%v] not found", account)
		return
	}

	ban_row := dbc.BanPlayers.GetRow(row.GetUniqueId())
	if ban_row != nil && ban_row.GetStartTime() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	if row.GetPassword() != password {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_PASSWORD_INVALID)
		log.Error("set_password_handler account[%v] password is invalid", account)
		return
	}

	row.SetPassword(new_password)

	response := &msg_client_message.S2CSetLoginPasswordResponse{
		Account:     account,
		Password:    password,
		NewPassword: new_password,
	}

	var err error
	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("set_password_handler marshal response error: %v", err.Error())
		return
	}

	return
}

func response_error(err_code int32, w http.ResponseWriter) {
	err_response := JsonResponseData{
		Code: err_code,
	}
	data, err := json.Marshal(err_response)
	if nil != err {
		log.Error("login_http_handler json mashal error")
		return
	}
	w.Write(data)
}

func _send_error(w http.ResponseWriter, msg_id, ret_code int32) {
	m := &msg_client_message.S2C_ONE_MSG{ErrorCode: ret_code}
	res2cli := &msg_client_message.S2C_MSG_DATA{MsgList: []*msg_client_message.S2C_ONE_MSG{m}}
	final_data, err := proto.Marshal(res2cli)
	if nil != err {
		log.Error("client_msg_handler marshal 1 client msg failed err(%s)", err.Error())
		return
	}

	data := final_data
	data = append(data, byte(0))

	iret, err := w.Write(data)
	if nil != err {
		log.Error("client_msg_handler write data 1 failed err[%s] ret %d", err.Error(), iret)
		return
	}
}

func client_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			debug.PrintStack()
		}
	}()

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_http_handler ReadAll err[%s]", err.Error())
		return
	}

	var msg msg_client_message.C2S_ONE_MSG
	err = proto.Unmarshal(data, &msg)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_http_handler proto Unmarshal err[%s]", err.Error())
		return
	}

	var err_code, msg_id int32
	if msg.MsgCode == int32(msg_client_message.C2SLoginRequest_ProtoID) {
		var login_msg msg_client_message.C2SLoginRequest
		err = proto.Unmarshal(msg.GetData(), &login_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SLoginRequest unmarshal err %v", err.Error())
			return
		}
		if login_msg.GetAcc() == "" {
			_send_error(w, 0, -1)
			log.Error("Acc is empty")
			return
		}
		msg_id = int32(msg_client_message.S2CLoginResponse_ProtoID)
		if db_use_new {
			//err_code, data = new_login_handler(login_msg.GetAcc(), login_msg.GetPassword(), login_msg.GetChannel())
		} else {
			err_code, data = login_handler(login_msg.GetAcc(), login_msg.GetPassword(), login_msg.GetChannel())
		}
	} else if msg.MsgCode == int32(msg_client_message.C2SRegisterRequest_ProtoID) {
		var register_msg msg_client_message.C2SRegisterRequest
		err = proto.Unmarshal(msg.GetData(), &register_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SRegisterRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message.S2CRegisterResponse_ProtoID)
		if db_use_new {
			//err_code, data = new_register_handler(register_msg.GetAccount(), register_msg.GetPassword(), register_msg.GetIsGuest())
		} else {
			err_code, data = register_handler(register_msg.GetAccount(), register_msg.GetPassword(), register_msg.GetIsGuest())
		}
	} else if msg.MsgCode == int32(msg_client_message.C2SSetLoginPasswordRequest_ProtoID) {
		var pass_msg msg_client_message.C2SSetLoginPasswordRequest
		err = proto.Unmarshal(msg.GetData(), &pass_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SSetLoginPasswordRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message.S2CSetLoginPasswordResponse_ProtoID)
		if db_use_new {
			//err_code, data = new_set_password_handler(pass_msg.GetAccount(), pass_msg.GetPassword(), pass_msg.GetNewPassword())
		} else {
			err_code, data = set_password_handler(pass_msg.GetAccount(), pass_msg.GetPassword(), pass_msg.GetNewPassword())
		}
	} else if msg.MsgCode == int32(msg_client_message.C2SGuestBindNewAccountRequest_ProtoID) {
		var bind_msg msg_client_message.C2SGuestBindNewAccountRequest
		err = proto.Unmarshal(msg.GetData(), &bind_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SGuestBindNewAccountRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message.S2CGuestBindNewAccountResponse_ProtoID)
		if db_use_new {
			//err_code, data = new_bind_new_account_handler(bind_msg.GetServerId(), bind_msg.GetAccount(), bind_msg.GetPassword(), bind_msg.GetNewAccount(), bind_msg.GetNewPassword(), bind_msg.GetNewChannel())
		} else {
			err_code, data = bind_new_account_handler(bind_msg.GetServerId(), bind_msg.GetAccount(), bind_msg.GetPassword(), bind_msg.GetNewAccount(), bind_msg.GetNewPassword(), bind_msg.GetNewChannel())
		}
	} else {
		if msg.MsgCode > 0 {
			_send_error(w, msg.MsgCode, int32(msg_client_message.E_ERR_PLAYER_MSG_ID_NOT_FOUND))
			log.Error("Unsupported msg %v", msg.MsgCode)
		} else {
			_send_error(w, msg.MsgCode, int32(msg_client_message.E_ERR_PLAYER_MSG_ID_INVALID))
			log.Error("Invalid msg %v", msg.MsgCode)
		}
		return
	}

	var resp_msg msg_client_message.S2C_ONE_MSG
	resp_msg.MsgCode = msg_id
	resp_msg.ErrorCode = err_code
	resp_msg.Data = data
	data, err = proto.Marshal(&resp_msg)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_msg_handler marshal 2 client msg failed err(%s)", err.Error())
		return
	}

	iret, err := w.Write(data)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_msg_handler write data 2 failed err[%s] ret %d", err.Error(), iret)
		return
	}
}
