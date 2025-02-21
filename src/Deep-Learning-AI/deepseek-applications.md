# DeepSeek 大模型应用

## 模型版本

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

### 显存+大内存低成本部署 DeepSeek-R1 (671B) 量化版

[低成本本地部署：4090单卡24G显存运行Deepseek R1 671B满血版](https://deepseek.csdn.net/67b6ab573c9cd21f4cb8d9ce.html)

## Citation

[DeepSeek-V3](https://github.com/deepseek-ai/DeepSeek-V3)
[DeepSeek-R1](https://github.com/deepseek-ai/DeepSeek-R1)
[deepseek-v3 量化版](https://ollama.com/library/deepseek-v3)
[deepseek-r1 量化版](https://ollama.com/library/deepseek-r1)