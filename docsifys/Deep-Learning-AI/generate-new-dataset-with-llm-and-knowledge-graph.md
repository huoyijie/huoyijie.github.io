# 基于LLM+KnowledgeGraph生成新的对话问答数据集

## LLM 模型

[通义千问 QwQ-32B](https://modelscope.cn/models/Qwen/QwQ-32B)，是最新开源的推理大模型，项目主页上报告其各项评测数据与DeepSeek-R1相当。个人电脑可以选择使用 [Ollama](https://ollama.com/)部署其Q4量化版本，需要20G显存。

```bash
# 先打开一个终端窗口运行 ollama 服务
$ ollama serve

# 再打开一个终端窗口下载 qwq:32b 量化模型
$ ollama pull qwq
```

下载完成后，可以通过 ollama cli 与 llm 交互，或者使用 LangChain 开发框架与之交互。

## LangChain 框架

[LangChain](https://python.langchain.com/docs/introduction/) is a framework for developing applications powered by large language models (LLMs).

## 知识图谱数据集

[海洋鱼类百科知识图谱数据](http://data.openkg.cn/dataset/ocean)

首先下载3个csv文件，该知识图谱有2个实体（鱼和栖息地），1个关系（鱼分布于栖息地），共计有32070种鱼、5379个栖息地以及114863个鱼与栖息地的分布关系。

## 知识图谱导入Neo4j并可视化

```python
import pandas as pd
from py2neo import Graph, Node, Relationship
from tqdm import tqdm

graph = Graph("bolt://localhost:7687", auth=("neo4j", ""))


def clear_db():
    query = "MATCH (n) DETACH DELETE n"
    graph.run(query)


num_places = 5379


def load_places_from_csv():
    place_df = pd.read_csv("data/place.csv")
    tx = graph.begin()
    for _, row in tqdm(place_df.iterrows(), desc="load places"):
        node = Node(row[":LABEL"],
                    id=row[":ID(区域)"],
                    **row.drop([":ID(区域)", ":LABEL"]).to_dict())
        tx.create(node)
    graph.commit(tx)


def load_fishes_from_csv():
    fish_df = pd.read_csv("data/fish.csv")
    tx = graph.begin()
    for _, row in tqdm(fish_df.iterrows(), desc="load fishes"):
        node = Node(row[":LABEL"],
                    id=num_places + row[":ID(鱼)"],
                    **row.drop([":ID(鱼)", ":LABEL"]).to_dict())
        tx.create(node)
    graph.commit(tx)


def load_re_from_csv():
    re_df = pd.read_csv("data/fishplacerelationship.csv")
    tx = graph.begin()
    for row in tqdm(re_df.iloc[:].itertuples(index=False),
                    desc="load relationship"):
        source = graph.nodes.match(id=num_places + row[1]).first()
        target = graph.nodes.match(id=row[2]).first()
        if source and target:
            rel = Relationship(source, row[0], target)
            tx.create(rel)
    graph.commit(tx)


if __name__ == "__main__":
    # loading minutes
    load_places_from_csv()
    load_fishes_from_csv()

    # 11w+, loading hours
    load_re_from_csv()
```

```SQL
MATCH p=()-[r:`分布`]->() RETURN p LIMIT 25
```

<svg xmlns="http://www.w3.org/2000/svg" width="595.9154052734375" height="698.5591430664062" viewBox="-292.5971374511719 -362.384765625 595.9154052734375 698.5591430664062"><title>Neo4j Graph Visualization</title><desc>Created using Neo4j (http://www.neo4j.com/)</desc><g class="layer relationships"><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(224.1597076504678)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 95.64136438601106 0.5 L 95.64136438601106 -0.5 L 25 -0.5 Z M 117.64115076296419 0.5 L 188.28251514897525 0.5 L 188.28251514897525 3.5 L 195.28251514897525 0 L 188.28251514897525 -3.5 L 188.28251514897525 -0.5 L 117.64115076296419 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="106.64125757448762" y="3" transform="rotate(180 106.64125757448762 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(503.7499563672841)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 68.38142686230154 0.5 L 68.38142686230154 -0.5 L 25 -0.5 Z M 90.38121323925466 0.5 L 133.7626401015562 0.5 L 133.7626401015562 3.5 L 140.7626401015562 0 L 133.7626401015562 -3.5 L 133.7626401015562 -0.5 L 90.38121323925466 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="79.3813200507781" y="3" transform="rotate(180 79.3813200507781 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(400.5949681846059)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 67.25611415465393 0.5 L 67.25611415465393 -0.5 L 25 -0.5 Z M 89.25590053160705 0.5 L 131.51201468626098 0.5 L 131.51201468626098 3.5 L 138.51201468626098 0 L 131.51201468626098 -3.5 L 131.51201468626098 -0.5 L 89.25590053160705 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="78.25600734313049" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(250.03073811831166)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 95.83462420588404 0.5 L 95.83462420588404 -0.5 L 25 -0.5 Z M 117.83441058283717 0.5 L 188.6690347887212 0.5 L 188.6690347887212 3.5 L 195.6690347887212 0 L 188.6690347887212 -3.5 L 188.6690347887212 -0.5 L 117.83441058283717 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="106.8345173943606" y="3" transform="rotate(180 106.8345173943606 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(469.7091180158796)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 69.21144602563214 0.5 L 69.21144602563214 -0.5 L 25 -0.5 Z M 91.21123240258527 0.5 L 135.4226784282174 0.5 L 135.4226784282174 3.5 L 142.4226784282174 0 L 135.4226784282174 -3.5 L 135.4226784282174 -0.5 L 91.21123240258527 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="80.2113392141087" y="3" transform="rotate(180 80.2113392141087 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(435.0117280174748)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 68.19023609358187 0.5 L 68.19023609358187 -0.5 L 25 -0.5 Z M 90.190022470535 0.5 L 133.38025856411687 0.5 L 133.38025856411687 3.5 L 140.38025856411687 0 L 133.38025856411687 -3.5 L 133.38025856411687 -0.5 L 90.190022470535 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="79.19012928205844" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(366.54325483895605)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 69.7716831023919 0.5 L 69.7716831023919 -0.5 L 25 -0.5 Z M 91.77146947934503 0.5 L 136.54315258173693 0.5 L 136.54315258173693 3.5 L 143.54315258173693 0 L 136.54315258173693 -3.5 L 136.54315258173693 -0.5 L 91.77146947934503 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="80.77157629086847" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(332.1555825212872)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 67.54223805951067 0.5 L 67.54223805951067 -0.5 L 25 -0.5 Z M 89.54202443646379 0.5 L 132.08426249597446 0.5 L 132.08426249597446 3.5 L 139.08426249597446 0 L 132.08426249597446 -3.5 L 132.08426249597446 -0.5 L 89.54202443646379 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="78.54213124798723" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(262.1505970929085)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 49.89778209086617 0.5 L 49.89778209086617 -0.5 L 25 -0.5 Z M 71.8975684678193 0.5 L 96.79535055868547 0.5 L 96.79535055868547 3.5 L 103.79535055868547 0 L 96.79535055868547 -3.5 L 96.79535055868547 -0.5 L 71.8975684678193 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="60.897675279342735" y="3" transform="rotate(180 60.897675279342735 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(297.7068886985978)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 70.284314254546 0.5 L 70.284314254546 -0.5 L 25 -0.5 Z M 92.28410063149913 0.5 L 137.56841488604513 0.5 L 137.56841488604513 3.5 L 144.56841488604513 0 L 137.56841488604513 -3.5 L 137.56841488604513 -0.5 L 92.28410063149913 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="81.28420744302257" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(213.64187106008035)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 48.84898244277567 0.5 L 48.84898244277567 -0.5 L 25 -0.5 Z M 70.8487688197288 0.5 L 94.69775126250447 0.5 L 94.69775126250447 3.5 L 101.69775126250447 0 L 94.69775126250447 -3.5 L 94.69775126250447 -0.5 L 70.8487688197288 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="59.848875631252234" y="3" transform="rotate(180 59.848875631252234 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(106.87291985664535 147.42047170106218) rotate(537.8365763557334)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 69.71618832563841 0.5 L 69.71618832563841 -0.5 L 25 -0.5 Z M 91.71597470259154 0.5 L 136.43216302822995 0.5 L 136.43216302822995 3.5 L 143.43216302822995 0 L 136.43216302822995 -3.5 L 136.43216302822995 -0.5 L 91.71597470259154 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="80.71608151411498" y="3" transform="rotate(180 80.71608151411498 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(433.3292719386424)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 77.5679513406258 0.5 L 77.5679513406258 -0.5 L 25 -0.5 Z M 99.56773771757892 0.5 L 152.13568905820472 0.5 L 152.13568905820472 3.5 L 159.13568905820472 0 L 152.13568905820472 -3.5 L 152.13568905820472 -0.5 L 99.56773771757892 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="88.56784452910236" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(326.22773810557214)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 64.96589654340687 0.5 L 64.96589654340687 -0.5 L 25 -0.5 Z M 86.96568292036 0.5 L 126.93157946376687 0.5 L 126.93157946376687 3.5 L 133.93157946376687 0 L 126.93157946376687 -3.5 L 126.93157946376687 -0.5 L 86.96568292036 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="75.96578973188343" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(363.55651264338)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 60.39752469816352 0.5 L 60.39752469816352 -0.5 L 25 -0.5 Z M 82.39731107511665 0.5 L 117.79483577328017 0.5 L 117.79483577328017 3.5 L 124.79483577328017 0 L 117.79483577328017 -3.5 L 117.79483577328017 -0.5 L 82.39731107511665 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="71.39741788664008" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(289.96698882859596)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 64.76564053184097 0.5 L 64.76564053184097 -0.5 L 25 -0.5 Z M 86.7654269087941 0.5 L 126.53106744063507 0.5 L 126.53106744063507 3.5 L 133.53106744063507 0 L 126.53106744063507 -3.5 L 126.53106744063507 -0.5 L 86.7654269087941 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="75.76553372031753" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(253.44003646441303)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 64.23942804329977 0.5 L 64.23942804329977 -0.5 L 25 -0.5 Z M 86.23921442025289 0.5 L 125.47864246355266 0.5 L 125.47864246355266 3.5 L 132.47864246355266 0 L 125.47864246355266 -3.5 L 125.47864246355266 -0.5 L 86.23921442025289 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="75.23932123177633" y="3" transform="rotate(180 75.23932123177633 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(180.52758144419823)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 65.31149134547431 0.5 L 65.31149134547431 -0.5 L 25 -0.5 Z M 87.31127772242743 0.5 L 127.62276906790174 0.5 L 127.62276906790174 3.5 L 134.62276906790174 0 L 127.62276906790174 -3.5 L 127.62276906790174 -0.5 L 87.31127772242743 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="76.31138453395087" y="3" transform="rotate(180 76.31138453395087 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(216.85133220716693)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 64.49505636688399 0.5 L 64.49505636688399 -0.5 L 25 -0.5 Z M 86.49484274383711 0.5 L 125.9898991107211 0.5 L 125.9898991107211 3.5 L 132.9898991107211 0 L 125.9898991107211 -3.5 L 125.9898991107211 -0.5 L 86.49484274383711 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="75.49494955536055" y="3" transform="rotate(180 75.49494955536055 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(402.1071043273036)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 76.81522936348954 0.5 L 76.81522936348954 -0.5 L 25 -0.5 Z M 98.81501574044266 0.5 L 150.6302451039322 0.5 L 150.6302451039322 3.5 L 157.6302451039322 0 L 150.6302451039322 -3.5 L 150.6302451039322 -0.5 L 98.81501574044266 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="87.8151225519661" y="3" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(504.63594614593023)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 65.81778854259045 0.5 L 65.81778854259045 -0.5 L 25 -0.5 Z M 87.81757491954357 0.5 L 128.63536346213402 0.5 L 128.63536346213402 3.5 L 135.63536346213402 0 L 128.63536346213402 -3.5 L 128.63536346213402 -0.5 L 87.81757491954357 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="76.81768173106701" y="3" transform="rotate(180 76.81768173106701 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g><g class="relationship" transform="translate(-103.98112566748188 -182.43803095221384) rotate(467.4256507437692)"><path class="b-outline" fill="#A5ABB6" stroke="none" d="M 25 0.5 L 60.23697366618573 0.5 L 60.23697366618573 -0.5 L 25 -0.5 Z M 82.23676004313886 0.5 L 117.47373370932459 0.5 L 117.47373370932459 3.5 L 124.47373370932459 0 L 117.47373370932459 -3.5 L 117.47373370932459 -0.5 L 82.23676004313886 -0.5 Z"/><text text-anchor="middle" pointer-events="none" font-size="8px" fill="#000000" x="71.2368668546623" y="3" transform="rotate(180 71.2368668546623 0)" font-family="Helvetica Neue, Helvetica, Arial, sans-serif">分布</text></g></g><g class="layer nodes"><g class="node" aria-label="graph-node5379" transform="translate(106.87291985664535,147.42047170106218)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#FFC454" stroke="#d7a013" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#2A2C34" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 尼罗口孵…</text></g><g class="node" aria-label="graph-node6" transform="translate(-51.1579108538433,-6.041715906541624)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 台湾</text></g><g class="node" aria-label="graph-node242" transform="translate(-26.805392198979174,245.43762193201474)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 冈比亚</text></g><g class="node" aria-label="graph-node158" transform="translate(231.03224559326043,253.81897159363183)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 尼罗河</text></g><g class="node" aria-label="graph-node2" transform="translate(31.510921153980888,-59.98105207859922)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 非洲</text></g><g class="node" aria-label="graph-node2588" transform="translate(50.410445408573146,305.0350085664132)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 下方艾伯…</text></g><g class="node" aria-label="graph-node2589" transform="translate(149.64378094110214,307.17429282853396)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 尼罗口</text></g><g class="node" aria-label="graph-node65" transform="translate(274.31820104879534,166.62651472591452)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 塞内加尔</text></g><g class="node" aria-label="graph-node400" transform="translate(251.95936728824793,70.78126521688287)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 三角州</text></g><g class="node" aria-label="graph-node373" transform="translate(89.2833661800136,19.831875318032832)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 以色列</text></g><g class="node" aria-label="graph-node507" transform="translate(185.71349892987072,-2.7048440133916416)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 贝努埃河</text></g><g class="node" aria-label="graph-node174" transform="translate(1.3949379232524788,77.2299060797887)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 佛塔河</text></g><g class="node" aria-label="graph-node496" transform="translate(-61.43918787925173,153.7787680033751)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 查德湖</text></g><g class="node" aria-label="graph-node5380" transform="translate(-103.98112566748188,-182.43803095221384)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#FFC454" stroke="#d7a013" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#2A2C34" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 莫三比克…</text></g><g class="node" aria-label="graph-node2590" transform="translate(28.131335650904415,-270.7870236756486)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 西河三…</text></g><g class="node" aria-label="graph-node508" transform="translate(45.525220052473806,-173.14580765479616)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 西河下游</text></g><g class="node" aria-label="graph-node249" transform="translate(-49.84614613413349,-331.43971998556384)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 非洲南部</text></g><g class="node" aria-label="graph-node1909" transform="translate(-148.86547613108849,-333.38476962391746)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 南非德兰…</text></g><g class="node" aria-label="graph-node803" transform="translate(-263.5971277453604,-183.90782195074715)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 奥歌亚湾</text></g><g class="node" aria-label="graph-node64" transform="translate(-230.40375421396436,-277.19101091809114)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 莫桑比克</text></g><g class="node" aria-label="graph-node419" transform="translate(-234.97782911666988,-89.46715551233537)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 凌波波河</text></g><g class="node" aria-label="graph-node1910" transform="translate(-148.74372068141062,-39.82419166364732)"><circle class="b-outline" cx="0" cy="0" r="25" fill="#4C8EDA" stroke="#2870c2" stroke-width="2px"/><text class="caption" text-anchor="middle" pointer-events="none" x="0" y="5" font-size="10px" fill="#FFFFFF" font-family="Helvetica Neue, Helvetica, Arial, sans-serif"> 雪利河…</text></g></g></svg>

## 导出为Parquet文件方便程序处理

```python
from py2neo import Graph

graph = Graph("bolt://localhost:7687", auth=("neo4j", ""))


def export_places():
    nodes_query = "MATCH (n:`区域`) RETURN n.id as id, n.name as name"
    nodes_df = graph.run(nodes_query).to_data_frame()
    nodes_df.to_parquet("data/place.parquet")


def export_fishes():
    nodes_query = ("MATCH (n:`鱼`) RETURN n.id as id, n.name as name,"
                   " n.bio_class_cn as bio_class, n.bio_family as bio_family,"
                   " n.bio_introduction as bio_introduction")
    nodes_df = graph.run(nodes_query).to_data_frame()
    nodes_df.to_parquet("data/fish.parquet")


def export_re():
    nodes_query = "MATCH (f)-[r:`分布`]->(p) RETURN f.id as src, p.id as target"
    nodes_df = graph.run(nodes_query).to_data_frame()
    nodes_df.to_parquet("data/fishplacerelationship.parquet")


if __name__ == "__main__":
    export_places()
    export_fishes()
    export_re()
```

## 基于知识图谱生成对话数据集

* 目前随机选择了100种鱼并生成了对话问答。
* 数据集主要生成思路是，以每种鱼的属性信息、分布地信息为上下文，以鱼为中心使用LLM生成3-5个问题，然后再以相同的上下文+问题输入到LLM，得到每个问题对应的回答。
* 为了使得生成的问题具有多样性，还尝试使用LLM随机生成不同的Role（具有不同的人格），以不同的视角提出问题。
* 鱼类属性和栖息地上下文构建比较简单，就是把实体和关系中的属性拼接起来。
* Garbage in/Garbage out，尽可能提高从知识图谱中提取的上下文质量可以改进生成质量
* 不断有针对性的改善Prompts，也可以提升生成质量

### Prompts

```python
role_prompt = PromptTemplate(
    input_variables=["example"],
    template=("你是一位在鱼类展览馆工作超过三十年的接待员，遇到过各种各样的不同年龄、性别、身份、职业、国籍、性格等的参观者。\n\n"
              "请随机提供一个参观者的一句话身份特征描述，可能是各种年龄、性别、身份、职业、国籍、性格等，"
              "要具有多样性，不要只生成和生物学有关的参观者，要符合常识和参观行为。\n\n"
              "示例：{example}\n\n不要包含任何不相关信息，不要输出markdown格式。"),
)

visitor_prompt = PromptTemplate(
    input_variables=["role", "context"],
    template=("{role}\n\n{context}\n\n"
              "请严格根据上下文提出3至5个以鱼为中心的关于鱼类特征及其分布地的问题，问题之间换行，不要编造任何信息，"
              "上下文中要包含问题的答案，不要包含答案，用中文提问，不要包含任何不相关信息，不要输出markdown格式。"),
)

expert_prompt = PromptTemplate(
    input_variables=["context", "question"],
    template=("你是一个在鱼类展览馆工作的生物学家，了解所有海洋和河流鱼类知识，"
              "对每位不同年龄、性别、国籍、性格的参观者的各种问题都能给出很好的回答。\n\n上下文：{context}\n\n"
              "问题：{question}\n\n请严格根据上下文回答问题，"
              "不要编造任何信息，上下文不能回答的问题回复不知道，用中文回答，不要提到上下文，"
              "不要包含任何不相关信息，不要输出markdown格式。"),
)
```

### 生成的单条数据示例

```json
{
    "qa": [
        {
            "question": "露斯塔野鲮的背鳍分支条数量范围是多少？",
            "answer": "露斯塔野鲮的背鳍分支条数量范围是12-14.5个。"
        },
        {
            "question": "这种鱼的头部轮廓具有什么显著形态特征？",
            "answer": "这种鱼的头部轮廓具有低而显著弓形的特征。"
        },
        {
            "question": "前背鳞片的数量范围在露斯塔野鲮中如何界定？",
            "answer": "露斯塔野鲮的前背鳞片数量范围为12-16个。"
        },
        {
            "question": "露斯塔野鲮的吻部结构有何特殊性？",
            "answer": "露斯塔野鲮的吻部结构特殊性在于没有侧叶的吻。"
        },
        {
            "question": "该鱼类主要分布在哪些亚洲国家或地区？",
            "answer": "该鱼类主要分布在孟加拉、尼泊尔、印度、缅甸和巴基斯坦。"
        }
    ],
    "role": "你是一位来自西班牙巴塞罗那的42岁自由插画师，习惯安静观察鱼群游动轨迹，总在速写本上记录光影变化，偶尔因过度专注撞到展示牌而慌张道歉。",
    "context": "鱼类：\nID:5450 露斯塔野鲮(南亚野鲮) 辐鳍鱼纲 None 背鳍有 12-141/2个分支条; 头部的低轮廓显着弓形; 短背鳍有前面的分支条短于头部; 12-16个前背的鳞片 ; 没有侧叶的吻.\n分布地：\n孟加拉 尼泊尔 亚洲 印度 缅甸 巴基斯坦"
}
```

### 数据集生成程序

```python
import pandas as pd
from collections import defaultdict
import itertools
import json
from typing import List
from tqdm import tqdm
from langchain_ollama import ChatOllama
from langchain.prompts import PromptTemplate


def read_places():
    return pd.read_parquet("data/place.parquet")


def read_fishes():
    return pd.read_parquet("data/fish.parquet")


def read_re():
    return pd.read_parquet("data/fishplacerelationship.parquet")


llm = ChatOllama(
    model="qwq",
    num_ctx=8192,
    temperature=0.6,
    top_p=0.95,
    top_k=30,
)

role_prompt = PromptTemplate(
    input_variables=["example"],
    template=("你是一位在鱼类展览馆工作超过三十年的接待员，遇到过各种各样的不同年龄、性别、身份、职业、国籍、性格等的参观者。\n\n"
              "请随机提供一个参观者的一句话身份特征描述，可能是各种年龄、性别、身份、职业、国籍、性格等，"
              "要具有多样性，不要只生成和生物学有关的参观者，要符合常识和参观行为。\n\n"
              "示例：{example}\n\n不要包含任何不相关信息，不要输出markdown格式。"),
)

visitor_prompt = PromptTemplate(
    input_variables=["role", "context"],
    template=("{role}\n\n{context}\n\n"
              "请严格根据上下文提出3至5个以鱼为中心的关于鱼类特征及其分布地的问题，问题之间换行，不要编造任何信息，"
              "上下文中要包含问题的答案，不要包含答案，用中文提问，不要包含任何不相关信息，不要输出markdown格式。"),
)

expert_prompt = PromptTemplate(
    input_variables=["context", "question"],
    template=("你是一个在鱼类展览馆工作的生物学家，了解所有海洋和河流鱼类知识，"
              "对每位不同年龄、性别、国籍、性格的参观者的各种问题都能给出很好的回答。\n\n上下文：{context}\n\n"
              "问题：{question}\n\n请严格根据上下文回答问题，"
              "不要编造任何信息，上下文不能回答的问题回复不知道，用中文回答，不要提到上下文，"
              "不要包含任何不相关信息，不要输出markdown格式。"),
)


def create_context(fish: dict, fish_places: List[dict]):
    fid, fname, bio_class, bio_family, bio_introduction = (
        fish["id"], fish["name"], fish["bio_class"], fish["bio_family"],
        fish["bio_introduction"])
    return (
        f"鱼类：\nID:{fid} {fname} {bio_class} {bio_family} {bio_introduction}" +
        "\n分布地：\n" + " ".join(fish_places))


def main():
    places = {row["id"]: row for _, row in read_places().iterrows()}
    fishes = {row["id"]: row for _, row in read_fishes().iterrows()}
    re = defaultdict(list)
    for _, row in read_re().iterrows():
        re[row["src"]].append(row["target"])

    with open("data/qa.jsonl", "w") as f:
        for id in tqdm(itertools.islice(fishes, 100), desc="RAGing"):
            fish = fishes[id]
            fish_places = [places[pid] for pid in re[id]]
            context = create_context(fish, [p["name"] for p in fish_places])

            final_role_prompt = role_prompt.format(
                example="你是一个来自中国上海的五年级男生，对各类鱼类生活习性都非常好奇，特别喜欢问有趣的问题。")
            role_response = llm.invoke(final_role_prompt)
            role = role_response.content.split("</think>\n\n")[1]

            final_visitor_prompt = visitor_prompt.format(
                role=role,
                context=context,
            )
            vistor_response = llm.invoke(final_visitor_prompt)
            questions = [
                q.strip() for q in vistor_response.content.split(
                    "</think>\n\n")[1].split("\n")
            ]

            answers = []
            for question in questions:
                final_expert_prompt = expert_prompt.format(
                    context=context,
                    question=question,
                )
                expert_response = llm.invoke(final_expert_prompt)
                answer = expert_response.content.split(
                    "</think>\n\n")[1].strip()
                answers.append(answer)

            record = {
                "qa": [{
                    "question": q,
                    "answer": a
                } for q, a in zip(questions, answers)],
                "role":
                role,
                "context":
                context,
            }
            print(json.dumps(record, indent=4, ensure_ascii=False))

            json_str = json.dumps(
                record,
                separators=(",", ":"),
                ensure_ascii=False,
            )
            f.write(json_str + "\n")


if __name__ == "__main__":
    main()
```