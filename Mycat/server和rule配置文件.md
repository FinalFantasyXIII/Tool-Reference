## server.xml & rule.xml

+ server.xml
    + user 标签
        ```
        <user name="test">
        <property name="password">test</property>
        <property name="schemas">TESTDB</property>
        <property name="readOnly">true</property>
        </user>
        ```
        + 这个标签主要用于定义登录mycat的用户和权限
        + schemas 表示用户可以访问到的逻辑库

    + system 标签
        > 这个标签内嵌套的所有property标签都与系统配置有关

        + defaultSqlParser 属性 : 这个属性用来指定默认的解析器

        + processors 属性 : 指定系统可用的线程数

        + processorBufferChunk 属性 : 指定每次分配Socket Direct Buffer的大小，默认是4096个字节

        + processorBufferPool 属性 : 指定bufferPool计算 比例值

        + processorBufferLocalPercent 属性 : 前面提到了ThreadLocalPool。这个属性就是用来控制分配这个pool的大小用的，但其也并不是一个准确的值，也是一个比例值。这个属性默认值为100

        + processorExecutor 属性 : 指定NIOProcessor上共享的businessExecutor固定线程池大小

        + sequnceHandlerType 属性 : 指定使用Mycat全局序列的类型。0为本地文件方式，1为数据库方式

        + 服务相关属性
            + bindIp : mycat服务监听的IP地址，默认值为0.0.0.0
            + serverPort : 定义mycat的使用端口，默认值为8066
            + managerPort : 定义mycat的管理端口，默认值为9066

+ rule.xml
    ```
    rule.xml里面就定义了我们对表进行拆分所涉及到的规则定义
    文件里面主要有tableRule和function这两个标签。在具体使用过程中可以按照需求添加tableRule和function
    ```
    + tableRule 标签
        ```
        <tableRule name="rule1">
            <rule>
                <columns>id</columns>
                <algorithm>func1</algorithm>
            </rule>
        </tableRule>
        ```
        + name 属性指定唯一的名字，用于标识不同的表规则
        + rule标签则指定对物理表中的哪一列进行拆分和使用什么路由算法
        + columns 内指定要拆分的列名字
        + algorithm 使用function标签中的name属性，连接表规则和具体路由算法。

    + function 标签
        ```
        <function name="hash-int" class="org.opencloudb.route.function.PartitionByFileMap">
            <property name="mapFile">partition-hash-int.txt</property>
        </function>
        ```
        + name 指定算法的名字
        + class 制定路由算法具体的类名字
        + property 为具体算法需要用到的一些属性
        
        
