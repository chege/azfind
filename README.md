# ⚡ azf — Fast Azure Resource Finder

A small, fast CLI for jumping to Azure resources.  
Type → find → open in Azure Portal.

## Features
- Instant search across all subscriptions  
- Local SQLite cache (`--sync`)  
- FZF-powered interactive picker  
- Shell completion (bash / zsh / fish / pwsh)  
- Zero noise, minimal dependencies

## Usage
```bash
azf
azf kvasir
azf --sync
azf --completion bash
```

## Install
```bash
go install github.com/chege/azfind@latest
brew install fzf
```

## Idea
A fast alternative to browsing the Azure Portal.  
Keep context, jump instantly. Minimal by design.

## License
MIT
