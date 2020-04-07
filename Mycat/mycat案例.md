## mycat 案例

+ 水平分库 + 读写分离
    ```
    <?xml version="1.0"?>
    <!DOCTYPE mycat:schema SYSTEM "schema.dtd">
    <mycat:schema xmlns:mycat="http://io.mycat/">
	<!--逻辑库表-->
	<schema name="account" checkSQLschema="false" sqlMaxLimit="100">
		<table name="user" dataNode="dn1,dn2" rule="mod-long" />
	</schema>
	<!--数据节点-->
	<dataNode name="dn1" dataHost="cluster1" database="account" />
	<dataNode name="dn2" dataHost="cluster2" database="account" />
	<!--Host节点-->
	<dataHost name="cluster1" maxCon="1000" minCon="10" balance="3"
			  writeType="1" dbType="mysql" dbDriver="native" switchType="1"  slaveThreshold="100">
		<heartbeat>select user()</heartbeat>
		<!-- can have multi write hosts -->
		<writeHost host="hostM1" url="192.168.31.220:3306" user="root"
			password="123">
			<readHost host="hostM2" url="192.168.31.220:3307" user="root"
				password="123"/>
		</writeHost>
	</dataHost>
	<dataHost name="cluster2" maxCon="1000" minCon="10" balance="3"
			  writeType="1" dbType="mysql" dbDriver="native" switchType="1"  slaveThreshold="100">
		<heartbeat>select user()</heartbeat>
		<!-- can have multi write hosts -->
		<writeHost host="hostM3" url="192.168.31.220:3308" user="root"
			password="123">
			<readHost host="hostM4" url="192.168.31.220:3309" user="root"
				password="123"/>
		</writeHost>
	</dataHost>
    </mycat:schema>
    ```
    + [hostM1(主) , hostM2(从)] -- [hostM3(主) , hostM4(从)]
    + 按mod-long规则，2个数据节点 dn1 dn2

+ 读写分离
    ```
    <?xml version="1.0"?>
    <!DOCTYPE mycat:schema SYSTEM "schema.dtd">
    <mycat:schema xmlns:mycat="http://io.mycat/">
	<!--配置数据表-->
	<schema name="account" checkSQLschema="false" sqlMaxLimit="100" dataNode="defaultDN">
	</schema>
	<!--配置分片关系-->
	<dataNode name="dn1" dataHost="cluster1" database="account" />
	<!--配置连接信息-->
	<dataHost name="cluster1" maxCon="1000" minCon="10" balance="3"
			  writeType="1" dbType="mysql" dbDriver="native" switchType="1"  slaveThreshold="100">
		<heartbeat>select user()</heartbeat>
		<!-- can have multi write hosts -->
		<writeHost host="hostM1" url="192.168.31.220:3306" user="root"
			password="123">
			<readHost host="hostM2" url="192.168.31.220:3307" user="root"
				password="123"/>
            <readHost host="hostM3" url="192.168.31.220:3308" user="root"
			password="123">
			<readHost host="hostM4" url="192.168.31.220:3309" user="root"
				password="123"/>
		</writeHost>
	</dataHost>
    </mycat:schema>
    ```
    + 一主三从，M1为主库，M2 M3 M4为从库

+ 水平分库 + ER表
    ```
    <?xml version="1.0"?>
    <!DOCTYPE mycat:schema SYSTEM "schema.dtd">
    <mycat:schema xmlns:mycat="http://io.mycat/">
	<!--配置数据表-->
	<schema name="account" checkSQLschema="false" sqlMaxLimit="100">
		<table name="user" primaryKey="id" dataNode="dn1,dn2" rule="mod-long">
			<childTable name="user_info" primaryKey="id" joinKey="user_id" parentKey="id">
			</childTable>
		</table>
	</schema>
	<!--配置分片关系-->
	<dataNode name="dn1" dataHost="cluster1" database="account" />
	<dataNode name="dn2" dataHost="cluster2" database="account" />
	<!--配置连接信息-->
	<dataHost name="cluster1" maxCon="1000" minCon="10" balance="3"
			  writeType="1" dbType="mysql" dbDriver="native" switchType="1"  slaveThreshold="100">
		<heartbeat>select user()</heartbeat>
		<!-- can have multi write hosts -->
		<writeHost host="hostM1" url="192.168.31.220:3306" user="root"
			password="123">
			<readHost host="hostM2" url="192.168.31.220:3307" user="root"
				password="123"/>
		</writeHost>
	</dataHost>
		<dataHost name="cluster2" maxCon="1000" minCon="10" balance="3"
			  writeType="1" dbType="mysql" dbDriver="native" switchType="1"  slaveThreshold="100">
		<heartbeat>select user()</heartbeat>
		<!-- can have multi write hosts -->
		<writeHost host="hostM3" url="192.168.31.220:3308" user="root"
			password="123">
			<readHost host="hostM4" url="192.168.31.220:3309" user="root"
				password="123"/>
		</writeHost>
	</dataHost>
    </mycat:schema>
    ```
    + 在childTable中，parentKey和JoinKey是对应关系
