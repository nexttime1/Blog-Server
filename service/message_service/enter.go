package message_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"errors"
	"fmt"
)

type MessageAddRequest struct {
	SendUserID uint   `json:"send_user_id"` // 发送人id
	RevUserID  uint   `json:"rev_user_id"`  // 接收人id
	Content    string `json:"content"`
}

type MessageRecordRequest struct {
	UserID uint `json:"user_id" binding:"required" msg:"请输入查询的用户id"`
}

func GetConversationID(SendUserID, RevUserID uint) string {
	if SendUserID > RevUserID {
		SendUserID, RevUserID = RevUserID, SendUserID
	}
	//从小到大排  1_3  3_1  选 1_3  当唯一的 房间id
	return fmt.Sprintf("%d_%d", SendUserID, RevUserID)
}

func MessageAddService(cr MessageAddRequest) error {
	var sendModel models.UserModel
	var msg string
	err := global.DB.Where("id = ?", cr.SendUserID).Take(&sendModel).Error
	if err != nil {
		msg = "未找到发送者id"
		return errors.New(msg)
	}
	var RevModel models.UserModel
	err = global.DB.Where("id = ?", cr.RevUserID).Take(&RevModel).Error
	if err != nil {
		msg = "未找到接收者id"
		return errors.New(msg)
	}
	conversationID := GetConversationID(cr.SendUserID, cr.RevUserID)
	messageModel := models.MessageModel{
		SendUserID:       sendModel.ID,
		SendUserNickName: sendModel.Nickname,
		SendUserAvatar:   sendModel.Avatar,
		RevUserID:        RevModel.ID,
		RevUserNickName:  RevModel.Nickname,
		RevUserAvatar:    RevModel.Avatar,
		IsRead:           false,
		Content:          cr.Content,
		ConversationID:   conversationID,
	}
	err = global.DB.Create(&messageModel).Error
	if err != nil {
		return err
	}
	// 在 Conversation表中 增加 相应参数
	var conversationModel models.ConversationModel
	var uid1, uid2 uint
	_, err = fmt.Sscanf(conversationID, "%d_%d", &uid1, &uid2)
	if err != nil {
		return fmt.Errorf("scanf 转换错误 %s", err)
	}
	userIDsStr := fmt.Sprintf("%d,%d", uid1, uid2)
	err = global.DB.Where("user_ids = ?", userIDsStr).Take(&conversationModel).Error
	if err == nil {
		//已经存在  只需修改最后 聊天
		global.DB.Model(&conversationModel).Update("last_msg_id", messageModel.ID)
	} else {
		//不存在  创建
		err = global.DB.Create(&models.ConversationModel{
			UserIDs:   []uint{uid1, uid2},
			LastMsgID: messageModel.ID,
		}).Error
		if err != nil {
			return fmt.Errorf("创建会话失败 %s", err)
		}
	}

	return nil
}
