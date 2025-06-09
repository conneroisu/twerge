# Dynamic Form Validation Use Case

## Overview
A comprehensive form validation system demonstrating Twerge's ability to handle complex state-dependent styling, real-time validation feedback, and progressive enhancement. Showcases dynamic error states, loading indicators, and accessibility features.

## Key Features
- Real-time field validation with visual feedback
- Multi-step form progression with state persistence
- Complex validation rules (async validation, dependencies)
- Accessibility-focused error messaging
- Loading states and submission feedback
- Conditional field visibility based on user input

## Twerge Benefits Demonstrated

### 1. Dynamic Field Validation States
```go
type FieldState struct {
    Value     string
    Touched   bool
    Validating bool
    Valid     *bool  // nil = not validated, true = valid, false = invalid
    Errors    []string
    Warnings  []string
}

type ValidationSeverity string

const (
    SeverityError   ValidationSeverity = "error"
    SeverityWarning ValidationSeverity = "warning"
    SeveritySuccess ValidationSeverity = "success"
)

templ FormField(field FieldState, fieldType string, label string, required bool) {
    <div class={twerge.It("space-y-2")}>
        @FieldLabel(label, required, getValidationSeverity(field))
        @FieldInput(field, fieldType)
        @FieldFeedback(field)
    </div>
}

templ FieldLabel(label string, required bool, severity ValidationSeverity) {
    <label class={twerge.It(
        "block text-sm font-medium transition-colors duration-200 " +
        getLabelClasses(severity)
    )}>
        { label }
        if required {
            <span class={twerge.It(
                "text-red-500 ml-1 " +
                twerge.If(severity == SeverityError, "animate-pulse", "")
            )}>*</span>
        }
    </label>
}

func getLabelClasses(severity ValidationSeverity) string {
    switch severity {
    case SeverityError:
        return "text-red-700 dark:text-red-400"
    case SeverityWarning:
        return "text-amber-700 dark:text-amber-400"
    case SeveritySuccess:
        return "text-green-700 dark:text-green-400"
    default:
        return "text-gray-700 dark:text-gray-300"
    }
}

templ FieldInput(field FieldState, fieldType string) {
    <div class={twerge.It("relative")}>
        <input 
            type={fieldType}
            value={field.Value}
            class={twerge.It(
                "w-full px-3 py-2 border rounded-lg shadow-sm " +
                "transition-all duration-200 focus:outline-none " +
                "placeholder-gray-400 dark:placeholder-gray-500 " +
                "bg-white dark:bg-gray-800 " +
                getInputClasses(field) +
                twerge.If(field.Validating, " pr-10", "") +
                getFocusRingClasses(getValidationSeverity(field))
            )}
            aria-invalid={field.Valid != nil && !*field.Valid}
            aria-describedby={getAriaDescribedBy(field)}
        />
        
        <!-- Validation Icons and Loading Spinner -->
        <div class={twerge.It(
            "absolute inset-y-0 right-0 flex items-center pr-3 " +
            "pointer-events-none"
        )}>
            if field.Validating {
                <div class={twerge.It(
                    "animate-spin h-4 w-4 border-2 border-blue-500 " +
                    "border-t-transparent rounded-full"
                )}></div>
            } else if field.Valid != nil {
                @ValidationIcon(*field.Valid, len(field.Warnings) > 0)
            }
        </div>
    </div>
}

func getInputClasses(field FieldState) string {
    if !field.Touched {
        return "border-gray-300 dark:border-gray-600 " +
               "focus:border-blue-500 dark:focus:border-blue-400"
    }

    severity := getValidationSeverity(field)
    switch severity {
    case SeverityError:
        return "border-red-500 dark:border-red-400 " +
               "bg-red-50 dark:bg-red-900/20 " +
               "text-red-900 dark:text-red-100"
    case SeverityWarning:
        return "border-amber-500 dark:border-amber-400 " +
               "bg-amber-50 dark:bg-amber-900/20 " +
               "text-amber-900 dark:text-amber-100"
    case SeveritySuccess:
        return "border-green-500 dark:border-green-400 " +
               "bg-green-50 dark:bg-green-900/20 " +
               "text-green-900 dark:text-green-100"
    default:
        return "border-gray-300 dark:border-gray-600 " +
               "focus:border-blue-500 dark:focus:border-blue-400"
    }
}

func getFocusRingClasses(severity ValidationSeverity) string {
    switch severity {
    case SeverityError:
        return " focus:ring-2 focus:ring-red-500 focus:ring-opacity-50"
    case SeverityWarning:
        return " focus:ring-2 focus:ring-amber-500 focus:ring-opacity-50"
    case SeveritySuccess:
        return " focus:ring-2 focus:ring-green-500 focus:ring-opacity-50"
    default:
        return " focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
    }
}

templ ValidationIcon(valid bool, hasWarnings bool) {
    if valid && !hasWarnings {
        <svg class={twerge.It("h-4 w-4 text-green-500 animate-scale-in")}>
            <!-- Checkmark icon -->
        </svg>
    } else if valid && hasWarnings {
        <svg class={twerge.It("h-4 w-4 text-amber-500 animate-scale-in")}>
            <!-- Warning icon -->
        </svg>
    } else {
        <svg class={twerge.It("h-4 w-4 text-red-500 animate-shake")}>
            <!-- X icon -->
        </svg>
    }
}
```

### 2. Progressive Multi-Step Form
```go
type FormStep struct {
    ID          string
    Title       string
    Description string
    Fields      []string
    Optional    bool
    Completed   bool
    Current     bool
    Accessible  bool
}

templ MultiStepForm(steps []FormStep, currentStep int) {
    <div class={twerge.It("max-w-4xl mx-auto")}>
        @StepIndicator(steps, currentStep)
        @StepContent(steps[currentStep])
        @StepNavigation(steps, currentStep)
    </div>
}

templ StepIndicator(steps []FormStep, currentStep int) {
    <nav class={twerge.It("mb-8")} aria-label="Form progress">
        <ol class={twerge.It(
            "flex items-center justify-between border border-gray-200 " +
            "dark:border-gray-700 rounded-lg p-4 bg-gray-50 dark:bg-gray-800"
        )}>
            for i, step := range steps {
                @StepIndicatorItem(step, i, currentStep, i < len(steps)-1)
            }
        </ol>
    </nav>
}

templ StepIndicatorItem(step FormStep, index int, currentStep int, hasNext bool) {
    <li class={twerge.It(
        "flex items-center " +
        twerge.If(hasNext, "flex-1", "")
    )}>
        <!-- Step Circle -->
        <div class={twerge.It(
            "flex items-center justify-center w-8 h-8 rounded-full " +
            "transition-all duration-300 " +
            getStepCircleClasses(step, index, currentStep)
        )}>
            if step.Completed {
                <svg class={twerge.It("w-4 h-4 text-white")}>
                    <!-- Checkmark -->
                </svg>
            } else {
                <span class={twerge.It(
                    "text-sm font-medium " +
                    getStepNumberClasses(step, index, currentStep)
                )}>
                    { strconv.Itoa(index + 1) }
                </span>
            }
        </div>

        <!-- Step Label -->
        <div class={twerge.It("ml-3 min-w-0 flex-1")}>
            <h3 class={twerge.It(
                "text-sm font-medium transition-colors duration-200 " +
                getStepTitleClasses(step, index, currentStep)
            )}>
                { step.Title }
            </h3>
            if step.Description != "" {
                <p class={twerge.It("text-xs text-gray-500 dark:text-gray-400 mt-1")}>
                    { step.Description }
                </p>
            }
        </div>

        <!-- Connector Line -->
        if hasNext {
            <div class={twerge.It(
                "flex-1 h-px mx-4 transition-colors duration-300 " +
                twerge.If(step.Completed, 
                    "bg-green-500", 
                    "bg-gray-300 dark:bg-gray-600"
                )
            )}></div>
        }
    </li>
}

func getStepCircleClasses(step FormStep, index int, currentStep int) string {
    if step.Completed {
        return "bg-green-500 text-white shadow-lg"
    } else if index == currentStep {
        return "bg-blue-600 text-white shadow-lg ring-2 ring-blue-500 ring-opacity-50"
    } else if step.Accessible {
        return "bg-white dark:bg-gray-700 border-2 border-gray-300 dark:border-gray-600 " +
               "hover:border-blue-500 dark:hover:border-blue-400 cursor-pointer"
    } else {
        return "bg-gray-200 dark:bg-gray-600 border-2 border-gray-200 dark:border-gray-600"
    }
}
```

### 3. Complex Field Dependencies and Conditional Logic
```go
type ConditionalField struct {
    FieldName   string
    Condition   string
    ShowWhen    interface{}
    RequiredWhen *string
}

templ ConditionalFieldGroup(fields []ConditionalField, formData map[string]interface{}) {
    <div class={twerge.It("space-y-6")}>
        for _, field := range fields {
            if shouldShowField(field, formData) {
                <div class={twerge.It(
                    "transform transition-all duration-300 ease-in-out " +
                    "animate-fade-in-up"
                )}>
                    @FieldByName(field.FieldName, isFieldRequired(field, formData))
                </div>
            }
        }
    </div>
}

// Example: Account type selection affecting visible fields
templ AccountTypeFields(accountType string, businessInfo BusinessInfo) {
    <!-- Base account information -->
    @FormField(FieldState{}, "email", "Email Address", true)
    @FormField(FieldState{}, "password", "Password", true)
    
    <!-- Conditional business fields -->
    if accountType == "business" {
        <div class={twerge.It(
            "space-y-4 p-4 border border-blue-200 dark:border-blue-800 " +
            "rounded-lg bg-blue-50 dark:bg-blue-900/20 " +
            "animate-fade-in"
        )}>
            <h4 class={twerge.It("font-medium text-blue-900 dark:text-blue-100")}>
                Business Information
            </h4>
            
            @FormField(FieldState{}, "text", "Company Name", true)
            @FormField(FieldState{}, "text", "Tax ID", true)
            
            <!-- Nested conditional: Company size affects additional fields -->
            @CompanySizeSelect(businessInfo.Size)
            
            if businessInfo.Size == "enterprise" {
                <div class={twerge.It(
                    "pl-4 border-l-2 border-blue-300 dark:border-blue-700 " +
                    "animate-slide-down"
                )}>
                    @FormField(FieldState{}, "text", "Dedicated Account Manager", false)
                    @FormField(FieldState{}, "number", "Estimated Annual Volume", true)
                </div>
            }
        </div>
    }
}
```

### 4. Real-time Validation with Debouncing
```go
templ AsyncValidationField(field FieldState, fieldName string, placeholder string) {
    <div class={twerge.It("space-y-2")}>
        <label class={twerge.It("block text-sm font-medium text-gray-700 dark:text-gray-300")}>
            { strings.Title(fieldName) }
        </label>
        
        <div class={twerge.It("relative")}>
            <input 
                type="text"
                placeholder={placeholder}
                value={field.Value}
                class={twerge.It(
                    "w-full px-3 py-2 border rounded-lg " +
                    "transition-all duration-200 focus:outline-none " +
                    "bg-white dark:bg-gray-800 " +
                    getAsyncValidationClasses(field)
                )}
                hx-post={"/validate/" + fieldName}
                hx-trigger="keyup changed delay:500ms"
                hx-target={".validation-" + fieldName}
                hx-indicator={".spinner-" + fieldName}
            />
            
            <!-- Loading indicator -->
            <div class={twerge.It(
                "absolute inset-y-0 right-0 flex items-center pr-3 " +
                "spinner-" + fieldName + " htmx-indicator"
            )}>
                <div class={twerge.It(
                    "animate-spin h-4 w-4 border-2 border-blue-500 " +
                    "border-t-transparent rounded-full"
                )}></div>
            </div>
        </div>
        
        <!-- Validation feedback area -->
        <div class={twerge.It("validation-" + fieldName)}>
            if field.Valid != nil {
                @AsyncValidationFeedback(field)
            }
        </div>
    </div>
}

func getAsyncValidationClasses(field FieldState) string {
    base := "border-gray-300 dark:border-gray-600 focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
    
    if field.Validating {
        return base + " pr-10 border-blue-300 dark:border-blue-500"
    }
    
    if field.Valid == nil {
        return base
    }
    
    if *field.Valid {
        return "border-green-500 dark:border-green-400 focus:border-green-500 " +
               "focus:ring-green-500 focus:ring-opacity-50 pr-10"
    }
    
    return "border-red-500 dark:border-red-400 focus:border-red-500 " +
           "focus:ring-red-500 focus:ring-opacity-50 pr-10"
}

templ AsyncValidationFeedback(field FieldState) {
    if len(field.Errors) > 0 {
        for _, error := range field.Errors {
            <div class={twerge.It(
                "flex items-start space-x-2 text-sm text-red-600 dark:text-red-400 " +
                "animate-fade-in"
            )}>
                <svg class={twerge.It("w-4 h-4 mt-0.5 flex-shrink-0")}>
                    <!-- Error icon -->
                </svg>
                <span>{ error }</span>
            </div>
        }
    } else if *field.Valid {
        <div class={twerge.It(
            "flex items-center space-x-2 text-sm text-green-600 dark:text-green-400 " +
            "animate-fade-in"
        )}>
            <svg class={twerge.It("w-4 h-4")}>
                <!-- Success icon -->
            </svg>
            <span>Looks good!</span>
        </div>
    }
}
```

### 5. Form Submission States and Feedback
```go
type SubmissionState struct {
    Submitting bool
    Success    bool
    Error      *string
    Progress   *int // for file uploads or multi-step processing
}

templ SubmitButton(state SubmissionState, disabled bool, label string) {
    <button 
        type="submit"
        disabled={disabled || state.Submitting}
        class={twerge.It(
            "relative w-full flex justify-center py-3 px-4 " +
            "border border-transparent rounded-lg text-sm font-medium " +
            "transition-all duration-200 focus:outline-none focus:ring-2 " +
            "focus:ring-offset-2 " +
            getSubmitButtonClasses(state, disabled)
        )}
    >
        if state.Submitting {
            @SubmittingContent(state.Progress)
        } else if state.Success {
            @SuccessContent()
        } else {
            { label }
        }
    </button>
}

func getSubmitButtonClasses(state SubmissionState, disabled bool) string {
    if disabled && !state.Submitting {
        return "bg-gray-300 text-gray-500 cursor-not-allowed " +
               "dark:bg-gray-600 dark:text-gray-400"
    }
    
    if state.Submitting {
        return "bg-blue-500 text-white cursor-wait " +
               "focus:ring-blue-500"
    }
    
    if state.Success {
        return "bg-green-600 text-white " +
               "focus:ring-green-500"
    }
    
    if state.Error != nil {
        return "bg-red-600 text-white hover:bg-red-700 " +
               "focus:ring-red-500"
    }
    
    return "bg-blue-600 text-white hover:bg-blue-700 " +
           "active:bg-blue-800 transform hover:scale-105 active:scale-95 " +
           "focus:ring-blue-500"
}

templ SubmittingContent(progress *int) {
    <div class={twerge.It("flex items-center space-x-2")}>
        <div class={twerge.It(
            "animate-spin h-4 w-4 border-2 border-white " +
            "border-t-transparent rounded-full"
        )}></div>
        <span>
            if progress != nil {
                Processing... { strconv.Itoa(*progress) }%
            } else {
                Submitting...
            }
        </span>
    </div>
}

templ SuccessContent() {
    <div class={twerge.It("flex items-center space-x-2 animate-bounce")}>
        <svg class={twerge.It("w-4 h-4")}>
            <!-- Checkmark icon -->
        </svg>
        <span>Success!</span>
    </div>
}

templ FormErrorSummary(errors []string) {
    if len(errors) > 0 {
        <div class={twerge.It(
            "mb-6 p-4 border border-red-200 dark:border-red-800 " +
            "rounded-lg bg-red-50 dark:bg-red-900/20 " +
            "animate-shake"
        )}>
            <div class={twerge.It("flex items-start space-x-3")}>
                <svg class={twerge.It("w-5 h-5 text-red-500 mt-0.5 flex-shrink-0")}>
                    <!-- Error icon -->
                </svg>
                <div>
                    <h3 class={twerge.It("text-sm font-medium text-red-800 dark:text-red-200")}>
                        Please correct the following errors:
                    </h3>
                    <ul class={twerge.It("mt-2 text-sm text-red-700 dark:text-red-300 space-y-1")}>
                        for _, error := range errors {
                            <li class={twerge.It("flex items-start space-x-1")}>
                                <span class={twerge.It("inline-block w-1 h-1 bg-red-500 rounded-full mt-2 flex-shrink-0")}></span>
                                <span>{ error }</span>
                            </li>
                        }
                    </ul>
                </div>
            </div>
        </div>
    }
}
```

## Performance Benefits

### State-Dependent CSS Optimization
```css
/* twerge:begin */
/* Input states consolidated */
.tw-1 { @apply w-full px-3 py-2 border rounded-lg shadow-sm transition-all duration-200 focus:outline-none placeholder-gray-400 dark:placeholder-gray-500 bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 focus:border-blue-500 dark:focus:border-blue-400; }
.tw-2 { @apply w-full px-3 py-2 border rounded-lg shadow-sm transition-all duration-200 focus:outline-none placeholder-gray-400 dark:placeholder-gray-500 bg-white dark:bg-gray-800 border-red-500 dark:border-red-400 bg-red-50 dark:bg-red-900/20 text-red-900 dark:text-red-100 pr-10; }
.tw-3 { @apply w-full px-3 py-2 border rounded-lg shadow-sm transition-all duration-200 focus:outline-none placeholder-gray-400 dark:placeholder-gray-500 bg-white dark:bg-gray-800 border-green-500 dark:border-green-400 bg-green-50 dark:bg-green-900/20 text-green-900 dark:text-green-100 pr-10; }

/* Animation states unified */
.tw-4 { @apply transform transition-all duration-300 ease-in-out animate-fade-in-up; }
.tw-5 { @apply space-y-4 p-4 border border-blue-200 dark:border-blue-800 rounded-lg bg-blue-50 dark:bg-blue-900/20 animate-fade-in; }
/* twerge:end */
```

### Validation State Benefits
- **Real-time feedback**: Smooth transitions between validation states
- **Loading indicators**: Consistent spinner and progress animations
- **Error handling**: Unified error styling across all field types
- **Accessibility**: Proper ARIA attributes without class conflicts

## Code Generation Integration

```go
//go:build ignore

package main

import (
    "github.com/conneroisu/twerge"
    "github.com/yourproject/forms/views"
    "github.com/yourproject/forms/types"
)

func main() {
    // Generate all validation states
    states := []types.FieldState{
        {Value: "", Touched: false, Valid: nil},
        {Value: "test", Touched: true, Valid: nil, Validating: true},
        {Value: "valid@email.com", Touched: true, Valid: &[]bool{true}[0]},
        {Value: "invalid", Touched: true, Valid: &[]bool{false}[0], Errors: []string{"Invalid format"}},
    }

    var components []templ.Component

    // Generate all field variations
    for _, state := range states {
        components = append(components,
            views.FormField(state, "text", "Sample Field", true),
            views.AsyncValidationField(state, "email", "Enter email"),
        )
    }

    // Generate form states
    submissionStates := []types.SubmissionState{
        {Submitting: false, Success: false},
        {Submitting: true, Progress: &[]int{45}[0]},
        {Submitting: false, Success: true},
        {Submitting: false, Error: &[]string{"Server error"}[0]},
    }

    for _, state := range submissionStates {
        components = append(components,
            views.SubmitButton(state, false, "Submit Form"),
        )
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

1. **Enhanced UX**: Smooth state transitions without jarring class changes
2. **Accessibility**: Consistent focus states and ARIA attributes
3. **Performance**: Optimized CSS for complex validation scenarios
4. **Maintainability**: Unified styling system across all form components
5. **Developer Experience**: Simple conditional styling without class conflict management