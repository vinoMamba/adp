---
name: adp-review
description: Use when auditing a client's ADP against the quality gates. Invoked as /adp-review <客户名称>. Reads 输出/客户名称-ADP.md and the knowledge pages, checks every quality-gate rule, reports gaps with file:location and a fix recommendation, optionally fixes them, and logs the audit. Also use when the user says "审查ADP", "检查ADP质量", "ADP合规检查", or wants to know what's missing or wrong before sharing the ADP.
---

# adp-review — 对照质量门槛审计 ADP

## 必读

开始前先读：

1. `references/核心原则.md`（「质量门槛」是审计清单）
2. `references/ADP方法论.md`（十段结构 + 判断边界）
3. `references/标准ADP输出模板.md`（图表规范 + 客情等级积分规范）

## 触发

```
/adp-review <客户名称>
```

## 工作流

### 1. 加载审计对象

读 `输出/<客户名称>-ADP.md` 和全部 `客户知识库/*.md`、`来源登记.md`。

### 2. 逐项核对质量门槛

对每条门槛，输出一个发现：

- **PASS** / **FAIL** / **WARN**（WARN = 边缘情况或证据不足）
- 失败时给出 `文件:位置` 和具体问题
- 给出修复建议

核对清单（来自 `references/核心原则.md` 的「质量门槛」，全部覆盖）：

1. 十段主结构完整且顺序正确。
2. 高风险事实（预算、金额、账号数、日期、角色、阶段、采购流程）有原始资料支撑；未核实的写"待确认"。
3. 最终 ADP 不含来源链接、来源名称、页码、截图位置、文件名、"待回溯原始来源"。
4. 第二节有 mermaid 关系图 + 图下总结段（不是只有图）。
5. 第三节是 Markdown 表格（不是字段提示）。
6. 未接触人物没有强客情结论。
7. 客情逐人 `-3`~`+3` 评分，证据不足写"待确认"。
8. 第五节「决策链客情分析」是 5 列简表、不画图。
9. 客情等级积分有活动空间/信息质量/立场程度或行为依据。
10. 有明确 PPL/报价/项目时，`客户知识库/机会与动机.md` 维护了「机会质量五维评估」。
11. 第七节是固定 3 列表（计划时间、产品、对接部门）。
12. 第八节只有已购买的钉钉产品。
13. 行动计划含目标人物、目的、负责人、时间、验证方式。
14. mermaid 只用基础 `graph`/`flowchart`；节点不含 emoji/图标；无独立图例图。

### 3. 汇总

给出总体结论（READY / NEEDS FIX / BLOCKED）和发现列表。READY = 全部 PASS 或仅 WARN；NEEDS FIX = 有 FAIL 但可修；BLOCKED = 缺关键资料（建议先 `/adp-ingest`）。

### 4. 记录审计（可选修复）

- 如果用户同意修复 FAIL 项，直接编辑 `输出/<客户名称>-ADP.md` 或对应知识页，修完重核对。
- 无论是否修复，都记录审计：

  ```bash
  adp log <客户名称> --action "质量审计" --judgement "<总体结论 + FAIL 项数>"
  ```

- 修复后状态变化时：

  ```bash
  adp status <客户名称> --state ready   # 审计通过
  adp status <客户名称> --state updating # 仍在修
  ```

## 注意

- 审计是对照门槛找问题，不是重写 ADP；除非用户同意，只报告不改文件。
- 不直接写 `metadata.json` / `更新日志.md`（用 CLI）。
- 发现系统性资料缺失时，建议 `/adp-ingest` 而不是在审计里硬补。
