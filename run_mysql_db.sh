docker run --name mysql-test \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=testdb \
  -e MYSQL_USER=testuser \
  -e MYSQL_PASSWORD=testpass \
  -p 3306:3306 \
  -d mysql:8

sleep 15

docker exec -i mysql-test mysql -utestuser -ptestpass testdb <<EOF
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100)
);
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com');
INSERT INTO users (name, email) VALUES ('Bob', 'bob@example.com');
EOF


docker exec -it mysql-test mysql -utestuser -ptestpass testdb -e "SELECT * FROM users;"
