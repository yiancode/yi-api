# GitHub Actions 工作流说明

## 🔄 自动同步上游仓库 (sync-upstream.yml)

### 功能说明

自动从上游仓库 [QuantumNous/new-api](https://github.com/QuantumNous/new-api) 同步最新代码到本 fork 仓库。

### 触发方式

#### 1. 自动触发 (定时任务)
- **执行时间**: 每天北京时间早上 9:00 (UTC 1:00)
- **自动执行**: 无需人工干预
- **冲突策略**: 默认采用上游版本 (`theirs`)

#### 2. 手动触发
在 GitHub 仓库页面操作：
1. 进入 `Actions` 标签页
2. 选择 `🔄 同步上游仓库` workflow
3. 点击 `Run workflow` 按钮
4. 选择冲突解决策略：
   - `theirs` (推荐) - 采用上游版本
   - `ours` - 保留我们的版本
5. 点击绿色的 `Run workflow` 按钮

### 冲突解决策略

#### `theirs` - 采用上游版本 (默认，推荐)
- ✅ 优点：保持与上游代码一致，获取最新功能和修复
- ⚠️ 注意：会覆盖本地对冲突文件的修改
- 📝 适用场景：大部分情况，尤其是配置文件和文档冲突

#### `ours` - 保留我们的版本
- ✅ 优点：保留本地的自定义修改
- ⚠️ 注意：可能错过上游的重要更新
- 📝 适用场景：对特定文件做了重要自定义且不想被覆盖

### 工作流程

```
1. 📥 检出代码
   ↓
2. ⚙️ 配置 Git
   ↓
3. 🔗 添加上游仓库并获取最新代码
   ↓
4. 📊 检查是否有新的提交
   ↓
5. 🔄 合并上游更新
   ├─ 无冲突 → 自动合并
   └─ 有冲突 → 应用解决策略
   ↓
6. 📤 推送到本仓库
   ↓
7. 📝 生成同步报告
```

### 查看同步结果

#### 方式 1: GitHub Actions 页面
1. 进入仓库的 `Actions` 标签页
2. 查看最新的 `🔄 同步上游仓库` 运行记录
3. 点击进入可查看详细日志和同步报告

#### 方式 2: 提交历史
查看仓库的提交记录，会看到类似这样的提交：
```
🔄 自动同步上游: 解决冲突 (策略: theirs)

自动同步来自 QuantumNous/new-api 的更新
冲突解决策略: theirs
冲突文件数: 2

🤖 由 GitHub Actions 自动执行
```

### 常见问题

#### Q1: 为什么同步失败了？
**A**: 可能的原因：
- 网络问题：GitHub Actions 访问上游仓库失败
- 复杂冲突：自动策略无法解决的复杂冲突
- 权限问题：GITHUB_TOKEN 权限不足

**解决方法**：
1. 查看 Actions 日志了解具体错误
2. 使用手动触发重试
3. 如果持续失败，需要本地手动解决冲突

#### Q2: 如何修改同步时间？
**A**: 编辑 `.github/workflows/sync-upstream.yml` 文件的 cron 表达式：
```yaml
schedule:
  - cron: '0 1 * * *'  # UTC 时间，每天 01:00
```

常用时间示例：
- `0 1 * * *` - 每天 UTC 01:00 (北京时间 09:00)
- `0 */12 * * *` - 每 12 小时
- `0 0 * * 1` - 每周一 UTC 00:00

#### Q3: 如何临时停止自动同步？
**A**: 有两种方式：
1. **禁用 workflow**: 在 Actions 页面点击 workflow 右上角的 `...` → `Disable workflow`
2. **删除 schedule**: 编辑 workflow 文件，注释掉 `schedule` 部分

#### Q4: 同步会覆盖我的自定义代码吗？
**A**:
- 如果你的修改与上游修改的是**不同文件**或**不同部分**，不会被覆盖
- 如果是**同一文件的同一位置**，会根据选择的策略处理：
  - `theirs`: 采用上游版本（会覆盖）
  - `ours`: 保留你的版本
- **建议**: 重要的自定义代码放在独立的文件中，避免冲突

### 权限要求

Workflow 使用 `GITHUB_TOKEN` 自动授权，具有以下权限：
- ✅ 读取仓库代码
- ✅ 推送提交到 main 分支
- ✅ 写入 workflow 运行状态

无需额外配置 Personal Access Token。

### 高级配置

#### 修改默认冲突策略
编辑 workflow 文件，修改 `STRATEGY` 的默认值：
```yaml
STRATEGY="${{ github.event.inputs.conflict_strategy || 'ours' }}"
```

#### 添加通知
可以添加通知步骤，在同步成功/失败时发送通知（如邮件、Slack、钉钉等）。

#### 同步到多个分支
复制 workflow 文件并修改目标分支，可以实现同步到 develop、staging 等分支。

## 其他 Workflows

### docker-image-*.yml
- Docker 镜像构建和发布

### electron-build.yml
- Electron 桌面应用构建

### release.yml
- 版本发布自动化

### sync-to-gitee.yml
- 同步到 Gitee 镜像仓库

---

💡 **提示**: 建议定期检查 Actions 执行情况，确保同步正常运行。
