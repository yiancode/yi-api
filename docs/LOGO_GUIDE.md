# Yi-API Logo 使用指南

## 📋 Logo 文件说明

Yi-API 品牌 Logo 基于 **"Quantum Relay"（量子中继）** 设计哲学，采用科技感十足的蓝紫渐变配色，完美体现了 AI API 网关的核心理念。

### 文件位置

所有 Logo 文件存储在 `web/public/` 目录：

```
web/public/
├── logo.png              # 主 Logo（深色版，1200×1200）
├── logo-light.png        # 浅色背景版 Logo（1200×1200）
├── favicon.ico           # 标准 ICO 格式 favicon（多尺寸）
├── favicon-256.png       # 256×256 PNG favicon
├── favicon-512.png       # 512×512 PNG favicon（Apple Touch Icon）
└── [备份文件]
    ├── logo.png.backup
    └── favicon.ico.backup
```

## 🎨 设计元素说明

### 核心符号

- **无限符号 (∞)** - 融合"8"的形态，象征无限的 API 接入能力
- **双网关节点** - 两个圆形代表"0"和网关概念，体现中继功能
- **电路纹理** - 六边形网格背景，暗示电路板和量子网络
- **蓝紫渐变** - 从光子蓝(#2D5AFF)到量子紫(#9146FF)
- **发光效果** - 多层光晕，营造能量传输视觉感受
- **"AI80"标识** - 简洁字体标签，呼应域名 api.ai80.vip

### 配色方案

**深色背景版：**
- 背景色：`#060A18` (深宇宙蓝)
- 主色调：`#2D69FF` → `#9141FF` (蓝紫渐变)
- 节点色：`#55AFFF` (电蓝)
- 文字色：`#CDD7F0` (柔光白)

**浅色背景版：**
- 背景色：`#F7F9FE` (纯净白蓝)
- 主色调：`#1E5AE1` → `#7D32E1` (深蓝紫)
- 节点色：`#2D7DF5` (鲜明蓝)
- 文字色：`#141E46` (深海军蓝)

## 💻 使用场景

### 1. 网站 Favicon

已在 `web/index.html` 中自动配置：

```html
<!-- Favicon 配置 -->
<link rel="icon" type="image/x-icon" href="/favicon.ico" />
<link rel="icon" type="image/png" sizes="256x256" href="/favicon-256.png" />
<link rel="icon" type="image/png" sizes="512x512" href="/favicon-512.png" />
<link rel="apple-touch-icon" sizes="512x512" href="/favicon-512.png" />
```

### 2. 代码中引用 Logo

项目使用 `getLogo()` 函数获取 Logo 路径（定义在 `web/src/helpers/utils.jsx`）：

```javascript
import { getLogo } from '@/helpers/utils';

// 使用主 Logo（深色版）
<img src={getLogo()} alt="Yi-API Logo" />

// 使用浅色版
<img src="/logo-light.png" alt="Yi-API Logo" />
```

### 3. 在 React 组件中使用

```jsx
// 自适应深色/浅色主题
const isDark = useTheme().mode === 'dark';
const logoSrc = isDark ? '/logo.png' : '/logo-light.png';

<img src={logoSrc} alt="Yi-API" style={{ width: 200 }} />
```

### 4. GitHub / 文档使用

**README.md 头部：**
```markdown
<p align="center">
  <img src="web/public/logo.png" alt="Yi-API Logo" width="200" />
</p>
```

**文档封面：**
- 使用 `logo.png`（深色背景场景）
- 使用 `logo-light.png`（白色/浅色背景）

### 5. 社交媒体头像

推荐使用：`favicon-512.png`（512×512 方形，适合各平台）

## 🔧 构建说明

Logo 文件会在前端构建时自动复制到 `web/dist/` 目录：

```bash
cd web
bun run build
```

构建后的文件：
```
web/dist/
├── logo.png
├── logo-light.png
├── favicon.ico
├── favicon-256.png
└── favicon-512.png
```

## 📐 尺寸规范

| 文件 | 尺寸 | 用途 |
|------|------|------|
| `logo.png` | 1200×1200 | 网站展示、文档 |
| `logo-light.png` | 1200×1200 | 浅色背景场景 |
| `favicon.ico` | 多尺寸 | 浏览器标签图标 |
| `favicon-256.png` | 256×256 | 标准 favicon |
| `favicon-512.png` | 512×512 | 高清 favicon / Apple Touch Icon |

## ⚠️ 使用注意事项

1. **保持比例**：Logo 为方形设计，缩放时保持 1:1 比例
2. **最小尺寸**：建议不小于 32×32 像素，以保持清晰度
3. **背景适配**：
   - 深色界面使用 `logo.png`
   - 浅色界面使用 `logo-light.png`
4. **不要修改**：避免拉伸、变形、改色或添加效果
5. **保持完整**：不要裁剪 Logo 的任何部分

## 🎯 品牌标准

### DO（推荐做法）✅

- 在深色背景上使用深色版 Logo
- 在浅色背景上使用浅色版 Logo
- 为 Logo 保留足够的呼吸空间（四周至少 Logo 高度的 10%）
- 保持 Logo 清晰可辨

### DON'T（避免做法）❌

- 不要改变 Logo 的颜色或渐变
- 不要旋转或倾斜 Logo
- 不要在复杂背景上使用（确保对比度）
- 不要添加阴影、描边等额外效果（Logo 自带发光效果）

## 📞 问题反馈

如需其他尺寸或格式的 Logo，请联系开发团队。

---

**设计哲学**：Quantum Relay
**创建日期**：2026-01-05
**版本**：v1.0
