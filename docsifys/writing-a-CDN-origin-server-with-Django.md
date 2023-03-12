# 基于 Django 搭建 CDN 回源服务

![CDN 回源服务](https://cdn.huoyijie.cn/uploads/2023/03/originsrv-upload-success.png)

[网站图片视频接入CDN](https://huoyijie.cn/gitbooks/huoyijie.cn/host-static-resources-of-website-with-tencent-CDN/latest/)

[快速搭建CDN回源服务器](https://huoyijie.cn/gitbooks/huoyijie.cn/build-a-CDN-back-to-source-server-with-nodejs/latest/)

我在很早之前有介绍过如何开通腾讯云 CDN 服务，以及如何搭建回源服务器，如果感兴趣可以先看一下。之前的回源服务器非常简单，就是后台启动一个 http server 服务，挂载指定目录。每次需要上传资源文件，可以执行 shell 脚本工具(或者手动 scp)，把指定文件远程上传到服务器的指定目录。然后就可以通过 cdn.huoyijie.cn 域名访问指定资源了。

但是每次执行脚本上传文件有点麻烦，也不方便资源文件管理，所以我这次想基于 Django 搭建一个新的回源服务。基于 Django Admin 的强大功能，几乎不需要写代码，就可以实现，而且支持用户及权限管理。

## 前置

* Python: v3.10+

* 安装 pip
  ```bash
  $ sudo apt install python3-pip
  ```

* 安装 python3-virtualenv
  ```bash
  $ sudo apt install python3-virtualenv
  ```

* Django: v4.1+

* 安装 gettext (可选)，用以支持 Django 多语言翻译
  ```bash
  $ sudo apt install gettext
  ```

使用 Sqlite3，不需要安装和配置数据库。

## 初始化项目

创建 Github 库 [originsrv](https://github.com/huoyijie/originsrv)，并克隆到本地。

创建项目并激活 venv
  ```bash
  $ cd originsrv/
  $ virtualenv -p python3 venv
  $ source venv/bin/activate
  ```

激活后 venv 检查 python 版本
  ```bash
  $ python -V
  Python 3.10.6
  ```

安装 Django
  ```bash
  $ python -m pip install Django
  # 验证 Django 是否安装成功
  $ python -m django --version
  4.1.7
  ```

通过 Django 脚手架快速生成项目代码
  ```bash
  $ django-admin startproject originsrv
  ```

在最开始已创建 originsrv 项目目录的情况下，直接运行脚手架命令会报错。所以这里有点 tricky，需要在其他地方运行脚手架命令，然后把生成的项目结构融合到最开始建好的 originsrv 目录。

之所以会麻烦一些，是因为我想通过搭建一个隔离的 python 运行环境 venv 来运行项目。

修改 originsrv/urls.py，增加 Django admin 相关页面路径
  ```python
  urlpatterns = [
      path('', admin.site.urls),
  ]
  ```

初始化 Django 内置数据库表
  ```bash
  $ python manage.py migrate
  Operations to perform:
    Apply all migrations: admin, auth, contenttypes, sessions
  Running migrations:
    Applying contenttypes.0001_initial... OK
    Applying auth.0001_initial... OK
    Applying admin.0001_initial... OK
    Applying admin.0002_logentry_remove_auto_add... OK
    Applying admin.0003_logentry_add_action_flag_choices... OK
    Applying contenttypes.0002_remove_content_type_name... OK
    Applying auth.0002_alter_permission_name_max_length... OK
    Applying auth.0003_alter_user_email_max_length... OK
    Applying auth.0004_alter_user_username_opts... OK
    Applying auth.0005_alter_user_last_login_null... OK
    Applying auth.0006_require_contenttypes_0002... OK
    Applying auth.0007_alter_validators_add_error_messages... OK
    Applying auth.0008_alter_user_username_max_length... OK
    Applying auth.0009_alter_user_last_name_max_length... OK
    Applying auth.0010_alter_group_name_max_length... OK
    Applying auth.0011_update_proxy_permissions... OK
    Applying auth.0012_alter_user_first_name_max_length... OK
    Applying sessions.0001_initial... OK
  ```

创建超级用户
  ```bash
  $ python manage.py createsuperuser --username=huoyijie --email=huoyijie@huoyijie.cn
  ```

编辑 originsrv/settings.py，调整显示简体中文，并调整到上海时区
  ```python
  LANGUAGE_CODE = 'zh-hans'

  TIME_ZONE = 'Asia/Shanghai'
  ```

初次运行
  ```bash
  $ cd originsrv
  $ python manage.py runserver
  Django version 4.1.7, using settings 'originsrv.settings'
  Starting development server at http://127.0.0.1:8000/
  Quit the server with CONTROL-C.
  ```

访问上方的 URL，应该会打开 Django 管理后台

## 增加 CDN 资源管理模块

创建 webapp
  ```bash
  $ python manage.py startapp webapp
  ```

增加模型 webapp/models.py
  ```python
  class Resource(models.Model):
    file = models.FileField(verbose_name=_('File'), upload_to='uploads/%Y/%m/')

    def __str__(self):
        return self.file.name

    class Meta:
        verbose_name = _('Resource')
        verbose_name_plural = _('Resources')
  ```

上面这段是最主要的代码了，Django admin 会根据上面的 Resource 模型，自动生成增删改查管理页面。

配置 originsrv/settings.py，增加新 app
  ```python
  INSTALLED_APPS = [
      'webapp.apps.WebappConfig',
      # ...
  ]
  ```

提取 message 文件并翻译
  ```bash
  $ cd originsrv/webapp
  $ mkdir locale
  $ django-admin makemessages -l zh_Hans
  ```

编辑 django.po 文件
  ```
  msgid "Origin server admin"
  msgstr "cdn.huoyijie.cn"

  msgid "Origin server administration"
  msgstr "cdn.huoyijie.cn"

  msgid "Dir"
  msgstr "目录"

  msgid "Name"
  msgstr "名称"

  msgid "File"
  msgstr "文件"

  msgid "Resource"
  msgstr "资源"

  msgid "Resources"
  msgstr "资源"
  ```

编译 messages 文件
  ```bash
  $ django-admin compilemessages
  ```

因为新增加了 Resource 模型，需要运行模型迁移
  ```bash
  $ python manage.py makemigrations webapp
  $ python manage.py migrate
  ```

其他主要是 settings.py 的配置调整，注意开发环境和显示部署需要不同的值，尤其是 SECRET_KEY 和 Debug。

## 线上部署

克隆代码库
  ```bash
  $ cd ~pywork
  $ git clone git@github.com:huoyijie/originsrv.git
  ```
检查前置条件，确保全部依赖已安装好，创建并 venv 隔离环境。

安装 gunicorn web 服务器
  ```bash
  $ python -m pip install gunicorn
  ```

迁移数据库
  ```bash
  $ python manage.py migrate
  ```

新建 super user
  ```bash
  $ python manage.py createsuperuser --username=huoyijie --email=huoyijie@huoyijie.cn
  ```

编译 messages 文件
  ```bash
  $ cd originsrc/webapp
  $ django-admin compilemessages
  ```

把静态文件收集到 originsrv/public/static 目录下
  ```bash
  $ python manage.py collectstatic
  ```

程序会把上传的资源文件放到 originsrv/public/media 目录下

配置 nginx，可分别读取 originsrv/public/static 和 originsrv/public/media 目录。前者是访问当前项目所需的静态文件。后者是访问用户上传的资源文件。

新建 systemd 服务 originsrv.service，编辑配置文件 /etc/systemd/system/originsrv.service
  ```conf
  [Unit]
  Description=OriginSrv

  [Service]
  User=ubuntu
  Group=ubuntu
  Type=idle
  Environment="SECRET_KEY=******"
  Environment="STATIC_URL=https://cdn.huoyijie.cn/static/"
  Environment="MEDIA_URL=https://cdn.huoyijie.cn/"
  Environment="CSRF_TRUSTED_ORIGINS=https://huoyijie.cn"
  ExecStart=/home/ubuntu/pywork/originsrv/venv/bin/gunicorn --bind 127.0.0.1:4001 originsrv.wsgi --chdir /home/ubuntu/pywork/originsrv
  Restart=always
  KillMode=process

  [Install]
  WantedBy=multi-user.target
  ```

启动 originsrv 服务
  ```
  $ sudo systemctl daemon-reload
  $ sudo systemctl enable --now originsrv
  ```

打开浏览器已经可以访问了。通过资源管理，可上传文件到 originsrc/public/media/uploads 目录下，而这个目录是会被 nginx 挂载的。简单来说，访问 https://cdn.huoyijie.cn/uploads/2023/03/originsrv-select-file.png 资源时，回源请求会发到 nginx 进程，nginx 会读取 originsrc/public/media/uploads 目录下的相应文件并返回。

## 运行

![CDN 回源服务](https://cdn.huoyijie.cn/uploads/2023/03/originsrv-select-file.png)

![CDN 回源服务](https://cdn.huoyijie.cn/uploads/2023/03/originsrv-upload-success.png)

## 引用
> [Django Docs](https://docs.djangoproject.com/zh-hans/4.1/)