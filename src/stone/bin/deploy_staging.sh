curl --header 'Content-Type: application/json' \
     --header 'Accept: application/json' \
     --header "Authorization: Bearer $1" \
     -d '{"template_id":7,"environment":"{\"TAG\":\"'$2'\"}"}' \
     http://ansible.58wallet.io/api/project/1/tasks