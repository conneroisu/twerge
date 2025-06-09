# Blog & Content Management Use Case

## Overview
A modern blog and content management system demonstrating Twerge's theming capabilities, typography handling, and dynamic content layout. Showcases complex theme switching, content-aware styling, and accessibility features.

## Key Features
- Dynamic theme switching (light/dark/auto)
- Rich typography with reading modes
- Content-aware layouts (articles, galleries, videos)
- Interactive comment system
- Author profiles with social links
- Tag-based content organization

## Twerge Benefits Demonstrated

### 1. Theme System with Context-Aware Styling
```go
type Theme string

const (
    ThemeLight Theme = "light"
    ThemeDark  Theme = "dark"
    ThemeAuto  Theme = "auto"
)

type ThemeContext struct {
    Current Theme
    System  Theme // detected system preference
}

templ BlogLayout(theme ThemeContext, readingMode bool) {
    <div class={twerge.It(
        "min-h-screen transition-colors duration-300 " +
        getThemeClasses(theme) +
        twerge.If(readingMode, " font-serif", " font-sans")
    )}>
        @Header(theme)
        <main class={twerge.It(
            "container mx-auto px-4 sm:px-6 lg:px-8 " +
            "max-w-4xl " +
            twerge.If(readingMode, "max-w-3xl", "max-w-4xl")
        )}>
            { children... }
        </main>
        @Footer(theme)
    </div>
}

func getThemeClasses(theme ThemeContext) string {
    switch theme.Current {
    case ThemeLight:
        return "bg-white text-gray-900"
    case ThemeDark:
        return "bg-gray-900 text-white"
    case ThemeAuto:
        return "bg-white dark:bg-gray-900 text-gray-900 dark:text-white"
    default:
        return "bg-white dark:bg-gray-900 text-gray-900 dark:text-white"
    }
}

templ ThemeToggle(current Theme) {
    <div class={twerge.It("relative")}>
        <button class={twerge.It(
            "flex items-center space-x-2 px-3 py-2 rounded-lg " +
            "bg-gray-100 hover:bg-gray-200 dark:bg-gray-800 dark:hover:bg-gray-700 " +
            "transition-colors duration-200 text-sm font-medium"
        )}>
            @ThemeIcon(current)
            <span class={twerge.It("hidden sm:inline")}>
                { string(current) }
            </span>
            <svg class={twerge.It("w-4 h-4 ml-1")}>
                <!-- Chevron down -->
            </svg>
        </button>
        
        <!-- Dropdown menu would be positioned here -->
    </div>
}
```

### 2. Content-Aware Article Layout
```go
type Article struct {
    ID          string
    Title       string
    Excerpt     string
    Content     string
    Author      Author
    PublishedAt time.Time
    Tags        []string
    Featured    bool
    ReadTime    int
    ContentType ContentType // article, gallery, video, podcast
}

type ContentType string

const (
    ContentTypeArticle ContentType = "article"
    ContentTypeGallery ContentType = "gallery"
    ContentTypeVideo   ContentType = "video"
    ContentTypePodcast ContentType = "podcast"
)

templ ArticleCard(article Article, layout string) {
    <article class={twerge.It(
        "group bg-white dark:bg-gray-800 rounded-xl overflow-hidden " +
        "shadow-sm hover:shadow-md transition-shadow duration-300 " +
        "border border-gray-200 dark:border-gray-700 " +
        getLayoutClasses(layout) +
        twerge.If(article.Featured, " ring-2 ring-blue-500 ring-opacity-30", "")
    )}>
        @ArticleMedia(article)
        
        <div class={twerge.It(
            "p-6 " +
            twerge.If(layout == "horizontal", "flex-1", "")
        )}>
            <!-- Content Type Badge -->
            @ContentTypeBadge(article.ContentType)
            
            <!-- Article Header -->
            <header class={twerge.It("space-y-3")}>
                <h2 class={twerge.It(
                    "font-bold text-gray-900 dark:text-white " +
                    "group-hover:text-blue-600 dark:group-hover:text-blue-400 " +
                    "transition-colors duration-200 " +
                    twerge.If(layout == "featured", "text-2xl lg:text-3xl", "text-xl")
                )}>
                    <a href={templ.SafeURL("/articles/" + article.ID)}>
                        { article.Title }
                    </a>
                </h2>
                
                @ArticleMeta(article, layout == "compact")
            </header>

            <!-- Article Excerpt -->
            if layout != "compact" {
                <p class={twerge.It(
                    "text-gray-600 dark:text-gray-300 leading-relaxed " +
                    twerge.If(layout == "featured", "text-lg", "text-base") +
                    " line-clamp-3"
                )}>
                    { article.Excerpt }
                </p>
            }

            <!-- Tags -->
            if len(article.Tags) > 0 && layout != "compact" {
                <div class={twerge.It("flex flex-wrap gap-2 pt-4")}>
                    for _, tag := range article.Tags[:min(3, len(article.Tags))] {
                        @TagBadge(tag, false)
                    }
                    if len(article.Tags) > 3 {
                        <span class={twerge.It("text-sm text-gray-500 dark:text-gray-400")}>
                            +{ strconv.Itoa(len(article.Tags) - 3) } more
                        </span>
                    }
                </div>
            }
        </div>
    </article>
}

func getLayoutClasses(layout string) string {
    switch layout {
    case "horizontal":
        return "flex flex-col sm:flex-row"
    case "featured":
        return "lg:col-span-2"
    case "compact":
        return ""
    default:
        return ""
    }
}

templ ContentTypeBadge(contentType ContentType) {
    <div class={twerge.It(
        "inline-flex items-center space-x-1 text-xs font-medium px-2 py-1 rounded-full " +
        getContentTypeClasses(contentType)
    )}>
        @ContentTypeIcon(contentType)
        <span class={twerge.It("capitalize")}>{ string(contentType) }</span>
    </div>
}

func getContentTypeClasses(contentType ContentType) string {
    switch contentType {
    case ContentTypeArticle:
        return "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300"
    case ContentTypeGallery:
        return "bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300"
    case ContentTypeVideo:
        return "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300"
    case ContentTypePodcast:
        return "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300"
    default:
        return "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
    }
}
```

### 3. Rich Typography and Reading Experience
```go
templ ArticleContent(article Article, readingMode bool, fontSize string) {
    <div class={twerge.It(
        "prose prose-lg max-w-none " +
        "prose-headings:text-gray-900 dark:prose-headings:text-white " +
        "prose-p:text-gray-700 dark:prose-p:text-gray-300 " +
        "prose-a:text-blue-600 dark:prose-a:text-blue-400 " +
        "prose-strong:text-gray-900 dark:prose-strong:text-white " +
        "prose-code:text-pink-600 dark:prose-code:text-pink-400 " +
        "prose-pre:bg-gray-100 dark:prose-pre:bg-gray-800 " +
        "prose-blockquote:border-l-blue-500 dark:prose-blockquote:border-l-blue-400 " +
        getReadingModeClasses(readingMode) +
        getFontSizeClasses(fontSize)
    )}>
        <!-- Article content with enhanced typography -->
        <div class={twerge.It("space-y-6")}>
            @ProcessMarkdown(article.Content)
        </div>
    </div>
}

func getReadingModeClasses(readingMode bool) string {
    if readingMode {
        return " font-serif leading-relaxed prose-xl prose-p:leading-loose " +
               "prose-headings:font-serif prose-headings:font-normal"
    }
    return " font-sans leading-normal prose-headings:font-sans prose-headings:font-bold"
}

func getFontSizeClasses(fontSize string) string {
    switch fontSize {
    case "small":
        return " prose-sm"
    case "large":
        return " prose-xl"
    case "extra-large":
        return " prose-2xl"
    default:
        return " prose-lg"
    }
}

templ ReadingControls(readingMode bool, fontSize string, theme Theme) {
    <div class={twerge.It(
        "fixed bottom-6 right-6 bg-white dark:bg-gray-800 rounded-full shadow-lg " +
        "border border-gray-200 dark:border-gray-700 p-2 " +
        "flex items-center space-x-2 transition-all duration-300 " +
        "hover:shadow-xl"
    )}>
        <!-- Font Size Controls -->
        <div class={twerge.It("flex items-center space-x-1")}>
            <button class={twerge.It(
                "p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 " +
                "transition-colors duration-200 " +
                twerge.If(fontSize == "small", "bg-blue-100 dark:bg-blue-900", "")
            )}>
                <span class={twerge.It("text-xs font-bold")}>A</span>
            </button>
            <button class={twerge.It(
                "p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 " +
                "transition-colors duration-200 " +
                twerge.If(fontSize == "medium", "bg-blue-100 dark:bg-blue-900", "")
            )}>
                <span class={twerge.It("text-sm font-bold")}>A</span>
            </button>
            <button class={twerge.It(
                "p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 " +
                "transition-colors duration-200 " +
                twerge.If(fontSize == "large", "bg-blue-100 dark:bg-blue-900", "")
            )}>
                <span class={twerge.It("text-base font-bold")}>A</span>
            </button>
        </div>

        <div class={twerge.It("w-px h-6 bg-gray-300 dark:bg-gray-600")}></div>

        <!-- Reading Mode Toggle -->
        <button class={twerge.It(
            "p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 " +
            "transition-colors duration-200 " +
            twerge.If(readingMode, "bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400", "")
        )}>
            📖
        </button>

        <div class={twerge.It("w-px h-6 bg-gray-300 dark:bg-gray-600")}></div>

        <!-- Share Button -->
        <button class={twerge.It(
            "p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 " +
            "transition-colors duration-200"
        )}>
            🔗
        </button>
    </div>
}
```

### 4. Interactive Comment System
```go
type Comment struct {
    ID        string
    Author    User
    Content   string
    CreatedAt time.Time
    Likes     int
    Replies   []Comment
    Liked     bool
    Editing   bool
}

templ CommentThread(comments []Comment, depth int) {
    <div class={twerge.It(
        "space-y-4 " +
        twerge.If(depth > 0, "ml-6 pl-4 border-l-2 border-gray-200 dark:border-gray-700", "")
    )}>
        for _, comment := range comments {
            @CommentItem(comment, depth)
            if len(comment.Replies) > 0 {
                @CommentThread(comment.Replies, depth+1)
            }
        }
    </div>
}

templ CommentItem(comment Comment, depth int) {
    <div class={twerge.It(
        "group bg-white dark:bg-gray-800 rounded-lg p-4 " +
        "border border-gray-200 dark:border-gray-700 " +
        "hover:border-gray-300 dark:hover:border-gray-600 " +
        "transition-all duration-200 " +
        twerge.If(comment.Editing, "ring-2 ring-blue-500 ring-opacity-50", "")
    )}>
        <!-- Comment Header -->
        <div class={twerge.It("flex items-start justify-between mb-3")}>
            <div class={twerge.It("flex items-center space-x-3")}>
                <img 
                    src={comment.Author.Avatar} 
                    alt={comment.Author.Name}
                    class={twerge.It("w-8 h-8 rounded-full")}
                />
                <div>
                    <h4 class={twerge.It("font-medium text-gray-900 dark:text-white")}>
                        { comment.Author.Name }
                    </h4>
                    <time class={twerge.It("text-xs text-gray-500 dark:text-gray-400")}>
                        { comment.CreatedAt.Format("Jan 2, 2006") }
                    </time>
                </div>
            </div>

            <!-- Comment Actions -->
            <div class={twerge.It(
                "flex items-center space-x-2 opacity-0 group-hover:opacity-100 " +
                "transition-opacity duration-200"
            )}>
                <button class={twerge.It(
                    "text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 " +
                    "p-1 rounded transition-colors duration-200"
                )}>
                    ⋯
                </button>
            </div>
        </div>

        <!-- Comment Content -->
        if comment.Editing {
            @CommentEditor(comment.Content, true)
        } else {
            <div class={twerge.It(
                "prose prose-sm max-w-none " +
                "prose-p:text-gray-700 dark:prose-p:text-gray-300 " +
                "prose-a:text-blue-600 dark:prose-a:text-blue-400"
            )}>
                { comment.Content }
            </div>
        }

        <!-- Comment Footer -->
        <div class={twerge.It("flex items-center justify-between mt-3 pt-3 border-t border-gray-100 dark:border-gray-700")}>
            <div class={twerge.It("flex items-center space-x-4")}>
                <button class={twerge.It(
                    "flex items-center space-x-1 text-sm " +
                    "transition-colors duration-200 " +
                    twerge.If(comment.Liked,
                        "text-red-500 hover:text-red-600",
                        "text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
                    )
                )}>
                    <span>{ twerge.If(comment.Liked, "❤️", "🤍") }</span>
                    <span>{ strconv.Itoa(comment.Likes) }</span>
                </button>

                if depth < 3 {
                    <button class={twerge.It(
                        "text-sm text-gray-500 hover:text-gray-700 " +
                        "dark:text-gray-400 dark:hover:text-gray-300 " +
                        "transition-colors duration-200"
                    )}>
                        Reply
                    </button>
                }
            </div>
        </div>
    </div>
}
```

### 5. Dynamic Tag and Category System
```go
templ TagCloud(tags []TagWithCount, selectedTags []string) {
    <div class={twerge.It(
        "bg-white dark:bg-gray-800 rounded-lg p-6 " +
        "border border-gray-200 dark:border-gray-700"
    )}>
        <h3 class={twerge.It("text-lg font-semibold text-gray-900 dark:text-white mb-4")}>
            Popular Tags
        </h3>
        
        <div class={twerge.It("flex flex-wrap gap-2")}>
            for _, tag := range tags {
                @TagBadge(tag.Name, slices.Contains(selectedTags, tag.Name))
            }
        </div>
    </div>
}

templ TagBadge(name string, selected bool) {
    <button class={twerge.It(
        "inline-flex items-center px-3 py-1 rounded-full text-sm font-medium " +
        "transition-all duration-200 hover:scale-105 " +
        twerge.If(selected,
            "bg-blue-600 text-white shadow-lg",
            "bg-gray-100 text-gray-700 hover:bg-gray-200 " +
            "dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
        )
    )}>
        #{ name }
        if selected {
            <span class={twerge.It("ml-1 text-blue-200")}>✕</span>
        }
    </button>
}
```

## Performance Benefits

### Theme-Aware CSS Generation
```css
/* twerge:begin */
/* Theme variations consolidated */
.tw-1 { @apply min-h-screen transition-colors duration-300 bg-white text-gray-900; }
.tw-2 { @apply min-h-screen transition-colors duration-300 bg-gray-900 text-white; }
.tw-3 { @apply min-h-screen transition-colors duration-300 bg-white dark:bg-gray-900 text-gray-900 dark:text-white; }

/* Typography modes unified */
.tw-4 { @apply min-h-screen transition-colors duration-300 bg-white dark:bg-gray-900 text-gray-900 dark:text-white font-serif; }
.tw-5 { @apply min-h-screen transition-colors duration-300 bg-white dark:bg-gray-900 text-gray-900 dark:text-white font-sans; }

/* Content layouts optimized */
.tw-6 { @apply group bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 border border-gray-200 dark:border-gray-700; }
.tw-7 { @apply group bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 border border-gray-200 dark:border-gray-700 flex flex-col sm:flex-row; }
/* twerge:end */
```

### Dynamic Content Benefits
- **Theme switching**: No class conflicts between theme variants
- **Typography scaling**: Unified reading experience controls
- **Content types**: Consistent styling across different media types
- **Interactive states**: Smooth transitions without performance overhead

## Code Generation Integration

```go
//go:build ignore

package main

import (
    "time"
    "github.com/conneroisu/twerge"
    "github.com/yourproject/blog/views"
    "github.com/yourproject/blog/types"
)

func main() {
    // Sample data for all theme and layout variations
    sampleArticle := types.Article{
        Title: "Sample Article",
        Excerpt: "This is a sample article excerpt for testing.",
        Content: "Sample content...",
        Author: types.Author{Name: "John Doe"},
        PublishedAt: time.Now(),
        Tags: []string{"go", "web", "programming"},
        Featured: true,
        ReadTime: 5,
        ContentType: types.ContentTypeArticle,
    }

    themes := []types.ThemeContext{
        {Current: types.ThemeLight},
        {Current: types.ThemeDark},
        {Current: types.ThemeAuto},
    }

    layouts := []string{"default", "featured", "horizontal", "compact"}
    readingModes := []bool{false, true}
    fontSizes := []string{"small", "medium", "large"}

    var components []templ.Component

    // Generate all theme combinations
    for _, theme := range themes {
        for _, readingMode := range readingModes {
            components = append(components, views.BlogLayout(theme, readingMode))
        }
    }

    // Generate all layout combinations
    for _, layout := range layouts {
        components = append(components, views.ArticleCard(sampleArticle, layout))
    }

    // Generate typography variations
    for _, readingMode := range readingModes {
        for _, fontSize := range fontSizes {
            components = append(components, 
                views.ArticleContent(sampleArticle, readingMode, fontSize),
                views.ReadingControls(readingMode, fontSize, types.ThemeLight),
            )
        }
    }

    if err := twerge.CodeGen(
        twerge.Default(),
        "styles/generated.go",
        "styles/input.css",
        "styles/classes.html",
        components...,
    ); err != nil {
        panic(err)
    }
}
```

## Expected Outcomes

1. **Theme Consistency**: Seamless switching between light/dark modes without style conflicts
2. **Reading Experience**: Enhanced typography with user-controllable reading preferences
3. **Content Flexibility**: Unified styling system across different content types
4. **Performance**: Optimized CSS bundle despite complex theming and layout variations
5. **Accessibility**: Consistent focus states and readable typography across all themes