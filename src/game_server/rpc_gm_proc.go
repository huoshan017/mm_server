package main

import (
	"errors"
	"fmt"
	"mm_server_new/libs/log"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	msg_server_message "mm_server_new/proto/gen_go/server_message"
	"mm_server_new/src/rpc_proto"
	"sync/atomic"
	"time"
)

// GM调用
type G2G_Proc struct {
}

func (this *G2G_Proc) Test(args *rpc_proto.GmTestCmd, result *rpc_proto.GmCommonResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	result.Res = 1

	log.Trace("@@@ G2G_Proc::Test %v", args)
	return nil
}

func (this *G2G_Proc) Anouncement(args *rpc_proto.GmAnouncementCmd, result *rpc_proto.GmCommonResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	if !system_chat_mgr.push_chat_msg(args.Content, args.RemainSeconds, 0, 0, "", 0) {
		err_str := fmt.Sprintf("@@@ G2G_Proc::Anouncement %v failed", args)
		return errors.New(err_str)
	}

	result.Res = 1

	log.Trace("@@@ G2G_Proc::Anouncement %v", args)
	return nil
}

func (this *G2G_Proc) SysMail(args *rpc_proto.GmSendSysMailCmd, result *rpc_proto.GmCommonResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	// 群发邮件
	if args.PlayerId <= 0 {
		row := dbc.SysMails.AddRow()
		if row == nil {
			log.Error("@@@ G2G_Proc::SysMail add new db row failed")
			result.Res = -1
		}
		result.Res = mail_has_subtype(args.MailTableID)
		if result.Res > 0 {
			row.SetTableId(args.MailTableID)
			row.AttachedItems.SetItemList(args.AttachItems)
			row.SetSendTime(int32(time.Now().Unix()))
			dbc.SysMailCommon.GetRow().SetCurrMailId(row.GetId())
		}
	} else {
		result.Res = RealSendMail(nil, args.PlayerId, MAIL_TYPE_SYSTEM, args.MailTableID, "", "", args.AttachItems, 0)
	}

	log.Trace("@@@ G2G_Proc::SysMail %v", args)
	return nil
}

func (this *G2G_Proc) PlayerInfo(args *rpc_proto.GmPlayerInfoCmd, result *rpc_proto.GmPlayerInfoResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	p := player_mgr.GetPlayerById(args.Id)
	if p == nil {
		result.Id = int32(msg_client_message.E_ERR_PLAYER_NOT_EXIST)
		return nil
	}

	result.Id = args.Id
	result.Account = p.db.GetAccount()
	result.UniqueId = p.db.GetUniqueId()
	result.CreateTime = p.db.Info.GetCreateUnix()
	result.IsLogin = atomic.LoadInt32(&p.is_login)
	result.LastLoginTime = p.db.Info.GetLastLogin()
	result.LogoutTime = p.db.Info.GetLastLogout()
	result.Level = p.db.GetLevel()
	result.VipLevel = p.db.Info.GetVipLvl()
	result.Gold = p.db.Info.GetGold()
	result.Diamond = p.db.Info.GetDiamond()
	result.CurStage = p.db.Info.GetCurPassMaxStage()
	items := p.db.Items.GetAllIndex()
	if items != nil {
		for _, item_id := range items {
			item_num, _ := p.db.Items.GetItemNum(item_id)
			result.Items = append(result.Items, []int32{item_id, item_num}...)
		}
	}
	cats := p.db.Cats.GetAllIndex()
	if cats != nil {
		for _, cat_id := range cats {
			cat_table_id, _ := p.db.Cats.GetCfgId(cat_id)
			coin_ability, _ := p.db.Cats.GetCoinAbility(cat_id)
			match_ability, _ := p.db.Cats.GetMatchAbility(cat_id)
			explore_ability, _ := p.db.Cats.GetExploreAbility(cat_id)
			result.Cats = append(result.Cats, []int32{cat_id, cat_table_id, coin_ability, match_ability, explore_ability}...)
		}
	}
	log.Trace("@@@ G2G_Proc::PlayerInfo %v %v", args, result)

	return nil
}

func (this *G2G_Proc) OnlinePlayerNum(args *rpc_proto.GmOnlinePlayerNumCmd, result *rpc_proto.GmOnlinePlayerNumResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	result.PlayerNum = []int32{conn_timer_wheel.GetCurrPlayerNum(), player_mgr.GetPlayersNum()}

	log.Trace("@@@ G2G_Proc::OnlinePlayerNum")

	return nil
}

func (this *G2G_Proc) MonthCardSend(args *rpc_proto.GmMonthCardSendCmd, result *rpc_proto.GmCommonResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	cards := pay_table_mgr.GetMonthCards()
	if cards == nil || len(cards) == 0 {
		log.Error("@@@ month cards is empty")
		result.Res = -1
		return nil
	}

	var found bool
	for i := 0; i < len(cards); i++ {
		if cards[i].BundleId == args.BundleId {
			found = true
			break
		}
	}

	if !found {
		log.Error("@@@ Not found month card with bundle id %v", args.BundleId)
		result.Res = -1
		return nil
	}

	p := player_mgr.GetPlayerById(args.PlayerId)
	if p == nil {
		log.Error("@@@ Month card send cant found player %v", args.PlayerId)
		result.Res = int32(msg_client_message.E_ERR_PLAYER_NOT_EXIST)
		return nil
	}

	res, _ := p._charge_with_bundle_id(args.ItemId, 0, args.BundleId, nil, nil, -1)
	if res < 0 {
		log.Error("@@@ Month card send with error %v", res)
		result.Res = res
		return nil
	}

	log.Trace("@@@ G2G_Proc::MonthCardSend %v", args)

	return nil
}

func (this *G2G_Proc) GetPlayerUniqueId(args *rpc_proto.GmGetPlayerUniqueIdCmd, result *rpc_proto.GmGetPlayerUniqueIdResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	if args.PlayerId > 0 {
		p := player_mgr.GetPlayerById(args.PlayerId)
		if p == nil {
			result.PlayerUniqueId = "Cant found player"
			log.Error("@@@ Get player %v cant found", args.PlayerId)
			return nil
		}

		result.PlayerUniqueId = p.db.GetUniqueId()
	}

	log.Trace("@@@ G2G_Proc::GetPlayerUniqueId %v", args)

	return nil
}

func (this *G2G_Proc) BanPlayer(args *rpc_proto.GmBanPlayerByUniqueIdCmd, result *rpc_proto.GmCommonResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	p := player_mgr.GetPlayerByUid(args.PlayerUniqueId)
	if p == nil {
		result.Res = int32(msg_client_message.E_ERR_PLAYER_NOT_EXIST)
		log.Error("@@@ Player cant get by unique id %v", args.PlayerUniqueId)
		return nil
	}

	p.OnLogout(true)

	row := dbc.BanPlayers.GetRow(args.PlayerUniqueId)
	if args.BanOrFree > 0 {
		now_time := time.Now()
		if row == nil {
			row = dbc.BanPlayers.AddRow(args.PlayerUniqueId)
			row.SetAccount(p.db.GetAccount())
			row.SetPlayerId(p.db.GetPlayerId())
		}
		row.SetStartTime(int32(now_time.Unix()))
		row.SetStartTimeStr(now_time.Format("2006-01-02 15:04:05"))
	} else {
		if row != nil {
			row.SetStartTime(0)
			row.SetStartTimeStr("")
		}
	}

	if args.PlayerId == p.db.GetPlayerId() {
		login_conn_mgr.Send(uint16(msg_server_message.MSGID_G2L_ACCOUNT_BAN), &msg_server_message.G2LAccountBan{
			UniqueId:  args.PlayerUniqueId,
			BanOrFree: args.BanOrFree,
			Account:   p.db.GetAccount(),
			PlayerId:  p.db.GetPlayerId(),
		})
	}

	log.Trace("@@@ G2G_Proc::BanPlayer %v", args)

	return nil
}
