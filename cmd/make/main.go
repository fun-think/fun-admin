package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func main() {
	// CLI flag definitions
	var (
		action string
		name   string
	)

	flag.StringVar(&action, "action", "", "Action to perform (make:resource, make:page)")
	flag.StringVar(&name, "name", "", "Name of the resource or page")
	flag.Parse()

	// validate required flags
	if action == "" {
		fmt.Println("Error: action is required")
		printUsage()
		os.Exit(1)
	}

	switch action {
	case "make:resource":
		if name == "" {
			fmt.Println("Error: name is required for make:resource")
			printUsage()
			os.Exit(1)
		}
		makeResource(name)
	case "make:page":
		if name == "" {
			fmt.Println("Error: name is required for make:page")
			printUsage()
			os.Exit(1)
		}
		makePage(name)
	default:
		fmt.Printf("Error: unknown action '%s'\n", action)
		printUsage()
		os.Exit(1)
	}
}

// makeResource creates a new resource scaffolding file.
func makeResource(name string) {
	fmt.Printf("Creating resource: %s\n", name)

	baseName, snakeName, kebabName, title, err := prepareNames(name)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if err := ensureDir(filepath.Join("internal", "resources")); err != nil {
		fmt.Printf("Failed to ensure resources directory: %v\n", err)
		return
	}

	filePath := filepath.Join("internal", "resources", fmt.Sprintf("%s_resource.go", snakeName))
	if err := createFile(filePath, buildResourceContent(baseName, title, kebabName)); err != nil {
		fmt.Printf("Error creating resource file: %v\n", err)
		return
	}

	fmt.Printf("Resource file created: %s\n", filePath)
	fmt.Println("Remember to register the resource in cmd/bootstrap/bootstrap.go:")
	fmt.Printf("  admin.GlobalResourceManager.Register(resources.New%sResource())\n", baseName)
}

// makePage creates a new custom admin page stub.
func makePage(name string) {
	fmt.Printf("Creating page: %s\n", name)

	baseName, snakeName, kebabName, title, err := prepareNames(name)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if err := ensureDir(filepath.Join("internal", "pages")); err != nil {
		fmt.Printf("Failed to ensure pages directory: %v\n", err)
		return
	}

	filePath := filepath.Join("internal", "pages", fmt.Sprintf("%s_page.go", snakeName))
	if err := createFile(filePath, buildPageContent(baseName, title, kebabName)); err != nil {
		fmt.Printf("Error creating page file: %v\n", err)
		return
	}

	fmt.Printf("Page file created: %s\n", filePath)
	fmt.Println("Remember to register the page in cmd/bootstrap/bootstrap.go:")
	fmt.Printf("  admin.GlobalResourceManager.RegisterPage(pages.New%sPage())\n", baseName)
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func createFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

func buildResourceContent(baseName, title, slug string) string {
	if slug == "" {
		slug = strings.ToLower(baseName)
	}

	typeName := baseName + "Resource"

	return fmt.Sprintf(`package resources

import "fun-admin/pkg/admin"

// %s defines the admin resource for %s.
type %s struct {
	admin.BaseResource
}

// New%sResource creates a new %s resource instance.
func New%sResource() *%s {
	return &%s{}
}

// GetTitle returns the display name of the resource.
func (r *%s) GetTitle() string {
	return "%s"
}

// GetSlug returns the resource slug that is used in API routes.
func (r *%s) GetSlug() string {
	return "%s"
}

// GetModel returns the underlying model for the resource.
func (r *%s) GetModel() interface{} {
	// TODO: return the appropriate model instance, e.g. &model.%s{}
	return nil
}

// GetFields returns the editable fields for the resource.
func (r *%s) GetFields() []admin.Field {
	return []admin.Field{}
}

// GetActions returns the actions that can be executed on the resource.
func (r *%s) GetActions() []admin.Action {
	return []admin.Action{}
}

// GetReadOnlyFields returns the readonly field names that should not be editable.
func (r *%s) GetReadOnlyFields() []string {
	return []string{}
}

// GetColumns returns the table columns displayed in the list view.
func (r *%s) GetColumns() []*admin.Column {
	return []*admin.Column{}
}

// GetFilters returns the filters available for the list view.
func (r *%s) GetFilters() []*admin.Filter {
	return []*admin.Filter{}
}
`, typeName, title, typeName, baseName, title, baseName, typeName, typeName, typeName, title, typeName, slug, typeName, baseName, typeName, typeName, typeName, typeName, typeName)
}

func buildPageContent(baseName, title, slug string) string {
	if slug == "" {
		slug = strings.ToLower(baseName)
	}

	pagePath := "/" + slug
	typeName := baseName + "Page"

	return fmt.Sprintf(`package pages

import "fun-admin/pkg/admin"

// %s defines a custom admin page.
type %s struct {
	*admin.BasePage
}

// New%sPage creates the %s admin page instance.
func New%sPage() *%s {
	page := admin.NewBasePage("%s", "%s", "%s")
	return &%s{
		BasePage: page,
	}
}
`, typeName, typeName, baseName, title, baseName, typeName, title, slug, pagePath, typeName)
}

func prepareNames(input string) (baseName, snakeName, kebabName, title string, err error) {
	words := splitWords(input)
	if len(words) == 0 {
		return "", "", "", "", fmt.Errorf("name must contain letters or numbers")
	}

	baseName = toPascal(words)
	if baseName == "" {
		return "", "", "", "", fmt.Errorf("invalid resource or page name")
	}

	snakeName = joinLower(words, "_")
	kebabName = joinLower(words, "-")
	title = strings.Join(capitalizeWords(words), " ")
	return baseName, snakeName, kebabName, title, nil
}

func splitWords(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var builder strings.Builder
	for _, r := range input {
		if r == '-' || r == '_' || unicode.IsSpace(r) {
			builder.WriteRune(' ')
			continue
		}
		builder.WriteRune(r)
	}

	var words []string
	for _, field := range strings.Fields(builder.String()) {
		words = append(words, splitCamelCase(field)...)
	}

	return words
}

func splitCamelCase(word string) []string {
	if word == "" {
		return nil
	}

	var parts []string
	runes := []rune(word)
	var buffer []rune

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
			parts = append(parts, string(buffer))
			buffer = buffer[:0]
		}
		buffer = append(buffer, r)
	}

	if len(buffer) > 0 {
		parts = append(parts, string(buffer))
	}

	return parts
}

func toPascal(words []string) string {
	var builder strings.Builder
	for _, word := range words {
		builder.WriteString(capitalize(word))
	}
	return builder.String()
}

func joinLower(words []string, sep string) string {
	lowered := make([]string, 0, len(words))
	for _, word := range words {
		if word == "" {
			continue
		}
		lowered = append(lowered, strings.ToLower(word))
	}
	return strings.Join(lowered, sep)
}

func capitalizeWords(words []string) []string {
	capitalized := make([]string, 0, len(words))
	for _, word := range words {
		capitalized = append(capitalized, capitalize(word))
	}
	return capitalized
}

func capitalize(word string) string {
	if word == "" {
		return ""
	}

	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/make/main.go -action=make:resource -name=ResourceName")
	fmt.Println("  go run cmd/make/main.go -action=make:page -name=PageName")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/make/main.go -action=make:resource -name=User")
	fmt.Println("  go run cmd/make/main.go -action=make:page -name=Dashboard")
}
