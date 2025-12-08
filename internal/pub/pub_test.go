package pub

import (
	"testing"

	"github.com/CarlosEduardoAD/broadcast-server/internal/sub"
	"github.com/stretchr/testify/assert"
)

func TestWhenGivenAWrongTypeOfSubsItShouldReturnError(t *testing.T) {
	var invalid any = "invalid"

	defer func() {
		if r := recover(); r != nil {
			t.Log("Test passed, panic was caught!")
		} else {
			t.Fatalf("expected panic but none occurred")
		}
	}()

	NewPublisher(invalid.([]sub.Subscriber))
}

func TestUserCreatesPublisherWithValidData(t *testing.T) {
	pub := NewPublisher([]sub.Subscriber{})

	if pub == nil {
		t.Fatalf("expected non-nil publisher")
	}

	assert.NotNil(t, pub.Subscribers, "expected non-nil subscribers slice")
	assert.Equal(t, 0, len(pub.Subscribers), "expected 0 subscribers on new publisher")
}

func TestUserInsertsANonSubscriberIntoPublisherShouldPanic(t *testing.T) {
	// Similar to the other panic test: attempt invalid type assertions
	var invalid any = 123
	pub := NewPublisher([]sub.Subscriber{})

	defer func() {
		if r := recover(); r != nil {
			t.Log("panic caught as expected")
		} else {
			t.Fatalf("expected panic but none occurred")
		}
	}()

	pub.Subscribe(invalid.(sub.Subscriber))

}

func TestUserSubscribesSuccessfully(t *testing.T) {
	pub := NewPublisher([]sub.Subscriber{})

	s := sub.Subscriber{Name: "alice"}

	subs := pub.Subscribe(s)

	if len(subs) != 1 {
		t.Fatalf("expected 1 subscriber after subscribe; got %d", len(subs))
	}

	assert.Equal(t, 1, len(subs), "expected 1 subscriber after subscribe")
	assert.Equal(t, s, subs[0], "expected subscriber to match the one added")
}

func TestUserRemovesNonExistentSubscriberShouldPanic(t *testing.T) {
	pub := NewPublisher([]sub.Subscriber{})

	toRemove := sub.Subscriber{Name: "bob"}

	defer func() {
		if r := recover(); r != nil {
			t.Log("panic caught as expected when removing non-existent subscriber")
		} else {
			t.Fatalf("expected panic when removing non-existent subscriber")
		}
	}()

	pub.Remove(toRemove)
}

func TestUserRemovesExistingSubscriberSuccessfully(t *testing.T) {
	s := sub.Subscriber{Name: "carol"}
	pub := NewPublisher([]sub.Subscriber{s})

	if len(pub.Subscribers) != 1 {
		t.Fatalf("setup failed: expected 1 subscriber; got %d", len(pub.Subscribers))
	}

	pub.Remove(s)

	assert.Equal(t, 0, len(pub.Subscribers), "expected 0 subscribers after removal")
}
