package project

type CommandAttr struct {
	Description       string   `json:"description"`
	Aliases           []string `json:"aliases"`
	Permission        string   `json:"permission"`
	PermissionMessage string   `json:"permission_message"`
	Usage             string   `json:"usage"`
}

type PermissionAttr struct {
	Description string          `json:"description"`
	Default     *bool           `json:"default"`
	Children    map[string]bool `json:"children"`
}

type PackageJSON struct {
	Name        string                     `json:"name"`
	Version     string                     `json:"version"`
	Commands    map[string]*CommandAttr    `json:"commands"`
	Permissions map[string]*PermissionAttr `json:"permissions"`
}
