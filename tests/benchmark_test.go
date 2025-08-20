package tests

import (
	"billing-service/internal/user"
	"billing-service/internal/user/domain"
	"billing-service/internal/user/event"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkCreateUser(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	testUser := domain.User{
		Balance:   100.0,
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userService.CreateUser(ctx, testUser)
		if err != nil {
			b.Fatalf("CreateUser failed: %v", err)
		}
	}
}

func BenchmarkGetUserByID(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 1000; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   100.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := domain.UserID((i % 1000) + 1)
		_, err := userService.GetUserByID(ctx, userID)
		if err != nil {
			b.Fatalf("GetUserByID failed: %v", err)
		}
	}
}

func BenchmarkCreditUserBalance(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 100; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   100.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := domain.UserID((i % 100) + 1)
		amount := float64(10 + (i % 90))
		err := userService.CreditUserBalance(ctx, userID, amount)
		if err != nil {
			b.Fatalf("CreditUserBalance failed: %v", err)
		}
	}
}

func BenchmarkDebitUserBalance(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 100; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   10000000.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := uint((i % 100) + 1)
		smsID := uint(i + 1)
		amount := float64(1 + (i % 10))

		eventData := event.UserBalanceEvent{
			UserID: userID,
			SMSID:  smsID,
			Amount: amount,
			Type:   event.SMSDebitEvent,
			Domain: event.SMS,
		}

		eventBytes, _ := json.Marshal(eventData)
		_, err := userService.DebitUserBalance(ctx, eventBytes)
		if err != nil {
			b.Fatalf("DebitUserBalance failed: %v", err)
		}
	}
}

func BenchmarkUpdateSMSStatus(b *testing.B) {
	mockRepo := NewMockUserRepo()
	mockRabbit := NewMockRabbit()

	userService := &MockUserService{
		repo:   mockRepo,
		rabbit: mockRabbit,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smsEvent := event.SMSUpdateEvent{
			SMSID:  uint(i + 1),
			Status: event.StatusSuccess,
			Domain: event.SMS,
		}

		err := userService.UpdateSMSStatus(ctx, smsEvent)
		if err != nil {
			b.Fatalf("UpdateSMSStatus failed: %v", err)
		}
	}
}

func BenchmarkConcurrentOperations(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 1000; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   1000.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			userID := domain.UserID((i % 1000) + 1)

			switch i % 3 {
			case 0:
				_, err := userService.GetUserByID(ctx, userID)
				if err != nil {
					b.Errorf("GetUserByID failed: %v", err)
				}
			case 1:
				err := userService.CreditUserBalance(ctx, userID, 10.0)
				if err != nil {
					b.Errorf("CreditUserBalance failed: %v", err)
				}
			case 2:
				eventData := event.UserBalanceEvent{
					UserID: uint(userID),
					SMSID:  uint(i + 1),
					Amount: 5.0,
					Type:   event.SMSDebitEvent,
					Domain: event.SMS,
				}
				eventBytes, _ := json.Marshal(eventData)
				_, err := userService.DebitUserBalance(ctx, eventBytes)
				if err != nil {
					b.Errorf("DebitUserBalance failed: %v", err)
				}
			}
			i++
		}
	})
}

func BenchmarkHighVolumeUserCreation(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			testUser := domain.User{
				Balance:   float64(100 + (i % 1000)),
				CreatedAt: time.Now(),
			}
			_, err := userService.CreateUser(ctx, testUser)
			if err != nil {
				b.Errorf("CreateUser failed: %v", err)
			}
			i++
		}
	})
}

func BenchmarkHighVolumeTransactions(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 10000; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   1000.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			userID := domain.UserID((i % 10000) + 1)

			if i%2 == 0 {
				amount := float64(1 + (i % 50))
				err := userService.CreditUserBalance(ctx, userID, amount)
				if err != nil {
					b.Errorf("CreditUserBalance failed: %v", err)
				}
			} else {
				eventData := event.UserBalanceEvent{
					UserID: uint(userID),
					SMSID:  uint(i + 1),
					Amount: float64(1 + (i % 20)),
					Type:   event.SMSDebitEvent,
					Domain: event.SMS,
				}
				eventBytes, _ := json.Marshal(eventData)
				_, _ = userService.DebitUserBalance(ctx, eventBytes)
			}
			i++
		}
	})
}

func BenchmarkStressTest(b *testing.B) {
	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 100000; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   1000000.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	var wg sync.WaitGroup
	concurrentUsers := 50

	b.ResetTimer()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < b.N/concurrentUsers; j++ {
				userID := domain.UserID((workerID*2000+j)%100000 + 1)

				switch j % 4 {
				case 0:
					_, _ = userService.GetUserByID(ctx, userID)
				case 1:
					_ = userService.CreditUserBalance(ctx, userID, 10.0)
				case 2:
					eventData := event.UserBalanceEvent{
						UserID: uint(userID),
						SMSID:  uint(j + 1),
						Amount: 0.5,
						Type:   event.SMSDebitEvent,
						Domain: event.SMS,
					}
					eventBytes, _ := json.Marshal(eventData)
					_, _ = userService.DebitUserBalance(ctx, eventBytes)
				case 3:
					newUser := domain.User{
						Balance:   1000.0,
						CreatedAt: time.Now(),
					}
					_, _ = userService.CreateUser(ctx, newUser)
				}
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkMemoryUsage(b *testing.B) {
	b.ReportAllocs()

	mockRepo := NewMockUserRepo()
	userService := user.NewService(mockRepo, nil)
	ctx := context.Background()

	users := make(map[uint]domain.User)
	for i := 1; i <= 1000; i++ {
		users[uint(i)] = domain.User{
			ID:        domain.UserID(i),
			Balance:   1000.0,
			APIKey:    domain.APIKey(fmt.Sprintf("api-key-%d", i)),
			CreatedAt: time.Now(),
		}
	}
	mockRepo.SetUsers(users)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := domain.UserID((i % 1000) + 1)

		user, _ := userService.GetUserByID(ctx, userID)
		_ = user
		err := userService.CreditUserBalance(ctx, userID, 1.0)
		if err != nil {
			b.Fatalf("Operation failed: %v", err)
		}
	}
}
