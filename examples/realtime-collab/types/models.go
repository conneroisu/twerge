package types

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Avatar        string    `json:"avatar"`
	Status        UserStatus `json:"status"`
	LastSeen      time.Time `json:"lastSeen"`
	Role          UserRole  `json:"role"`
	Permissions   []string  `json:"permissions"`
	JoinedAt      time.Time `json:"joinedAt"`
	IsOnline      bool      `json:"isOnline"`
	CurrentRoom   string    `json:"currentRoom,omitempty"`
	ActiveCursor  *Cursor   `json:"activeCursor,omitempty"`
}

type UserStatus string

const (
	StatusOnline    UserStatus = "online"
	StatusAway      UserStatus = "away"
	StatusBusy      UserStatus = "busy"
	StatusOffline   UserStatus = "offline"
	StatusInvisible UserStatus = "invisible"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
	RoleEditor    UserRole = "editor"
	RoleViewer    UserRole = "viewer"
	RoleGuest     UserRole = "guest"
)

type Project struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Owner       User               `json:"owner"`
	Members     []ProjectMember    `json:"members"`
	Documents   []Document         `json:"documents"`
	Status      ProjectStatus      `json:"status"`
	Priority    ProjectPriority    `json:"priority"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	DeadlineAt  *time.Time         `json:"deadlineAt,omitempty"`
	Tags        []string           `json:"tags"`
	Settings    ProjectSettings    `json:"settings"`
	Analytics   ProjectAnalytics   `json:"analytics"`
	Activity    []ActivityEvent    `json:"activity"`
}

type ProjectMember struct {
	User        User               `json:"user"`
	Role        ProjectRole        `json:"role"`
	JoinedAt    time.Time          `json:"joinedAt"`
	Permissions []string           `json:"permissions"`
	IsActive    bool               `json:"isActive"`
}

type ProjectRole string

const (
	ProjectRoleOwner        ProjectRole = "owner"
	ProjectRoleAdmin        ProjectRole = "admin"
	ProjectRoleCollaborator ProjectRole = "collaborator"
	ProjectRoleReviewer     ProjectRole = "reviewer"
	ProjectRoleObserver     ProjectRole = "observer"
)

type ProjectStatus string

const (
	ProjectStatusDraft      ProjectStatus = "draft"
	ProjectStatusActive     ProjectStatus = "active"
	ProjectStatusOnHold     ProjectStatus = "on_hold"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusArchived   ProjectStatus = "archived"
	ProjectStatusCancelled  ProjectStatus = "cancelled"
)

type ProjectPriority string

const (
	PriorityLow      ProjectPriority = "low"
	PriorityMedium   ProjectPriority = "medium"
	PriorityHigh     ProjectPriority = "high"
	PriorityCritical ProjectPriority = "critical"
)

type ProjectSettings struct {
	IsPublic              bool     `json:"isPublic"`
	AllowGuestEditing     bool     `json:"allowGuestEditing"`
	RequireApproval       bool     `json:"requireApproval"`
	EnableVersionControl  bool     `json:"enableVersionControl"`
	EnableComments        bool     `json:"enableComments"`
	EnableNotifications   bool     `json:"enableNotifications"`
	AutoSaveInterval      int      `json:"autoSaveInterval"`
	MaxConcurrentEditors  int      `json:"maxConcurrentEditors"`
	AllowedFileTypes      []string `json:"allowedFileTypes"`
	MaxFileSize           int64    `json:"maxFileSize"`
}

type ProjectAnalytics struct {
	TotalViews           int       `json:"totalViews"`
	TotalEdits           int       `json:"totalEdits"`
	TotalComments        int       `json:"totalComments"`
	TotalCollaborators   int       `json:"totalCollaborators"`
	LastActivityAt       time.Time `json:"lastActivityAt"`
	AverageSessionTime   int       `json:"averageSessionTime"`
	PeakConcurrentUsers  int       `json:"peakConcurrentUsers"`
	DocumentCount        int       `json:"documentCount"`
	CompletionRate       float64   `json:"completionRate"`
}

type Document struct {
	ID              string           `json:"id"`
	ProjectID       string           `json:"projectId"`
	Title           string           `json:"title"`
	Content         string           `json:"content"`
	Type            DocumentType     `json:"type"`
	Status          DocumentStatus   `json:"status"`
	Author          User             `json:"author"`
	Editors         []User           `json:"editors"`
	Comments        []Comment        `json:"comments"`
	Versions        []DocumentVersion `json:"versions"`
	Tags            []string         `json:"tags"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	LastEditedBy    *User            `json:"lastEditedBy,omitempty"`
	IsLocked        bool             `json:"isLocked"`
	LockHolder      *User            `json:"lockHolder,omitempty"`
	WordCount       int              `json:"wordCount"`
	ReadTime        int              `json:"readTime"`
	Settings        DocumentSettings `json:"settings"`
}

type DocumentType string

const (
	DocumentTypeText     DocumentType = "text"
	DocumentTypeMarkdown DocumentType = "markdown"
	DocumentTypeCode     DocumentType = "code"
	DocumentTypeJSON     DocumentType = "json"
	DocumentTypeYAML     DocumentType = "yaml"
	DocumentTypeXML      DocumentType = "xml"
	DocumentTypeCSV      DocumentType = "csv"
)

type DocumentStatus string

const (
	DocumentStatusDraft     DocumentStatus = "draft"
	DocumentStatusReview    DocumentStatus = "review"
	DocumentStatusPublished DocumentStatus = "published"
	DocumentStatusArchived  DocumentStatus = "archived"
)

type DocumentSettings struct {
	IsPublic         bool   `json:"isPublic"`
	AllowComments    bool   `json:"allowComments"`
	RequireApproval  bool   `json:"requireApproval"`
	SyntaxHighlight  string `json:"syntaxHighlight,omitempty"`
	LineNumbers      bool   `json:"lineNumbers"`
	WordWrap         bool   `json:"wordWrap"`
	Theme            string `json:"theme"`
}

type DocumentVersion struct {
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	Content   string    `json:"content"`
	Author    User      `json:"author"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	Size      int64     `json:"size"`
	Changes   []Change  `json:"changes"`
}

type Change struct {
	Type     ChangeType `json:"type"`
	Line     int        `json:"line"`
	Column   int        `json:"column"`
	Length   int        `json:"length"`
	Content  string     `json:"content"`
	Author   User       `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

type ChangeType string

const (
	ChangeTypeInsert ChangeType = "insert"
	ChangeTypeDelete ChangeType = "delete"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeMove   ChangeType = "move"
)

type Comment struct {
	ID         string          `json:"id"`
	DocumentID string          `json:"documentId"`
	Author     User            `json:"author"`
	Content    string          `json:"content"`
	Position   CommentPosition `json:"position"`
	Replies    []Comment       `json:"replies"`
	IsResolved bool            `json:"isResolved"`
	ResolvedBy *User           `json:"resolvedBy,omitempty"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	Reactions  []Reaction      `json:"reactions"`
}

type CommentPosition struct {
	Line      int `json:"line"`
	Column    int `json:"column"`
	EndLine   int `json:"endLine"`
	EndColumn int `json:"endColumn"`
}

type Reaction struct {
	Type   ReactionType `json:"type"`
	User   User         `json:"user"`
	AddedAt time.Time   `json:"addedAt"`
}

type ReactionType string

const (
	ReactionLike     ReactionType = "like"
	ReactionDislike  ReactionType = "dislike"
	ReactionLove     ReactionType = "love"
	ReactionLaugh    ReactionType = "laugh"
	ReactionSurprise ReactionType = "surprise"
	ReactionAngry    ReactionType = "angry"
)

type ActivityEvent struct {
	ID          string           `json:"id"`
	Type        ActivityType     `json:"type"`
	Actor       User             `json:"actor"`
	Target      ActivityTarget   `json:"target"`
	Description string           `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time        `json:"createdAt"`
	IsImportant bool             `json:"isImportant"`
}

type ActivityType string

const (
	ActivityProjectCreated   ActivityType = "project_created"
	ActivityProjectUpdated   ActivityType = "project_updated"
	ActivityProjectDeleted   ActivityType = "project_deleted"
	ActivityMemberAdded      ActivityType = "member_added"
	ActivityMemberRemoved    ActivityType = "member_removed"
	ActivityMemberRoleChanged ActivityType = "member_role_changed"
	ActivityDocumentCreated  ActivityType = "document_created"
	ActivityDocumentUpdated  ActivityType = "document_updated"
	ActivityDocumentDeleted  ActivityType = "document_deleted"
	ActivityCommentAdded     ActivityType = "comment_added"
	ActivityCommentResolved  ActivityType = "comment_resolved"
	ActivityUserJoined       ActivityType = "user_joined"
	ActivityUserLeft         ActivityType = "user_left"
)

type ActivityTarget struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Cursor struct {
	UserID    string    `json:"userId"`
	DocumentID string   `json:"documentId"`
	Line      int       `json:"line"`
	Column    int       `json:"column"`
	Selection *Selection `json:"selection,omitempty"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Selection struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndLine     int `json:"endLine"`
	EndColumn   int `json:"endColumn"`
}

type WebSocketMessage struct {
	Type      MessageType     `json:"type"`
	RoomID    string          `json:"roomId"`
	UserID    string          `json:"userId"`
	Data      interface{}     `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	MessageID string          `json:"messageId"`
}

type MessageType string

const (
	MessageTypeJoinRoom        MessageType = "join_room"
	MessageTypeLeaveRoom       MessageType = "leave_room"
	MessageTypeUserJoined      MessageType = "user_joined"
	MessageTypeUserLeft        MessageType = "user_left"
	MessageTypeUserTyping      MessageType = "user_typing"
	MessageTypeCursorMove      MessageType = "cursor_move"
	MessageTypeDocumentChange  MessageType = "document_change"
	MessageTypeCommentAdded    MessageType = "comment_added"
	MessageTypeCommentResolved MessageType = "comment_resolved"
	MessageTypePresenceUpdate  MessageType = "presence_update"
	MessageTypeSystemAlert     MessageType = "system_alert"
	MessageTypeError           MessageType = "error"
)

type Room struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"projectId"`
	DocumentID  string    `json:"documentId"`
	Users       []User    `json:"users"`
	MaxUsers    int       `json:"maxUsers"`
	CreatedAt   time.Time `json:"createdAt"`
	LastActivity time.Time `json:"lastActivity"`
	IsActive    bool      `json:"isActive"`
}

type Notification struct {
	ID          string            `json:"id"`
	UserID      string            `json:"userId"`
	Type        NotificationType  `json:"type"`
	Title       string            `json:"title"`
	Message     string            `json:"message"`
	Data        map[string]interface{} `json:"data"`
	IsRead      bool              `json:"isRead"`
	IsDismissed bool              `json:"isDismissed"`
	Priority    NotificationPriority `json:"priority"`
	CreatedAt   time.Time         `json:"createdAt"`
	ReadAt      *time.Time        `json:"readAt,omitempty"`
	ExpiresAt   *time.Time        `json:"expiresAt,omitempty"`
}

type NotificationType string

const (
	NotificationTypeInfo     NotificationType = "info"
	NotificationTypeSuccess  NotificationType = "success"
	NotificationTypeWarning  NotificationType = "warning"
	NotificationTypeError    NotificationType = "error"
	NotificationTypeSystem   NotificationType = "system"
)

type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityMedium   NotificationPriority = "medium"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"
)

type SearchResult struct {
	Type        SearchResultType `json:"type"`
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	URL         string           `json:"url"`
	Score       float64          `json:"score"`
	Highlights  []string         `json:"highlights"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time        `json:"createdAt"`
}

type SearchResultType string

const (
	SearchResultTypeProject  SearchResultType = "project"
	SearchResultTypeDocument SearchResultType = "document"
	SearchResultTypeUser     SearchResultType = "user"
	SearchResultTypeComment  SearchResultType = "comment"
)

type FilterOptions struct {
	ProjectIDs   []string             `json:"projectIds,omitempty"`
	UserIDs      []string             `json:"userIds,omitempty"`
	Status       []string             `json:"status,omitempty"`
	Priority     []ProjectPriority    `json:"priority,omitempty"`
	Tags         []string             `json:"tags,omitempty"`
	DateRange    *DateRange           `json:"dateRange,omitempty"`
	SearchQuery  string               `json:"searchQuery,omitempty"`
	SortBy       string               `json:"sortBy,omitempty"`
	SortOrder    SortOrder            `json:"sortOrder,omitempty"`
	Limit        int                  `json:"limit,omitempty"`
	Offset       int                  `json:"offset,omitempty"`
}

type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type DashboardStats struct {
	TotalProjects        int                    `json:"totalProjects"`
	ActiveProjects       int                    `json:"activeProjects"`
	TotalUsers           int                    `json:"totalUsers"`
	OnlineUsers          int                    `json:"onlineUsers"`
	TotalDocuments       int                    `json:"totalDocuments"`
	TotalComments        int                    `json:"totalComments"`
	RecentActivity       []ActivityEvent        `json:"recentActivity"`
	ProjectsByStatus     map[ProjectStatus]int  `json:"projectsByStatus"`
	ProjectsByPriority   map[ProjectPriority]int `json:"projectsByPriority"`
	UserActivityChart    []ChartDataPoint       `json:"userActivityChart"`
	DocumentCreationChart []ChartDataPoint      `json:"documentCreationChart"`
	TopProjects          []Project              `json:"topProjects"`
	TopUsers             []User                 `json:"topUsers"`
	SystemHealth         SystemHealth           `json:"systemHealth"`
}

type ChartDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Date  time.Time `json:"date"`
	Color string  `json:"color,omitempty"`
}

type SystemHealth struct {
	CPUUsage      float64   `json:"cpuUsage"`
	MemoryUsage   float64   `json:"memoryUsage"`
	DiskUsage     float64   `json:"diskUsage"`
	ActiveConnections int   `json:"activeConnections"`
	Uptime        int64     `json:"uptime"`
	LastChecked   time.Time `json:"lastChecked"`
	Status        HealthStatus `json:"status"`
}

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusWarning   HealthStatus = "warning"
	HealthStatusCritical  HealthStatus = "critical"
	HealthStatusDown      HealthStatus = "down"
)