package code

import "fmt"

type NotificationService struct{}

func (ns *NotificationService) NotifyUser(user *User, meeting *Meeting) {
	fmt.Printf("%v, meeting %v is updated", user.name, meeting.id)
}
