#!/bin/sh
commit=$(git log -1 --pretty=%h)
output='./git.go'
cat <<EOF > $output
package falcon
const (
	COMMIT = "$commit"
	CHANGELOG = \`
EOF

git log --format='* %cd %aN%n- (%h) %s%d%n' --date=local | grep 'feature\|bugfix\|change' | sed 's/[0-9]+:[0-9]+:[0-9]+ //' >> $output

cat <<'EOF' >> $output
`
)
EOF
