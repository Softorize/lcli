# lcli - LinkedIn CLI

A command-line tool for interacting with the LinkedIn API. Manage your profile, posts, comments, reactions, media uploads, organization pages, and analytics directly from the terminal.

## Prerequisites

1. **Go 1.25+** installed
2. A **LinkedIn Developer Application** with:
   - OAuth 2.0 client ID and client secret
   - Redirect URI set to `http://localhost:8484/callback`
   - Required scopes: `r_liteprofile`, `r_emailaddress`, `w_member_social`, `r_organization_social`, `rw_organization_admin` (as needed)

## Installation

### From source

```bash
git clone https://github.com/Softorize/lcli.git
cd lcli
make build
# Binary is at ./bin/lcli
```

### Install to GOPATH

```bash
make install
```

### Go install

```bash
go install github.com/Softorize/lcli/cmd/lcli@latest
```

## Quick Start

```bash
# 1. Configure your LinkedIn app credentials
lcli config setup --client-id YOUR_ID --client-secret YOUR_SECRET

# 2. Authenticate via OAuth
lcli auth login

# 3. View your profile
lcli profile me

# 4. Create a post
lcli post create --text "Hello from lcli!"
```

## Commands

### Authentication

```bash
lcli auth login              # OAuth login via browser
lcli auth login --port 9090  # Use custom callback port
lcli auth logout             # Remove stored credentials
lcli auth status             # Show auth status and token expiry
```

### Configuration

```bash
lcli config setup --client-id ID --client-secret SECRET
```

### Profile

```bash
lcli profile me                           # Your profile
lcli profile me --output json             # JSON output
lcli profile view --id PERSON_ID          # View another profile
```

### Posts

```bash
lcli post create --text "Hello!"                        # Public post
lcli post create --text "Hi" --visibility CONNECTIONS    # Connections only
lcli post create --text "Look!" --image photo.jpg       # Post with image
lcli post create --text "Watch!" --video clip.mp4       # Post with video
lcli post list                                          # List recent posts
lcli post list --count 20 --start 0                     # Paginated
lcli post get URN                                       # Get single post
lcli post delete URN --confirm                          # Delete post
```

### Comments

```bash
lcli comment create --post POST_URN --text "Nice post!"
lcli comment list --post POST_URN
lcli comment list --post POST_URN --count 20
lcli comment delete COMMENT_URN --confirm
```

### Reactions

```bash
lcli reaction like POST_URN                        # Default LIKE
lcli reaction like POST_URN --type CELEBRATE        # Other types
lcli reaction unlike POST_URN
lcli reaction list POST_URN
```

Reaction types: `LIKE`, `CELEBRATE`, `SUPPORT`, `LOVE`, `INSIGHTFUL`, `FUNNY`

### Media

```bash
lcli media upload photo.jpg              # Auto-detects image type
lcli media upload video.mp4              # Auto-detects video type
lcli media upload file.bin --type image  # Manual type override
```

Supported formats: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp` (image), `.mp4`, `.mov`, `.avi`, `.wmv`, `.webm` (video)

### Organizations

```bash
lcli org info --id 12345                 # By numeric ID
lcli org info --vanity company-name      # By vanity name
lcli org followers --org ORG_URN         # Follower stats
lcli org stats --org ORG_URN             # Page view stats
```

### Analytics

```bash
lcli analytics post POST_URN            # Post engagement metrics
lcli analytics views                     # Profile/network size
```

### Shell Completions

```bash
lcli completion bash    # Bash completions
lcli completion zsh     # Zsh completions

# Install completions
lcli completion bash > /etc/bash_completion.d/lcli
lcli completion zsh > "${fpath[1]}/_lcli"
```

### Other

```bash
lcli version    # Print version, commit, build date
lcli help       # Show usage
```

## Output Formats

All read commands support `--output` with three formats:

| Format  | Flag             | Description                  |
|---------|------------------|------------------------------|
| Table   | `--output table` | Aligned columns (default)    |
| JSON    | `--output json`  | Pretty-printed JSON          |
| YAML    | `--output yaml`  | YAML output                  |

## Configuration

Configuration is stored in `~/.config/lcli/`:

- `config.yaml` - Client credentials and settings
- `tokens.json` - OAuth tokens (auto-managed)

## Development

```bash
make build       # Build binary
make test        # Run tests with race detector
make test-cover  # Tests with coverage report
make lint        # Run golangci-lint
make vet         # Run go vet
make fmt         # Format code
make check       # Run fmt + vet + test
make clean       # Remove build artifacts
```

## License

MIT
