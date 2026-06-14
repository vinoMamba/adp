---
name: adp-ingest
description: Use when ingesting NEW raw materials into an EXISTING client workspace. Invoked as /adp-ingest <客户名称> followed by the new materials (拜访纪要/方案报价/CRM记录/公开调研/系统资料, pasted text, or file paths). Saves originals under 原始资料/, registers each via `adp source`, extracts into the knowledge pages, and logs the update. Also use when the user says "摄入资料", "补充拜访纪要", "更新客户信息", or hands over new client documents.
---

# adp-ingest — 摄入新资料，更新知识页

## 必读

开始前先读：

1. `references/核心原则.md`
2. `references/资料提取指南.md`（提取口径，含机会质量五维评估提取）
3. `references/事实溯源与Grounding.md`（高风险事实回源规则）

## 触发

```
/adp-ingest <客户名称>
<附上的新资料>
```

## 前置检查

- 读 `AGENTS.md`、`客户知识库/索引.md`、`客户知识库/来源登记.md`，了解已有资料和当前判断。
- 标记工作区进入更新中：

  ```bash
  adp status <客户名称> --state updating
  ```

## 工作流

### 1. 归档原始资料

把每份新资料原文**不改写**地放进 `原始资料/` 对应子目录：

| 资料类型 | 目标目录 |
|---|---|
| 行业报告 / 公开新闻 / 招标信息 | `原始资料/公开调研/` |
| 拜访纪要 / 会议记录 | `原始资料/拜访纪要/` |
| 方案 / 报价 / 商务文件 | `原始资料/方案报价/` |
| CRM 导出 / 商机记录 | `原始资料/CRM记录/` |
| 系统截图 / 账号导出 / 产品配置 | `原始资料/系统资料/` |

### 2. 登记来源（每份资料一条）

```bash
adp source <客户名称> --origin "<资料标识或文件名>" --type <类型> --authority <权威层级> [--date <日期>] [--page <影响页面>] [--key "<关键字段/数字>"] [--note "<备注>"]
```

类型：公开调研 / 拜访纪要 / 方案报价 / CRM记录 / 系统资料。
权威层级示例：一手（客户确认/系统导出）、二手（纪要/转述）、三手（公开报道/推测）。

### 3. 提取并更新知识页

按 `references/资料提取指南.md` 的口径，更新：

- `客户知识库/客户画像.md`
- `客户知识库/现状.md`
- `客户知识库/人物与决策链.md`
- `客户知识库/机会与动机.md`（出现明确 PPL/报价/项目时，必须维护「机会质量五维评估」）
- `客户知识库/行动计划.md`

高风险事实（预算、成本、账号数、金额、日期、角色、采购流程、阶段）必须回溯到本轮原始资料或工具结果；找不到原始来源时在该字段写"待回溯原始来源"或"待确认"，不要硬编。

### 4. 记录更新

```bash
adp log <客户名称> --action "摄入<资料概述>" --judgement "<一句话判断变化，如新增决策人 / 修正预算>"
```

## 注意

- 只改 `客户知识库/` 下的知识页和 `原始资料/`；**不要**直接编辑 `metadata.json` / `更新日志.md` / `来源登记.md`（用上面的 CLI）。
- **不要**在本 skill 里改 `输出/<客户名称>-ADP.md`——资料摄入完成后提示用户运行 `/adp-generate <客户名称>` 重新生成 ADP。
- 摄入多份资料时，逐份归档 + 登记 + 提取，最后统一写一条更新日志。
