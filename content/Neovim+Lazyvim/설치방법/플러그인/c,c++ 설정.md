---
aliases:
  - c/c++ 설정
---

```bash
# C / C++ 컴파일러 및 디버거
sudo dnf install -y gcc gcc-c++ cmake clang clang-tools-extra gdb lldb

# 개발 헤더
sudo dnf install -y glibc-devel libstdc++-devel

# LLVM (clangd LSP가 내부적으로 사용)
sudo dnf install -y llvm

# 확인
gcc --version
g++ --version
clang --version
```



nvim 내부에서 아래 내용 입력
```
:MasonInstall clangd clang-format codelldb rust-analyzer
```