<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/video.js@8.0.4/dist/video-js.min.css">
<script src="https://cdn.jsdelivr.net/npm/video.js@8.0.4/dist/video.min.js"></script>
<script>
    window.HELP_IMPROVE_VIDEOJS = false
</script>

# 基于 React 实现经典俄罗斯方块小游戏

## Github

[tetris 代码地址](https://github.com/huoyijie/tetris)

## 更新日志

* Round 1 (2023-10-31)

<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/11/tetris-v1.webm" type="video/webm">
</video><br>

先实现一个最简单的 I 型四格拼板，可以进行左移、右移、下移、旋转操作，包含边界碰撞检测。把7种四格拼板封装成 React 组件，通过 x、y 属性控制四格拼板的位置，移动拼板就是更新其 x、y 值。通过 rotate 属性控制四格拼板的形态，旋转 0/90/180/270 度。后面继续开发 J、L、O、S、T、Z 型四格拼板。

* Round 2 (2023-11-02)

实现了当前拼板掉落、旋转、移动时的碰撞检测函数，当前拼板显示为红色，每秒自动降落一格直到发生碰撞，然后冻结当前拼板，再随机产生一个新的拼板。