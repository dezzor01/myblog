package templates

import "embed"

//go:embed index.html post.html admin_login.html
var FS embed.FS
