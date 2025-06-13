//go:build ignore

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/conneroisu/twerge"
	"github.com/conneroisu/twerge/examples/realtime-collab/data"
	"github.com/conneroisu/twerge/examples/realtime-collab/types"
	"github.com/conneroisu/twerge/examples/realtime-collab/views"
)

func main() {
	log.Println("🎨 Generating optimized classes for realtime collaboration dashboard...")
	
	// Initialize data store to get realistic sample data
	store := data.NewStore()
	
	// Get sample data for component generation
	demoUser := store.GetDemoUser()
	stats := store.GetDashboardStats()
	projects := store.GetUserProjects(demoUser.ID)
	notifications := store.GetUserNotifications(demoUser.ID)
	
	// Create comprehensive component instances covering all states
	log.Println("📊 Generating dashboard components...")
	dashboardComponents := []templ.Component{
		// Dashboard with different user states and data scenarios
		views.DashboardLayout(demoUser, stats, notifications),
		views.DashboardLayout(demoUser, stats, []types.Notification{}), // Empty notifications
		views.DashboardLayout(createOfflineUser(), stats, notifications), // Offline user
		views.DashboardLayout(createBusyUser(), createMinimalStats(), notifications), // Minimal data
		
		// Navigation states
		views.NavigationSidebar(demoUser, stats),
		views.NavigationSidebar(createOfflineUser(), createMinimalStats()),
		
		// System health widget in different states
		views.SystemHealthWidget(createHealthySystem()),
		views.SystemHealthWidget(createWarningSystem()),
		views.SystemHealthWidget(createCriticalSystem()),
		
		// Top bar with different notification counts
		views.TopBar(demoUser, notifications),
		views.TopBar(demoUser, []types.Notification{}),
		views.TopBar(demoUser, createManyNotifications()),
		
		// Dashboard content with various data states
		views.DashboardContent(stats, demoUser),
		views.DashboardContent(createMinimalStats(), demoUser),
		views.DashboardContent(createMaximalStats(), demoUser),
		
		// Welcome hero for different users
		views.WelcomeHero(demoUser),
		views.WelcomeHero(createNewUser()),
		views.WelcomeHero(createPowerUser()),
		
		// Stats overview with different values
		views.StatsOverview(stats),
		views.StatsOverview(createMinimalStats()),
		views.StatsOverview(createMaximalStats()),
		
		// Activity feed states
		views.RecentActivityFeed(stats.RecentActivity),
		views.RecentActivityFeed([]types.ActivityEvent{}), // Empty state
		views.RecentActivityFeed(createMaximalActivity()),
		
		// Project status charts with different distributions
		views.ProjectStatusChart(stats),
		views.ProjectStatusChart(createSkewedStats()),
		
		// Contributor lists in different configurations
		views.TopContributors(stats.TopUsers),
		views.TopContributors(createMixedOnlineUsers()),
		views.TopContributors(createAllOfflineUsers()),
		
		// Active projects grid with different project states
		views.ActiveProjects(projects),
		views.ActiveProjects(createVariedProjects()),
		views.ActiveProjects([]types.Project{}), // Empty state
		
		// Notification center states
		views.NotificationCenter(notifications),
		views.NotificationCenter(createManyNotifications()),
		views.NotificationCenter([]types.Notification{}),
	}
	
	log.Println("📄 Generating project view components...")
	projectComponents := []templ.Component{
		// Project views with different project types and member configurations
		views.ProjectView(projects[0], store.GetProjectDocuments(projects[0].ID), store.GetProjectMembers(projects[0].ID), demoUser),
		views.ProjectView(createLargeProject(), createManyDocuments(), createManyMembers(), demoUser),
		views.ProjectView(createMinimalProject(), []types.Document{}, createMinimalMembers(), demoUser),
		
		// Project sidebar states
		views.ProjectSidebar(projects[0], store.GetProjectDocuments(projects[0].ID), store.GetProjectMembers(projects[0].ID), demoUser),
		views.ProjectSidebar(createLargeProject(), createManyDocuments(), createManyMembers(), demoUser),
		
		// Project header configurations
		views.ProjectHeader(projects[0], demoUser),
		views.ProjectHeader(createUrgentProject(), demoUser),
		views.ProjectHeader(createCompletedProject(), demoUser),
		
		// Project content with different data loads
		views.ProjectContent(projects[0], store.GetProjectDocuments(projects[0].ID), store.GetProjectMembers(projects[0].ID)),
		views.ProjectContent(createLargeProject(), createManyDocuments(), createManyMembers()),
		views.ProjectContent(createMinimalProject(), []types.Document{}, createMinimalMembers()),
		
		// Project overview states
		views.ProjectOverview(projects[0]),
		views.ProjectOverview(createUrgentProject()),
		views.ProjectOverview(createCompletedProject()),
		
		// Quick stats variations
		views.QuickStats(projects[0], store.GetProjectDocuments(projects[0].ID), store.GetProjectMembers(projects[0].ID)),
		views.QuickStats(createLargeProject(), createManyDocuments(), createManyMembers()),
		
		// Document lists in different states
		views.RecentDocuments(store.GetProjectDocuments(projects[0].ID)),
		views.RecentDocuments(createManyDocuments()),
		views.RecentDocuments([]types.Document{}),
		
		// Timeline with different activity levels
		views.ProjectTimeline(stats.RecentActivity),
		views.ProjectTimeline(createMaximalActivity()),
		views.ProjectTimeline([]types.ActivityEvent{}),
		
		// Team member configurations
		views.TeamMembers(store.GetProjectMembers(projects[0].ID)),
		views.TeamMembers(createManyMembers()),
		views.TeamMembers(createMinimalMembers()),
		
		// Project tasks in various states
		views.ProjectTasks(projects[0]),
		views.ProjectTasks(createUrgentProject()),
	}
	
	// Combine all components
	allComponents := append(dashboardComponents, projectComponents...)
	
	log.Printf("🔄 Processing %d component instances...", len(allComponents))
	
	// Generate optimized classes, CSS, and HTML
	twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		"input.css",
		"classes/classes.html",
		allComponents...,
	)
	
	log.Println("✅ Code generation complete!")
	log.Println("📁 Generated files:")
	log.Println("   • classes/classes.go - Optimized Go class mappings")
	log.Println("   • input.css - CSS with @apply directives")
	log.Println("   • classes/classes.html - HTML for TailwindCSS purging")
	log.Println("")
	log.Println("🎯 Next steps:")
	log.Println("   1. Run: tailwindcss -i input.css -o _static/dist/styles.css --minify")
	log.Println("   2. Run: templ generate")
	log.Println("   3. Run: go run main.go")
	log.Println("   4. Open: http://localhost:8080")
}

// Helper functions to create various component states

func createOfflineUser() types.User {
	return types.User{
		ID:       "user-offline",
		Username: "offline.user",
		Email:    "offline@example.com",
		Avatar:   "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150",
		Status:   types.StatusOffline,
		Role:     types.RoleViewer,
		IsOnline: false,
		JoinedAt: time.Now().AddDate(0, -1, 0),
	}
}

func createBusyUser() types.User {
	return types.User{
		ID:       "user-busy",
		Username: "busy.developer",
		Email:    "busy@example.com",
		Avatar:   "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150",
		Status:   types.StatusBusy,
		Role:     types.RoleEditor,
		IsOnline: true,
		JoinedAt: time.Now().AddDate(0, -2, 0),
	}
}

func createNewUser() types.User {
	return types.User{
		ID:       "user-new",
		Username: "new.joiner",
		Email:    "new@example.com",
		Avatar:   "https://images.unsplash.com/photo-1494790108755-2616b612b786?w=150",
		Status:   types.StatusOnline,
		Role:     types.RoleGuest,
		IsOnline: true,
		JoinedAt: time.Now().Add(-24 * time.Hour),
	}
}

func createPowerUser() types.User {
	return types.User{
		ID:       "user-power",
		Username: "super.admin",
		Email:    "admin@example.com",
		Avatar:   "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150",
		Status:   types.StatusOnline,
		Role:     types.RoleAdmin,
		IsOnline: true,
		JoinedAt: time.Now().AddDate(-2, 0, 0),
	}
}

func createMinimalStats() types.DashboardStats {
	return types.DashboardStats{
		TotalProjects:         1,
		ActiveProjects:        1,
		TotalUsers:           2,
		OnlineUsers:          1,
		TotalDocuments:       1,
		TotalComments:        0,
		RecentActivity:       []types.ActivityEvent{},
		ProjectsByStatus:     map[types.ProjectStatus]int{types.ProjectStatusDraft: 1},
		ProjectsByPriority:   map[types.ProjectPriority]int{types.PriorityLow: 1},
		UserActivityChart:    []types.ChartDataPoint{},
		DocumentCreationChart: []types.ChartDataPoint{},
		TopProjects:          []types.Project{createMinimalProject()},
		TopUsers:             []types.User{createNewUser()},
		SystemHealth:         createHealthySystem(),
	}
}

func createMaximalStats() types.DashboardStats {
	return types.DashboardStats{
		TotalProjects:    50,
		ActiveProjects:   32,
		TotalUsers:      150,
		OnlineUsers:     45,
		TotalDocuments:  1250,
		TotalComments:   3456,
		RecentActivity:  createMaximalActivity(),
		ProjectsByStatus: map[types.ProjectStatus]int{
			types.ProjectStatusActive:    32,
			types.ProjectStatusDraft:     8,
			types.ProjectStatusCompleted: 7,
			types.ProjectStatusOnHold:    3,
		},
		ProjectsByPriority: map[types.ProjectPriority]int{
			types.PriorityCritical: 5,
			types.PriorityHigh:     15,
			types.PriorityMedium:   20,
			types.PriorityLow:      10,
		},
		TopProjects: createVariedProjects(),
		TopUsers:    createMixedOnlineUsers(),
		SystemHealth: createWarningSystem(),
	}
}

func createSkewedStats() types.DashboardStats {
	return types.DashboardStats{
		TotalProjects:  25,
		ActiveProjects: 5,
		ProjectsByStatus: map[types.ProjectStatus]int{
			types.ProjectStatusActive:   5,
			types.ProjectStatusDraft:    15,
			types.ProjectStatusArchived: 5,
		},
		ProjectsByPriority: map[types.ProjectPriority]int{
			types.PriorityCritical: 15,
			types.PriorityHigh:     8,
			types.PriorityMedium:   2,
		},
		SystemHealth: createCriticalSystem(),
	}
}

func createHealthySystem() types.SystemHealth {
	return types.SystemHealth{
		CPUUsage:          15.2,
		MemoryUsage:       45.8,
		DiskUsage:         32.1,
		ActiveConnections: 12,
		Uptime:            86400 * 30,
		LastChecked:       time.Now(),
		Status:            types.HealthStatusHealthy,
	}
}

func createWarningSystem() types.SystemHealth {
	return types.SystemHealth{
		CPUUsage:          75.5,
		MemoryUsage:       82.3,
		DiskUsage:         68.9,
		ActiveConnections: 89,
		Uptime:            86400 * 7,
		LastChecked:       time.Now(),
		Status:            types.HealthStatusWarning,
	}
}

func createCriticalSystem() types.SystemHealth {
	return types.SystemHealth{
		CPUUsage:          95.8,
		MemoryUsage:       98.2,
		DiskUsage:         94.7,
		ActiveConnections: 150,
		Uptime:            86400,
		LastChecked:       time.Now(),
		Status:            types.HealthStatusCritical,
	}
}

func createManyNotifications() []types.Notification {
	notifications := []types.Notification{}
	for i := 0; i < 15; i++ {
		notifications = append(notifications, types.Notification{
			ID:       fmt.Sprintf("notif-%d", i),
			UserID:   "user-1",
			Type:     []types.NotificationType{types.NotificationTypeInfo, types.NotificationTypeWarning, types.NotificationTypeSuccess}[i%3],
			Title:    fmt.Sprintf("Notification %d", i+1),
			Message:  fmt.Sprintf("This is notification message %d with some content", i+1),
			Priority: []types.NotificationPriority{types.NotificationPriorityLow, types.NotificationPriorityMedium, types.NotificationPriorityHigh}[i%3],
			IsRead:   i%4 == 0,
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		})
	}
	return notifications
}

func createMaximalActivity() []types.ActivityEvent {
	activities := []types.ActivityEvent{}
	for i := 0; i < 20; i++ {
		activities = append(activities, types.ActivityEvent{
			ID:   fmt.Sprintf("activity-%d", i),
			Type: []types.ActivityType{
				types.ActivityDocumentCreated,
				types.ActivityDocumentUpdated,
				types.ActivityCommentAdded,
				types.ActivityProjectUpdated,
				types.ActivityMemberAdded,
			}[i%5],
			Actor: []types.User{createNewUser(), createPowerUser(), createBusyUser()}[i%3],
			Target: types.ActivityTarget{
				Type: "document",
				ID:   fmt.Sprintf("doc-%d", i),
				Name: fmt.Sprintf("Document %d", i),
			},
			Description: fmt.Sprintf("performed action %d", i),
			CreatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
			IsImportant: i%5 == 0,
		})
	}
	return activities
}

func createVariedProjects() []types.Project {
	return []types.Project{
		createUrgentProject(),
		createLargeProject(),
		createMinimalProject(),
		createCompletedProject(),
	}
}

func createUrgentProject() types.Project {
	deadline := time.Now().AddDate(0, 0, 3)
	return types.Project{
		ID:          "proj-urgent",
		Name:        "Critical Security Update",
		Description: "Urgent security patches and vulnerability fixes",
		Owner:       createPowerUser(),
		Status:      types.ProjectStatusActive,
		Priority:    types.PriorityCritical,
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * time.Hour),
		DeadlineAt:  &deadline,
		Tags:        []string{"security", "urgent", "patch"},
		Members:     createManyMembers(),
		Analytics: types.ProjectAnalytics{
			TotalViews:          2890,
			TotalEdits:          456,
			TotalComments:       123,
			TotalCollaborators:  8,
			LastActivityAt:      time.Now().Add(-30 * time.Minute),
			PeakConcurrentUsers: 6,
			DocumentCount:       25,
			CompletionRate:      85.5,
		},
	}
}

func createLargeProject() types.Project {
	return types.Project{
		ID:          "proj-large",
		Name:        "Enterprise Platform Redesign",
		Description: "Complete overhaul of the enterprise platform with new architecture, microservices, and modern UI",
		Owner:       createPowerUser(),
		Status:      types.ProjectStatusActive,
		Priority:    types.PriorityHigh,
		CreatedAt:   time.Now().AddDate(0, -6, 0),
		UpdatedAt:   time.Now().Add(-2 * time.Hour),
		Tags:        []string{"enterprise", "architecture", "microservices", "ui", "backend", "frontend"},
		Members:     createManyMembers(),
		Documents:   createManyDocuments(),
		Analytics: types.ProjectAnalytics{
			TotalViews:          15670,
			TotalEdits:          2340,
			TotalComments:       567,
			TotalCollaborators:  15,
			LastActivityAt:      time.Now().Add(-15 * time.Minute),
			PeakConcurrentUsers: 12,
			DocumentCount:       89,
			CompletionRate:      62.3,
		},
	}
}

func createMinimalProject() types.Project {
	return types.Project{
		ID:          "proj-minimal",
		Name:        "Quick Fix",
		Description: "Small bug fix",
		Owner:       createNewUser(),
		Status:      types.ProjectStatusDraft,
		Priority:    types.PriorityLow,
		CreatedAt:   time.Now().Add(-6 * time.Hour),
		UpdatedAt:   time.Now().Add(-4 * time.Hour),
		Tags:        []string{"bugfix"},
		Members:     createMinimalMembers(),
		Analytics: types.ProjectAnalytics{
			TotalViews:          23,
			TotalEdits:          5,
			TotalComments:       1,
			TotalCollaborators:  1,
			LastActivityAt:      time.Now().Add(-4 * time.Hour),
			PeakConcurrentUsers: 1,
			DocumentCount:       1,
			CompletionRate:      10.0,
		},
	}
}

func createCompletedProject() types.Project {
	return types.Project{
		ID:          "proj-completed",
		Name:        "Q3 Marketing Campaign",
		Description: "Successful marketing campaign for Q3 product launch",
		Owner:       createBusyUser(),
		Status:      types.ProjectStatusCompleted,
		Priority:    types.PriorityMedium,
		CreatedAt:   time.Now().AddDate(0, -4, 0),
		UpdatedAt:   time.Now().AddDate(0, -1, 0),
		Tags:        []string{"marketing", "campaign", "q3", "completed"},
		Members:     createManyMembers(),
		Analytics: types.ProjectAnalytics{
			TotalViews:          5670,
			TotalEdits:          890,
			TotalComments:       234,
			TotalCollaborators:  12,
			LastActivityAt:      time.Now().AddDate(0, -1, 0),
			PeakConcurrentUsers: 8,
			DocumentCount:       45,
			CompletionRate:      100.0,
		},
	}
}

func createManyDocuments() []types.Document {
	documents := []types.Document{}
	docTypes := []types.DocumentType{
		types.DocumentTypeMarkdown,
		types.DocumentTypeCode,
		types.DocumentTypeJSON,
		types.DocumentTypeYAML,
	}
	statuses := []types.DocumentStatus{
		types.DocumentStatusDraft,
		types.DocumentStatusReview,
		types.DocumentStatusPublished,
	}
	
	for i := 0; i < 12; i++ {
		documents = append(documents, types.Document{
			ID:        fmt.Sprintf("doc-many-%d", i),
			ProjectID: "proj-large",
			Title:     fmt.Sprintf("Document %d - Comprehensive Guide", i+1),
			Content:   fmt.Sprintf("This is the content for document %d with detailed information...", i+1),
			Type:      docTypes[i%len(docTypes)],
			Status:    statuses[i%len(statuses)],
			Author:    []types.User{createNewUser(), createPowerUser(), createBusyUser()}[i%3],
			Editors:   []types.User{createNewUser(), createPowerUser()},
			CreatedAt: time.Now().Add(-time.Duration(i*12) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(i*2) * time.Hour),
			WordCount: 1500 + i*200,
			ReadTime:  8 + i,
		})
	}
	return documents
}

func createManyMembers() []types.ProjectMember {
	users := []types.User{
		createNewUser(),
		createPowerUser(),
		createBusyUser(),
		createOfflineUser(),
	}
	roles := []types.ProjectRole{
		types.ProjectRoleOwner,
		types.ProjectRoleAdmin,
		types.ProjectRoleCollaborator,
		types.ProjectRoleReviewer,
		types.ProjectRoleObserver,
	}
	
	members := []types.ProjectMember{}
	for i, user := range users {
		members = append(members, types.ProjectMember{
			User:        user,
			Role:        roles[i%len(roles)],
			JoinedAt:    time.Now().Add(-time.Duration(i*24) * time.Hour),
			Permissions: []string{"read", "write", "comment"},
			IsActive:    i%4 != 3, // Some inactive members
		})
	}
	return members
}

func createMinimalMembers() []types.ProjectMember {
	return []types.ProjectMember{
		{
			User:        createNewUser(),
			Role:        types.ProjectRoleOwner,
			JoinedAt:    time.Now().Add(-6 * time.Hour),
			Permissions: []string{"all"},
			IsActive:    true,
		},
	}
}

func createMixedOnlineUsers() []types.User {
	return []types.User{
		createPowerUser(),
		createBusyUser(),
		createNewUser(),
		createOfflineUser(),
	}
}

func createAllOfflineUsers() []types.User {
	users := []types.User{}
	for i := 0; i < 5; i++ {
		user := createOfflineUser()
		user.ID = fmt.Sprintf("offline-%d", i)
		user.Username = fmt.Sprintf("offline.user.%d", i)
		users = append(users, user)
	}
	return users
}