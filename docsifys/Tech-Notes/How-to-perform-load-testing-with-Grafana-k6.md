# 如何基于 Grafana k6 实现 API 负载测试

我当前正在做的一个 Web 项目要交付了，客户希望可以支持最少 200 用户同时进行频繁的交互。由于该系统功能上的特殊性，用户的核心高频操作中，数据库写操作占比非常大，主要依赖数据库事务来确保并发下的数据一致性。

我们选择了具有2核CPU、4GB内存的基本配置 ECS 云主机， 但当前系统最少需要部署几台云主机呢？要想知道答案，靠拍脑袋不是不行，但是不太科学，一旦上线后，系统运营一段时间后，性能不理想影响用户体验了，损伤的也是开发团队的技术信誉。最好的做法是能够模拟真实流量进行负载测试，弄清楚当前系统的性能情况，通过反复对比测试，找到可以支撑 200 用户同时高频访问的最少主机数量。

## Grafana k6

当下有非常多的负载测试工具，但我对 Golang 是有偏好的，所以选择了更新的基于 Golang 实现的开源负载测试工具 [Grafana k6](https://k6.io/)。它的主要特点有 3 个:

* 可以基于 JavaScript ES2015/ES6 编写自动化测试脚本
* 通过定义 Checks、Thresholds，可以编写面向目标、对自动化更友好的测试脚本
* 对开发者更友好的 API

## Load Test

为了更准确评估系统性能，我们基于一台2核4GB主机搭建了真实线上系统，并安装了 mysql 数据库（还未采买独立的 mysql 数据库实例）。同时准备了同机房内的另一台云主机作为客户端来运行自动化测试脚本。

为了后续可以更方便的反复运行自动化负载测试，我创建了一个单独的负载测试项目，存放脚本、准备数据等。后面的代码仅仅是示例，非项目代码。

### 安装 k6

[安装指南](https://k6.io/docs/get-started/installation/)

### 准备测试数据

* seed.sql

下面的 sql 制造了一个虚拟测试用户，以及为它开通了3个产品服务。

```sql
INSERT INTO user(id, username, password) VALUES
(10001, 'test10001', '******');

INSERT INTO asset(user_id, prod_id, qty) VALUES
(10001, 1, 1000);
INSERT INTO asset(user_id, prod_id, qty) VALUES
(10001, 2, 1000);
INSERT INTO asset(user_id, prod_id, qty) VALUES
(10001, 3, 1000);
```

* data.json

下面的 json 配置了自动化测试脚本可以使用的虚拟用户，当运行自动化脚本时，会按照数组索引顺序依次使用对应虚拟用户来登录并操作核心 API。

```json
{
  "users": [
    "test10001"
  ]
}
```

* genseed.mjs

上面的 seed.sql 和 data.json 示例仅造了一个测试用户，但是目标是 200 个测试用户，如果手动准备 seed.js 和 data.json 还是比较麻烦的，我们可以再编写一个脚本用来生成 seed.js 和 data.json 文件。

genseed.mjs 的内容也非常简单，是一个执行 200 次的循环，每个循环内，构造一个插入用户的 sql 语句，构造为其开通产品服务的 sql 语句。

```js
import fs from 'fs';

let sql1 = `INSERT INTO user(id, username, password) VALUES\n`;

let sql2 = `INSERT INTO asset(user_id, prod_id, qty) VALUES\n`;

// test users
let users = '';

for (let i = 1; i <= 200; i++) {
  let id = 10000 + i;
  sql1 += `(${id}, 'test${id}', '******')`;

  for (let j = 1; j <= 3; j++) {
    sql2 += `(${id}, ${j}, 1000)`;

    if (j < 6) {
      sql2 += ',\n';
    }
  }

  users += `"test${id}"`;

  if (i < 200) {
    sql1 += ',\n';
    sql2 += ',\n';
    users += ',\n    ';
  } else {
    sql1 += ';\n\n';
    sql2 += ';\n\n';
  }
}

fs.writeFileSync('./seed.sql', sql1 + sql2);

// data.json
let data = `{
  "users": [
    ${users}
  ]
}`;

fs.writeFileSync('./data.json', data);
```

### 编写自动化测试脚本

```js
// script.js

import http from 'k6/http';
import { Trend, Rate } from 'k6/metrics';
import { check, group, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import { vu } from 'k6/execution';

// 请求 host，根据实际情况进行修改
const host = 'http://localhost:8080';
// 加载测试用户列表
const users = new SharedArray('users', function () {
  return JSON.parse(open('./data.json')).users;
});

/** 重点测试的接口列表 */
// 接口 1
export const APIOneTrendRTT = new Trend('api_one_rtt');
export const APIOneRateContentOK = new Rate('api_one_ok');
// 接口 2
export const APITwoTrendRTT = new Trend('api_two_rtt');
export const APITwoRateContentOK = new Rate('api_two_ok');

// 自定义选项配置
export const options = {
  // 定义所测试接口的目标阈值
  thresholds: {
    'api_one_ok': ['rate>0.999'],
    'api_one_rtt': ['p(99)<300', 'p(90)<250', 'avg<200', 'med<150'],
    'api_two_ok': ['rate>0.999'],
    'api_two_rtt': ['p(99)<300', 'p(90)<250', 'avg<200', 'med<150'],
  },

  // 定义测试完成后 summary 显示的字段
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)', 'count'],

  // 定义测试场景，配置了200个虚拟用户，针对当前自动化测试脚本，确保对每个虚拟测试用户跑且仅跑一次
  scenarios: {
    contacts: {
      vus: 200,
      executor: 'per-vu-iterations',
      iterations: 1,
      maxDuration: '10m',
    },
  },
};

// 模拟用户行为
export default function () {
  var token;
  var list;

  group('signin', () => {

    // 根据本次测试中虚拟用户序号获取测试用户
    const user = users[vu.idInTest - 1];
    const data = {
      username: user,
      password: '******',
    }

    let res = http.post(`${host}/**/oauth/token`, JSON.stringify(data), {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json;charset=UTF-8',
      },
      tags: { name: 'grantToken' },
    });
    check(res, { 'status was 200': (r) => r.status == 200 });
    const { accessToken } = JSON.parse(res.body);
    token = accessToken;
    sleep(1);
  });

  group('home', () => {
    let res = http.get(`${host}/**/list`, {
      headers: {
        'Accept': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      tags: { name: 'list' },
    });
    check(res, { 'list was 200': (r) => r.status == 200 });
    list = JSON.parse(res.body);
    sleep(1);
  });

  group('start', () => {
    list.forEach(({ dataType }) => {
      let res = http.post(`${host}/**/user-data-list`, JSON.stringify({ dataType }), {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json;charset=UTF-8',
          'Authorization': `Bearer ${token}`,
        },
        tags: { name: 'user-data-list' },
      });
      check(res, { 'user-data-list was 200': (r) => r.status == 200 });
      const { id: objId } = JSON.parse(res.body);
      sleep(1);

      let lastId = 0;
      let hasNext;
      do {
        let res = http.get(`${host}/**/user-data-list/${objId}/apione?lastId=${lastId}`, {
          headers: {
            'Accept': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          tags: { name: 'apione' },
        });
        check(res, { 'apione was 200': (r) => r.status == 200 });
        APIOneTrendRTT.add(res.timings.duration);
        APIOneRateContentOK.add(res.status === 200);
        if (res.status != 200) {
          console.log(token, objId, lastId, res.status, res.body);
        } else {
          const details = JSON.parse(res.body);
          lastId = details.lastId;
          hasNext = details.hasNext;
          sleep(1);

          const submitData = dataType !== 6 ? [1] : [1, 1];
          details.data.forEach(({ id, dataType, subList }) => {
            if (dataType <= 3) {
              let res = http.post(`${host}/**/user-data-list/apitwo/${id}`, JSON.stringify({ submitData }), {
                headers: {
                  'Accept': 'application/json',
                  'Content-Type': 'application/json;charset=UTF-8',
                  'Authorization': `Bearer ${token}`,
                },
                tags: { name: 'apitwo' },
              });
              lastId = id;
              check(res, { 'apitwo was 200': (r) => r.status == 200 });
              if (res.status != 200) {
                console.log(id, res.status, res.body);
              }
              APITwoTrendRTT.add(res.timings.duration);
              APITwoRateContentOK.add(res.status == 200);
              // 等 1s
              sleep(1);
            } else {
              subList.forEach(({ id }) => {
                let res = http.post(`${host}/**/user-data-list/apitwo/${id}`, JSON.stringify({ submitData }), {
                  headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json;charset=UTF-8',
                    'Authorization': `Bearer ${token}`,
                  },
                  tags: { name: 'apitwo' },
                });
                lastId = id;
                check(res, { 'apitwo was 200': (r) => r.status == 200 });
                if (res.status != 200) {
                  console.log(id, res.status, res.body);
                }
                APITwoTrendRTT.add(res.timings.duration);
                APITwoRateContentOK.add(res.status == 200);
                // 等 1s
                sleep(1);
              });
            }
          });
        }
      } while (hasNext);

      // 等 3s
      sleep(3);
    });
  });
}
```

### 运行自动化测试脚本

```bash
# 首先生成数据准备文件
$ node genseed.mjs
  seed.sql and data.json generated...

# 通过任意数据库管理工具执行 seed.sql 插入测试数据
# 然后运行脚本，等待片刻
$ k6 run --console-output "loadtest.log" script.js
```

上面的测试脚本模拟了 200 个虚拟测试用户，同时登录网站并进行频繁核心 API 调用，我们的项目里，实际上这200个用户的大量的核心操作是在3分钟内完成的，每秒的请求处理数（http_reqs）在 150 左右。下面是其中一次的测试报告，可以重点关注 api_one_ok、api_one_rtt、api_two_ok、api_two_rtt、http_reqs 等几个指标，其他指标的解释说明在 k6 官方文档中都有介绍。

```bash
          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: script.js
     output: -

  scenarios: (100.00%) 1 scenario, 200 max VUs, 10m30s max duration (incl. graceful stop):
           * contacts: 1 iterations for each of 200 VUs (maxDuration: 10m0s, gracefulStop: 30s)


     █ signin

       ✓ status was 200

     █ home

       ✓ list was 200

     █ startTest

       ✓ user-data-list was 200
       ✓ apione was 200
       ✓ apitwo was 200

     checks.........................: 100.00% ✓ 28600      ✗ 0    
     data_received..................: 17 MB   87 kB/s
     data_sent......................: 8.2 MB  43 kB/s
     group_duration.................: avg=1m3s      min=1s       med=13.84s    max=2m53s       p(95)=2m52s      p(99)=2m52s      count=600  
     http_req_blocked...............: avg=72.39µs   min=1.07µs   med=2.25µs    max=18.78ms     p(95)=5.31µs     p(99)=337.86µs   count=28600
     http_req_connecting............: avg=69.22µs   min=0s       med=0s        max=18.71ms     p(95)=0s         p(99)=299.9µs    count=28600
     http_req_duration..............: avg=196.57ms  min=3.2ms    med=50.43ms   max=19.81s      p(95)=293.99ms   p(99)=2.96s      count=28600
       { expected_response:true }...: avg=196.57ms  min=3.2ms    med=50.43ms   max=19.81s      p(95)=293.99ms   p(99)=2.96s      count=28600
     http_req_failed................: 0.00%   ✓ 0          ✗ 28600
     http_req_receiving.............: avg=37.41µs   min=10.33µs  med=31.41µs   max=21.34ms     p(95)=64.08µs    p(99)=124.86µs   count=28600
     http_req_sending...............: avg=29.98µs   min=4.62µs   med=14.3µs    max=2.68ms      p(95)=73.18µs    p(99)=323.57µs   count=28600
     http_req_tls_handshaking.......: avg=0s        min=0s       med=0s        max=0s          p(95)=0s         p(99)=0s         count=28600
     http_req_waiting...............: avg=196.5ms   min=3.16ms   med=50.38ms   max=19.81s      p(95)=293.96ms   p(99)=2.96s      count=28600
     http_reqs......................: 28600   148.881731/s
     iteration_duration.............: avg=3m9s      min=3m4s     med=3m9s      max=3m12s       p(95)=3m11s      p(99)=3m12s      count=200  
     iterations.....................: 200     1.041131/s
   ✓ api_one_ok.....................: 100.00% ✓ 3000       ✗ 0    
   ✗ api_one_rtt....................: avg=132.02598 min=6.582471 med=76.277697 max=1215.379137 p(95)=445.060647 p(99)=978.078818 count=3000 
   ✓ api_two_ok.....................: 100.00% ✓ 24000      ✗ 0    
   ✗ api_two_rtt....................: avg=74.208904 min=5.349855 med=48.139855 max=1221.654595 p(95)=175.69172  p(99)=712.478799 count=24000
     vus............................: 7       min=7        max=200
     vus_max........................: 200     min=200      max=200


running (03m12.1s), 000/200 VUs, 200 complete and 0 interrupted iterations
contacts ✓ [======================================] 200 VUs  03m12.1s/10m0s  200/200 iters, 1 per VU
ERRO[0194] thresholds on metrics 'api_one_rtt, api_two_rtt' have been crossed
```