package main

import (
	msg_client_message "mm_server_new/proto/gen_go/client_message"
)

func (this *Player) guide_data() {
	response := &msg_client_message.S2CGuideDataResponse{
		Data: this.db.GuideData.GetData(),
	}
	this.Send(uint16(msg_client_message.S2CGuideDataResponse_ProtoID), response)
}
