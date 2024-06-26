# AI人脸面部表情识别医学应用研究

## 1.基本介绍

面部表情是人类传达情感状态和意图的最强大、自然和普遍的信号之一。由于其在社交机器人、医疗治疗、驾驶员疲劳监测以及许多其他人机交互系统中的实际重要性，自动面部表情分析已经成为许多研究的焦点。

随着面部表情识别（FER）从实验室控制条件向具有挑战性的自然环境过渡，以及深度学习技术在各个领域取得的最新成功，深度神经网络越来越多地被利用来学习用于自动识别面部表情。

目前面部表情识别最主要的困难，除了缺乏规模大、质量高的数据集外。另一个主要难点是同一个表情，可能差距比较大，而不同的表情之间可能差别又非常细微。

举个例子，同一个人随着开心程度的不同，表达出的表情是有变化的，而不同的人表达开心的方式也可能有差别，有时甚至差别比较大（其他表情也是同理），有的人笑哭了，到底是开心还是悲伤呢？

再举个例子，不同的表情也可能差距很小，开心的大笑和平静的表情区分度很大，特征明显。但是微微一笑和平静的表情就差别很小了，有时人类肉眼难辨。再比如遇到意外状况后，是惊讶多一些，还是恐惧多一些？再比如蔑视的看着其他人，是讨厌多一些还是愤怒多一些？有些情况下，人类仅通过图片也是难以分辨。

这些情况对 AI 模型来说也是比较难学习到的，想要尽可能准确识别不明显的表情特征，需要提供尽可能多的训练数据，越多越好，能够覆盖到这些细微的表情特征差别。

## 2.数据集 (Dataset)

### 2.1面部表情识别在数据集和方法方面的演变

![The evolution of facial expression recognition in terms of datasets
and methods](https://cdn.huoyijie.cn/uploads/2024/04/datasets.png)

早期的研究主要集中在实验室可控条件下的数据集（如 CK+），图片数量少，一般几百到上千张，被拍摄的对象一般只有几十到上百人，样本多样性差，一般都是正面、无遮挡、灯光良好、分辨率高的照片。

后面出现了一些通过互联网收集整理而来的较大规模的数据集（如 FER2013/RAF-DB/AffectNet），这些数据集试图解决样本多样性问题。样本基本可以覆盖不同年龄、性别、种族、头部姿势、光照条件、部分遮挡（眼镜、头发、手等）、照片分辨率等自然条件。

FER2013/RAF-DB 数据集规模依然较小，图片样本数不足3万，与一般的图片分类任务动辄几十万、上百万的规模相比很小。AffectNet 数据集有近四十万张图片构成，这几个是目前研究面部表情识别最常用的数据集。

### 2.2 现有数据集主要问题

1. 样本数量少/多样性方面不足
2. 表情间样本数量不平衡
3. 标柱错误
4. 缺乏医疗场景下的数据集

#### 2.2.1 样本数量少/多样性方面不足

导致训练数据不足，模型训练不充分，经常会出现过拟合，模型泛化性能不差。除了覆盖不同年龄、性别、种族、头部姿势、光照条件、部分遮挡（眼镜、头发、手等）、照片分辨率等自然条件，还需要结合具体场景，尽可能提供符合实际应用场景多样性的样本集合。

#### 2.2.2 表情间样本数量不平衡

模型倾向预测为样本数量更多的表情，使得样本占比低的表情识别准确率低，模型泛化性能差。

#### 2.2.3 存在标注错误

典型的是 FER2013 数据集，可能是由于没有通过标注团队协作标注的原因，导致有很多图片标注的表情是不对的。这会导致模型训练过程中出现困惑，无法学习到真实表情特征，导致模型性能不佳。所以在这个数据集上训练的模型，准确率通常很低。后面发现这个问题后，由众包数据标注团队进行了重新标注，数据集升级为 FER+，在这个数据集上训练的模型准确率大大提升。

#### 2.2.4 缺乏医疗场景下的数据集

目前缺乏医疗场景下的公开的数据集，现有的数据集图片都是生活场景下的图片，其上训练出的模型在医疗场景下估计并不能直接应用，患者卧床拍摄的医疗场景照片和生活场景照片差别大。如果自行收集医疗场景下的数据集，也可以考虑公开给学术领域研究，也是非常有价值的事情。

### 2.3 注意事项

1. 确定医院实际问题目标

一般来说，模型可以解决什么问题是由数据集定义的，所以一定要沟通好医院场景的实际问题，才能决定收集哪些以及怎样收集样本。

要确定好需要识别的表情分类，有些表情特征明显（如大笑），不但人类容易分辨，模型也相对容易学习。不明显或与其他表情区分度不高的表情，人类不容易分辨，模型的识别成功率也会相对低。比如惊讶和恐惧，有时差别很细微，人类和模型都可能误判。这也是目前通过人脸图片识别表情最主要的困难。

一旦数据集确定后，后面训练模型时，只能尽可能提升在数据集上的准确率，尽可能提高模型的泛化能力。模型的泛化能力高，意味着在它没见过的图片上也可以很好的预测其表情，但是这种泛化能力是有限制的，越接近训练数据集的样本照片，预测准确的可能性越高，与训练数据集差别越大的样本照片，预测准确的可能性会随着差距变大而降低。

如果医疗场景下要识别的表情特征是非常细微的，比如类似镇痛等级识别这种，从数据收集角度和训练模型进行识别的角度来说，都是比较有挑战的。这种和目前的人脸表情识别任务差距还是很大的，也许更类似于通过实时拍摄人脸视频，检测司机疲劳驾驶的任务。因为镇痛通常会持续一段时间，通过分析视频帧，更可能捕捉到这种细微表情特征变化，而仅通过一帧图片，可能难以学习到不同镇痛等级的表情特征。

2. 避免标注错误

为了防止标注错误，一般每张图片都会由多人同时标注。所以如果是自己收集标注数据，特别要注意这个问题。如果数据集出现标注不准确的情况，训练出的模型往往性能都不会好。

3. 数据隐私问题

如果是自己拍摄照片，可能要与患者及家属沟通好相关照片授权问题，避免纠纷。

4. 拍摄方式

摄像头如何放置，可以拍摄到人脸正面、头部姿势较正、清晰度好、亮度好的照片，以及如何在一段时间内收集到足够多的照片

## REFERENCES

> [Deep Facial Expression Recognition: A Survey](https://arxiv.org/abs/1804.08348)