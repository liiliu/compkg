package livekit

// server.go
import (
	"context"
	"fmt"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

var roomClient *lksdk.RoomServiceClient
var apiKey string
var apiSecret string

func Initial() {
	apiKey = config.GetString("liveKit.apiKey")
	apiSecret = config.GetString("liveKit.apiSecret")
	//roomClient = lksdk.NewRoomServiceClient(config.GetString("liveKit.host"), apiKey, apiSecret)
}

// GetJoinToken 生成客户端加入令牌
func GetJoinToken(room, identity, name string) (token string, err error) {
	at := auth.NewAccessToken(apiKey, apiSecret)
	canUpdateOwnMetadata := true
	grant := &auth.VideoGrant{
		RoomJoin:             true,
		Room:                 room,
		CanUpdateOwnMetadata: &canUpdateOwnMetadata,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour).
		SetName(name)

	token, err = at.ToJWT()
	return
}

// CreateRoom 创建房间
func CreateRoom(roomName string) error {
	room, err := roomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            roomName,
		EmptyTimeout:    10 * 60, // 10 minutes
		MaxParticipants: 20,
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("PutFromFile Error: %s", err.Error()))
		return err
	}

	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("CreateRoom Success: %s", room.String()))
	return nil
}

// ListRooms 列出房间
func ListRooms() error {
	rooms, err := roomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("ListRooms Error: %s", err.Error()))
		return err

	}
	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("ListRooms Success: %s", rooms.String()))
	return nil
}

// DeleteRoom 删除房间
func DeleteRoom(roomName string) error {
	_, err := roomClient.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
		Room: roomName,
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("DeleteRoom Error: %s", err.Error()))
		return err

	}
	return nil
}

// ListParticipants 列出参与者
func ListParticipants(roomName string) (participantInfos []*livekit.ParticipantInfo, err error) {
	res, err := roomClient.ListParticipants(context.Background(), &livekit.ListParticipantsRequest{
		Room: roomName,
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("ListParticipants Error: %s", err.Error()))
		return
	}
	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("ListParticipants Success: %s", res.String()))
	participantInfos = res.Participants
	return
}

// GetParticipant  获取参与者信息
func GetParticipant(roomName, identity string) error {
	res, err := roomClient.GetParticipant(context.Background(), &livekit.RoomParticipantIdentity{
		Room:     roomName,
		Identity: identity,
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("GetParticipant Error: %s", err.Error()))
		return err
	}
	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("GetParticipant Success: %s", res.String()))
	return nil
}

// UpdateParticipant  更新参与者信息
func UpdateParticipant(roomName, identity string) error {
	// Promotes an audience member to a speaker
	res, err := roomClient.UpdateParticipant(context.Background(), &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Permission: &livekit.ParticipantPermission{
			CanSubscribe:   true, //可以订阅
			CanPublish:     true, //可以发布
			CanPublishData: true, //可以发布数据
		},
		Metadata: "{}",
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("UpdateParticipant Error: %s", err.Error()))
		return err
	}
	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("UpdateParticipant Success: %s", res.String()))
	return nil
}

// RemoveParticipant 移除参与者
func RemoveParticipant(roomName, identity string) error {
	res, err := roomClient.RemoveParticipant(context.Background(), &livekit.RoomParticipantIdentity{
		Room:     roomName,
		Identity: identity,
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("RemoveParticipant Error: %s", err.Error()))
	}
	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("RemoveParticipant Success: %s", res.String()))
	return nil
}

// MutePublishedTrack  静音/取消静音特定轨道
func MutePublishedTrack(roomName, identity, trackSid string) error {
	res, err := roomClient.MutePublishedTrack(context.Background(), &livekit.MuteRoomTrackRequest{
		Room:     roomName,
		Identity: identity,
		TrackSid: trackSid,
		Muted:    true, //静音
	})
	if err != nil {
		logger.Error(common.LogTagLiveKitError, fmt.Sprintf("MutePublishedTrack Error: %s", err.Error()))

	}
	logger.Info(common.LogTagLiveKitError, fmt.Sprintf("MutePublishedTrack Success: %s", res.String()))
	return nil
}
