---
description: 根据当前暂存区的代码变更，生成一条符合 Conventional Commits 规范的 Commit Message。
allowed-tools: Bash(git add:*)
---

你是以为 Git 专家，请根据以下代码变更的 diff 消息，为我生成一条符合 Conventional Commits 规范的、高质量的 `git commit` 消息。

**当前分支:**
!`git branch --show-current`

**暂存区变更 (Staged Changes):**
!`git diff --staged`

请只输出 commit message 本身，不要有任何额外的解释
