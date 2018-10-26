cd /home/ld-sgdev/workspace/open-falcon
tar -xzvf /home/ld-sgdev/zhufeng_li/workspace/src/github.com/open-falcon/falcon-plus/*gz -C .
find . -name cfg.json |xargs sed -i 's/root:@tcp(127.0.0.1:3306)/root:123456@tcp(127.0.0.1:3306)/g'
sed -i   's#"max_conns": 100,#"max_conns": 1000,#' api/config/cfg.json
sed -i   's#"max_idle": 100,#"max_idle": 1000,#' api/config/cfg.json
find . -name '*.json' |xargs sed -i 's#"log_level": "debug"#"log_level": "warn"#g'
find . -name '*.json' |xargs sed -i 's#"debug": true,#"debug": false,#g'
