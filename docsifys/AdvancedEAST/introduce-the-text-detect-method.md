## AdvancedEAST 文本检测原理

![AdvancedEAST 后置处理](AdvancedEAST/image/AdvancedEast.detect-method.jpg)

输出层说明：

* 1位score map, 是否在文本框内
* 2位vertex code，是否属于文本框边界像素以及是头还是尾
* 4位geo，是边界像素可以预测的2个顶点坐标。其中所有像素构成了文本框形状，然后只用边界像素去预测回归顶点坐标。边界像素定义为黄色和绿色框内部所有像素，是用所有的边界像素预测值的加权平均来预测头或尾的短边两端的两个顶点。头和尾部分边界像素分别预测2个顶点，最后得到4个顶点坐标