---
name: adp-ask
description: Use when answering questions about a specific client. Invoked as /adp-ask <客户名称> followed by the question on the next line. Read-only — reads the knowledge base and raw materials to answer, never edits files. Also use when the user says "这个客户…", "问一下客户情况", or asks anything answerable from an existing client workspace.
---

# adp-ask — 回答关于客户的问题（只读）

## 必读

开始前先读：

1. `references/核心原则.md`
2. `references/事实溯源与Grounding.md`（回答高风险事实时的回源要求）

## 触发

```
/adp-ask <客户名称>
<问题>
```

## 工作流

1. 读 `AGENTS.md` 和 `客户知识库/索引.md` 建立整体认知。
2. 读 `客户知识库/来源登记.md`，了解资料覆盖面和权威层级。
3. 根据问题相关性，读对应的知识页：

   | 问题方向 | 优先读 |
   |---|---|
   | 客户背景 / 业务 | `客户画像.md` |
   | 现网 / 已购产品 / 对接记录 | `现状.md` |
   | 组织 / 关键人 / 决策 / 客情 | `人物与决策链.md` |
   | 机会 / PPL / 动机 / 报价 | `机会与动机.md` |
   | 行动 / 关键行为 / 风险 | `行动计划.md` |

4. 涉及高风险事实（预算、金额、账号数、日期、角色、阶段、采购流程）时，必须回溯到 `原始资料/` 对应文件或本轮最新工具结果核对；只凭 `客户知识库/*.md` 不足以定论时，明确说"知识库记录为 X，未在原始资料中核实"。
5. 回答时区分：**事实**（带来源）、**判断**（带依据）、**假设**（标"待验证"）。不把三者混着说。

## 注意

- **只读**：不修改任何文件，不调用任何写入 CLI。
- 信息不足时如实说"待确认"或"资料中未覆盖"，不要编造。
- 如果问题指向需要更新知识库或重新生成 ADP，提示用户运行 `/adp-ingest` 或 `/adp-generate`。
