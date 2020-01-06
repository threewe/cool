package room

import (
	"sync"
)
type RoomSetUp struct {
	Roomid int
	RoomName string
	CreteUserId int
	RoomUserId int
	PeopleNum int
	Config interface{}
	GameType int
	GameName string
	WatchNum int
}
type Room struct {
	roomId int // 房间号
	roomName string
	createUserId int // 房间创建者
	roomUserId int // 房间归属者
	seats [] interface{} //房间座位
	history sync.Map
	cards [] int // 牌
	config interface{} // 配置
	gameType int // 游戏类型
	gameName string // 游戏名
	maxWatchNum int // 最大观战人数
	peopleNum int // 人数
	Begin bool // 房间游戏是否开始
	OfReady int // 准备人数
}
type CreateRoomConfig struct {
	PeopleNum int16 // 人数
}
// 创建房间
func CreateRoom(up *RoomSetUp) *Room {
	room := &Room{

	}
	//rooms.Store(up.RoomName, room) // 存储房间实例
	return room
}
func (ro *Room) Init(up *RoomSetUp) {
	ro.config =       up.Config
	ro.gameType=      up.GameType
	ro.gameName=      up.GameName
	ro.maxWatchNum =  up.WatchNum
	ro.peopleNum =    up.PeopleNum
	ro.roomId =       up.Roomid
	ro.roomName =     up.RoomName
	ro.createUserId = up.CreteUserId
	ro.roomUserId =   up.RoomUserId
	ro.seats =        make([]interface{}, up.PeopleNum)
	ro.history = sync.Map{}
	ro.cards = nil
	ro.Begin = false
	ro.OfReady = 0
}
// 获取房间号
func (ro *Room) GetRoomId() int {
	return ro.roomId
}
// 获取创建者id
func (ro *Room) GetCreateUserId() int {
	return ro.createUserId
}
// 获取房间所属者
func (ro *Room) GetRoomUserId() int {
	return ro.roomUserId
}

func (ro *Room) GetUser(userId int) (interface{}, bool){
	user, ok := ro.history.Load(userId)
	if !ok {return nil, ok}
	return user, ok
}





