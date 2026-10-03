package code

import "fmt"

type NotificationService struct{}

func (ns *NotificationService) NotifyUser(user *User, meetingTitle string) {
	fmt.Printf("%v, meeting %v is updated", user.name, meetingTitle)
}
