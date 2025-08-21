package models

import (
	"sync"
	"time"
)

// Store provides in-memory storage for bouquets, users, channels, and providers
type Store struct {
	bouquets      map[int]Bouquet
	users         map[int]User
	channels      map[int]Channel
	providers     map[int]Provider
	nextBouquetID int
	nextUserID    int
	nextChannelID int
	nextProviderID int
	mutex         sync.RWMutex
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	store := &Store{
		bouquets:       make(map[int]Bouquet),
		users:          make(map[int]User),
		channels:       make(map[int]Channel),
		providers:      make(map[int]Provider),
		nextBouquetID:  1,
		nextUserID:     1,
		nextChannelID:  1,
		nextProviderID: 1,
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
	
	// Sample channels
	s.channels[1] = Channel{
		ID:        1,
		Name:      "BBC One",
		Manifest:  "https://manifest.bbc.co.uk/bbc1/manifest.mpd",
		KeyKid:    "bbc1-key-001",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.channels[2] = Channel{
		ID:        2,
		Name:      "BBC Two",
		Manifest:  "https://manifest.bbc.co.uk/bbc2/manifest.mpd",
		KeyKid:    "bbc2-key-001",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.channels[3] = Channel{
		ID:        3,
		Name:      "ITV",
		Manifest:  "https://manifest.itv.com/itv1/manifest.mpd",
		KeyKid:    "itv1-key-001",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.nextChannelID = 4
	
	// Sample providers
	s.providers[1] = Provider{
		ID:          1,
		Name:        "BBC",
		Description: "British Broadcasting Corporation",
		URL:         "https://www.bbc.co.uk",
		APIKey:      "bbc-api-key-123",
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.providers[2] = Provider{
		ID:          2,
		Name:        "Sky",
		Description: "Sky Television Services",
		URL:         "https://www.sky.com",
		APIKey:      "sky-api-key-456",
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.nextProviderID = 3
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

// Channel operations
func (s *Store) GetAllChannels() []Channel {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	channels := make([]Channel, 0, len(s.channels))
	for _, channel := range s.channels {
		channels = append(channels, channel)
	}
	return channels
}

func (s *Store) GetChannel(id int) (Channel, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	channel, exists := s.channels[id]
	return channel, exists
}

func (s *Store) CreateChannel(channel Channel) Channel {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	channel.ID = s.nextChannelID
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()
	s.channels[channel.ID] = channel
	s.nextChannelID++
	return channel
}

func (s *Store) UpdateChannel(channel Channel) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.channels[channel.ID]; !exists {
		return false
	}
	channel.UpdatedAt = time.Now()
	s.channels[channel.ID] = channel
	return true
}

func (s *Store) DeleteChannel(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.channels[id]; !exists {
		return false
	}
	delete(s.channels, id)
	return true
}

// Provider operations
func (s *Store) GetAllProviders() []Provider {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	providers := make([]Provider, 0, len(s.providers))
	for _, provider := range s.providers {
		providers = append(providers, provider)
	}
	return providers
}

func (s *Store) GetProvider(id int) (Provider, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	provider, exists := s.providers[id]
	return provider, exists
}

func (s *Store) CreateProvider(provider Provider) Provider {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	provider.ID = s.nextProviderID
	provider.CreatedAt = time.Now()
	provider.UpdatedAt = time.Now()
	s.providers[provider.ID] = provider
	s.nextProviderID++
	return provider
}

func (s *Store) UpdateProvider(provider Provider) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.providers[provider.ID]; !exists {
		return false
	}
	provider.UpdatedAt = time.Now()
	s.providers[provider.ID] = provider
	return true
}

func (s *Store) DeleteProvider(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.providers[id]; !exists {
		return false
	}
	delete(s.providers, id)
	return true
}

// Global store instance
var GlobalStore = NewStore()