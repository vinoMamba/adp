---
name: adp-generate
description: Use when generating or iterating the standard 10-section ADP output (输出/客户名称-ADP.md) from the knowledge base. Invoked as /adp-generate <客户名称>. Reads all 客户知识库/ pages, produces or updates the single 输出/客户名称-ADP.md against the 标准ADP输出模板, self-checks against the quality gates, and marks the client ready. Also use when the user says "生成ADP", "更新ADP", "刷新ADP输出", "重新生成ADP", or after /adp-ingest has added new materials.
---

# adp-generate — 基于知识库迭代 ADP 输出

## 必读

开始前先读：

1. `references/核心原则.md`
2. `references/标准ADP输出模板.md`（十段骨架 + mermaid 图表规范 + 客情等级积分规范）
3. `references/ADP方法论.md`（十段结构定义 + 判断边界）

## 触发

```
/adp-generate <客户名称>
```

## 前置检查

- 读全部知识页：`客户画像.md`、`现状.md`、`人物与决策链.md`、`机会与动机.md`、`行动计划.md`、`来源登记.md`。
- 读现有 `输出/<客户名称>-ADP.md`（如果已存在，**在原结构上迭代**，不新建文件）。
- 标记进入更新中：

  ```bash
  adp status <客户名称> --state updating
  ```

## 工作流

### 1. 生成 / 迭代 ADP

按 `references/ADP方法论.md` 的十段结构，在 `输出/<客户名称>-ADP.md` 上生成或迭代：

1. 一、客户简介
2. 二、人员关系图谱（含 mermaid `graph TD` 关系图 + 图下总结段）
3. 三、数字化现状（必须 Markdown 表格）
4. 四、核心 KP 个人
5. 五、决策链及客情关系分析（子段：决策链分析 mermaid 图 / 决策链客情分析 5 列简表 / 项目对接情况）
6. 六、五大关键行为
7. 七、近期规划（产品线）（固定 3 列表）
8. 八、已购买产品（仅钉钉产品）
9. 九、客户群
10. 十、风险&方案

每一节优先引用 `客户知识库/` 已有内容；信息不足时写"待确认"。

### 2. 自检（对照 `references/核心原则.md` 的「质量门槛」）

逐条核对，重点关注：

- 十段主结构未变。
- 高风险事实已回源；未核实的统一写"待确认"（不要写"待回溯原始来源"这类过程表述）。
- 第二节有图 + 图下总结段。
- 第三节是表格。
- 第五节「决策链客情分析」是 5 列简表、逐人 `-3`~`+3` 评分、不画图。
- 第七节是固定 3 列表。
- 第八节只有钉钉产品。
- 最终 ADP 不含任何来源链接、页码、文件名、截图位置。
- mermaid 只用基础 `graph`/`flowchart`，节点不用 emoji/图标。

### 3. 更新状态

```bash
adp status <客户名称> --state ready --model "<本轮使用的模型>"
adp log <客户名称> --action "迭代 ADP 输出" --judgement "<一句话总结本轮变化>"
```

## 注意

- 一个客户**只有一份** `输出/<客户名称>-ADP.md`，永不创建版本化副本。
- 不直接写 `metadata.json` / `更新日志.md`（用 CLI）。
- 如果发现知识页缺失关键信息，先提示用户用 `/adp-ingest` 补资料，再生成；不要凭空编造。
