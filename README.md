# RuleIniter

简易动态流媒体解锁规则生成器，通过检测流媒体解锁状态生成*Ray路由规则，适用于V2bX、XrayR或其他XrayR分支。

## 配置文件
请参考仓库中默认配置文件及下方示例修改配置文件
``` json
// 示例用，非标准json
{
    "RoutePath": "/etc/V2bX/route.json", // 路由配置文件所在路径
    "OutTag": "unlock", // OutboundTag, 需自行配置Outbound
    "CheckGlobal": true, // 检查跨国平台
    "CheckRegion": false, // 检查跨国平台地区与指定的地区一致性
    "MediaList": { // 指定不同地区要检测的流媒体平台
        "global": { // 指定跨国平台
			"Hotstar",
			...
		}
        "jp": { // 2位国家标识, 可随意填写, 但会造成CheckRegion失效
			"DMM", // 流媒体平台标识
			"Abema",
			...
		},
		...
    },
    "MatchRuleList": { // 指定流媒体平台的路由匹配规则
        "Abema": {
            "Domain": [ // 域名匹配规则, 与*Ray一致
                "geosite:abema",
                ...
            ],
            "Ip": { //Ip匹配规则, 与*Ray一致
                "114.114.514.514",
                ...
            }
        },
        ...
    }
```
## 运行
从Releases下载可执行文件或自行编译，然后参考下方命令运行
``` bash
./RuleIniter -conf ./config.json -region jp
```

## Thanks
[nkeonkeo/MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest)