package friendsservices

import (
	"fmt"
	"time"

	friendsapperr "github.com/suryansh74/chat_app/internal/friends/apperr"
	friendsdomain "github.com/suryansh74/chat_app/internal/friends/domain"
	friendsrepositories "github.com/suryansh74/chat_app/internal/friends/repositories"
	notificationservices "github.com/suryansh74/chat_app/internal/notification/services"
	wshandler "github.com/suryansh74/chat_app/internal/ws"
	"github.com/suryansh74/chat_app/pkg/logger"
	emailadapters "github.com/suryansh74/chat_app/shared/email/adapters"
	"github.com/suryansh74/chat_app/shared/token"
)

type FriendsServicePort interface {
	GetFriends(userID string) ([]friendsdomain.FriendListItem, error)
	SendFriendRequest(fromUserID, toUserID string) error
	AcceptFriendRequest(requestID, userID string) error
	RejectFriendRequest(requestID, userID string) error
	SearchByEmail(query string) ([]friendsdomain.Friend, error)
	SearchFriends(userID, query string) ([]friendsdomain.FriendListItem, error)
	RemoveFriend(userID, friendID string) error
}

type UserInfo struct {
	ID    string
	Name  string
	Email string
}

type friendsService struct {
	friendRepo      friendsrepositories.FriendRepositoryPort
	notificationSvc notificationservices.NotificationServicePort
	emailSender     emailadapters.EmailSenderPort
	passetoMaker    token.Maker
	wsHub           *wshandler.Hub
}

func NewFriendsService(
	friendRepo friendsrepositories.FriendRepositoryPort,
	notificationSvc notificationservices.NotificationServicePort,
	emailSender emailadapters.EmailSenderPort,
	passetoMaker token.Maker,
	wsHub *wshandler.Hub,
) FriendsServicePort {
	return &friendsService{
		friendRepo:      friendRepo,
		notificationSvc: notificationSvc,
		emailSender:     emailSender,
		passetoMaker:    passetoMaker,
		wsHub:           wsHub,
	}
}

func (s *friendsService) GetFriends(userID string) ([]friendsdomain.FriendListItem, error) {
	return s.friendRepo.GetFriendsByUserID(userID)
}

func (s *friendsService) SendFriendRequest(fromUserID, toUserID string) error {
	isFriend, err := s.friendRepo.IsFriend(fromUserID, toUserID)
	if err != nil {
		return err
	}
	if isFriend {
		return friendsapperr.ErrAlreadyFriends
	}

	// Get sender and receiver info
	senderInfo, err := s.friendRepo.GetUserInfo(fromUserID)
	if err != nil {
		logger.Log.Error("SendFriendRequest: failed to get sender info", "error", err.Error())
		return err
	}

	receiverInfo, err := s.friendRepo.GetUserInfo(toUserID)
	if err != nil {
		logger.Log.Error("SendFriendRequest: failed to get receiver info", "error", err.Error())
		return err
	}

	// Create in-app notification
	err = s.notificationSvc.CreateFriendRequest(fromUserID, toUserID)
	if err != nil {
		return err
	}

	// Send real-time notification via WebSocket
	if s.wsHub != nil {
		s.wsHub.BroadcastNewNotification(toUserID, "FRIEND_REQUEST", senderInfo.Name+" has sent you a friend request")
	}

	// Send email notification
	subject := "New Friend Request"
	body := fmt.Sprintf("Hello,\n\n%s (%s) has sent you a friend request.\n\nAccept or reject the request in the app.\n\n- ChatApp", senderInfo.Name, senderInfo.Email)

	if err := s.emailSender.SendEmail(receiverInfo.Email, subject, body); err != nil {
		logger.Log.Error("SendFriendRequest: failed to send email", "error", err.Error(), "to", receiverInfo.Email)
	}

	return nil
}

func (s *friendsService) AcceptFriendRequest(requestID, userID string) error {
	notification, err := s.notificationSvc.GetNotificationByID(requestID)
	if err != nil {
		return err
	}

	if notification.ToUserID != userID {
		return friendsapperr.ErrUnauthorized
	}

	if notification.Type != "FRIEND_REQUEST" {
		return friendsapperr.ErrInvalidRequestType
	}

	// Create friend relationship both ways
	err = s.friendRepo.CreateFriend(userID, notification.FromUserID)
	if err != nil {
		logger.Log.Error("AcceptFriendRequest: failed to create friend", "error", err.Error())
		return err
	}

	err = s.friendRepo.CreateFriend(notification.FromUserID, userID)
	if err != nil {
		logger.Log.Error("AcceptFriendRequest: failed to create friend reverse", "error", err.Error())
		return err
	}

	// Mark notification as read
	err = s.notificationSvc.MarkAsRead(requestID)
	if err != nil {
		logger.Log.Error("AcceptFriendRequest: failed to mark as read", "error", err.Error())
	}

	// Get acceptor info for notifications
	acceptorInfo, _ := s.friendRepo.GetUserInfo(userID)

	// Create acceptance notification for the sender
	err = s.notificationSvc.CreateFriendAccepted(userID, notification.FromUserID)
	if err != nil {
		logger.Log.Error("AcceptFriendRequest: failed to create notification", "error", err.Error())
	}

	// Send real-time notification to requester
	if s.wsHub != nil {
		acceptorName := "A user"
		if acceptorInfo != nil {
			acceptorName = acceptorInfo.Name
		}
		s.wsHub.BroadcastNewNotification(notification.FromUserID, "FRIEND_ACCEPTED", acceptorName+" accepted your friend request")
	}

	// Send email to the request sender
	requesterInfo, _ := s.friendRepo.GetUserInfo(notification.FromUserID)
	if requesterInfo != nil && acceptorInfo != nil {
		subject := "Friend Request Accepted"
		body := fmt.Sprintf("Hello,\n\n%s (%s) has accepted your friend request.\n\nYou can now chat with them in the app!\n\n- ChatApp", acceptorInfo.Name, acceptorInfo.Email)

		if err := s.emailSender.SendEmail(requesterInfo.Email, subject, body); err != nil {
			logger.Log.Error("AcceptFriendRequest: failed to send email", "error", err.Error())
		}
	}

	return nil
}

func (s *friendsService) RejectFriendRequest(requestID, userID string) error {
	notification, err := s.notificationSvc.GetNotificationByID(requestID)
	if err != nil {
		return err
	}

	if notification.ToUserID != userID {
		return friendsapperr.ErrUnauthorized
	}

	// Get rejector info
	rejectorInfo, _ := s.friendRepo.GetUserInfo(userID)
	senderID := notification.FromUserID

	// Delete the notification
	err = s.notificationSvc.DeleteNotification(requestID)
	if err != nil {
		return err
	}

	// Create rejection notification for the sender
	err = s.notificationSvc.CreateFriendRejected(userID, senderID)
	if err != nil {
		logger.Log.Error("RejectFriendRequest: failed to create rejection notification", "error", err.Error())
	}

	// Send real-time notification
	if s.wsHub != nil {
		rejectorName := "A user"
		if rejectorInfo != nil {
			rejectorName = rejectorInfo.Name
		}
		s.wsHub.BroadcastNewNotification(senderID, "FRIEND_REJECTED", rejectorName+" rejected your friend request")
	}

	return nil
}

func (s *friendsService) SearchByEmail(query string) ([]friendsdomain.Friend, error) {
	return s.friendRepo.GetFriendsByEmail(query)
}

func (s *friendsService) SearchFriends(userID, query string) ([]friendsdomain.FriendListItem, error) {
	friends, err := s.friendRepo.GetFriendsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var filtered []friendsdomain.FriendListItem
	for _, f := range friends {
		if query == "" || containsIgnoreCase(f.FriendName, query) || containsIgnoreCase(f.FriendEmail, query) {
			filtered = append(filtered, f)
		}
	}

	return filtered, nil
}

func (s *friendsService) RemoveFriend(userID, friendID string) error {
	err := s.friendRepo.DeleteFriend(userID, friendID)
	if err != nil {
		return err
	}
	return s.friendRepo.DeleteFriend(friendID, userID)
}

func containsIgnoreCase(s, substr string) bool {
	sLower := ""
	substrLower := ""
	for _, r := range s {
		sLower += string(r)
	}
	for _, r := range substr {
		substrLower += string(r)
	}
	return len(sLower) > 0 && len(substrLower) > 0 && containsSubstring(sLower, substrLower)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	_ = time.Second
}
