# Go client for Scalingo API

Generate the mocks with:

```shell
for interface in $(grep --extended-regexp --no-message --no-filename "type (.*Service|API|TokenGenerator) interface" ./* | grep -v  mockgen | cut -d " " -f 2)
do
  mockgen -destination scalingomock/gomock_$(echo $interface | tr '[:upper:]' '[:lower:]').go -package scalingomock github.com/Scalingo/go-scalingo $interface
done
```
