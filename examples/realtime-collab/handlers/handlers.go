package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/conneroisu/twerge/examples/realtime-collab/data"
	"github.com/conneroisu/twerge/examples/realtime-collab/types"
	"github.com/conneroisu/twerge/examples/realtime-collab/views"
	"github.com/conneroisu/twerge/examples/realtime-collab/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	dataStore *data.Store
	wsHub     *websocket.Hub
}

func NewHandlers(store *data.Store, hub *websocket.Hub) *Handlers {
	return &Handlers{
		dataStore: store,
		wsHub:     hub,
	}
}

func (h *Handlers) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(h.corsMiddleware)

	// Static files
	fileServer := http.FileServer(http.Dir("./_static/"))
	r.Handle("/_static/*", http.StripPrefix("/_static/", fileServer))

	// WebSocket endpoint
	r.Get("/ws", h.handleWebSocket)

	// Dashboard routes
	r.Get("/", h.handleDashboard)
	r.Get("/dashboard", h.handleDashboard)

	// Project routes
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", h.handleProjectsList)
		r.Post("/", h.handleCreateProject)
		r.Get("/{projectID}", h.handleProjectView)
		r.Put("/{projectID}", h.handleUpdateProject)
		r.Delete("/{projectID}", h.handleDeleteProject)
		r.Get("/{projectID}/members", h.handleProjectMembers)
		r.Post("/{projectID}/members", h.handleAddProjectMember)
		r.Delete("/{projectID}/members/{userID}", h.handleRemoveProjectMember)
	})

	// Document routes
	r.Route("/documents", func(r chi.Router) {
		r.Get("/", h.handleDocumentsList)
		r.Post("/", h.handleCreateDocument)
		r.Get("/{documentID}", h.handleDocumentView)
		r.Put("/{documentID}", h.handleUpdateDocument)
		r.Delete("/{documentID}", h.handleDeleteDocument)
		r.Get("/{documentID}/collaborators", h.handleDocumentCollaborators)
		r.Post("/{documentID}/comments", h.handleAddComment)
		r.Get("/{documentID}/comments", h.handleGetComments)
	})

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.handleUsersList)
		r.Get("/{userID}", h.handleUserProfile)
		r.Put("/{userID}", h.handleUpdateUser)
		r.Get("/{userID}/projects", h.handleUserProjects)
		r.Get("/{userID}/activity", h.handleUserActivity)
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/search", h.handleSearch)
		r.Get("/stats", h.handleStats)
		r.Get("/notifications", h.handleNotifications)
		r.Post("/notifications/{notificationID}/read", h.handleMarkNotificationRead)
		r.Get("/online-users", h.handleOnlineUsers)
		r.Get("/system-health", h.handleSystemHealth)
		r.Post("/upload", h.handleFileUpload)
		r.Get("/export/{projectID}", h.handleExportProject)
	})

	// Real-time API routes
	r.Route("/realtime", func(r chi.Router) {
		r.Get("/presence", h.handlePresence)
		r.Post("/cursor", h.handleCursorUpdate)
		r.Post("/typing", h.handleTypingIndicator)
		r.Get("/room/{roomID}/users", h.handleRoomUsers)
		r.Post("/broadcast", h.handleBroadcast)
	})

	return r
}

// Dashboard handlers
func (h *Handlers) handleDashboard(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	stats := h.dataStore.GetDashboardStats()
	notifications := h.dataStore.GetUserNotifications(user.ID)

	component := views.DashboardLayout(user, stats, notifications)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Project handlers
func (h *Handlers) handleProjectsList(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	projects := h.dataStore.GetUserProjects(user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *Handlers) handleProjectView(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	user := h.getCurrentUser(r)

	project, err := h.dataStore.GetProject(projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	documents := h.dataStore.GetProjectDocuments(projectID)
	members := h.dataStore.GetProjectMembers(projectID)

	component := views.ProjectView(project, documents, members, user)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	
	var projectData struct {
		Name        string                `json:"name"`
		Description string                `json:"description"`
		Priority    types.ProjectPriority `json:"priority"`
		Tags        []string              `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&projectData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	project := types.Project{
		ID:          h.dataStore.GenerateID(),
		Name:        projectData.Name,
		Description: projectData.Description,
		Owner:       user,
		Status:      types.ProjectStatusDraft,
		Priority:    projectData.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        projectData.Tags,
		Members:     []types.ProjectMember{
			{
				User:        user,
				Role:        types.ProjectRoleOwner,
				JoinedAt:    time.Now(),
				Permissions: []string{"all"},
				IsActive:    true,
			},
		},
		Settings: types.ProjectSettings{
			IsPublic:              false,
			AllowGuestEditing:     false,
			RequireApproval:       true,
			EnableVersionControl:  true,
			EnableComments:        true,
			EnableNotifications:   true,
			AutoSaveInterval:      30,
			MaxConcurrentEditors:  10,
			AllowedFileTypes:      []string{"md", "txt", "json", "yaml"},
			MaxFileSize:           10485760, // 10MB
		},
	}

	if err := h.dataStore.CreateProject(project); err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	// Broadcast project creation event
	activity := types.ActivityEvent{
		ID:          h.dataStore.GenerateID(),
		Type:        types.ActivityProjectCreated,
		Actor:       user,
		Target:      types.ActivityTarget{Type: "project", ID: project.ID, Name: project.Name},
		Description: "created a new project",
		CreatedAt:   time.Now(),
		IsImportant: true,
	}

	h.broadcastActivity(activity)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *Handlers) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	user := h.getCurrentUser(r)

	project, err := h.dataStore.GetProject(projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Check permissions
	if !h.dataStore.HasProjectPermission(user.ID, projectID, "edit") {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.dataStore.UpdateProject(projectID, updateData); err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	// Broadcast update event
	activity := types.ActivityEvent{
		ID:          h.dataStore.GenerateID(),
		Type:        types.ActivityProjectUpdated,
		Actor:       user,
		Target:      types.ActivityTarget{Type: "project", ID: project.ID, Name: project.Name},
		Description: "updated the project",
		CreatedAt:   time.Now(),
		IsImportant: false,
	}

	h.broadcastActivity(activity)

	w.WriteHeader(http.StatusOK)
}

// Document handlers
func (h *Handlers) handleDocumentView(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "documentID")
	user := h.getCurrentUser(r)

	document, err := h.dataStore.GetDocument(documentID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Check permissions
	if !h.dataStore.HasDocumentPermission(user.ID, documentID, "read") {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handlers) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	
	var docData struct {
		ProjectID string               `json:"projectId"`
		Title     string               `json:"title"`
		Content   string               `json:"content"`
		Type      types.DocumentType   `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&docData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check project permissions
	if !h.dataStore.HasProjectPermission(user.ID, docData.ProjectID, "create_document") {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	document := types.Document{
		ID:          h.dataStore.GenerateID(),
		ProjectID:   docData.ProjectID,
		Title:       docData.Title,
		Content:     docData.Content,
		Type:        docData.Type,
		Status:      types.DocumentStatusDraft,
		Author:      user,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		WordCount:   len(docData.Content),
		ReadTime:    calculateReadTime(docData.Content),
		Settings: types.DocumentSettings{
			IsPublic:        false,
			AllowComments:   true,
			RequireApproval: true,
			LineNumbers:     true,
			WordWrap:        true,
			Theme:           "default",
		},
	}

	if err := h.dataStore.CreateDocument(document); err != nil {
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	// Broadcast document creation event
	activity := types.ActivityEvent{
		ID:          h.dataStore.GenerateID(),
		Type:        types.ActivityDocumentCreated,
		Actor:       user,
		Target:      types.ActivityTarget{Type: "document", ID: document.ID, Name: document.Title},
		Description: "created a new document",
		CreatedAt:   time.Now(),
		IsImportant: false,
	}

	h.broadcastActivity(activity)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

// WebSocket handler
func (h *Handlers) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	h.wsHub.HandleWebSocket(w, r, user)
}

// API handlers
func (h *Handlers) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]types.SearchResult{})
		return
	}

	user := h.getCurrentUser(r)
	results := h.dataStore.Search(query, user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handlers) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := h.dataStore.GetDashboardStats()
	wsStats := h.wsHub.GetStats()

	// Merge stats
	combined := map[string]interface{}{
		"dashboard":    stats,
		"websocket":    wsStats,
		"last_updated": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combined)
}

func (h *Handlers) handleNotifications(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	notifications := h.dataStore.GetUserNotifications(user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *Handlers) handleOnlineUsers(w http.ResponseWriter, r *http.Request) {
	users := h.wsHub.GetOnlineUsers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handlers) handleSystemHealth(w http.ResponseWriter, r *http.Request) {
	health := h.dataStore.GetSystemHealth()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Real-time handlers
func (h *Handlers) handlePresence(w http.ResponseWriter, r *http.Request) {
	onlineUsers := h.wsHub.GetOnlineUsers()
	
	presence := map[string]interface{}{
		"online_users": onlineUsers,
		"total_count":  len(onlineUsers),
		"timestamp":    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(presence)
}

func (h *Handlers) handleCursorUpdate(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	
	var cursorData struct {
		DocumentID string `json:"documentId"`
		Line       int    `json:"line"`
		Column     int    `json:"column"`
		Selection  *struct {
			StartLine   int `json:"startLine"`
			StartColumn int `json:"startColumn"`
			EndLine     int `json:"endLine"`
			EndColumn   int `json:"endColumn"`
		} `json:"selection,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&cursorData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cursor := &types.Cursor{
		UserID:     user.ID,
		DocumentID: cursorData.DocumentID,
		Line:       cursorData.Line,
		Column:     cursorData.Column,
		Color:      generateUserColor(user.ID),
		UpdatedAt:  time.Now(),
	}

	if cursorData.Selection != nil {
		cursor.Selection = &types.Selection{
			StartLine:   cursorData.Selection.StartLine,
			StartColumn: cursorData.Selection.StartColumn,
			EndLine:     cursorData.Selection.EndLine,
			EndColumn:   cursorData.Selection.EndColumn,
		}
	}

	// Broadcast cursor update
	msg := types.WebSocketMessage{
		Type:      types.MessageTypeCursorMove,
		RoomID:    fmt.Sprintf("document:%s", cursorData.DocumentID),
		UserID:    user.ID,
		Data:      cursor,
		Timestamp: time.Now(),
	}

	h.wsHub.BroadcastToRoom(msg.RoomID, msg)

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleRoomUsers(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	users := h.wsHub.GetRoomUsers(roomID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Utility methods
func (h *Handlers) getCurrentUser(r *http.Request) types.User {
	// In a real application, this would extract the user from JWT token or session
	return h.dataStore.GetDemoUser()
}

func (h *Handlers) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) broadcastActivity(activity types.ActivityEvent) {
	msg := types.WebSocketMessage{
		Type:      types.MessageTypeSystemAlert,
		Data:      activity,
		Timestamp: time.Now(),
	}

	data, _ := json.Marshal(msg)
	h.wsHub.Broadcast <- data
}

func calculateReadTime(content string) int {
	// Approximate reading time: 200 words per minute
	wordCount := len(content) / 5 // Rough estimate
	return (wordCount / 200) + 1
}

func generateUserColor(userID string) string {
	colors := []string{
		"#3B82F6", "#EF4444", "#10B981", "#F59E0B",
		"#8B5CF6", "#EC4899", "#06B6D4", "#84CC16",
		"#F97316", "#6366F1", "#14B8A6", "#F43F5E",
	}
	
	// Simple hash to pick a consistent color for each user
	hash := 0
	for _, char := range userID {
		hash = hash*31 + int(char)
	}
	
	return colors[hash%len(colors)]
}

// Additional handlers would include file upload, export, member management, etc.
// These are simplified for brevity but would include proper validation,
// error handling, and security measures in a production system.

func (h *Handlers) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) handleProjectMembers(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	members := h.dataStore.GetProjectMembers(projectID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func (h *Handlers) handleAddProjectMember(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) handleRemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) handleDocumentsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]types.Document{})
}

func (h *Handlers) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) handleUpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleDocumentCollaborators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]types.User{})
}

func (h *Handlers) handleAddComment(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) handleGetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]types.Comment{})
}

func (h *Handlers) handleUsersList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]types.User{})
}

func (h *Handlers) handleUserProfile(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handlers) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleUserProjects(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	projects := h.dataStore.GetUserProjects(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *Handlers) handleUserActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]types.ActivityEvent{})
}

func (h *Handlers) handleMarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) handleExportProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"export": "data"})
}

func (h *Handlers) handleTypingIndicator(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleBroadcast(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}