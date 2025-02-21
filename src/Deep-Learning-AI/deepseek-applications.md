# DeepSeek 大模型应用

## 模型版本

DeepSeek V3/R1 应用了 MoE、FP8、PTX二进制优化等，并联合算法、训练框架和硬件团队一起进行了工程上的优化，大幅提升了训练效率，降低了训练成本。

* DeepSeek-V3-Base (预训练 2.664M H800 GPU hours)
* DeepSeek-V3-Base + RL => DeepSeek-R1-Zero (实验了仅通过增加 RL，就大大增强了推理能力)
* DeepSeek-V3-Base + 两阶段 RL + 两阶段 SFT => DeepSeek-R1
* DeepSeek-V3-Base + 蒸馏 R1(后训练 0.1M GPU hours) => DeepSeek-V3

### DeepSeek-V3 (671B)

* 对标 gpt-4o，推理能力一般，无蒸馏小模型。
* 原生 FP8 训练，在支持 FP8 的 GPU 上，约需 700G 显存。不支持 FP8 的（如 A100），则需要转为 BF16，约需 1.4T 显存。推理设备建议 two H20*8 nodes,two H200*8 nodes,four A100*8 nodes (单卡都为 80G 显存)

### DeepSeek-R1 (671B)

* 对标 gpt-o1，推理能力强，有蒸馏小模型。
* 推理配置与 V3 相同

### DeepSeek-R1 蒸馏版

基于 qwen/llama 开源小参数模型，使用 DeepSeek-R1 蒸馏得到，适合垂类场景应用部署，服务器端可考虑 32B/70B 版本。

### DeepSeek 量化版

进一步降低显存占用，DeepSeek-V3/R1 (671B) 量化版大小为 404G，1个 A100*8 节点可运行（单卡80G显存）。ollama 中可下载，且可运行单节点运行。

### 本地私有部署 DeepSeek-R1 (671B) 量化版

* [单机8卡A800(80G显存)运行 Deepseek-R1 671B满血量化版](https://www.cnblogs.com/zhayujie/p/18719199)
* [显存+大内存低成本本地部署：4090单卡24G显存运行 Deepseek-R1 671B满血量化版](https://deepseek.csdn.net/67b6ab573c9cd21f4cb8d9ce.html)

## 应用场景

辅助导诊、医疗、病例生成等基于医疗系统内结构化信息，结合业务更容易落地。
辅助病例生成对原有系统依赖性强，非独立性功能，可先不做。
医疗知识库搭建，需要处理各种格式文件，扫描图片需要 OCR 识别，目前复杂格式（如表格）的 OCR 识别能力不足，效果不佳，也不需要重复建设。

**建议先做辅助导诊和诊疗**

### 辅助导诊（*）
* 使用者患者及家属
* 自助导诊/咨询/信息查询(公众号)，类似电商客服机器人，支持语音对讲，可转人工
* 医院内，结构化、半结构化、格式良好的诊疗文档等数据，可建索引进行 RAG 语义或关键字查询

### 辅助诊疗（*）
* 使用者医生
* 诊疗问答，可自动查询患者基本信息、实时监控信息、检查报告等自动诊断和给出建议
* 医疗系统内，结构化、半结构化、格式良好的诊疗文档等数据，可建索引进行 RAG 语义或关键字查询
* 支持直接输入LLM或者先检索相关数据再输入LLM

### 辅助病例生成
* 使用者医生
* 问诊对话录音
* 多人语音分类与语音识别 (*)
* 查询患者相关诊疗信息，辅助生成病例文本

### 通用医疗知识库(RAG)
全国医疗系统存在一个公共平台即可，因此意义不大

* 医疗知识问答
* 文献检索/翻译

### 辅助办公
* 略

## Citation

* [DeepSeek-V3](https://github.com/deepseek-ai/DeepSeek-V3)
* [DeepSeek-R1](https://github.com/deepseek-ai/DeepSeek-R1)
* [deepseek-v3 量化版](https://ollama.com/library/deepseek-v3)
* [deepseek-r1 量化版](https://ollama.com/library/deepseek-r1)