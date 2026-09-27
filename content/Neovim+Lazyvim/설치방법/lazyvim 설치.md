# Lazyvim 설치
```bash
## 기존 Neovim 설정 백업 (있다면)
mv ~/.config/nvim ~/.config/nvim.bak 2>/dev/null
mv ~/.local/share/nvim ~/.local/share/nvim.bak 2>/dev/null
mv ~/.cache/nvim ~/.cache/nvim.bak 2>/dev/null

## LazyVim starter 클론
git clone https://github.com/LazyVim/starter ~/.config/nvim

## .git 제거 (본인 설정으로 관리하기 위해)
rm -rf ~/.config/nvim/.git
```

# 이후 진행
nvim 으로 lazyvim 플러그인 설치
이후 :q 후 다시 nvim 으로 실행
