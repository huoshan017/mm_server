package main

import (
	"mm_server_new/libs/log"
	"mm_server_new/libs/utils"
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	"mm_server_new/src/common"
	"sync"
)

type RankList struct {
	rank_list *utils.CommonRankingList
	item_pool *sync.Pool
	root_node utils.SkiplistNode
}

func (this *RankList) Init(root_node utils.SkiplistNode) {
	this.root_node = root_node
	this.rank_list = utils.NewCommonRankingList(this.root_node, 100000)
	this.item_pool = &sync.Pool{
		New: func() interface{} {
			return this.root_node.New()
		},
	}
}

func (this *RankList) GetItemByKey(key interface{}) (item utils.SkiplistNode) {
	return this.rank_list.GetByKey(key)
}

func (this *RankList) GetRankByKey(key interface{}) int32 {
	return this.rank_list.GetRank(key)
}

func (this *RankList) GetItemByRank(rank int32) (item utils.SkiplistNode) {
	return this.rank_list.GetByRank(rank)
}

func (this *RankList) SetValueByKey(key interface{}, value interface{}) {
	this.rank_list.SetValueByKey(key, value)
}

func (this *RankList) RankNum() int32 {
	return this.rank_list.GetLength()
}

// 获取排名项
func (this *RankList) GetItemsByRange(key interface{}, start_rank, rank_num int32) (rank_items []utils.SkiplistNode, self_rank int32, self_value interface{}) {
	start_rank, rank_num = this.rank_list.GetRankRange(start_rank, rank_num)
	if start_rank == 0 {
		log.Error("Get rank list range with [%v,%v] failed", start_rank, rank_num)
		return nil, 0, nil
	}

	nodes := make([]interface{}, rank_num)
	for i := int32(0); i < rank_num; i++ {
		nodes[i] = this.item_pool.Get().(utils.SkiplistNode)
	}

	num := this.rank_list.GetRangeNodes(start_rank, rank_num, nodes)
	if num == 0 {
		log.Error("Get rank list nodes failed")
		return nil, 0, nil
	}

	rank_items = make([]utils.SkiplistNode, num)
	for i := int32(0); i < num; i++ {
		rank_items[i] = nodes[i].(utils.SkiplistNode)
	}

	self_rank, self_value = this.rank_list.GetRankAndValue(key)
	return

}

// 获取最后的几个排名
func (this *RankList) GetLastRankRange(rank_num int32) (int32, int32) {
	return this.rank_list.GetLastRankRange(rank_num)
}

// 更新排行榜
func (this *RankList) UpdateItem(item utils.SkiplistNode) bool {
	if !this.rank_list.Update(item) {
		log.Error("Update rank item[%v] failed", item)
		return false
	}
	return true
}

// 删除指定值
func (this *RankList) DeleteItem(key interface{}) bool {
	return this.DeleteItem(key)
}

var root_rank_item = []utils.SkiplistNode{
	nil,                             // 0
	&common.PlayerInt32RankItem{},   // 1
	&common.PlayerInt32RankItem{},   // 2
	&common.PlayerCatOuqiRankItem{}, // 3
	&common.PlayerInt32RankItem{},   // 4
}

type RankListManager struct {
	rank_lists []*RankList
	rank_map   map[int32]*RankList
	locker     *sync.RWMutex
}

var rank_list_mgr RankListManager

func (this *RankListManager) Init() {
	this.rank_lists = make([]*RankList, common.RANK_LIST_TYPE_MAX)
	for i := int32(1); i < common.RANK_LIST_TYPE_MAX; i++ {
		if int(i) >= len(root_rank_item) {
			break
		}
		this.rank_lists[i] = &RankList{}
		this.rank_lists[i].Init(root_rank_item[i])
	}
	this.rank_map = make(map[int32]*RankList)
	this.locker = &sync.RWMutex{}
}

func (this *RankListManager) GetRankList(rank_type int32) (rank_list *RankList) {
	if int(rank_type) >= len(this.rank_lists) {
		return nil
	}
	if this.rank_lists[rank_type] == nil {
		return nil
	}
	return this.rank_lists[rank_type]
}

func (this *RankListManager) GetItemByKey(rank_type int32, key interface{}) (item utils.SkiplistNode) {
	if int(rank_type) >= len(this.rank_lists) {
		return nil
	}
	if this.rank_lists[rank_type] == nil {
		return nil
	}
	return this.rank_lists[rank_type].GetItemByKey(key)
}

func (this *RankListManager) GetRankByKey(rank_type int32, key interface{}) int32 {
	if int(rank_type) >= len(this.rank_lists) {
		return -1
	}
	if this.rank_lists[rank_type] == nil {
		return -1
	}
	return this.rank_lists[rank_type].GetRankByKey(key)
}

func (this *RankListManager) GetItemByRank(rank_type, rank int32) (item utils.SkiplistNode) {
	if int(rank_type) >= len(this.rank_lists) {
		return nil
	}
	if this.rank_lists[rank_type] == nil {
		return nil
	}
	return this.rank_lists[rank_type].GetItemByRank(rank)
}

func (this *RankListManager) GetItemsByRange(rank_type int32, key interface{}, start_rank, rank_num int32) (rank_items []utils.SkiplistNode, self_rank int32, self_value interface{}) {
	if int(rank_type) >= len(this.rank_lists) {
		return nil, 0, nil
	}
	if this.rank_lists[rank_type] == nil {
		return nil, 0, nil
	}
	return this.rank_lists[rank_type].GetItemsByRange(key, start_rank, rank_num)
}

func (this *RankListManager) GetLastRankRange(rank_type, rank_num int32) (int32, int32) {
	if int(rank_type) >= len(this.rank_lists) {
		return -1, -1
	}
	if this.rank_lists[rank_type] == nil {
		return -1, -1
	}
	return this.rank_lists[rank_type].GetLastRankRange(rank_num)
}

func (this *RankListManager) UpdateItem(rank_type int32, item utils.SkiplistNode) bool {
	if int(rank_type) >= len(this.rank_lists) {
		return false
	}
	if this.rank_lists[rank_type] == nil {
		return false
	}
	return this.rank_lists[rank_type].UpdateItem(item)
}

func (this *RankListManager) DeleteItem(rank_type int32, key interface{}) bool {
	if int(rank_type) >= len(this.rank_lists) {
		return false
	}
	if this.rank_lists[rank_type] == nil {
		return false
	}
	return this.rank_lists[rank_type].DeleteItem(key)
}

func (this *RankListManager) GetRankList2(rank_type int32) (rank_list *RankList) {
	this.locker.RLock()
	rank_list = this.rank_map[rank_type]
	if rank_list == nil {
		rank_list = &RankList{}
		this.rank_map[rank_type] = rank_list
	}
	this.locker.RUnlock()
	return
}

func (this *RankListManager) GetItemByKey2(rank_type int32, key interface{}) (item utils.SkiplistNode) {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return
	}
	return rank_list.GetItemByKey(key)
}

func (this *RankListManager) GetRankByKey2(rank_type int32, key interface{}) int32 {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return 0
	}
	return rank_list.GetRankByKey(key)
}

func (this *RankListManager) GetItemByRank2(rank_type, rank int32) (item utils.SkiplistNode) {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return
	}
	return rank_list.GetItemByRank(rank)
}

func (this *RankListManager) GetItemsByRange2(rank_type int32, key interface{}, start_rank, rank_num int32) (rank_items []utils.SkiplistNode, self_rank int32, self_value interface{}) {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return
	}
	return rank_list.GetItemsByRange(key, start_rank, rank_num)
}

func (this *RankListManager) GetLastRankRange2(rank_type, rank_num int32) (int32, int32) {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return -1, -1
	}
	return rank_list.GetLastRankRange(rank_num)
}

func (this *RankListManager) UpdateItem2(rank_type int32, item utils.SkiplistNode) bool {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return false
	}
	return rank_list.UpdateItem(item)
}

func (this *RankListManager) DeleteItem2(rank_type int32, key interface{}) bool {
	rank_list := this.GetRankList2(rank_type)
	if rank_list == nil {
		return false
	}
	return rank_list.DeleteItem(key)
}

func transfer_nodes_to_rank_items(rank_type int32, start_rank int32, items []utils.SkiplistNode) (rank_items []*msg_client_message.RankItemInfo) {
	var item *common.PlayerInt32RankItem
	for i := int32(0); i < int32(len(items)); i++ {
		item = (items[i]).(*common.PlayerInt32RankItem)
		if item == nil {
			continue
		}
		var rank_item = &msg_client_message.RankItemInfo{}
		if rank_type == common.RANK_LIST_TYPE_STAGE_TOTAL_SCORE {

		} else if rank_type == common.RANK_LIST_TYPE_CHARM {

		} else if rank_type == common.RANK_LIST_TYPE_CAT_OUQI {

		} else if rank_type == common.RANK_LIST_TYPE_BE_ZANED {

		} else {
			log.Error("invalid rank type[%v] transfer nodes to rank items", rank_type)
			return
		}
		rank_items = append(rank_items, rank_item)
	}
	return
}
