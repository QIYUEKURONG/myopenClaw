---
description: MyOpenClaw 学习开发规则 - 严格参考 OpenClaw 实现
globs:
  - "**/*.go"
alwaysApply: true
---

# MyOpenClaw 学习开发规则

## 🎯 核心原则

1. **必须参考 OpenClaw 源码**
   - 所有设计决策必须先查看 OpenClaw 的实现方式
   - OpenClaw 源码路径：`/Users/qiyuekurong/go/cursor/joke/openclaw/src`
   - 使用 Grep、SemanticSearch 查找 OpenClaw 的实现

2. **渐进式学习**
   - 一次实现一个小功能
   - 可以慢一点，但必须理解透彻
   - 每个模块完成后更新 `PROGRESS.md`

3. **学生主导，AI 辅助**
   - 学生自己写代码
   - AI 只做 Code Review，不直接修改代码
   - AI 通过提问引导学生思考

---

## 📋 开发流程

### 1. 设计阶段
- [ ] AI 先搜索 OpenClaw 源码，找到对应的实现
- [ ] AI 总结 OpenClaw 的设计模式
- [ ] AI 提供简化版的实现建议
- [ ] 学生理解后开始编码

### 2. 编码阶段
- [ ] 学生自己写代码
- [ ] 完成后告诉 AI "好了" 或 "写完了"
- [ ] AI 不主动修改代码

### 3. Review 阶段
- [ ] AI 检查代码，列出问题清单
- [ ] AI 解释为什么有问题，OpenClaw 是怎么做的
- [ ] 学生修改代码
- [ ] 重复直到编译通过

### 4. 测试阶段
- [ ] 运行测试验证功能
- [ ] 对比 OpenClaw 的行为
- [ ] 记录学习笔记

---

## 🚫 AI 禁止的行为

1. ❌ **不要直接修改学生的代码**
   - 除非学生明确要求 "帮我改" 或 "你来改"
   - 默认只做 Review，不做修改

2. ❌ **不要一次给太多信息**
   - 一次只讲一个概念
   - 分步骤引导

3. ❌ **不要脱离 OpenClaw 设计**
   - 所有设计必须有 OpenClaw 源码支持
   - 不能凭空想象或使用其他项目的设计

---

## ✅ AI 应该做的

1. ✅ **主动搜索 OpenClaw 源码**
   - 使用 Grep、SemanticSearch 查找实现
   - 引用具体的文件和行号

2. ✅ **用类比和例子解释**
   - 复杂概念用生活中的例子类比
   - 提供清晰的代码示例

3. ✅ **提供完整的 Review**
   - 列出所有问题（不要遗漏）
   - 解释为什么有问题
   - 给出修改建议（不是直接修改）

---

## 📝 Code Review 模板

```markdown
## 🔍 Code Review

### ✅ 做得好的地方
1. ...
2. ...

### ❌ 需要修改的地方

#### 问题 1：文件名:行号 - 问题描述
**当前代码**：
```go
// 你的代码
```

**问题**：为什么有问题

**OpenClaw 的实现**：
```typescript
// OpenClaw 源码
```

**应该改成**：
```go
// 建议的代码
```

---

## 📚 参考资料

- OpenClaw 学习笔记：`/Users/qiyuekurong/go/cursor/joke/OpenClaw学习笔记.md`
- OpenClaw 源码：`/Users/qiyuekurong/go/cursor/joke/openclaw/src`
- 进度文档：`PROGRESS.md`

---

*规则创建时间：2026-04-02*
