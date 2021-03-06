package main

import (
	"mm_server_new/libs/log"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	"mm_server_new/src/share_data"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	DEFAULT_PLAYER_ARRAY_MAX  = 1
	PLAYER_ARRAY_MAX_ADD_STEP = 1
)

type PlayerManager struct {
	uid2players        map[string]*Player
	uid2players_locker *sync.RWMutex

	id2players        map[int32]*Player
	id2players_locker *sync.RWMutex
}

var player_mgr PlayerManager

func (this *PlayerManager) Init() bool {
	this.uid2players = make(map[string]*Player)
	this.uid2players_locker = &sync.RWMutex{}
	this.id2players = make(map[int32]*Player)
	this.id2players_locker = &sync.RWMutex{}
	return true
}

func (this *PlayerManager) GetPlayersNum() int32 {
	this.uid2players_locker.RLock()
	defer this.uid2players_locker.RUnlock()
	return int32(len(this.uid2players))
}

func (this *PlayerManager) GetPlayerById(id int32) *Player {
	this.id2players_locker.Lock()
	defer this.id2players_locker.Unlock()

	return this.id2players[id]
}

func (this *PlayerManager) GetAllPlayers() []*Player {
	this.id2players_locker.RLock()
	defer this.id2players_locker.RUnlock()

	ret_ps := make([]*Player, 0, len(this.id2players))
	for _, p := range this.id2players {
		ret_ps = append(ret_ps, p)
	}

	return ret_ps
}

func (this *PlayerManager) Add2IdMap(p *Player) {
	if nil == p {
		log.Error("Player_agent_mgr Add2IdMap p nil !")
		return
	}
	this.id2players_locker.Lock()
	defer this.id2players_locker.Unlock()

	if nil != this.id2players[p.Id] {
		log.Error("PlayerManager Add2IdMap already have player(%d)", p.Id)
	}

	this.id2players[p.Id] = p
}

func (this *PlayerManager) RemoveFromIdMap(id int32) {
	this.id2players_locker.Lock()
	defer this.id2players_locker.Unlock()

	cur_p := this.id2players[id]
	if nil != cur_p {
		delete(this.id2players, id)
	}

	return
}

func (this *PlayerManager) Add2UidMap(unique_id string, p *Player) {
	if unique_id == "" {
		return
	}

	this.uid2players_locker.Lock()
	defer this.uid2players_locker.Unlock()

	if this.uid2players[unique_id] != nil {
		log.Warn("UniqueId %v already added", unique_id)
		return
	}

	this.uid2players[unique_id] = p
}

func (this *PlayerManager) RemoveFromUidMap(unique_id string) {
	this.uid2players_locker.Lock()
	defer this.uid2players_locker.Unlock()

	delete(this.uid2players, unique_id)
}

func (this *PlayerManager) GetPlayerByUid(unique_id string) *Player {
	this.uid2players_locker.RLock()
	defer this.uid2players_locker.RUnlock()

	return this.uid2players[unique_id]
}

func (this *PlayerManager) PlayerLogout(p *Player) {
	if nil == p {
		log.Error("PlayerManager PlayerLogout p nil !")
		return
	}

	//this.RemoveFromAccMap(p.Account)
	this.RemoveFromUidMap(p.UniqueId)

	p.OnLogout(true)
}

func (this *PlayerManager) OnTick() {

}

//==============================================================================
func (this *PlayerManager) RegMsgHandler() {
	if !config.DisableTestCommand {
		msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2S_TEST_COMMAND_ProtoID), C2STestCommandHandler)
	}

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SReconnectRequest_ProtoID), C2SReconnectHandler)

	msg_handler_mgr.SetMsgHandler(uint16(msg_client_message.C2SEnterGameRequest_ProtoID), C2SEnterGameRequestHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SLeaveGameRequest_ProtoID), C2SLeaveGameRequestHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SHeartbeat_ProtoID), C2SHeartbeatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPlayerChangeNameRequest_ProtoID), C2SPlayerChangeNameHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPlayerChangeHeadRequest_ProtoID), C2SPlayerChangeHeadHandler)

	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPlayerInfoRequest_ProtoID), C2SPlayerInfoHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetInfo_ProtoID), C2SGetInfoHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetItemInfos_ProtoID), C2SGetItemInfosHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetDepotBuildingInfos_ProtoID), C2SGetDepotBuildingInfosHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetCatInfos_ProtoID), C2SGetCatInfosHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetStageInfos_ProtoID), C2SGetStageInfosHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetOptions_ProtoID), C2SGetOptionsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSaveOptions_ProtoID), C2SSaveOptionsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChgName_ProtoID), C2SChgNameHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChangeHead_ProtoID), C2SChangeHeadHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SZanPlayer_ProtoID), C2SZanPlayerHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SUseItem_ProtoID), C2SUseItemHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSellItem_ProtoID), C2SSellItemHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SComposeCat_ProtoID), C2SComposeCatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SItemResource_ProtoID), C2SItemResourceHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SShopItems_ProtoID), C2SShopItemsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SBuyShopItem_ProtoID), C2SBuyShopItemHandler)

	// ???
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatFeed_ProtoID), C2SFeedCatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatUpgradeStar_ProtoID), C2SCatUpgradeStarHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatSkillLevelUp_ProtoID), C2SCatSkillLevelUpHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatRenameNick_ProtoID), C2SCatRenameNickHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatLock_ProtoID), C2SCatLockHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatDecompose_ProtoID), C2SCatDecomposeHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPlayerCatInfo_ProtoID), C2SPlayerCatInfoHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetMakingFormulaBuildings_ProtoID), C2SGetMakingFormulaBuildingsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SExchangeBuildingFormula_ProtoID), C2SExchangeBuildingFormulaHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMakeFormulaBuilding_ProtoID), C2SMakeFormulaBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SBuyMakeBuildingSlot_ProtoID), C2SBuyMakeBuildingSlotHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSpeedupMakeBuilding_ProtoID), C2SSpeedupMakeBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetCompletedFormulaBuilding_ProtoID), C2SGetCompletedFormulaBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCancelMakingFormulaBuilding_ProtoID), C2SCancelMakingFormulaBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetFormulas_ProtoID), C2SGetFormulasHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetCrops_ProtoID), C2SGetCropsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPlantCrop_ProtoID), C2SPlantCropHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SHarvestCrop_ProtoID), C2SHarvestCropHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCropSpeedup_ProtoID), C2SSpeedupCropHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SHarvestCrops_ProtoID), C2SHarvestCropsHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetCatHouseInfo_ProtoID), C2SGetCatHousesInfoHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseAddCat_ProtoID), C2SCatHouseAddCatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseRemoveCat_ProtoID), C2SCatHouseRemoveCatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseStartLevelup_ProtoID), C2SCatHouseStartLevelupHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseSpeedLevelup_ProtoID), C2SCatHouseSpeedLevelupHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSellCatHouse_ProtoID), C2SCatHouseSellHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseProduceGold_ProtoID), C2SCatHouseProduceGoldHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseGetGold_ProtoID), C2SCatHouseGetGoldHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SCatHouseSetDone_ProtoID), C2SCatHouseSetDoneHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetHandbook_ProtoID), C2SGetHandbookHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetHead_ProtoID), C2SGetHeadHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetSuitHandbookReward_ProtoID), C2SGetSuitHandbookRewardHandler)

	// ?????????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SRankListRequest_ProtoID), C2SRankingListHandler)

	// ????????????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChatMsgPullRequest_ProtoID), C2SWorldChatMsgPullHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChatRequest_ProtoID), C2SWorldChatSendHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SHeartbeat_ProtoID), C2SHeartbeatHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPullFosterData_ProtoID), C2SPullFosterDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterEquipCard_ProtoID), C2SFosterEquipCardHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterUnequipCard_ProtoID), C2SFosterUnequipCardHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterSetCat_ProtoID), C2SFosterSetCatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterOutCat_ProtoID), C2SFosterOutCatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterSetCat2Friend_ProtoID), C2SFosterSetCat2FriendHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetPlayerFosterCats_ProtoID), C2SGetPlayerFosterCatsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SPullFosterCatsWithFriend_ProtoID), C2SPullFosterDataWithFriendHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterCardCompose_ProtoID), C2SFosterCardComposeHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFosterGetEmptySlotFriends_ProtoID), C2SFosterGetEmptySlotFriendsHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChapterUnlock_ProtoID), C2SChapterUnlockHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetCurHelpReqPIds_ProtoID), C2SGetCurHelpReqPIdsHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SUnlockArea_ProtoID), C2SUnlockAreaHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetAreasInfos_ProtoID), C2SGetAreasInfosHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMailSendRequest_ProtoID), C2SMailSendHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMailListRequest_ProtoID), C2SMailListHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMailDetailRequest_ProtoID), C2SMailDetailHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMailGetAttachedItemsRequest_ProtoID), C2SMailGetAttachedItemsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMailDeleteRequest_ProtoID), C2SMailDeleteHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFriendSearch_ProtoID), C2SFriendSearchHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SAddFriendByPId_ProtoID), C2SAddFriendByIdHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SRefuseFriend_ProtoID), C2SRefuseAddFriendHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SRemoveFriend_ProtoID), C2SFriendRemoveHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetFriendList_ProtoID), C2SGetFriendListHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGiveFriendPoints_ProtoID), C2SGiveFriendPointsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetFriendPoints_ProtoID), C2SGetFriendPointsHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFriendChat_ProtoID), C2SFriendChatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFriendGetUnreadMessageNum_ProtoID), C2SFriendGetUnreadMessageNumHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFriendPullUnreadMessage_ProtoID), C2SFriendPullUnreadMessageHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFriendConfirmUnreadMessage_ProtoID), C2SFriendConfirmUnreadMessageHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SAgreeFriend_ProtoID), C2SAddFriendAgreeHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SOpenFriendChest_ProtoID), C2SOpenFriendChestHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SDraw_ProtoID), C2SDrawHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SStageBegin_ProtoID), C2SStagePassBeginHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SStagePass_ProtoID), C2SStagePassHandler)
	//msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.), C2SDayBuyTiLiHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetBuildingInfos_ProtoID), C2SGetBuildingInfosHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSetBuilding_ProtoID), C2SSetBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetBackBuilding_ProtoID), C2SGetBackBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSellBuilding_ProtoID), C2SSellBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SRemoveBlock_ProtoID), C2SRemoveBlockHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SOpenMapChest_ProtoID), C2SOpenMapChestHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMoveBuilding_ProtoID), C2SMoveBuildingHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChgBuildingDir_ProtoID), C2SChgBuildingDirHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SVisitPlayer_ProtoID), C2SVisitPlayerHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSurfaceDataRequest_ProtoID), C2SSurfaceDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSurfaceUpdateRequest_ProtoID), C2SSurfaceUpdateHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetAllExpedition_ProtoID), C2SGetAllExpeditionHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChgExpedition_ProtoID), C2SChgExpeditionHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SStopExpedition_ProtoID), C2SStopExpeditionHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetExpeditionReward_ProtoID), C2SGetExpeditionRewardHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SStartExpedition_ProtoID), C2SStartExpeditionHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChgExpeditionResult_ProtoID), C2SChgExpeditionResultHandler)

	// ?????????
	/*msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message_id.MSGID_C2S_GOLD_HAND_DATA_REQUEST), C2SGoldHandDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message_id.MSGID_C2S_TOUCH_GOLD_REQUEST), C2STouchGoldHandler)*/

	// ??????
	/*msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message_id.MSGID_C2S_SHOP_DATA_REQUEST), C2SShopDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message_id.MSGID_C2S_SHOP_BUY_ITEM_REQUEST), C2SShopBuyItemHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message_id.MSGID_C2S_SHOP_REFRESH_REQUEST), C2SShopRefreshHandler)*/

	// ?????????
	/*msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message_id.MSGID_C2S_RANK_LIST_REQUEST), C2SRankListHandler)*/

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2STaskDataRequest_ProtoID), C2STaskDataHanlder)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2STaskRewardRequest_ProtoID), C2SGetTaskRewardHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChatRequest_ProtoID), C2SChatHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChatMsgPullRequest_ProtoID), C2SChatPullMsgHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSignDataRequest_ProtoID), C2SSignDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSignAwardRequest_ProtoID), C2SSignAwardHandler)

	// ?????????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSevenDaysDataRequest_ProtoID), C2SSevenDaysDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSevenDaysAwardRequest_ProtoID), C2SSevenDaysAwardHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChargeDataRequest_ProtoID), C2SChargeDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChargeRequest_ProtoID), C2SChargeHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SChargeFirstAwardRequest_ProtoID), C2SChargeFirstAwardHandler)

	// ????????????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SRedPointStatesRequest_ProtoID), C2SRedPointStatesHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGuideDataSaveRequest_ProtoID), C2SGuideDataSaveHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SActivityDataRequest_ProtoID), C2SActivityDataHandler)

	// ??????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFocusDataRequest_ProtoID), C2SFocusDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFocusPlayerRequest_ProtoID), C2SFocusPlayerHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SFocusPlayerCancalRequest_ProtoID), C2SUnfocusPlayerHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMyPictureDataRequest_ProtoID), C2SMyPictureDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SMyPictureSetRequest_ProtoID), C2SMyPictureSetHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSpaceDataRequest_ProtoID), C2SSpaceDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSpaceSetGenderRequest_ProtoID), C2SSpaceSetGenderHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSpaceFashionSaveRequest_ProtoID), C2SSpaceFashionSaveHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSpaceFashionDataRequest_ProtoID), C2SSpaceFashionDataHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SSpaceCatUnlockRequest_ProtoID), C2SSpaceOtherCatUnlockHandler)

	// ?????????
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetAllActivityInfos_ProtoID), C2SGetAllActivityInfosHandler)
	msg_handler_mgr.SetPlayerMsgHandler(uint16(msg_client_message.C2SGetActivityReward_ProtoID), C2SGetActivityRewardHandler)
}

func C2SEnterGameRequestHandler(msg_data []byte) (int32, *Player) {
	var p *Player
	var req msg_client_message.C2SEnterGameRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("unmarshal msg failed err(%s) !", err.Error())
		return -1, p
	}

	uid := login_token_mgr.GetUidByAccount(req.GetAcc())
	if uid == "" {
		log.Error("PlayerEnterGameHandler account[%v] cant get", req.GetAcc())
		return int32(msg_client_message.E_ERR_PLAYER_TOKEN_ERROR), p
	}

	row := dbc.BanPlayers.GetRow(uid)
	if row != nil && row.GetStartTime() > 0 {
		log.Error("Player unique id %v be banned", uid)
		return int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED), p
	}

	var is_new bool
	p = player_mgr.GetPlayerByUid(uid)
	if nil == p {
		global_row := dbc.Global.GetRow()
		player_id := global_row.GetNextPlayerId()
		pdb := dbc.Players.AddRow(player_id)
		if nil == pdb {
			log.Error("player_db_to_msg AddRow pid(%d) failed !", player_id)
			return -1, p
		}
		pdb.SetUniqueId(uid)
		pdb.SetAccount(req.GetAcc())
		//pdb.SetCurrReplyMsgNum(0)
		p = new_player(player_id, uid, req.GetAcc(), "", pdb)
		player_mgr.Add2IdMap(p)
		player_mgr.Add2UidMap(uid, p)
		p.OnCreate()
		is_new = true
		log.Info("player_db_to_msg new player(%d) !", player_id)
	} else {
		p.Account = req.GetAcc()
		pdb := dbc.Players.GetRow(p.Id)
		if pdb != nil {
			//pdb.SetCurrReplyMsgNum(0)
		}
	}

	p.OnLogin()
	p.send_enter_game(req.Acc, p.Id)
	p.send_data_on_login(is_new)
	p.notify_enter_complete()

	log.Info("PlayerEnterGameHandler account[%s]", req.GetAcc())

	return 1, p
}

func C2SLeaveGameRequestHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SLeaveGameRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s) !", err.Error())
		return -1
	}
	p.OnLogout(true)
	return 1
}

func C2SHeartbeatHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SHeartbeat
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s) !", err.Error())
		return -1
	}

	if p.IsOffline() {
		log.Error("Player[%v] is offline", p.Id)
		return int32(msg_client_message.E_ERR_PLAYER_IS_OFFLINE)
	}

	// ??????????????????
	p.CheckSystemMail()

	// ??????
	p.check_and_pull_chat()

	response := &msg_client_message.S2CHeartbeat{
		SysTime: int32(time.Now().Unix()),
	}
	p.Send(uint16(msg_client_message.S2CHeartbeat_ProtoID), response)

	return 1
}

func C2SPlayerChangeNameHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SPlayerChangeNameRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s)!", err.Error())
		return -1
	}
	if len(req.GetNewName()) > int(global_config.MaxNameLen) {
		log.Error("Player[%v] change new name[%v] is too long", p.Id, req.GetNewName())
		return int32(msg_client_message.E_ERR_PLAYER_NAME_TOO_LONG)
	}
	if p.db.GetName() != "" {
		if global_config.ChgNameCost != nil && len(global_config.ChgNameCost) > 0 {
			/*if p.get_diamond() < global_config.ChgNameCost[0] {
				return int32(msg_client_message.E_ERR_PLAYER_DIAMOND_NOT_ENOUGH)
			}
			p.add_diamond(-global_config.ChgNameCost[0])*/
		}
	}
	p.db.SetName(req.GetNewName())
	p.rpc_player_base_info_update()
	p.Send(uint16(msg_client_message.S2CPlayerChangeNameResponse_ProtoID), &msg_client_message.S2CPlayerChangeNameResponse{
		NewName: req.GetNewName(),
	})

	return 1
}

func C2SPlayerChangeHeadHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SPlayerChangeHeadRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s)!", err.Error())
		return -1
	}
	return p.change_head(req.GetNewHead())
}

func C2SRedPointStatesHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SRedPointStatesRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s)!", err.Error())
		return -1
	}
	return p.send_red_point_states(req.GetModules())
}

func C2SGuideDataSaveHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SGuideDataSaveRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s)!", err.Error())
		return -1
	}
	p.db.GuideData.SetData(req.GetData())
	response := &msg_client_message.S2CGuideDataSaveResponse{
		Data: req.GetData(),
	}
	p.Send(uint16(msg_client_message.S2CGuideDataSaveResponse_ProtoID), response)
	log.Debug("Player[%v] guide save %v", p.Id, req.GetData())
	return 1
}

func (p *Player) reconnect() int32 {
	uid := p.db.GetUniqueId()
	row := dbc.BanPlayers.GetRow(uid)
	if row != nil && row.GetStartTime() > 0 {
		log.Error("Player unique id %v be banned", uid)
		return int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
	}

	new_token := share_data.GenerateAccessToken(uid)
	login_token_mgr.SetToken(uid, new_token, p.Id)
	conn_timer_wheel.Remove(p.Id)
	atomic.StoreInt32(&p.is_login, 1)

	response := &msg_client_message.S2CReconnectResponse{
		NewToken: new_token,
	}
	p.Send(uint16(msg_client_message.S2CReconnectResponse_ProtoID), response)

	p.send_items()

	log.Trace("Player[%v] reconnected, new token %v", p.Id, new_token)
	return 1
}

func C2SReconnectHandler(p *Player, msg_data []byte) int32 {
	var req msg_client_message.C2SReconnectRequest
	err := proto.Unmarshal(msg_data, &req)
	if err != nil {
		log.Error("Unmarshal msg failed err(%s)!", err.Error())
		return -1
	}

	return p.reconnect()
}
