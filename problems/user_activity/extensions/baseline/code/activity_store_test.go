package code

import (
	"testing"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

var (
	OneHrBack        = time.Now().Add(-time.Hour)
	OneAndHalfHrBack = time.Now().Add(-1 * time.Hour).Add(-30 * time.Minute)
	TwoHrBack        = time.Now().Add(-2 * time.Hour)
)

func TestStore_InvalidInput(t *testing.T) {
	as := NewActivityStore(DefaultNumberOfMinutesInWindow)
	err := as.Store(time.Now(), 0, enum.LOGIN)

	if err == nil {
		t.Errorf("%s | storing activity should not have been allowed", t.Name())
	}
}

func TestStore_GetCountInTimerange(t *testing.T) {
	as := NewActivityStore(DefaultNumberOfMinutesInWindow)

	userID := uint(1)
	activity := enum.LOGIN
	if err := as.Store(OneAndHalfHrBack, userID, activity); err != nil {
		t.Errorf("%s | error in storing user %d activity %v at time %v", t.Name(), userID, activity, OneAndHalfHrBack)
	}

	if err := as.Store(OneAndHalfHrBack.Add(time.Minute), userID, activity); err != nil {
		t.Errorf("%s | error in storing user %d activity %v at time %v", t.Name(), userID, activity, OneAndHalfHrBack)
	}

	expectedActivityCount := uint(2)
	if count := as.GetCountInTimerange(TwoHrBack, OneHrBack, userID, activity); count != expectedActivityCount {
		t.Errorf("%s | unexpected activity count for user %d activity %v between time %v - %v. Expected %v, got %v",
			t.Name(), userID, activity, TwoHrBack, OneHrBack, expectedActivityCount, count)
	}
}
