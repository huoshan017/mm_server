package main

import (
	"mm_server_new/libs/log"
	"mm_server_new/libs/rpc"
	"mm_server_new/src/share_data"
)

type GameRpcClient struct {
	server_idx int32
	server_id  int32
	server_ip  string
	rpc_client *rpc.Client
}

// 通过ServerId对应rpc客户端
func GetRpcClientByServerId(server_id int32) *rpc.Client {
	server_info := server_list.GetServerById(server_id)
	if server_info == nil {
		log.Error("get server info by server_id[%v] from failed", server_id)
		return nil
	}
	r := server.game_rpc_clients[server_id]
	if r == nil {
		log.Error("通过ServerID[%v]获取rpc客户端失败", server_id)
		return nil
	}
	return r.rpc_client
}

// 通过玩家ID对应大厅的rpc客户端
func GetRpcClientByPlayerId(player_id int32) *rpc.Client {
	server_id := share_data.GetServerIdByPlayerId(player_id)
	return GetRpcClientByServerId(server_id)
}

// 通过源玩家ID和目标玩家ID获得跨服rpc客户端
func GetCrossRpcClientByPlayerId(from_player_id, to_player_id int32) *rpc.Client {
	from_server_id := share_data.GetServerIdByPlayerId(from_player_id)
	to_server_id := share_data.GetServerIdByPlayerId(to_player_id)
	if !server_list.IsSameCross(from_server_id, to_server_id) {
		return nil
	}
	return GetRpcClientByServerId(to_server_id)
}

// 通过源玩家ID和目标公会ID获得跨服rpc客户端
func GetCrossRpcClientByGuildId(from_player_id, to_guild_id int32) *rpc.Client {
	from_server_id := share_data.GetServerIdByPlayerId(from_player_id)
	to_server_id := share_data.GetServerIdByGuildId(to_guild_id)
	if !server_list.IsSameCross(from_server_id, to_server_id) {
		return nil
	}
	return GetRpcClientByServerId(to_server_id)
}
