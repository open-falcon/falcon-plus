if [ "$1" == "create" ]; then
  echo "team_create_test.go"
  go test -v team_create_test.go -test.run TestTeamCreate
fi

echo "team_info_test.go"
go test -v team_info_test.go -test.run TestTeamGetInfo
echo "team_update_test.go"
go test -v team_update_test.go -test.run TestTeamUpdate
echo "team_delete_test.go"
go test -v team_delete_test.go -test.run TestDeleteTeam
