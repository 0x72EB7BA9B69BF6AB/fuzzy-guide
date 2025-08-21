package models

import (
	"sync"
	"time"
)

// Store provides in-memory storage for bouquets and users
type Store struct {
	bouquets    map[int]Bouquet
	users       map[int]User
	nextBouquetID int
	nextUserID   int
	mutex       sync.RWMutex
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	store := &Store{
		bouquets:      make(map[int]Bouquet),
		users:         make(map[int]User),
		nextBouquetID: 1,
		nextUserID:    1,
	}
	
	// Add some sample data
	store.initSampleData()
	return store
}

// initSampleData adds initial sample data
func (s *Store) initSampleData() {
	now := time.Now()
	
	// Sample bouquets
	s.bouquets[1] = Bouquet{
		ID:          1,
		Name:        "Basic Package",
		Description: "Essential channels for everyday viewing",
		Channels: []Channel{
			{Name: "BBC One", Manifest: "https://manifest.bbc.co.uk/bbc1/manifest.mpd", KeyKid: "bbc1-key-001"},
			{Name: "BBC Two", Manifest: "https://manifest.bbc.co.uk/bbc2/manifest.mpd", KeyKid: "bbc2-key-001"},
			{Name: "ITV", Manifest: "https://manifest.itv.com/itv1/manifest.mpd", KeyKid: "itv1-key-001"},
			{Name: "Channel 4", Manifest: "https://manifest.channel4.com/c4/manifest.mpd", KeyKid: "c4-key-001"},
		},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.bouquets[2] = Bouquet{
		ID:          2,
		Name:        "Premium Package", 
		Description: "Complete entertainment experience with sports and movies",
		Channels: []Channel{
			{Name: "BBC One", Manifest: "https://manifest.bbc.co.uk/bbc1/manifest.mpd", KeyKid: "bbc1-key-001"},
			{Name: "BBC Two", Manifest: "https://manifest.bbc.co.uk/bbc2/manifest.mpd", KeyKid: "bbc2-key-001"},
			{Name: "Sky Sports", Manifest: "https://manifest.sky.com/sports/manifest.mpd", KeyKid: "sky-sports-key-001"},
			{Name: "Sky Movies", Manifest: "https://manifest.sky.com/movies/manifest.mpd", KeyKid: "sky-movies-key-001"},
			{Name: "Discovery", Manifest: "https://manifest.discovery.com/main/manifest.mpd", KeyKid: "discovery-key-001"},
		},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.nextBouquetID = 3
	
	// Sample users
	s.users[1] = User{
		ID:        1,
		Username:  "admin",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Role:      "Administrator",
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users[2] = User{
		ID:        2,
		Username:  "user1",
		Email:     "user1@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "User",
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.nextUserID = 3
}

// Bouquet operations
func (s *Store) GetAllBouquets() []Bouquet {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	bouquets := make([]Bouquet, 0, len(s.bouquets))
	for _, bouquet := range s.bouquets {
		bouquets = append(bouquets, bouquet)
	}
	return bouquets
}

func (s *Store) GetBouquet(id int) (Bouquet, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	bouquet, exists := s.bouquets[id]
	return bouquet, exists
}

func (s *Store) CreateBouquet(bouquet Bouquet) Bouquet {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	bouquet.ID = s.nextBouquetID
	bouquet.CreatedAt = time.Now()
	bouquet.UpdatedAt = time.Now()
	s.bouquets[bouquet.ID] = bouquet
	s.nextBouquetID++
	return bouquet
}

func (s *Store) UpdateBouquet(bouquet Bouquet) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.bouquets[bouquet.ID]; !exists {
		return false
	}
	bouquet.UpdatedAt = time.Now()
	s.bouquets[bouquet.ID] = bouquet
	return true
}

func (s *Store) DeleteBouquet(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.bouquets[id]; !exists {
		return false
	}
	delete(s.bouquets, id)
	return true
}

// User operations
func (s *Store) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *Store) GetUser(id int) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	user, exists := s.users[id]
	return user, exists
}

func (s *Store) CreateUser(user User) User {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	user.ID = s.nextUserID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	s.nextUserID++
	return user
}

func (s *Store) UpdateUser(user User) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.users[user.ID]; !exists {
		return false
	}
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	return true
}

func (s *Store) DeleteUser(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.users[id]; !exists {
		return false
	}
	delete(s.users, id)
	return true
}

// Global store instance
var GlobalStore = NewStore()