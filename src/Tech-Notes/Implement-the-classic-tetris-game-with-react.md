<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/video.js@8.0.4/dist/video-js.min.css">
<script src="https://cdn.jsdelivr.net/npm/video.js@8.0.4/dist/video.min.js"></script>
<script>
    window.HELP_IMPROVE_VIDEOJS = false
</script>

# 基于 React 实现经典俄罗斯方块小游戏

## Github

[tetris 代码地址](https://github.com/huoyijie/tetris)

## 实现思路

俄罗斯方块游戏最重要的部分就是四格拼板 (Tetromino)，共有 7 种基本形状如下：

### 七种四格拼板 (Tetromino)

```bash
# 坐标轴方向
#     y
#     |
# x ——|——>
#     v

# 方块 (Cube)
# []

# 7 种四格拼板
# I - [[0, 0], [0, 1], [0, 2], [0, 3]]
#
# [1][2][3][4]
#     *

[][][][]

# J - [[0, 0], [0, 1], [1, 1], [2, 1]]
#
# [1]
# [2][3][4]
#     *

[]
[][][]

# L - [[0, 1], [1, 1], [2, 1], [2, 0]]
#
#       [4]
# [1][2][3]
#     *

    []
[][][]

# O - [[0, 0], [0, 1], [1, 1], [1, 0]]
#
# [1][4]
# [2][3]

[][]
[][]

# S - [[0, 1], [1, 1], [1, 0], [2, 0]]
#
#    [3][4]
# [1][2]
#     *

  [][]
[][]

# T - [[0, 0], [1, 0], [1, 1], [2, 0]]
#
#     *
# [1][2][4]
#    [3]

[][][]
  []

# Z - [[0, 0], [1, 0], [1, 1], [2, 1]]
#
# [1][2]
#    [3][4]
#     *

[][]
  [][]
```

那么怎样用数字形式表达四格拼板形状呢？我采用的方式是矩阵，也就是二维数组。在上面每种拼板上方的注释中给出了相应的二维数组，注释里 [] 里面的数字和二维数组的索引(+1)是对应的，有 * 标记的是整体形状的中心块，当形状按照顺时针旋转时，就是围绕中心块进行。O 型是上下左右完全对称图形，所以不需要响应旋转操作，所以没有 * 标记。

### 底板 (Board) 坐标系

```bash
     0 1 2 3 4 5 6 7 8 9      ...      18 19
   +----------------------------------------+
 0 | . . . . . . . . . . . . . . . . . . . .|
 1 | . . . . . . . . . . . . . . . . . . . .|
 2 | . . . . . . . . . . . . . . . . . . . .|
 3 | . . . . . . . O O O . . . . . . . . . .|
 4 | . . . . . . . . O . . . . . . . . . . .|
 5 | . . . . . . . . . . . . . . . . . . . .|
 6 | . . . . . . . . . . . . . . . . . . . .|
 7 | . . . . . . . . . . . . . . . . . . . .|
 8 | . . . . . . . . . . . . . . . . . . . .|
 9 | . . . . . . . . . . . . . . . . . . . .|
10 | . . . . . . . . . . . . . . . . . . . .|
11 | . . . . . . . . . . . . . . . . . . . .|
12 | . . . . . . . . . . . . . . . . . . . .|
13 | . . . . . . . . . . . . . . . . . . . .|
14 | . . . . . . . . . . . . . . . . . . . .|
15 | . . . . . . . . . . . . . . . . . . . .|
16 | . . . . . . . . . . . . . . . . . . . .|
17 | . . . . . . . . . . . . . . . . . . . .|
18 | . . . . . . . . . . . . . . . . . . . .|
19 | . . . . . . . . . . . . . . . . . . . .|
20 | . . . . . . . . . . . . . . . . . . . .|
21 | . . . . . . . . . . . . . . . . . . . .|
22 | . . . . . . . . . . . . . . . . . . . .|
23 | . . . . . . . . . . . . . . . . . . . .|
24 | . . . . . . . . . . . . . . . . . . . .|
   +----------------------------------------+
```

每个四格拼板作为一个整体会有一个坐标 [x, y]，表示自己在上方 Board 中的位置，这个坐标与四格拼板内每个方块 (Cube) 的坐标加在一起可以确定每个方块 (Cube) 在 Board 中的坐标。

如上方的 T 型拼板的坐标是 [7, 3]，而 T 型拼板内 4 个方块的坐标是 [[0, 0], [1, 0], [1, 1], [2, 1]]，T 型拼板坐标分别加上 4 个方块坐标，就可以得到每个方块在 Board 中的坐标，并绘制相应的方块，组成 T 型拼板。

### 操作四格拼板



## 更新日志

* Round 1 (2023-10-31)
<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v1.webm" type="video/webm">
</video><br>

先实现一个最简单的 I 型四格拼板，可以进行左移、右移、下移、旋转操作，包含边界碰撞检测。把7种四格拼板封装成 React 组件，通过 x、y 属性控制四格拼板的位置，移动拼板就是更新其 x、y 值。通过 rotate 属性控制四格拼板的形态，旋转 0/90/180/270 度。后面继续开发 J、L、O、S、T、Z 型四格拼板。

* Round 2 (2023-11-02)
<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v2.webm" type="video/webm">
</video><br>

实现了当前拼板掉落、旋转、移动时的碰撞检测函数，当前拼板显示为红色，每秒自动降落一格直到发生碰撞，然后冻结当前拼板，再随机产生一个新的拼板。

* Round 3 (2023-11-03 )
<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v3.webm" type="video/webm">
</video><br>

已支持 I、J、L、O、S、T、Z 型四格拼板，没有行消除早晚要玩完啊！！！接下来把行消除功能加上。

* Round 4 (2023-11-06)
<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v4.mp4" type="video/webm">
</video><br>

已支持 I、J、L、O、S、T、Z 型四格拼板，今天终于把行消除功能加上了，试玩了一下，刚开始还强装镇定，到后面一直不来 I 型长块，心好慌啊！！！

* Round 5 (2023-11-07)
<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v5.mkv" type="video/webm">
</video><br>

今天把显示当前拼板掉落位置、显示得分、消除行数、提前预览下一块拼板等功能加上了。显示当前拼板掉落位置功能是我最喜欢的，之前没有这个功能，为了不串行，眼睛都要看花了！！！

* Round 6 (2023-11-08)
<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v6.mkv" type="video/webm">
</video><br>

今天把游戏计时和难度等级功能加上了，解决了一个组件 infinite render 的 bug。今天的状态不错，采取了比较有效的消除策略，玩了 35 分钟累积了 1480 分还再继续！！！主要功能已经差不多实现好了，后面考虑加入一些有趣的小道具，比如可以用获得的炸弹炸掉不好的地方，再比如可以在慌乱时使用暂停几秒的道具等等。