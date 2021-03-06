package tables

import (
	"encoding/xml"
	"io/ioutil"
	"mm_server_new/libs/log"
	"mm_server_new/src/server_config"
)

type XmlSysMsgItem struct {
	Id            int32 `xml:"Id,attr"`
	Type          int32 `xml:"Type,attr"`
	DescriptionId int32 `xml:"DescriptionId,attr"`
	NewsTicker    int32 `xml:"OutPut,attr"`
}

type XmlSysMsgConfig struct {
	Items []XmlSysMsgItem `xml:"item"`
}

type SysMsgTableMgr struct {
	Map   map[int32]*XmlSysMsgItem
	Array []*XmlSysMsgItem
}

func (this *SysMsgTableMgr) Init(table_file string) bool {
	if table_file == "" {
		table_file = "SysMessage.xml"
	}
	file_path := server_config.GetGameDataPathFile(table_file)
	data, err := ioutil.ReadFile(file_path)
	if nil != err {
		log.Error("CropTableMgr read file err[%s] !", err.Error())
		return false
	}

	tmp_cfg := &XmlSysMsgConfig{}
	err = xml.Unmarshal(data, tmp_cfg)
	if nil != err {
		log.Error("SysMsgTableMgr xml Unmarshal failed error [%s] !", err.Error())
		return false
	}

	if this.Map == nil {
		this.Map = make(map[int32]*XmlSysMsgItem)
	}

	if this.Array == nil {
		this.Array = make([]*XmlSysMsgItem, 0)
	}

	tmp_len := int32(len(tmp_cfg.Items))

	var tmp_item *XmlSysMsgItem
	for idx := int32(0); idx < tmp_len; idx++ {
		tmp_item = &tmp_cfg.Items[idx]

		this.Map[tmp_item.Id] = tmp_item
		this.Array = append(this.Array, tmp_item)
	}

	return true
}
