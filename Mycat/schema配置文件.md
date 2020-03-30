## Schema.xml

+ scheam 标签
    ```
    <schema name="TESTDB" checkSQLschema="false" sqlMaxLimit="100">
    ...
    </schema>

    schema 标签用于定义MyCat实例中的逻辑库，MyCat可以有多个逻辑库，每个逻辑库都有自己的相关配置。可以使用 schema 标签来划分这些不同的逻辑库。
    如果不配置 schema 标签，所有的表配置，会属于同一个默认的逻辑库。
    ```
    + schema 属性
        + dataNode
            ```
            该属性用于绑定逻辑库到某个具体的database上，如果定义了这个属性，那么这个逻辑库就不能工作在分库分表模式下了。也就是说对这个逻辑库的所有操作会直接作用到绑定的dataNode上，这个schema就可以用作读写分离和主从切换
            <schema name="USERDB" checkSQLschema="false" sqlMaxLimit="100" dataNode="dn1"
                <!—这里不能配置任何逻辑表信息-->
            </schema>
            那么现在USERDB就绑定到dn1所配置的具体database上，可以直接访问这个database。当然该属性只能配置绑定到一个database上，不能绑定多个dn
            ```
        + checkSQLschema
            ```
            当该值设置为 true 时，如果我们执行语句select * from TESTDB.travelrecord;则MyCat会把语句修改为select * from
            travelrecord;。即把表示schema的字符去掉，避免发送到后端数据库执行时报（ERROR 1146 (42S02): Table ‘testdb.travelrecord’ doesn’t exist）。
            不过，即使设置该值为 true ，如果语句所带的是并非是schema指定的名字，例如：select * from db1.travelrecord; 那么
            MyCat并不会删除db1这个字段，如果没有定义该库的话则会报错，所以在提供SQL语句的最好是不带这个字段。
            ```
        + sqlMaxLimit : 为不带limit操作的查询语句提供limit操作
            > 需要注意的是，如果运行的schema为非拆分库的，那么该属性不会生效。需要手动添加limit语句。

    + table 标签
        ```
        <table name="travelrecord" dataNode="dn1,dn2,dn3" rule="auto-sharding-long" ></table>
        ```
        + name 属性 : 定义逻辑表的表名，这个名字就如同在数据库中执行create table命令指定的名字一样，同个schema标签中定义的名字必须唯一

        + dataNode 属性
            ```
            定义这个逻辑表所属的dataNode, 该属性的值需要和dataNode标签中name属性的值相互对应。如果需要定义的dn过多可以使用如下的方法减少配置：
            <table name="travelrecord" dataNode="multipleDn$0-99,multipleDn2$100-199" rule="auto-sharding-long" ></table>
            <dataNode name="multipleDn" dataHost="localhost1" database="db$0-99" ></dataNode>
            <dataNode name="multipleDn2" dataHost="localhost1" database=" db$0-99" ></dataNode>
            这里需要注意的是database属性所指定的真实database name需要在后面添加一个，例如上面的例子中，我需要在真实的mysql上建立名称为dbs0到dbs99的database。
            ```
        + rule 属性 : 该属性用于指定逻辑表要使用的规则名字，规则名字在rule.xml中定义，必须与tableRule标签中name属性属性值一一对应

        + ruleRequired 属性 : 该属性用于指定表是否绑定分片规则，如果配置为true，但没有配置具体rule的话 ，程序会报错

        + primaryKey 属性 : 该逻辑表对应真实表的主键,主要是缓存相关，建议配置

        + type 属性 : 值设为gobal的话就是全局表，默认普通表

        + autoIncrement 属性 : true，false

        + needAddLimit属性 : limit相关，默认为true,false禁用

    + childTable 标签
        ```
        childTable标签用于定义E-R分片的子表。通过标签上的属性与父表进行关联
        ```
        + name 属性 : 定义子表的表名

        + joinKey 属性 : 插入子表的时候会使用这个列的值查找父表存储的数据节点

        + parentKey 属性
            ```
            属性指定的值一般为与父表建立关联关系的列名。程序首先获取joinkey的值，再通过parentKey属性指定的列名产生查询语句，通过执行该语句得到父表存储在哪个分片上。从而确定子表存储的位置。
            ```
        + primaryKey 属性 : 同name属性

        + needAddLimit 属性 : 同name属性


+ dataNode标签
    ```
    <dataNode name="dn1" dataHost="lch3307" database="db1" ></dataNode>

    dataNode 标签定义了MyCat中的数据节点，也就是我们通常说所的数据分片。一个 dataNode 标签就是一个独立的数据分片。
    例子中所表述的意思为：使用名字为lch3307数据库实例上的db1物理数据库，这就组成一个数据分片，最后，我们使用名字dn1标识这个分片。
    ```
    + name 属性 :  定义数据节点的名字，这个名字需要是唯一的，我们需要在table标签上应用这个名字，来建立表与分片对应的关系

    + dataHost 属性 : 该属性用于定义该分片属于哪个数据库实例的，属性值是引用dataHost标签上定义的name属性

    + database 属性 : 该属性用于定义该分片属性哪个具体数据库实例上的具体库，因为这里使用两个纬度来定义分片，就是：实例+具体的库。因为每个库上建立的表和表结构是一样的。所以这样做就可以轻松的对表进行水平拆分

+ dataHost 标签
    ```
    作为Schema.xml中最后的一个标签，该标签在mycat逻辑库中也是作为最底层的标签存在，直接定义了具体的数据库实例、读写分离配置和心跳语句
    <dataHost name="localhost1" maxCon="1000" minCon="10" balance="0" writeType="0" dbType="mysql" dbDriver="native">
        <heartbeat>select user()</heartbeat>
        <!-- can have multi write hosts -->
        <writeHost host="hostM1" url="localhost:3306" user="root" password="123456">
            <!-- can have multi read hosts -->
            <!-- <readHost host="hostS1" url="localhost:3306" user="root" password="123456"/> -->
        </writeHost>
        <!-- <writeHost host="hostM2" url="localhost:3316" user="root" password="123456"/> -->
    </dataHost>
    ```
    + name 属性 : 唯一标识dataHost标签，供上层的标签使用

    + maxCon 属性 : 指定每个读写实例连接池的最大连接。也就是说，标签内嵌套的writeHost、readHost标签都会使用这个属性的值来实例化出连接池的最大连接数

    + minCon 属性 : 指定每个读写实例连接池的最小连接，初始化连接池的大小

    + balance 属性
        ```
        负载均衡类型，目前的取值有3种：
        1. balance=“0”, 所有读操作都发送到当前可用的writeHost上
        2. balance=“1”，所有读操作都随机的发送到readHost
        3. balance=“2”，所有读操作都随机的在writeHost、readhost上分发
        ```
    + writeType 属性
        ```
        负载均衡类型，目前的取值有3种：
        1. writeType=“0”, 所有写操作都发送到可用的writeHost上。
        2. writeType=“1”，所有写操作都随机的发送到readHost。
        3. writeType=“2”，所有写操作都随机的在writeHost、readhost分上发。
        ```
    + dbType 属性 : 指定后端连接的数据库类型，目前支持二进制的mysql协议，还有其他使用JDBC连接的数据库。例如：mongodb、oracle、spark等

    + dbDriver 属性 : 指定连接后端数据库使用的Driver，目前可选的值有native和JDBC。使用native的话，因为这个值执行的是二进制的mysql协议，所以可以使用mysql和maridb。其他类型的数据库则需要使用JDBC驱动来支持

    + heartbeat 标签
        ```
        这个标签内指明用于和后端数据库进行心跳检查的语句。例如,MYSQL可以使用select user()，Oracle可以使用select 1 from dual等
        这个标签还有一个connectionInitSql属性，主要是当使用Oracla数据库时，需要执行的初始化SQL语句就这个放到这里面来。例如：alter session set nls_date_format='yyyy-mm-dd hh24:mi:ss'
        ```
    + writeHost标签、readHost标签
        ```
        这两个标签都指定后端数据库的相关配置给mycat，用于实例化后端连接池。唯一不同的是，writeHost指定写实例、readHost指定读实例，组着这些读写实例来满足系统的要求。
        在一个dataHost内可以定义多个writeHost和readHost。但是，如果writeHost指定的后端数据库宕机，那么这个writeHost绑定的所有readHost都将不可用。另一方面，由于这个writeHost宕机系统会自动的检测到，并切换到备用的writeHost上去。
        ```
        + host 属性 : 用于标识不同实例

        + url 属性 : 后端实例连接地址，如果是使用native的dbDriver，则一般为address:port这种形式。用JDBC或其他的dbDriver，则需要特殊指定。当使用JDBC时则可以这么写：jdbc:mysql://localhost:3306/

        + user 属性 : root

        + password 属性 : 123



