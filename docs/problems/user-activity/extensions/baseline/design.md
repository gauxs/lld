---
title: User activity
problem: user_activity
extension: baseline
prev:
  text: requirements
  link: /problems/user-activity/extensions/baseline/requirements
next:
  text: codebase
  link: /problems/user-activity/extensions/baseline/codebase
---

## Core types

- `Activitystore` — Store all the activities with time as the first class dimensions
- `Record` — A single record present in Activitystore. This holds the users activity data
- `User` — An user.
- `Activity` — An activity. This will be simple enum.
- `ActivityTracker` — Holds ActivityStore and provides activity tracking functionality.

## API

### ActivityStore

```text
Store(timestamp, userID, activityID)
  Store record for timestamp, user and activity

GetCountInTimerange(startTimestamp, endTimestamp, userID, activityID)
  Get the count of record in time range by userID and activityID

GetRateInTimerange(startTimestamp, endTimestamp, userID, activityID)
  Get the rate in time range by userID and activityID

GetDistinctCountInTimerange(startTimestamp, endTimestamp, activityID)
  Get the distinct count of record in time range by activityID

GetInTimerangeByActivityID(startTimestamp, endTimestamp, activityID)
  Get records in time range by activityID

GetTopKInTimerange(startTimestamp, endTimestamp, activityID)
  Get top K frequency records by userID in time range for activityID

GetInTimerangeByUserID(startTimestamp, endTimestamp, userID)
  Get records in time range by userID
```

### ActivityTracker

```text
StoreUserActivity(timestamp, userID, activityID)
  Store users activity for timestamp

GetActivityCountForUserInTimerange(startTimestamp, endTimestamp, userID, activityID)
  Get activity count for user in time range

GetActivityRateForUserInTimerange(startTimestamp, endTimestamp, userID, activityID)
  Get activity rate for user in time range

GetDistinctUserCountByActivityInTimerange(startTimestamp, endTimestamp, activityID)
  Get distinct user count by activity in time range

GetUsersByActivityInTimerange(startTimestamp, endTimestamp, activityID)
  Get users by activity in time range

GetTopKUsersForActivityInTimerange(startTimestamp, endTimestamp, activityID)
  Get topK users for activity in time range

GetUsersActivitySummaryInTimerange(startTimestamp, endTimestamp, userID)
  Get users activity summary in time range
```

## Relationships

1. Activitystore stores Record
2. User represents a user
3. Activity will be an enum representing an activity
4. ActivityTracker holds ActivityStore.
