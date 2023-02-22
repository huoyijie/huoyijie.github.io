## AdvancedEAST 后置处理过程

运行模型预测得到输出后，还需要根据模型输出计算得出最终的文本框box坐标，这个过程就是后置处理。先看一张图:

![AdvancedEAST 后置处理](AdvancedEAST/image/AdvancedEast.post-processing.jpg)

* 由预测矩阵根据配置阈值得出激活像素集合
* 左右邻接像素集合生成region list集合
* 上下邻接region list组成region group（文本框激活区域）集合
* 遍历每个region group，生成其头和尾边界像素集合，
* 根据头和尾边界像素预测的到顶点Delta值与该边界像素坐标值计算顶点坐标，每个顶点的所有预测值的加权平均值作为最后的预测坐标值，并输出score