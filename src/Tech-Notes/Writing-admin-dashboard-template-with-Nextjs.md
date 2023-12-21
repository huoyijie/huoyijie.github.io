<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/video.js@8.0.4/dist/video-js.min.css">
<script src="https://cdn.jsdelivr.net/npm/video.js@8.0.4/dist/video.min.js"></script>
<script>
    window.HELP_IMPROVE_VIDEOJS = false
</script>

# 基于 OpenLDAP、React、MUI组件库、Next.js、Serverless 等技术实现一个管理后台模板（二）

[上篇文章](https://huoyijie.cn/docsifys/Tech-Notes/How-to-AuthN-user-with-OpenLDAP)主要介绍了 OpenLDAP 的安装、部署，这篇文章会讲一下，如何搭建基于 Next.js 搭建项目、前后端代码实现和如何把项目免费部署到 Vercel 云。

## 在线体验 Demo

[代码地址](https://github.com/huoyijie/tech-notes-code/tree/master/user-auth-with-openldap)

[在线体验 Demo](https://ldap-auth.vercel.app/)
用户名: huoyijie
密码:123456

![Sign in](https://cdn.huoyijie.cn/uploads/2023/12/signin.png)

![Dashboard](https://cdn.huoyijie.cn/uploads/2023/12/dashboard.png)

<br><video id="video-1" class="video-js" controls muted preload="auto" width="1080" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/uploads/2023/12/react-admin.webm" type="video/webm">
</video><br>

## 部署图

![部署图](https://cdn.huoyijie.cn/uploads/2023/12/deployment.png)

图中红色部分属于应用状态，主要包括 LDAP 服务和数据库服务。可以看到数据库是虚线，因为这个项目仅仅是 Demo，暂时不涉及数据库操作。实际中，需要自己购买数据库服务或者自己搭建数据库服务器。LDAP 服务搭建在了我自己的服务器上，代码免费跑在 Vercel 云上，是无状态的。

## 接上篇已建 Next.js 项目

