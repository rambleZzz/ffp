# ffp
FFP - FastFingerPrint的简称，意为快速WEB指纹识别工具. 一款基于GO语言编写的支持调用httpx、observerWard指纹识别结果、cdn检测、无cdn ip及c段去重统计、ip归属地查询于一体的快速自动化资产指纹识别工具.

## 环境
go1.20开发  
目前已编译MAC amd64位环境 、windows64位环境 、linux64位环境
其它环境请自行编译

## 功能
```
1、首先会对扫描目标进行存在检测，统计所有不存活的目标写入excel结果文件便于查看，不进行接下来的扫描检测。
2、存活的目标首先会调用httpx(调用的原生go代码，不是调用的可执行文件)对目标站点识别，获取以下信息
 {"目标", "标题", "应用服务器","wappalyzer指纹", "网址", "是否有CDN", "A记录", "CNAME值", "主机", "端口", "类型", "状态码", "网站协议", "跳转URL"}

3、调用cdnCheck模块二次确认目标使用CDN情况，获取以下信息
 {"是否使用CDN","CDN名称","CNAME值"}
 
4、调用ObserverWard(可执行文件，需要自行配置路径)，进行web指纹识别。
5、调用iconhash模块进行favicon图标hash识别。
6、对所有无CDN的IP进行IP归属地检测。
7、整理所有无CDN的IP C段地址并去重，方便手动收集到更多的相关资产。
```




## 介绍：
```
./ffp -h


    ________________
   / ____/ ____/ __ \
  / /_  / /_  / /_/ /
 / __/ / __/ / ____/
/_/   /_/   /_/       Ver:1.0

https://github.com/rambleZzz/ffp
FFP v1.0   (FastFingerPrint)  Dev:go1.20


NAME:
   FFP - FastFingerPrint - 一款基于GO语言编写的支持调用httpx、observerWard指纹识别结果、cdn检测、无cdn ip及c段去重统计、ip归属地查询于一体的快速自动化资产指纹识别工具.

USAGE:
   ffp -u example.com (单一目标扫描)
   ffp -f target.txt   (批量目标扫描)

   -noe表示不输出到excel，此参数默认为false,不加此参数会输出,如果-noe不会输出至excel
   -t参数为线程默认为5
   -u 或-f参数为必填项(-u和-f不可同时)，其他参数为可选项，更多参数请参考GLOBAL OPTIONS

VERSION:
   v1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value     从txt文件中读取内容 如:-f target.txt
   --url value, -u value      target参数为需要检测的目标 如:-u example.com
   --threads value, -t value  扫描线程数，如：-t 5 (default: 5)
   --notO2Excel, --noe        是否不输出到excel, 如：-noe 不输出 (default: false)
   --help, -h                 show help
   --version, -v              print the version
 

```  
首次执行会在当前目录自动生成config.yaml配置文件，由于ObserverWard调用的原始文件，不同操作系统请自己自行配置不同的文件路径

```
# AutoScan Yaml Config
thirdparty: # 路径必须配置在当前目录下，不可自定义当前执行文件以外的其他目录
  ObserverWardPath: /thirdparty/observerWard/observer_ward # observer_ward 可执行文件
  ObserverWardDir: /thirdparty/observerWard/ # observerWard 所在目录
  WebFingerprintPath: /thirdparty/observerWard/web_fingerprint_v3.json # fingerprintHub指纹库
  GeoLite2Path: /thirdparty/cdnCheck/GeoLite2-ASN.mmdb # ip相关
 ``` 

## 使用帮助:
#### 1、直接目标扫描模式：  
```
./ffp -u www.example.com

扫描目标格式如下都可，会自动识别判断：  
36.x.x.x
36.x.x.x:8080
36.x.x.x:8080/erp
www.example.com
http://www.example.com
www.example:8080
https://www.example:com:7777/erp

```
#### 2、从文本中批量扫描模式：
````

./ffp -f 1.txt
扫描目标格式如下,直接写入txt文件：
36.x.x.x
36.x.x.x:8080
36.x.x.x:8080/erp
www.example.com
http://www.example.com
www.example:8080
https://www.example:com:7777/erp

````
#### 3、其他参数
```

(请先使用masscan扫描 masscan -p 1-65535 -oJ test.json)
./ffp -f 1.txt -t 10 -noe

-t 代表线程，不填写默认为5
-noe 表示不输出结果到excel，只输出到终端

````
#### 运行截图
、、、

![ffp1.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp1.jpg)   

![ffp2.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp2.jpg)   

![ffp3.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp3.jpg) 

### 结果输出
#### 1、终端输出
运行过程中自动输出到终端
#### 2、excel结果输出
运行完默认输出到excel，默认会输出, 不加任何参数也会输出,如果-noe 则表示不输出至excel
#### 结果截图

![ffp4.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp4.jpg)

![ffp5.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp5.jpg) 

![ffp6.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp6.jpg) 

![ffp7.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp7.jpg) 

![ffp8.png](https://github.com/rambleZzz/ffp/blob/main/images/ffp8.jpg) 

### 参考链接   

https://github.com/projectdiscovery/httpx   
https://github.com/0x727/ObserverWard   
https://github.com/0x727/FingerprintHub   
https://github.com/hanc00l/nemo_go   
https://github.com/Becivells/iconhash   
https://github.com/damit5/cdnCheck_go   
https://github.com/wgpsec/ENScan_GO


