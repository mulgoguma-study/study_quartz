# <font color="#ff0000">처음에는 이거부터 설정하고 개별 lua 설정</font>



#### 현재는 go, vue3, python, rust, c, c++만 해놓음

방법: `~/.config/nvim/lua/config/lazy.lua` 수정

bash

```bash
nvim ~/.config/nvim/lua/config/lazy.lua
```

`spec` 항목에 아래 추가:

lua

```lua
{ import = "lazyvim.plugins.extras.lang.go" },
```

최종적으로 이런 형태:

lua

````lua
require("lazy").setup({
  spec = {
    { "LazyVim/LazyVim", import = "lazyvim.plugins" },
    
    { import = "lazyvim.plugins.extras.lang.go" }, 
    { import = "lazyvim.plugins.extras.lang.vue" }, 
    { import = "lazyvim.plugins.extras.lang.svelte" },
    { import = "lazyvim.plugins.extras.lang.python" }, 
    {"lazyvim.plugins.extras.lang.clangd" }, -- c/c++
    { import = "lazyvim.plugins.extras.lang.rust" },
    { import = "plugins" },
  },
  defaults = { lazy = false, version = false },
  ...
})
```

저장 후 nvim 재시작하면 Go 관련 플러그인 자동 설치됨.

---

## 8. gopls 및 Go 툴 설치 (Mason으로 자동 설치됨)

LazyVim Go extra를 활성화하면 Mason이 자동으로 설치하지만, 수동으로도 가능:
```
## nvim 내부에서
:MasonInstall gopls goimports golangci-lint-langserver delve
````

또는 Go 커맨드로:

bash

```bash
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/cmd/goimports@latest
```