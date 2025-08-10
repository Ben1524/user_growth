# 数据库配置放在环境变量中
````
export USER_GROWTH_CONFIG='{"Db":{"Engine":"mysql","Username":"root","Password":"12345678","Host":"localhost","Port":3306,"Database":"user_growth","Charset":"utf8","ShowSql":true,"MaxIdleConns":2,"MaxOpenConns":10,"ConnMaxLifetime":30},"Cache":{}}'
````