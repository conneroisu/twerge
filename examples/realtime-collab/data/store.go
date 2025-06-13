package data

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/conneroisu/twerge/examples/realtime-collab/types"
)

// Store provides in-memory data storage for the collaboration system
type Store struct {
	users         map[string]types.User
	projects      map[string]types.Project
	documents     map[string]types.Document
	comments      map[string][]types.Comment
	notifications map[string][]types.Notification
	activities    map[string][]types.ActivityEvent
	
	// Indexes for efficient queries
	userProjects    map[string][]string // userID -> projectIDs
	projectDocuments map[string][]string // projectID -> documentIDs
	projectMembers   map[string][]types.ProjectMember
	
	mu sync.RWMutex
}

// NewStore creates a new data store with sample data
func NewStore() *Store {
	store := &Store{
		users:            make(map[string]types.User),
		projects:         make(map[string]types.Project),
		documents:        make(map[string]types.Document),
		comments:         make(map[string][]types.Comment),
		notifications:    make(map[string][]types.Notification),
		activities:       make(map[string][]types.ActivityEvent),
		userProjects:     make(map[string][]string),
		projectDocuments: make(map[string][]string),
		projectMembers:   make(map[string][]types.ProjectMember),
	}
	
	store.initializeSampleData()
	return store
}

// Initialize the store with comprehensive sample data
func (s *Store) initializeSampleData() {
	now := time.Now()
	
	// Create sample users
	users := []types.User{
		{
			ID:       "user-1",
			Username: "alice.dev",
			Email:    "alice@example.com",
			Avatar:   "https://images.unsplash.com/photo-1494790108755-2616b612b786?w=150",
			Status:   types.StatusOnline,
			Role:     types.RoleAdmin,
			IsOnline: true,
			JoinedAt: now.AddDate(0, -6, 0),
		},
		{
			ID:       "user-2",
			Username: "bob.designer",
			Email:    "bob@example.com",
			Avatar:   "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150",
			Status:   types.StatusOnline,
			Role:     types.RoleEditor,
			IsOnline: true,
			JoinedAt: now.AddDate(0, -4, 0),
		},
		{
			ID:       "user-3",
			Username: "carol.pm",
			Email:    "carol@example.com",
			Avatar:   "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150",
			Status:   types.StatusAway,
			Role:     types.RoleModerator,
			IsOnline: false,
			JoinedAt: now.AddDate(0, -3, 0),
		},
		{
			ID:       "user-4",
			Username: "david.analyst",
			Email:    "david@example.com",
			Avatar:   "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150",
			Status:   types.StatusOnline,
			Role:     types.RoleViewer,
			IsOnline: true,
			JoinedAt: now.AddDate(0, -2, 0),
		},
		{
			ID:       "user-5",
			Username: "eve.writer",
			Email:    "eve@example.com",
			Avatar:   "https://images.unsplash.com/photo-1517841905240-472988babdf9?w=150",
			Status:   types.StatusBusy,
			Role:     types.RoleEditor,
			IsOnline: true,
			JoinedAt: now.AddDate(0, -1, 0),
		},
	}
	
	for _, user := range users {
		s.users[user.ID] = user
	}
	
	// Create sample projects
	projects := []types.Project{
		{
			ID:          "proj-1",
			Name:        "Design System Overhaul",
			Description: "Complete redesign of the company design system with new components, tokens, and documentation.",
			Owner:       users[0],
			Status:      types.ProjectStatusActive,
			Priority:    types.PriorityHigh,
			CreatedAt:   now.AddDate(0, -2, -15),
			UpdatedAt:   now.Add(-2 * time.Hour),
			DeadlineAt:  &[]time.Time{now.AddDate(0, 1, 0)}[0],
			Tags:        []string{"design", "frontend", "components"},
			Settings: types.ProjectSettings{
				IsPublic:              true,
				AllowGuestEditing:     false,
				RequireApproval:       true,
				EnableVersionControl:  true,
				EnableComments:        true,
				EnableNotifications:   true,
				AutoSaveInterval:      30,
				MaxConcurrentEditors:  5,
				AllowedFileTypes:      []string{"md", "figma", "sketch"},
				MaxFileSize:           52428800, // 50MB
			},
			Analytics: types.ProjectAnalytics{
				TotalViews:          1247,
				TotalEdits:          89,
				TotalComments:       23,
				TotalCollaborators:  4,
				LastActivityAt:      now.Add(-30 * time.Minute),
				AverageSessionTime:  1800, // 30 minutes
				PeakConcurrentUsers: 3,
				DocumentCount:       8,
				CompletionRate:      65.5,
			},
		},
		{
			ID:          "proj-2",
			Name:        "API Documentation Portal",
			Description: "Interactive API documentation with examples, authentication guides, and SDK references.",
			Owner:       users[1],
			Status:      types.ProjectStatusActive,
			Priority:    types.PriorityMedium,
			CreatedAt:   now.AddDate(0, -1, -20),
			UpdatedAt:   now.Add(-4 * time.Hour),
			Tags:        []string{"documentation", "api", "developer-tools"},
			Settings: types.ProjectSettings{
				IsPublic:              true,
				AllowGuestEditing:     true,
				RequireApproval:       false,
				EnableVersionControl:  true,
				EnableComments:        true,
				EnableNotifications:   true,
				AutoSaveInterval:      60,
				MaxConcurrentEditors:  10,
				AllowedFileTypes:      []string{"md", "yaml", "json"},
				MaxFileSize:           10485760, // 10MB
			},
			Analytics: types.ProjectAnalytics{
				TotalViews:          892,
				TotalEdits:          156,
				TotalComments:       45,
				TotalCollaborators:  6,
				LastActivityAt:      now.Add(-1 * time.Hour),
				AverageSessionTime:  2100, // 35 minutes
				PeakConcurrentUsers: 4,
				DocumentCount:       12,
				CompletionRate:      78.2,
			},
		},
		{
			ID:          "proj-3",
			Name:        "User Research Findings",
			Description: "Compilation of user interviews, surveys, and usability testing results for Q4 2024.",
			Owner:       users[2],
			Status:      types.ProjectStatusOnHold,
			Priority:    types.PriorityLow,
			CreatedAt:   now.AddDate(0, -3, -5),
			UpdatedAt:   now.Add(-24 * time.Hour),
			Tags:        []string{"research", "ux", "analytics"},
			Settings: types.ProjectSettings{
				IsPublic:              false,
				AllowGuestEditing:     false,
				RequireApproval:       true,
				EnableVersionControl:  false,
				EnableComments:        true,
				EnableNotifications:   false,
				AutoSaveInterval:      120,
				MaxConcurrentEditors:  3,
				AllowedFileTypes:      []string{"md", "pdf", "xlsx"},
				MaxFileSize:           104857600, // 100MB
			},
			Analytics: types.ProjectAnalytics{
				TotalViews:          234,
				TotalEdits:          67,
				TotalComments:       12,
				TotalCollaborators:  3,
				LastActivityAt:      now.Add(-24 * time.Hour),
				AverageSessionTime:  900, // 15 minutes
				PeakConcurrentUsers: 2,
				DocumentCount:       5,
				CompletionRate:      45.0,
			},
		},
		{
			ID:          "proj-4",
			Name:        "Mobile App Prototype",
			Description: "Interactive prototype for the new mobile application with user flows and animations.",
			Owner:       users[3],
			Status:      types.ProjectStatusDraft,
			Priority:    types.PriorityCritical,
			CreatedAt:   now.Add(-72 * time.Hour),
			UpdatedAt:   now.Add(-6 * time.Hour),
			DeadlineAt:  &[]time.Time{now.AddDate(0, 0, 14)}[0],
			Tags:        []string{"mobile", "prototype", "ux"},
			Settings: types.ProjectSettings{
				IsPublic:              false,
				AllowGuestEditing:     false,
				RequireApproval:       true,
				EnableVersionControl:  true,
				EnableComments:        true,
				EnableNotifications:   true,
				AutoSaveInterval:      15,
				MaxConcurrentEditors:  8,
				AllowedFileTypes:      []string{"figma", "sketch", "principle"},
				MaxFileSize:           209715200, // 200MB
			},
			Analytics: types.ProjectAnalytics{
				TotalViews:          567,
				TotalEdits:          234,
				TotalComments:       78,
				TotalCollaborators:  5,
				LastActivityAt:      now.Add(-6 * time.Hour),
				AverageSessionTime:  3600, // 60 minutes
				PeakConcurrentUsers: 5,
				DocumentCount:       15,
				CompletionRate:      25.7,
			},
		},
	}
	
	for _, project := range projects {
		s.projects[project.ID] = project
	}
	
	// Create project members
	s.createProjectMembers(projects, users)
	
	// Create sample documents
	s.createSampleDocuments(projects, users, now)
	
	// Create sample activities
	s.createSampleActivities(projects, users, now)
	
	// Create sample notifications
	s.createSampleNotifications(users, now)
}

func (s *Store) createProjectMembers(projects []types.Project, users []types.User) {
	// Project 1 members
	s.projectMembers["proj-1"] = []types.ProjectMember{
		{User: users[0], Role: types.ProjectRoleOwner, JoinedAt: projects[0].CreatedAt, IsActive: true},
		{User: users[1], Role: types.ProjectRoleCollaborator, JoinedAt: projects[0].CreatedAt.Add(24 * time.Hour), IsActive: true},
		{User: users[2], Role: types.ProjectRoleReviewer, JoinedAt: projects[0].CreatedAt.Add(48 * time.Hour), IsActive: true},
		{User: users[4], Role: types.ProjectRoleCollaborator, JoinedAt: projects[0].CreatedAt.Add(72 * time.Hour), IsActive: true},
	}
	
	// Project 2 members  
	s.projectMembers["proj-2"] = []types.ProjectMember{
		{User: users[1], Role: types.ProjectRoleOwner, JoinedAt: projects[1].CreatedAt, IsActive: true},
		{User: users[0], Role: types.ProjectRoleAdmin, JoinedAt: projects[1].CreatedAt.Add(12 * time.Hour), IsActive: true},
		{User: users[3], Role: types.ProjectRoleCollaborator, JoinedAt: projects[1].CreatedAt.Add(36 * time.Hour), IsActive: true},
		{User: users[4], Role: types.ProjectRoleReviewer, JoinedAt: projects[1].CreatedAt.Add(60 * time.Hour), IsActive: true},
	}
	
	// Project 3 members
	s.projectMembers["proj-3"] = []types.ProjectMember{
		{User: users[2], Role: types.ProjectRoleOwner, JoinedAt: projects[2].CreatedAt, IsActive: true},
		{User: users[3], Role: types.ProjectRoleCollaborator, JoinedAt: projects[2].CreatedAt.Add(24 * time.Hour), IsActive: true},
		{User: users[4], Role: types.ProjectRoleObserver, JoinedAt: projects[2].CreatedAt.Add(48 * time.Hour), IsActive: false},
	}
	
	// Project 4 members
	s.projectMembers["proj-4"] = []types.ProjectMember{
		{User: users[3], Role: types.ProjectRoleOwner, JoinedAt: projects[3].CreatedAt, IsActive: true},
		{User: users[0], Role: types.ProjectRoleAdmin, JoinedAt: projects[3].CreatedAt.Add(6 * time.Hour), IsActive: true},
		{User: users[1], Role: types.ProjectRoleCollaborator, JoinedAt: projects[3].CreatedAt.Add(12 * time.Hour), IsActive: true},
		{User: users[2], Role: types.ProjectRoleReviewer, JoinedAt: projects[3].CreatedAt.Add(18 * time.Hour), IsActive: true},
		{User: users[4], Role: types.ProjectRoleCollaborator, JoinedAt: projects[3].CreatedAt.Add(24 * time.Hour), IsActive: true},
	}
	
	// Update user projects index
	for projectID, members := range s.projectMembers {
		for _, member := range members {
			s.userProjects[member.User.ID] = append(s.userProjects[member.User.ID], projectID)
		}
	}
}

func (s *Store) createSampleDocuments(projects []types.Project, users []types.User, now time.Time) {
	documents := []types.Document{
		{
			ID:          "doc-1",
			ProjectID:   "proj-1",
			Title:       "Design Tokens Specification",
			Content:     "# Design Tokens\n\nThis document outlines the color palette, typography scale, spacing system, and component tokens for our design system...",
			Type:        types.DocumentTypeMarkdown,
			Status:      types.DocumentStatusPublished,
			Author:      users[0],
			Editors:     []types.User{users[1], users[2]},
			CreatedAt:   now.Add(-48 * time.Hour),
			UpdatedAt:   now.Add(-2 * time.Hour),
			WordCount:   1250,
			ReadTime:    6,
		},
		{
			ID:          "doc-2",
			ProjectID:   "proj-1",
			Title:       "Component Library Guidelines",
			Content:     "# Component Guidelines\n\nStandards for creating reusable UI components...",
			Type:        types.DocumentTypeMarkdown,
			Status:      types.DocumentStatusReview,
			Author:      users[1],
			Editors:     []types.User{users[0]},
			CreatedAt:   now.Add(-36 * time.Hour),
			UpdatedAt:   now.Add(-4 * time.Hour),
			WordCount:   2100,
			ReadTime:    11,
		},
		{
			ID:          "doc-3",
			ProjectID:   "proj-2",
			Title:       "Authentication API",
			Content:     "# Authentication\n\n## Overview\nOur API uses OAuth 2.0 for authentication...",
			Type:        types.DocumentTypeMarkdown,
			Status:      types.DocumentStatusPublished,
			Author:      users[1],
			Editors:     []types.User{users[3], users[4]},
			CreatedAt:   now.Add(-72 * time.Hour),
			UpdatedAt:   now.Add(-1 * time.Hour),
			WordCount:   3200,
			ReadTime:    16,
		},
		{
			ID:          "doc-4",
			ProjectID:   "proj-4",
			Title:       "User Flow Specifications",
			Content:     "# Mobile App User Flows\n\nDetailed user journey mapping for the mobile application...",
			Type:        types.DocumentTypeMarkdown,
			Status:      types.DocumentStatusDraft,
			Author:      users[3],
			Editors:     []types.User{users[0], users[1], users[2]},
			CreatedAt:   now.Add(-24 * time.Hour),
			UpdatedAt:   now.Add(-30 * time.Minute),
			WordCount:   1800,
			ReadTime:    9,
		},
	}
	
	for _, doc := range documents {
		s.documents[doc.ID] = doc
		s.projectDocuments[doc.ProjectID] = append(s.projectDocuments[doc.ProjectID], doc.ID)
	}
}

func (s *Store) createSampleActivities(projects []types.Project, users []types.User, now time.Time) {
	activities := []types.ActivityEvent{
		{
			ID:          "activity-1",
			Type:        types.ActivityDocumentUpdated,
			Actor:       users[0],
			Target:      types.ActivityTarget{Type: "document", ID: "doc-1", Name: "Design Tokens Specification"},
			Description: "updated the design tokens documentation",
			CreatedAt:   now.Add(-2 * time.Hour),
			IsImportant: false,
		},
		{
			ID:          "activity-2",
			Type:        types.ActivityCommentAdded,
			Actor:       users[1],
			Target:      types.ActivityTarget{Type: "document", ID: "doc-2", Name: "Component Library Guidelines"},
			Description: "added a comment on component naming conventions",
			CreatedAt:   now.Add(-4 * time.Hour),
			IsImportant: false,
		},
		{
			ID:          "activity-3",
			Type:        types.ActivityProjectUpdated,
			Actor:       users[2],
			Target:      types.ActivityTarget{Type: "project", ID: "proj-1", Name: "Design System Overhaul"},
			Description: "updated project timeline and milestones",
			CreatedAt:   now.Add(-6 * time.Hour),
			IsImportant: true,
		},
		{
			ID:          "activity-4",
			Type:        types.ActivityMemberAdded,
			Actor:       users[3],
			Target:      types.ActivityTarget{Type: "project", ID: "proj-4", Name: "Mobile App Prototype"},
			Description: "added eve.writer as a collaborator",
			CreatedAt:   now.Add(-8 * time.Hour),
			IsImportant: false,
		},
	}
	
	// Store activities globally and per project
	for _, activity := range activities {
		s.activities["global"] = append(s.activities["global"], activity)
		if activity.Target.Type == "project" {
			s.activities[activity.Target.ID] = append(s.activities[activity.Target.ID], activity)
		}
	}
}

func (s *Store) createSampleNotifications(users []types.User, now time.Time) {
	notifications := []types.Notification{
		{
			ID:       "notif-1",
			UserID:   users[0].ID,
			Type:     types.NotificationTypeInfo,
			Title:    "New comment on Design Tokens",
			Message:  "Bob added feedback on the color palette section",
			Priority: types.NotificationPriorityMedium,
			IsRead:   false,
			CreatedAt: now.Add(-1 * time.Hour),
		},
		{
			ID:       "notif-2",
			UserID:   users[0].ID,
			Type:     types.NotificationTypeSuccess,
			Title:    "Document published",
			Message:  "Authentication API documentation is now live",
			Priority: types.NotificationPriorityLow,
			IsRead:   true,
			CreatedAt: now.Add(-3 * time.Hour),
			ReadAt:   &[]time.Time{now.Add(-2 * time.Hour)}[0],
		},
		{
			ID:       "notif-3",
			UserID:   users[1].ID,
			Type:     types.NotificationTypeWarning,
			Title:    "Approaching deadline",
			Message:  "Mobile App Prototype deadline is in 2 weeks",
			Priority: types.NotificationPriorityHigh,
			IsRead:   false,
			CreatedAt: now.Add(-30 * time.Minute),
		},
	}
	
	for _, notif := range notifications {
		s.notifications[notif.UserID] = append(s.notifications[notif.UserID], notif)
	}
}

// User methods
func (s *Store) GetDemoUser() types.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users["user-1"] // Return Alice as the demo user
}

func (s *Store) GetUser(userID string) (types.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if user, exists := s.users[userID]; exists {
		return user, nil
	}
	return types.User{}, fmt.Errorf("user not found")
}

// Project methods
func (s *Store) GetProject(projectID string) (types.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if project, exists := s.projects[projectID]; exists {
		return project, nil
	}
	return types.Project{}, fmt.Errorf("project not found")
}

func (s *Store) GetUserProjects(userID string) []types.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	projectIDs := s.userProjects[userID]
	var projects []types.Project
	
	for _, projectID := range projectIDs {
		if project, exists := s.projects[projectID]; exists {
			projects = append(projects, project)
		}
	}
	
	// Sort by last updated
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].UpdatedAt.After(projects[j].UpdatedAt)
	})
	
	return projects
}

func (s *Store) GetProjectMembers(projectID string) []types.ProjectMember {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.projectMembers[projectID]
}

func (s *Store) CreateProject(project types.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.projects[project.ID] = project
	s.userProjects[project.Owner.ID] = append(s.userProjects[project.Owner.ID], project.ID)
	
	return nil
}

func (s *Store) UpdateProject(projectID string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	project, exists := s.projects[projectID]
	if !exists {
		return fmt.Errorf("project not found")
	}
	
	// Apply updates (simplified - in production would use reflection or specific setters)
	project.UpdatedAt = time.Now()
	s.projects[projectID] = project
	
	return nil
}

// Document methods
func (s *Store) GetDocument(documentID string) (types.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if document, exists := s.documents[documentID]; exists {
		return document, nil
	}
	return types.Document{}, fmt.Errorf("document not found")
}

func (s *Store) GetProjectDocuments(projectID string) []types.Document {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	documentIDs := s.projectDocuments[projectID]
	var documents []types.Document
	
	for _, docID := range documentIDs {
		if doc, exists := s.documents[docID]; exists {
			documents = append(documents, doc)
		}
	}
	
	// Sort by last updated
	sort.Slice(documents, func(i, j int) bool {
		return documents[i].UpdatedAt.After(documents[j].UpdatedAt)
	})
	
	return documents
}

func (s *Store) CreateDocument(document types.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.documents[document.ID] = document
	s.projectDocuments[document.ProjectID] = append(s.projectDocuments[document.ProjectID], document.ID)
	
	return nil
}

// Notification methods
func (s *Store) GetUserNotifications(userID string) []types.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	notifications := s.notifications[userID]
	
	// Sort by created date (newest first)
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].CreatedAt.After(notifications[j].CreatedAt)
	})
	
	return notifications
}

// Search functionality
func (s *Store) Search(query string, userID string) []types.SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []types.SearchResult
	query = strings.ToLower(query)
	
	// Search projects
	for _, project := range s.projects {
		if strings.Contains(strings.ToLower(project.Name), query) ||
		   strings.Contains(strings.ToLower(project.Description), query) {
			results = append(results, types.SearchResult{
				Type:        types.SearchResultTypeProject,
				ID:          project.ID,
				Title:       project.Name,
				Description: project.Description,
				URL:         fmt.Sprintf("/projects/%s", project.ID),
				Score:       0.9,
				CreatedAt:   project.CreatedAt,
			})
		}
	}
	
	// Search documents
	for _, document := range s.documents {
		if strings.Contains(strings.ToLower(document.Title), query) ||
		   strings.Contains(strings.ToLower(document.Content), query) {
			results = append(results, types.SearchResult{
				Type:        types.SearchResultTypeDocument,
				ID:          document.ID,
				Title:       document.Title,
				Description: truncateString(document.Content, 150),
				URL:         fmt.Sprintf("/documents/%s", document.ID),
				Score:       0.8,
				CreatedAt:   document.CreatedAt,
			})
		}
	}
	
	// Search users
	for _, user := range s.users {
		if strings.Contains(strings.ToLower(user.Username), query) ||
		   strings.Contains(strings.ToLower(user.Email), query) {
			results = append(results, types.SearchResult{
				Type:        types.SearchResultTypeUser,
				ID:          user.ID,
				Title:       user.Username,
				Description: user.Email,
				URL:         fmt.Sprintf("/users/%s", user.ID),
				Score:       0.7,
				CreatedAt:   user.JoinedAt,
			})
		}
	}
	
	// Sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	return results
}

// Dashboard and analytics
func (s *Store) GetDashboardStats() types.DashboardStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	stats := types.DashboardStats{
		TotalProjects:    len(s.projects),
		TotalUsers:       len(s.users),
		TotalDocuments:   len(s.documents),
		ProjectsByStatus: make(map[types.ProjectStatus]int),
		ProjectsByPriority: make(map[types.ProjectPriority]int),
	}
	
	// Count active projects and aggregate stats
	for _, project := range s.projects {
		if project.Status == types.ProjectStatusActive {
			stats.ActiveProjects++
		}
		stats.ProjectsByStatus[project.Status]++
		stats.ProjectsByPriority[project.Priority]++
		stats.TotalComments += project.Analytics.TotalComments
	}
	
	// Count online users
	for _, user := range s.users {
		if user.IsOnline {
			stats.OnlineUsers++
		}
	}
	
	// Get recent activity
	if activities, exists := s.activities["global"]; exists {
		// Sort by created date and take the most recent 10
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].CreatedAt.After(activities[j].CreatedAt)
		})
		
		maxItems := 10
		if len(activities) < maxItems {
			maxItems = len(activities)
		}
		stats.RecentActivity = activities[:maxItems]
	}
	
	// Get top projects (by recent activity)
	var topProjects []types.Project
	for _, project := range s.projects {
		topProjects = append(topProjects, project)
	}
	sort.Slice(topProjects, func(i, j int) bool {
		return topProjects[i].Analytics.LastActivityAt.After(topProjects[j].Analytics.LastActivityAt)
	})
	
	maxProjects := 5
	if len(topProjects) < maxProjects {
		maxProjects = len(topProjects)
	}
	stats.TopProjects = topProjects[:maxProjects]
	
	// Get top users (by online status and role)
	var topUsers []types.User
	for _, user := range s.users {
		topUsers = append(topUsers, user)
	}
	sort.Slice(topUsers, func(i, j int) bool {
		if topUsers[i].IsOnline != topUsers[j].IsOnline {
			return topUsers[i].IsOnline
		}
		return topUsers[i].Role < topUsers[j].Role // Admin roles first
	})
	
	maxUsers := 5
	if len(topUsers) < maxUsers {
		maxUsers = len(topUsers)
	}
	stats.TopUsers = topUsers[:maxUsers]
	
	// Mock system health
	stats.SystemHealth = types.SystemHealth{
		CPUUsage:          25.7,
		MemoryUsage:       68.3,
		DiskUsage:         45.2,
		ActiveConnections: stats.OnlineUsers,
		Uptime:            86400 * 7, // 7 days
		LastChecked:       time.Now(),
		Status:            types.HealthStatusHealthy,
	}
	
	return stats
}

func (s *Store) GetSystemHealth() types.SystemHealth {
	return types.SystemHealth{
		CPUUsage:          25.7,
		MemoryUsage:       68.3,
		DiskUsage:         45.2,
		ActiveConnections: 5,
		Uptime:            86400 * 7, // 7 days
		LastChecked:       time.Now(),
		Status:            types.HealthStatusHealthy,
	}
}

// Permission methods
func (s *Store) HasProjectPermission(userID, projectID, permission string) bool {
	members := s.projectMembers[projectID]
	for _, member := range members {
		if member.User.ID == userID && member.IsActive {
			// Simplified permission check - in production would be more granular
			return member.Role == types.ProjectRoleOwner || 
				   member.Role == types.ProjectRoleAdmin ||
				   (member.Role == types.ProjectRoleCollaborator && permission != "delete")
		}
	}
	return false
}

func (s *Store) HasDocumentPermission(userID, documentID, permission string) bool {
	document, exists := s.documents[documentID]
	if !exists {
		return false
	}
	
	// Check if user has project permission
	return s.HasProjectPermission(userID, document.ProjectID, permission)
}

// Utility methods
func (s *Store) GenerateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}