package code

import (
	"fmt"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

func Execute() {
	activityTracker := NewActivityTracker(24 * 60)
	for {
		fmt.Println("\nEnter operation:\n1. STORE <userID> <activity_id>\n2. EXIT")

		var op int
		_, err := fmt.Scan(&op)
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			continue
		}

		switch op {
		case 1:
			var userID uint
			var activityID int // Assuming your enum is an integer type underneath

			fmt.Println("Enter space-separated <userID> and <activityID>:")
			_, err := fmt.Scan(&userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters. Try again.")
				continue
			}

			// Capture the current timestamp for the store function
			currentTime := time.Now()

			// Cast the int input to your enum type
			activityTracker.StoreUserActivity(currentTime, userID, enum.Activity(activityID))
			fmt.Printf("Successfully stored activity %d for user %d at %s\n", activityID, userID, currentTime.Format("15:04"))

		case 2:
			fmt.Println("Exiting application...")
			return

		default:
			fmt.Println("Unknown operation number. Please select 1 or 2.")
		}
	}
}
