# Skill

## 项目介绍

基于 golang 的 skill 第三方开源库，用于搭建 agent 应用开发中的 skill 模块，让 skill 不止在本地使用，更可以自己开发代码应用于生产项目。

## 构成

概念基本符合 claude code 的理念

> resources

* reference：额外的参考文献
* script：需要用的脚本
* asset：暂时未知有什么用

> skill

* 大模型真正使用的技能，技能中包含使用说明、resources

## ext

* 契合 cloudwego 的 eino 框架：
  * 将 skill 封装成 eino 的 tool 概念，一共有3类 tool，分别是：
    * {skill_name}()：该 tool 的 name 是 {skill_name} , desc 是 skill 的 desc ，result 是 skill 的 body【举例：创建了一个名为 skill_create 的 skill ，它包含有 references：workflows、output_patterns；scripts：init_skill、quick_validate，那么这个tool 就叫做 skill_create】。
    * read_reference(reference_name string)：该 tool 传入 reference_name ，result 则为 reference 的 body 【举例：传入skill_create/workflows，result 为 skill_create 包含的参考文献 workflows 的 body】
    * use_script(script_name string, args string) string：该 tool 传入 script_name 和参数，执行脚本并返回脚本的执行结果【举例：传入 skill_create/init_skill 和 `{"skill_name":"skill_create"}`,执行skill初始化的脚本，返回初始化的结果的json格式字符串】