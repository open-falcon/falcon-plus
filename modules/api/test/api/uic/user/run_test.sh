if [ "$1" == "create" ]; then
  echo "user_create_test.go"
  go test -v user_create_test.go -test.run TestUserCreate
  go test -v user_create_test.go -test.run TestSingUpDisableCreateUser
fi

echo "user_login_test.go"
go test -v user_login_test.go -test.run TestUserLogin
echo "user_session_test.go"
go test -v user_session_test.go -test.run TestUserSession
go test -v user_session_test.go -test.run TestUserGetInfo
