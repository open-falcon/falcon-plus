cf=$PWD
if [[ $# -eq 1 ]]; then
  cf=$1
fi

echo "working dir: $cf"
for entry in `ls ${cf}/*.sql`
do
  cat $entry | mysql -h 0.0.0.0 -u root --password=password
done
