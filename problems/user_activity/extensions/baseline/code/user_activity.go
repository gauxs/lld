package code

import (
	"fmt"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

func Execute() {
	activityTracker := NewActivityTracker(24 * 60)
	for {
		fmt.Println("\n--- Select Operation ---")
		fmt.Println("1. Store User Activity")
		fmt.Println("2. Get Activity Count For User In Timerange")
		fmt.Println("3. Get Activity Rate For User In Timerange")
		fmt.Println("4. Get Distinct User Count By Activity In Timerange")
		fmt.Println("5. Get Users By Activity In Timerange")
		fmt.Println("6. Get Top K Users For Activity In Timerange")
		fmt.Println("7. Get Users Activity Summary In Timerange")
		fmt.Println("8. EXIT")
		fmt.Print("Enter option (1-8): ")

		var op int
		_, err := fmt.Scan(&op)
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			continue
		}

		switch op {
		case 1: // StoreUserActivity
			var userID uint
			var activityID int
			fmt.Println("Enter space-separated <userID> and <activityID> (1=LOGIN, 2=LOGOUT):")
			_, err := fmt.Scan(&userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			currentTime := time.Now()
			activityTracker.StoreUserActivity(currentTime, userID, enum.Activity(activityID))
			fmt.Printf("Stored activity %d for user %d at %s\n", activityID, userID, currentTime.Format("15:04"))

		case 2: // GetActivityCountForUserInTimerange
			var startOffset, endOffset int
			var userID uint
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <userID> <activityID>:")
			fmt.Println("(e.g., '-10 0 1 1' checks from 10 mins ago until now)")
			_, err := fmt.Scan(&startOffset, &endOffset, &userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if count, err := activityTracker.GetActivityCountForUserInTimerange(startTime, endTime, userID, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Resulting Count: %d\n", count)
			}

		case 3: // GetActivityRateForUserInTimerange
			var startOffset, endOffset int
			var userID uint
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <userID> <activityID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if rate, err := activityTracker.GetActivityRateForUserInTimerange(startTime, endTime, userID, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Resulting Rate: %.4f\n", rate)
			}

		case 4: // GetDistinctUserCountByActivityInTimerange
			var startOffset, endOffset int
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <activityID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if distinctCount, err := activityTracker.GetDistinctUserCountByActivityInTimerange(startTime, endTime, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Distinct User Count: %d\n", distinctCount)
			}

		case 5: // GetUsersByActivityInTimerange
			var startOffset, endOffset int
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <activityID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if users, err := activityTracker.GetUsersByActivityInTimerange(startTime, endTime, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Matching User IDs: %v\n", users)
			}

		case 6: // GetTopKUsersForActivityInTimerange
			var startOffset, endOffset int
			var activityID int
			var k uint
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <activityID> <k>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &activityID, &k)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if topUsers, err := activityTracker.GetTopKUsersForActivityInTimerange(startTime, endTime, enum.Activity(activityID), k); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Top %d Users: %v\n", k, topUsers)
			}

		case 7: // GetUsersActivitySummaryInTimerange
			var startOffset, endOffset int
			var userID uint
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <userID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &userID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if summary, err := activityTracker.GetUsersActivitySummaryInTimerange(startTime, endTime, userID); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Activity History Summary: %v\n", summary)
			}

		case 8: // Exit
			fmt.Println("Exiting application...")
			return

		default:
			fmt.Println("Unknown operation number. Please select between 1 and 8.")
		}
	}
}
