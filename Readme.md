通过GO的x.509库和PKI库构建一个易于管理和部署的企业内部证书颁发系统。
 1. 初始化根证书
 2. 创建中间证书 (自签发证书，如果再多一级中间证书的话分发会更麻烦。所以功能去掉了)
 3. 颁发Web服务器证书
 4. 颁发个人证书
 5. 代码签名证书
 6. 吊销证书 (提供空的吊销证书列表)
 7. 格式转换 (尚未实现)


Start OpenLDAP
docker run -d -p 389:389 -p 636:636 --env LDAP_ORGANISATION="easyca" --env LDAP_DOMAIN="easyca.local" --env LDAP_ADMIN_PASSWORD="pass@word1"     --name openldap     --hostname openldap-host   --network bridge     osixia/openldap
